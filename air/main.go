package air

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"golang.org/x/crypto/ssh/terminal"

	"os"
)

type Report struct {
	Email Email
}

const (
	targetsFileName = "targets.yml"
	reportFileName  = "report.yml"
	filtersFileName = "filters.yml"

	DefaultMaxReportAge = 60
)

type AppConfig struct {
	Debug        bool
	TargetsFile  string
	FiltersFile  string
	ReportFile   string
	ConfigPath   string
	MaxReportAge int
	filters      filters
	targets      targets
	report       Report
	OutputDir    string
}

func (appConfig *AppConfig) load() {
	loaded := appConfig

	loaded.targets = loadTargets(appConfig.ConfigPath, appConfig.Debug)
	loaded.filters = loadFilters(appConfig.ConfigPath, appConfig.Debug)
	loaded.report = loadReportConfig(appConfig.ConfigPath, appConfig.Debug)
	*appConfig = *loaded
}

func (ar *accountsResults) hasFindings() bool {
	for _, r := range *ar {
		for _, rr := range r.regionResults {
			for _, rtr := range rr.regionTemplateResults {
				for _, rtrr := range rtr.runs {
					if len(rtrr.findings) > 0 {
						return true
					}
				}
			}
		}
	}
	return false
}

func Run(appConfig AppConfig) error {
	var err error
	var accountsResults accountsResults
	var initialSess *session.Session
	initialSess, err = session.NewSession()
	if err != nil {
		return err
	}
	var tems targetErrorsMaps
	appConfig.load()
	if len(appConfig.targets) > 0 {
		accountsResults, tems, err = processMultipleAccounts(initialSess, appConfig.targets, appConfig.MaxReportAge)
	} else {
		accountsResults, tems, err = processSingleAccount(initialSess, appConfig.MaxReportAge)
	}
	clearConsoleLine()

	// if we have results and filters defined, then apply filters
	if accountsResults != nil && accountsResults.hasFindings() {
		if len(appConfig.filters) > 0 {
			accountsResults.filter(appConfig.filters)
		}
		// if we still have results, then output to spreadsheet
		var reportPath string
		if len(accountsResults) > 0 {
			reportPath, err = generateSpreadsheet(accountsResults, appConfig.OutputDir)
			if err != nil {
				fmt.Println("failed to generate spreadsheet:", err)
				os.Exit(1)
			}
			if !reflect.DeepEqual(appConfig.report.Email, Email{}) {
				if err = emailReport(initialSess, reportPath, appConfig.report.Email, false); err != nil {
					return err
				}
			}
		}
	} else {
		log.Print("No findings found.")
		fmt.Println("No findings found.")
	}

	if anyTargetErrors(tems) {
		fmt.Printf("Errors encountered during processing...\n\n")
		for _, t := range tems {
			if len(t.errors) > 0 {
				fmt.Printf("Account: %s (%s)\n", t.target.ID, t.target.Alias)
				for _, e := range t.errors {
					fmt.Printf("  Issue: %s\n", e.desc)
					fmt.Printf("  Detail: %s\n", e.err)
				}
			}
		}
	}

	return err
}

func processMultipleAccounts(sess *session.Session, targets targets, maxReportAge int) (accountsResults accountsResults, tems targetErrorsMaps, err error) {
	for _, target := range targets {
		var tem targetErrorsMap
		tem.target = target
		var shortAccountOutput string
		if target.Alias != "" {
			shortAccountOutput = target.Alias
		} else {
			shortAccountOutput = target.ID
		}
		statusOutput := fmt.Sprintf("Processing: [%s]...", shortAccountOutput)
		statusOutput = padToWidth(statusOutput, true)
		width, _, _ := terminal.GetSize(0)
		if len(statusOutput) == width {
			fmt.Printf(statusOutput[0:width-3] + "   \r")
		} else {
			fmt.Print(statusOutput)
		}

		var creds *credentials.Credentials
		creds, err = getAssumeRoleCreds(getAssumeRoleCredsInput{
			Sess:       sess,
			AccountID:  target.ID,
			RoleName:   target.RoleName,
			ExternalID: target.RoleExternalID,
		})
		if err != nil {
			aErr := annotatedError{
				err:  err,
				desc: fmt.Sprintf("failed to assume role: %s", genRoleArn(target.ID, target.RoleName)),
			}
			tem.errors = append(tem.errors, aErr)
			if isUnrecoverable(err) {
				tems = append(tems, tem)
				continue
			}
		}
		var accountOutput accountResults
		accountOutput.accountID = target.ID
		accountOutput.accountAlias = target.Alias
		var perRegionResults []regionResult
		inspectorRegions := getAllInspectorRegions()

		perRegionResults, err = processAllRegions(creds, inspectorRegions, maxReportAge)
		if err != nil {
			aErr := annotatedError{
				err:  err,
				desc: fmt.Sprintf("failed to get region results"),
			}
			tem.errors = append(tem.errors, aErr)
			if isUnrecoverable(err) {
				tems = append(tems, tem)
				continue
			}
		}
		accountOutput.regionResults = perRegionResults
		accountsResults = append(accountsResults, accountOutput)
		tems = append(tems, tem)
	}
	return accountsResults, tems, err
}

func processSingleAccount(sess *session.Session, maxReportAge int) (accountsResults accountsResults, tems targetErrorsMaps, err error) {
	inspectorRegions := getAllInspectorRegions()
	var tem targetErrorsMap
	svc := iam.New(sess)
	stsSvc := sts.New(sess)
	accountID := getAccountID(stsSvc)
	accountAlias := getAccountAlias(svc)
	sessCreds, err := sess.Config.Credentials.Get()
	if err != nil {
		os.Exit(1)
	}
	var accountOutput accountResults
	accountOutput.accountID = accountID
	accountOutput.accountAlias = accountAlias
	var shortAccountOutput string
	if accountAlias != "" {
		shortAccountOutput = accountAlias
	} else {
		shortAccountOutput = accountID
	}
	var perRegionResults []regionResult
	creds := credentials.NewStaticCredentials(sessCreds.AccessKeyID,
		sessCreds.SecretAccessKey, sessCreds.SessionToken)
	statusOutput := fmt.Sprintf("Processing: [%s]...", shortAccountOutput)
	statusOutput = padToWidth(statusOutput, true)
	width, _, _ := terminal.GetSize(0)
	if len(statusOutput) == width {
		fmt.Printf(statusOutput[0:width-3] + "   \r")
	} else {
		fmt.Print(statusOutput)
	}

	perRegionResults, err = processAllRegions(creds, inspectorRegions, maxReportAge)
	if err != nil {
		aErr := annotatedError{
			err:  err,
			desc: fmt.Sprintf("failed to get region results"),
		}
		tem.errors = append(tem.errors, aErr)
		if isUnrecoverable(err) {
			tems = append(tems, tem)
		}
	}
	accountOutput.regionResults = perRegionResults
	accountsResults = append(accountsResults, accountOutput)
	tems = append(tems, tem)
	return accountsResults, tems, err
}

func anyTargetErrors(tems targetErrorsMaps) bool {
	for _, t := range tems {
		if len(t.errors) > 0 {
			return true
		}
	}
	return false
}

func isUnrecoverable(err error) bool {
	switch {
	case strings.HasPrefix(err.Error(), "ExpiredToken"):
		return true
	case strings.Contains(err.Error(), "NoCredentialProviders"):
		return true
	}
	return false
}

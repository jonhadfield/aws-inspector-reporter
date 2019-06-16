package air

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/aws/aws-sdk-go/service/inspector/inspectoriface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/snwfdhmp/errlog"
)

type run struct {
	runArn   string
	findings findings
}

type regionTemplateResult struct {
	templateArn  string
	templateName string
	runs         []run
}

type regionTemplateResults []regionTemplateResult

type regionResult struct {
	region                string
	regionTemplateResults []regionTemplateResult
}

type accountResults struct {
	accountID     string
	accountAlias  string
	regionResults []regionResult
}

type accountsResults []accountResults

func getAssessmentTargetsArns(svc inspectoriface.InspectorAPI) ([]*string, error) {
	var err error
	lato := &inspector.ListAssessmentTargetsOutput{}
	var assessmentTargetArns []*string

	for {
		lati := &inspector.ListAssessmentTargetsInput{
			MaxResults: ptrToInt64(10),
			NextToken:  lato.NextToken,
		}
		lato, err = svc.ListAssessmentTargets(lati)

		if err != nil {
			return assessmentTargetArns, err
		}
		assessmentTargetArns = append(assessmentTargetArns, lato.AssessmentTargetArns...)
		if lato.NextToken == nil {
			break
		}
	}
	return assessmentTargetArns, err
}

func getAssessmentTemplatesArns(svc inspectoriface.InspectorAPI, targetArns []*string) ([]*string, error) {
	var err error

	var assessmentTemplateArns []*string
	lato := &inspector.ListAssessmentTemplatesOutput{}
	for {
		lati := &inspector.ListAssessmentTemplatesInput{
			AssessmentTargetArns: targetArns,
			MaxResults:           ptrToInt64(10),
			NextToken:            lato.NextToken,
		}
		lato, err = svc.ListAssessmentTemplates(lati)

		if errlog.Debug(err) { // will debug & pass if err != nil, will ignore if err == nil
			return assessmentTemplateArns, err
		}
		assessmentTemplateArns = append(assessmentTemplateArns, lato.AssessmentTemplateArns...)
		if lato.NextToken == nil {
			break
		}
	}
	return assessmentTemplateArns, err
}

type annotatedError struct {
	err  error
	desc string
}

func processAllRegions(creds *credentials.Credentials, inspectorRegions []string, maxReportAge int) (results []regionResult, err error) {
	GetRegionResults := func(ctx context.Context, regions []string) ([]regionResult, error) {
		g, ctx := errgroup.WithContext(ctx)

		perRegionResults := make([]regionResult, len(regions))
		for i, region := range regions {
			i := i
			region := region
			g.Go(func() error {
				// TODO: catch errors
				var rtr regionResult
				rtr.region = region
				var sess *session.Session
				sess, err = session.NewSession(&aws.Config{Credentials: creds, Region: &region})
				var regionTemplateResults regionTemplateResults
				svc := inspector.New(sess)
				regionTemplateResults, err = getRegionTemplateResults(svc, maxReportAge)
				rtr.regionTemplateResults = append(rtr.regionTemplateResults, regionTemplateResults...)
				if err == nil {
					perRegionResults[i] = rtr
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return nil, err
		}
		return perRegionResults, nil
	}

	results, err = GetRegionResults(context.Background(), inspectorRegions)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	return results, err
}

func getLatestAssessmentTemplateRuns(svc inspectoriface.InspectorAPI, templateArns []*string) ([]*string, error) {
	laro := &inspector.ListAssessmentRunsOutput{}
	var err error
	var assessmentRunArns []*string
	for {
		lari := &inspector.ListAssessmentRunsInput{
			AssessmentTemplateArns: templateArns,
			MaxResults:             ptrToInt64(10),
			NextToken:              laro.NextToken,
		}

		laro, err = svc.ListAssessmentRuns(lari)

		if errlog.Debug(err) { // will debug & pass if err != nil, will ignore if err == nil
			return assessmentRunArns, err
		}
		assessmentRunArns = append(assessmentRunArns, laro.AssessmentRunArns...)
		if laro.NextToken == nil {
			break
		}
	}
	return assessmentRunArns, err
}

func getAssessmentRunDetails(svc inspectoriface.InspectorAPI, assessmentRunArns []*string) ([]*inspector.AssessmentRun, error) {
	var assessmentRunDetails []*inspector.AssessmentRun
	var err error
	for i := 0; i <= len(assessmentRunArns)-1; i += 10 {
		var last int
		if i+10 > len(assessmentRunArns) {
			last = len(assessmentRunArns)
		} else {
			last = i + 10
		}
		dardi := &inspector.DescribeAssessmentRunsInput{
			AssessmentRunArns: assessmentRunArns[i:last],
		}
		var dardo *inspector.DescribeAssessmentRunsOutput
		dardo, err = svc.DescribeAssessmentRuns(dardi)
		if err != nil {
			log.Fatal(err)
			return assessmentRunDetails, err
		}
		assessmentRunDetails = append(assessmentRunDetails, dardo.AssessmentRuns...)
	}
	return assessmentRunDetails, err
}

func listFindingArns(svc inspectoriface.InspectorAPI, assRunArn *string) ([]*string, error) {
	var findingsArns []*string
	var err error
	var nextToken *string
	for {
		lfi := &inspector.ListFindingsInput{
			AssessmentRunArns: []*string{assRunArn},
			MaxResults:        ptrToInt64(100),
			NextToken:         nextToken,
		}
		var lfo *inspector.ListFindingsOutput
		lfo, err = svc.ListFindings(lfi)
		if err != nil {
			return findingsArns, err
		}
		findingsArns = append(findingsArns, lfo.FindingArns...)
		if lfo.NextToken != nil {
			nextToken = lfo.NextToken
		} else {
			return findingsArns, err
		}
	}
}

func getRegionTemplateResults(svc inspectoriface.InspectorAPI, maxReportAge int) (results regionTemplateResults, err error) {
	// list assessment targets
	var assTargetArns []*string
	assTargetArns, err = getAssessmentTargetsArns(svc)
	if err != nil || len(assTargetArns) == 0 {
		return
	}
	// Output templates
	var assTemplatesArns []*string
	assTemplatesArns, err = getAssessmentTemplatesArns(svc, assTargetArns)

	for _, assTemplateArn := range assTemplatesArns {
		var result regionTemplateResult
		result.templateArn = *assTemplateArn
		var assessmentRunArns []*string

		aTa := []*string{assTemplateArn}
		// Get latest assessment runs for templates
		assessmentRunArns, err = getLatestAssessmentTemplateRuns(svc, aTa)
		if err != nil {
			return
		}

		if len(assessmentRunArns) == 0 {
			continue
		}

		// Get latest assessment run details
		var assessmentRunsDetails []*inspector.AssessmentRun
		assessmentRunsDetails, err = getAssessmentRunDetails(svc, assessmentRunArns)
		if err != nil {
			return
		}
		type assessmentRunToTime struct {
			runArn string
			name   string
			time   time.Time
		}

		// loop through assessment runs and get the latest for each template
		latestRunPerTemplate := make(map[string]assessmentRunToTime)
		for _, ard := range assessmentRunsDetails {
			// ignore any runs that are older than max report age
			timeNow := time.Now()
			timeSinceReportCompleted := timeNow.Sub(*ard.CompletedAt)
			timeMaxReportAge := time.Duration(maxReportAge) * (24 * time.Hour)
			if timeSinceReportCompleted > timeMaxReportAge {
				//log.Printf("Ignoring report as time now is: %s, report completed at: %s and that's: %s duration and maxReportAge is: %s", time.Now(), *ard.CompletedAt,
				//	timeSinceReportCompleted, timeMaxReportAge)
				continue
			}
			if ard.CompletedAt != nil && (latestRunPerTemplate[*ard.AssessmentTemplateArn].time.IsZero() ||
				latestRunPerTemplate[*ard.AssessmentTemplateArn].time.Before(*ard.StartedAt)) {
				latestRunPerTemplate[*ard.AssessmentTemplateArn] = assessmentRunToTime{
					time:   *ard.StartedAt,
					runArn: *ard.Arn,
					name:   *ard.Name,
				}
			}
		}
		var latestRunsArns []*string
		for _, v := range latestRunPerTemplate {
			latestRunsArns = append(latestRunsArns, ptrToStr(v.runArn))
		}

		// Get template name
		dtni := inspector.DescribeAssessmentTemplatesInput{
			AssessmentTemplateArns: aTa,
		}
		var dtno *inspector.DescribeAssessmentTemplatesOutput
		dtno, err = svc.DescribeAssessmentTemplates(&dtni)
		if dtno != nil {
			result.templateName = *dtno.AssessmentTemplates[0].Name
		} else {
			result.templateName = "-"
		}

		for _, runArn := range latestRunsArns {
			var resultRun run
			resultRun.runArn = *runArn
			var findingArns []*string
			findingArns, err = listFindingArns(svc, runArn)
			if err != nil {
				log.Fatal(err)
			}

			var findings findings
			findings, err = describeFindings(svc, findingArns)
			resultRun.findings = append(resultRun.findings, findings...)
			result.templateArn = *assTemplateArn

			result.runs = append(result.runs, resultRun)
		}
		results = append(results, result)
	}
	return results, err

}

type findings []finding

func (in *findings) copy(aFs []*inspector.Finding) {
	nFindings := make(findings, 0, len(aFs))
	for _, f := range aFs {
		nFinding := transformFinding(f)
		nFindings = append(nFindings, nFinding)
	}
	*in = nFindings
}

func (in *findings) propagateRulesPackageNames(svc inspectoriface.InspectorAPI) error {
	rulesPackages, err := getRulesPackages(svc)
	updated := make(findings, 0, len(*in))
	for _, f := range *in {
		f.rulePackageName = rulesPackages[*f.ServiceAttributes.RulesPackageArn]
		updated = append(updated, f)
	}
	*in = updated
	return err
}

func getRulesPackages(svc inspectoriface.InspectorAPI) (map[string]string, error) {
	rulesPackagesLookup := make(map[string]string)
	var nextToken string
	var rpArns []*string
	var lrpi inspector.ListRulesPackagesInput
	for {
		if nextToken == "" {
			lrpi = inspector.ListRulesPackagesInput{}
		} else {
			lrpi = inspector.ListRulesPackagesInput{
				NextToken: ptrToStr(nextToken),
			}
		}
		lrpo, err := svc.ListRulesPackages(&lrpi)
		if err != nil {
			return rulesPackagesLookup, err
		}

		rpArns = append(rpArns, lrpo.RulesPackageArns...)
		if lrpo.NextToken != nil && *lrpo.NextToken != "" {
			nextToken = *lrpo.NextToken
		} else {
			break
		}
	}

	dri := &inspector.DescribeRulesPackagesInput{
		RulesPackageArns: rpArns,
	}
	dro, err := svc.DescribeRulesPackages(dri)
	if dro != nil && dro.RulesPackages != nil {
		for _, rp := range dro.RulesPackages {
			rulesPackagesLookup[*rp.Arn] = *rp.Name
		}
	}
	return rulesPackagesLookup, err

}

func describeFindings(svc inspectoriface.InspectorAPI, findingsArns []*string) (findings, error) {
	var err error
	var results findings
	for i := 0; i <= len(findingsArns); i += 100 {
		var last int
		if i+100 > len(findingsArns) {
			last = len(findingsArns)
		} else {
			last = i + 100
		}
		dfi := inspector.DescribeFindingsInput{
			FindingArns: findingsArns[i:last],
		}
		var dfo *inspector.DescribeFindingsOutput
		dfo, err = svc.DescribeFindings(&dfi)
		if err != nil {
			return results, err
		}
		ufs := findings{}
		ufs.copy(dfo.Findings)
		results = append(results, ufs...)
	}
	err = results.propagateRulesPackageNames(svc)
	return results, err
}

func getAllInspectorRegions() (result []string) {
	nonStandardRegions := []string{"cn-north-1", "cn-northwest-1", "us-gov-east-1", "us-gov-west-1"}
	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	resMaps := make([]map[string]endpoints.Region, 0, len(partitions))
	for _, p := range partitions {
		resMap, _ := endpoints.RegionsForService(endpoints.DefaultPartitions(), p.ID(), "inspector")
		resMaps = append(resMaps, resMap)
	}
	keys := make([]string, 0, len(resMaps))
	for _, resMap := range resMaps {
		for _, ra := range resMap {
			if !stringInSlice(ra.ID(), nonStandardRegions) {
				keys = append(keys, ra.ID())
			}
		}
	}
	result = keys
	return
}

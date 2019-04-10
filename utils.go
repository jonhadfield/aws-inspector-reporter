package air

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

func getInstanceName(finding finding) string {
	for _, t := range finding.AssetAttributes.Tags {
		if *t.Key == "Name" {
			return *t.Value
		}
	}
	return "-"
}
func deleteFile(path string) (err error) {
	err = os.Remove(path)
	return
}
func PadToWidth(input, char string, inputLengthOverride int, trimToWidth bool) (output string) {
	// Split string into lines
	var lines []string
	var newLines []string
	if strings.Contains(input, "\n") {
		lines = strings.Split(input, "\n")
	} else {
		lines = []string{input}
	}
	var paddingSize int
	for i, line := range lines {
		width, _, _ := terminal.GetSize(0)
		if width == -1 {
			width = 80
		}
		// No padding for a line that already meets or exceeds console width
		var length int
		if inputLengthOverride > 0 {
			length = inputLengthOverride
		} else {
			length = len(line)
		}

		switch {
		case length >= width:
			if trimToWidth {
				output = line[0:width]
			} else {
				output = input
			}
			return
		case i == len(lines)-1:
			if inputLengthOverride != 0 {
				paddingSize = width - inputLengthOverride
			} else {
				paddingSize = width - len(line)
			}
			if paddingSize >= 1 {
				newLines = append(newLines, fmt.Sprintf("%s%s\r", line, strings.Repeat(char, paddingSize)))
			} else {
				newLines = append(newLines, fmt.Sprintf("%s\r", line))
			}
		default:
			var suffix string
			newLines = append(newLines, fmt.Sprintf("%s%s%s\n", line, strings.Repeat(char, paddingSize), suffix))

		}
	}
	return strings.Join(newLines, "")
}
func OutputError(err error) {
	output := PadToWidth(fmt.Sprintf("error: %v\n", err), " ", 0, false)
	_, _ = fmt.Fprintf(os.Stderr, output)
}

func getAccountID(sess *session.Session) (id string) {
	stsSvc := sts.New(sess)
	callerID, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	credsNotFoundMessage := "credentials not found\nsee: https://docs.aws.amazon.com/cli/" +
		"latest/userguide/cli-chap-getting-started.html#cli-quick-configuration"

	switch {
	case err != nil:
		if awsErr, okBPA2 := errors.Cause(err).(awserr.Error); okBPA2 {
			if strings.Contains(awsErr.Message(), "non-User credentials") {
				// not using user creds, so need to try a different method
			} else if awsErr.Code() == "NoCredentialProviders" {
				err = errors.New(credsNotFoundMessage)
				OutputError(err)
				os.Exit(1)
			} else if awsErr.Code() == "ExpiredToken" {
				err = errors.New("temporary credentials have expired")
				OutputError(err)
				os.Exit(1)
			} else if strings.Contains(awsErr.Message(), "security token included in the request is invalid") {
				err = errors.New("specified credentials have an invalid security token")
				OutputError(err)
				os.Exit(1)
			} else {
				fmt.Println(fmt.Sprintf("unhandled exception using specified credentials: %s", awsErr.Message()))
			}
		}
	case callerID.Arn == nil:
		err = errors.New("credentials not found\nsee: https://docs.aws.amazon.com/cli/" +
			"latest/userguide/cli-chap-getting-started.html#cli-quick-configuration")
		OutputError(err)
		os.Exit(1)
	default:
		id = *callerID.Account
		return
	}
	return id
}

func getAccountAlias(sess *session.Session) (alias string) {
	// try to get the account alias
	iamSvc := iam.New(sess)
	var getAliasOutput *iam.ListAccountAliasesOutput
	var err error
	getAliasOutput, err = iamSvc.ListAccountAliases(&iam.ListAccountAliasesInput{})
	if err != nil {
		fmt.Println("missing \"iam:ListAccountAliases\" permission so unable to retrieve alias")
	} else if len(getAliasOutput.AccountAliases) > 0 {
		alias = *getAliasOutput.AccountAliases[0]
	}
	return
}
func formatTitle(in string) string {
	result := strings.Split(in, "\n")
	var titleLines []string
	for _, line := range result {
		if len(strings.TrimSpace(line)) > 0 {
			titleLines = append(titleLines, strings.TrimSpace(line))
		}
	}
	return strings.Join(titleLines, "\r\n")
}
func formatDescription(in string) string {
	result := strings.Split(in, "\n")
	var descriptionLines []string
	for _, line := range result {
		// Strip 'Description' prefix
		if strings.HasPrefix(line, "Description") {
			line = line[11:]
		}
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) > 0 {
			descriptionLines = append(descriptionLines, trimmedLine)
		}
	}
	return strings.Join(descriptionLines, "\r\n")
}

func formatRecommendation(in string) string {
	result := strings.Split(in, "\n")
	var recLines []string
	for _, line := range result {
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) > 0 {
			recLines = append(recLines, trimmedLine)
		}
	}
	return strings.Join(recLines, "\r\n")
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func clearConsoleLine() {
	fmt.Printf("%s", PadToWidth("", " ", 0, false))
}

func ptrToInt64(in int64) *int64 {
	return &in
}

func ptrToStr(in string) *string {
	return &in
}

func ptrToTime(in time.Time) *time.Time {
	return &in
}

func ptrToBool(in bool) *bool {
	return &in
}

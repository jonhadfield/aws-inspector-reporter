package air

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func loadFilters(configPath string, debug bool) (filters filters) {
	filtersPath := ensureTrailingSlash(configPath) + filtersFileName
	var err error

	// try loading from s3
	if strings.HasPrefix(configPath, "s3://") {
		// split config path (minus prefix)
		parts := strings.Split(configPath[5:], "/")
		key := ensureTrailingSlash(strings.Join(parts[1:], "/")) + filtersFileName
		sess := session.Must(session.NewSession())
		var region string
		region, err = s3manager.GetBucketRegion(context.Background(), sess, parts[0], "us-east-1")
		sess = session.Must(session.NewSession(&aws.Config{Region: ptrToStr(region)}))
		svc := s3.New(sess)
		input := &s3.GetObjectInput{
			Bucket: aws.String(parts[0]),
			Key:    aws.String(key),
		}
		var goo *s3.GetObjectOutput
		goo, err = svc.GetObject(input)
		if err != nil && debug {
			fmt.Printf("failed to load filters from s3://%s/%s\n", parts[0], key)
		}
		if err == nil {
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(goo.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = yaml.Unmarshal(buf.Bytes(), &filters)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	} else {
		// try loading from filesystem
		if _, err = os.Stat(filtersPath); err == nil {
			_, openErr := os.Open(filtersPath)
			if openErr != nil && debug {
				fmt.Println(err)
			}
			filtersFileContent, readErr := ioutil.ReadFile(filtersPath)
			if readErr != nil && debug {
				fmt.Println(err)
			}
			filters, err = parseFiltersFileContent(filtersFileContent)
			if err != nil && debug {
				fmt.Println(err)
			}
		} else if debug {
			fmt.Println(err)
		}
	}
	return
}

// try loading report configuration from envvars, then from path provided
func loadReportConfig(configPath string, debug bool) (reportConfig Report) {
	var err error
	reportFilePath := ensureTrailingSlash(configPath) + reportFileName
	// try loading from envvars (only AWS SES Supported so far)
	upper := strings.ToUpper
	if upper(os.Getenv("AIR_EMAIL_PROVIDER")) == "SES" {
		if os.Getenv("AIR_EMAIL_AWS_REGION") != "" &&
			os.Getenv("AIR_EMAIL_SOURCE") != "" &&
			os.Getenv("AIR_EMAIL_RECIPIENTS") != "" &&
			os.Getenv("AIR_EMAIL_SUBJECT") != "" {
			recipients := strings.Split(os.Getenv("AIR_EMAIL_RECIPIENTS"), ",")
			email := Email{
				Provider:   "ses",
				Region:     os.Getenv("AIR_EMAIL_AWS_REGION"),
				Source:     os.Getenv("AIR_EMAIL_SOURCE"),
				Recipients: recipients,
				Subject:    os.Getenv("AIR_EMAIL_SUBJECT"),
			}
			reportConfig.Email = email
			return
		}
	}
	// try loading from s3
	if strings.HasPrefix(configPath, "s3://") {
		// split config path (minus prefix)
		parts := strings.Split(configPath[5:], "/")
		key := ensureTrailingSlash(strings.Join(parts[1:], "/")) + reportFileName
		sess := session.Must(session.NewSession())
		var region string
		region, err = s3manager.GetBucketRegion(context.Background(), sess, parts[0], "us-east-1")
		sess = session.Must(session.NewSession(&aws.Config{Region: ptrToStr(region)}))
		svc := s3.New(sess)
		input := &s3.GetObjectInput{
			Bucket: aws.String(parts[0]),
			Key:    aws.String(key),
		}
		var goo *s3.GetObjectOutput
		goo, err = svc.GetObject(input)
		if err != nil && debug {
			fmt.Printf("failed to load report from s3://%s/%s\n", parts[0], key)
			return
		}
		if goo == nil || goo.Body == nil {
			return
		}
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(goo.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = yaml.Unmarshal(buf.Bytes(), &reportConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}
	// try loading from file
	if _, err = os.Stat(reportFilePath); err == nil {
		_, err = os.Open(reportFilePath)
		if err != nil && debug {
			fmt.Println(err)
		}
		var reportFileContent []byte
		reportFileContent, err = ioutil.ReadFile(reportFilePath)
		if err != nil && debug {
			fmt.Println(err)
		}
		err = yaml.Unmarshal(reportFileContent, &reportConfig)
		if err != nil && debug {
			fmt.Println(err)
		}
	} else if debug {
		fmt.Println(err)
	}
	return
}

func loadTargets(configPath string, debug bool) (targets targets) {
	var err error
	// try loading from s3
	if strings.HasPrefix(configPath, "s3://") {
		// split config path (minus prefix)
		parts := strings.Split(configPath[5:], "/")
		key := ensureTrailingSlash(strings.Join(parts[1:], "/")) + targetsFileName
		sess := session.Must(session.NewSession())
		var region string
		region, err = s3manager.GetBucketRegion(context.Background(), sess, parts[0], "us-east-1")
		sess = session.Must(session.NewSession(&aws.Config{Region: ptrToStr(region)}))
		svc := s3.New(sess)
		input := &s3.GetObjectInput{
			Bucket: aws.String(parts[0]),
			Key:    aws.String(key),
		}
		var goo *s3.GetObjectOutput
		goo, err = svc.GetObject(input)
		if err != nil && debug {
			fmt.Printf("failed to load targets from s3://%s/%s\n", parts[0], key)
		}
		if err == nil {
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(goo.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = yaml.Unmarshal(buf.Bytes(), &targets)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	} else {
		// try loading from filesystem
		targets, err = readTargets(ensureTrailingSlash(configPath) + targetsFileName)
		if err != nil && debug {
			fmt.Println(err)
		}
	}
	return
}

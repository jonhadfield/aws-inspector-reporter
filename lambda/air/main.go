package main

import (
	air2 "github.com/jonhadfield/aws-inspector-reporter/air"
	"os"
	"strconv"

	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(cwe events.CloudWatchEvent) error {
	log.Printf("Processing Lambda cwe: %s\n", cwe.Time)
	var debug bool
	if os.Getenv("AIR_DEBUG") != "" {
		debug = true
	}
	var err error
	var maxReportAge int
	if os.Getenv("AIR_MAX_REPORT_AGE") != "" {
		maxReportAge, err = strconv.Atoi(os.Getenv("AIR_MAX_REPORT_AGE"))
		if err != nil {
			maxReportAge = air2.DefaultMaxReportAge
		}
	}

	err = air2.Run(air2.AppConfig{
		Debug:        debug,
		ConfigPath:   os.Getenv("AIR_CONFIG_PATH"),
		MaxReportAge: maxReportAge,
		OutputDir:    "/tmp",
	})
	if err != nil {
		log.Printf("error: %+v\n", err)
	}
	return err
}

func main() {
	lambda.Start(Handler)
}

package main

import (
	"os"

	air2 "github.com/jonhadfield/aws-inspector-reporter/air"

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

	err := air2.Run(air2.AppConfig{
		Debug:     debug,
		OutputDir: "/tmp",
	})
	if err != nil {
		log.Printf("error: %+v\n", err)
	}
	return err
}

func main() {
	lambda.Start(Handler)
}

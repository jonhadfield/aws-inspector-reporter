package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	air2 "github.com/jonhadfield/aws-inspector-reporter/air"

	"github.com/urfave/cli"
)

// overwritten at build time
var version, versionOutput, tag, sha, buildDate string

func main() {
	msg, display, err := startCLI(os.Args)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}
	if display && msg != "" {
		fmt.Println(msg)
	}
	os.Exit(0)
}

func startCLI(args []string) (msg string, display bool, err error) {
	if tag != "" && buildDate != "" {
		versionOutput = fmt.Sprintf("[%s-%s] %s UTC", tag, sha, buildDate)
	} else {
		versionOutput = version
	}

	app := cli.NewApp()
	app.EnableBashCompletion = true

	app.Name = "air"
	app.Version = versionOutput
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Jon Hadfield",
			Email: "jon@lessknown.co.uk",
		},
	}
	app.HelpName = "-"
	app.Usage = "AWS Inspector Reporter"
	app.Description = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config-path", Usage: "load configuration files from filesystem path or AWS S3 using s3://...", Value: "config/"},
		cli.StringFlag{Name: "output", Usage: "report output directory"},
		cli.IntFlag{Name: "max-report-age", Usage: "max age (in days) of reports to check", Value: air2.DefaultMaxReportAge},
		cli.BoolFlag{Name: "debug"},
	}

	app.Action = func(c *cli.Context) error {
		_ = air2.Run(air2.AppConfig{
			Debug:        c.Bool("debug"),
			ConfigPath:   c.String("config-path"),
			MaxReportAge: c.Int("max-report-age"),
			OutputDir:    strings.Trim(c.String("output"), " "),
		})

		return nil
	}
	return msg, display, app.Run(args)
}

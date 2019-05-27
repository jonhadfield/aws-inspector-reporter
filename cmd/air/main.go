package main

import (
	"fmt"
	"os"
	"sort"
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
		cli.BoolFlag{Name: "debug"},
		cli.StringFlag{Name: "targets", Usage: "load targets from `FILE`", Value: "config/targets.yml"},
		cli.StringFlag{Name: "filters", Usage: "load filters from `FILE`", Value: "config/filters.yml"},
		cli.StringFlag{Name: "report", Usage: "load report configuration from `FILE`", Value: "config/report.yml"},
		cli.StringFlag{Name: "output", Usage: "output directory"},
	}

	app.Action = func(c *cli.Context) error {
		_ = air2.Run(air2.AppConfig{
			Debug:       c.Bool("debug"),
			TargetsFile: c.String("targets"),
			FiltersFile: c.String("filters"),
			ReportFile:  c.String("report"),
			OutputDir:   strings.Trim(c.String("output"), " "),
		})

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	return msg, display, app.Run(args)
}

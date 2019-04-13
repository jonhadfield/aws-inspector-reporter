package air

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Filter struct {
	TitleMatch string `yaml:"title-match"`
	Severity   string `yaml:"severity"`
	Comment    string `yaml:"comment"`
}

type Filters []Filter

func loadFilters(filtersPath string, debug bool) (filters Filters) {
	var err error
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
	return
}

func parseFiltersFileContent(content []byte) (filters Filters, err error) {
	err = errors.WithStack(yaml.Unmarshal(content, &filters))
	return
}

func (ar *accountsResults) filter(filters Filters) {
	var filteredResults accountsResults
	for _, res := range *ar {
		filteredResult := res
		var filteredRegionResults []regionResult
		for _, rres := range res.regionResults {
			filteredRegionResult := rres
			var filteredRegionTemplateResults regionTemplateResults
			for _, rtr := range rres.regionTemplateResults {
				filtersRegionTemplateResult := rtr
				var filteredRuns []run
				for _, run := range rtr.runs {
					filteredRun := run
					var filteredFindings findings
					for _, f := range run.findings {
						filteredFinding := filterFinding(f, filters)
						filteredFindings = append(filteredFindings, filteredFinding)
					}
					filteredRun.findings = filteredFindings
					filteredRuns = append(filteredRuns, filteredRun)
				}
				filtersRegionTemplateResult.runs = filteredRuns
				filteredRegionTemplateResults = append(filteredRegionTemplateResults, filtersRegionTemplateResult)
			}
			filteredRegionResult.regionTemplateResults = filteredRegionTemplateResults
			filteredRegionResults = append(filteredRegionResults, filteredRegionResult)

		}
		filteredResult.regionResults = filteredRegionResults
		filteredResults = append(filteredResults, filteredResult)
	}
	*ar = filteredResults
}

func filterFinding(finding finding, filters Filters) (out finding) {
	out = finding
	for _, f := range filters {
		if f.TitleMatch != "" {
			r := regexp.MustCompile(f.TitleMatch)
			if r.MatchString(*finding.Title) {
				out.Severity = ptrToStr("ignored")
				out.comment = f.Comment
				return out
			}
		}
	}
	return finding
}

type finding struct {
	inspector.Finding
	rulePackageName string
	comment         string
}

func transformFinding(aF *inspector.Finding) (out finding) {
	out.Service = aF.Service
	out.ServiceAttributes = aF.ServiceAttributes
	out.Severity = aF.Severity
	out.AssetAttributes = aF.AssetAttributes
	out.Attributes = aF.Attributes
	out.Title = aF.Title
	out.CreatedAt = aF.CreatedAt
	out.Recommendation = aF.Recommendation
	out.Description = aF.Description
	out.UserAttributes = aF.UserAttributes
	out.AssetType = aF.AssetType
	out.Confidence = aF.Confidence
	out.Id = aF.Id
	out.IndicatorOfCompromise = aF.IndicatorOfCompromise
	out.SchemaVersion = aF.SchemaVersion
	out.NumericSeverity = aF.NumericSeverity
	out.UpdatedAt = aF.UpdatedAt
	return
}

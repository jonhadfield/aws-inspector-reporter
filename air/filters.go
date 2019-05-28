package air

import (
	"regexp"

	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type filter struct {
	TitleMatch string `yaml:"title-match"`
	Severity   string `yaml:"severity"`
	Comment    string `yaml:"comment"`
}

type filters []filter

func parseFiltersFileContent(content []byte) (filters filters, err error) {
	err = errors.WithStack(yaml.Unmarshal(content, &filters))
	return
}

func (ar *accountsResults) filter(filters filters) {
	if filterMaps == nil {
		filterMaps = make(map[string]*regexp.Regexp)
	}
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

var filterMaps map[string]*regexp.Regexp

func getCompiledRegex(regexString string) *regexp.Regexp {
	for k, v := range filterMaps {
		if k == regexString {
			return v
		}
	}
	newRegex := regexp.MustCompile(regexString)
	filterMaps[regexString] = newRegex
	return newRegex
}

func filterFinding(finding finding, filters filters) (out finding) {
	out = finding
	for _, f := range filters {
		if f.TitleMatch != "" {
			r := getCompiledRegex(f.TitleMatch)
			if r.MatchString(*finding.Title) {
				out.Severity = ptrToStr(f.Severity)
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

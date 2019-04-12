package air

import (
	"strings"
	"testing"

	a "github.com/jonhadfield/aws-inspector-reporter/airtest"
	"github.com/stretchr/testify/assert"
)

// Test with full set of results

func TestGetLatestAssessmentTemplateRunsComplete(t *testing.T) {
	m := &a.MockInspectorClient1{}
	output, _ := getLatestAssessmentTemplateRuns(m, nil)
	assert.Len(t, output, 11)
	assert.True(t, strings.HasSuffix(*output[0], "2j38BEoa"))
	assert.True(t, strings.HasSuffix(*output[10], "2j38BEok"))
}

func TestGetAssessmentTargetsArnsComplete(t *testing.T) {
	m := &a.MockInspectorClient1{}
	output, _ := getAssessmentTargetsArns(m)
	assert.Len(t, output, 11)
	assert.True(t, strings.HasSuffix(*output[0], "0-EgrdrY3A"))
	assert.True(t, strings.HasSuffix(*output[10], "0-EgrdrY3K"))
}

func TestGetAssessmentTemplatesArnsComplete(t *testing.T) {
	m := &a.MockInspectorClient1{}
	output, _ := getAssessmentTemplatesArns(m, nil)
	assert.Len(t, output, 11)
	assert.True(t, strings.HasSuffix(*output[0], "0-3gLCoEvA"))
	assert.True(t, strings.HasSuffix(*output[10], "0-3gLCoEvK"))
}

func TestGetRegionTemplateResultsComplete(t *testing.T) {
	m := &a.MockInspectorClient1{}
	results, err := getRegionTemplateResults(m)
	assert.NoError(t, err)
	assert.Len(t, results, 11)
}

func TestGetAllInspectorRegionsComplete(t *testing.T) {
	results := getAllInspectorRegions()
	assert.NotEmpty(t, results)
	assert.Contains(t, results, "eu-west-1")
	assert.Len(t, results, 10)
}

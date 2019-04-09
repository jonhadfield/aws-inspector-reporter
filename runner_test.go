package air

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/aws/aws-sdk-go/service/inspector/inspectoriface"
	"github.com/stretchr/testify/assert"
)

type MockInspectorClient struct {
	inspectoriface.InspectorAPI
}

func (m *MockInspectorClient) ListAssessmentRuns(input *inspector.ListAssessmentRunsInput) (*inspector.ListAssessmentRunsOutput, error) {
	runArns := []*string{
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEoa"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEob"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEoc"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEod"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEoe"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEof"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEog"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEoh"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEoi"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEoj"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEok"),
	}
	return &inspector.ListAssessmentRunsOutput{
		AssessmentRunArns: runArns,
	}, nil
}

func (m *MockInspectorClient) ListAssessmentTargets(input *inspector.ListAssessmentTargetsInput) (*inspector.ListAssessmentTargetsOutput, error) {
	targetArns := []*string{
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3B"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3C"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3D"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3E"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3F"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3G"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3H"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3I"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3J"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3K"),
	}
	return &inspector.ListAssessmentTargetsOutput{
		AssessmentTargetArns: targetArns,
	}, nil
}

func (m *MockInspectorClient) ListAssessmentTemplates(input *inspector.ListAssessmentTemplatesInput) (*inspector.ListAssessmentTemplatesOutput, error) {
	templateArns := []*string{
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvA"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvB"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvC"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvD"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvE"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvF"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvG"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvH"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvI"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvJ"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-EgrdrY3A/template/0-3gLCoEvK"),
	}
	return &inspector.ListAssessmentTemplatesOutput{
		AssessmentTemplateArns: templateArns,
	}, nil
}

func TestGetLatestAssessmentTemplateRuns(t *testing.T) {
	m := &MockInspectorClient{}
	output, _ := getLatestAssessmentTemplateRuns(m, nil)
	assert.Len(t, output, 11)
	assert.True(t, strings.HasSuffix(*output[0], "2j38BEoa"))
	assert.True(t, strings.HasSuffix(*output[10], "2j38BEok"))
}

func TestGetAssessmentTargetsArns(t *testing.T) {
	m := &MockInspectorClient{}
	output, _ := getAssessmentTargetsArns(m)
	assert.Len(t, output, 11)
	assert.True(t, strings.HasSuffix(*output[0], "0-EgrdrY3A"))
	assert.True(t, strings.HasSuffix(*output[10], "0-EgrdrY3K"))
}

func TestGetAssessmentTemplatesArns(t *testing.T) {
	m := &MockInspectorClient{}
	output, _ := getAssessmentTemplatesArns(m, nil)
	assert.Len(t, output, 11)
	assert.True(t, strings.HasSuffix(*output[0], "0-3gLCoEvA"))
	assert.True(t, strings.HasSuffix(*output[10], "0-3gLCoEvK"))
}

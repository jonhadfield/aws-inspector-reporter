package air

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/inspector"
	a "github.com/jonhadfield/aws-inspector-reporter/air/airtest"
	"github.com/stretchr/testify/assert"
)

func TestGetInstanceName(t *testing.T) {
	finding := finding{
		Finding: inspector.Finding{
			AssetAttributes: &inspector.AssetAttributes{
				Tags: []*inspector.Tag{
					{Key: ptrToStr("Service"), Value: ptrToStr("myService")},
					{Key: ptrToStr("Name"), Value: ptrToStr("test")},
				},
			},
		},
	}
	assert.True(t, getInstanceName(finding) == "test")
}

func TestFormatTitle(t *testing.T) {
	input := "One\nTwo Three\nFour"
	expected := "One\r\nTwo Three\r\nFour"
	assert.Equal(t, expected, formatTitle(input))
}

func TestFormatDescription(t *testing.T) {
	input := "One\nTwo Three\nDescription Four\n  Five  "
	expected := "One\r\nTwo Three\r\nFour\r\nFive"
	assert.Equal(t, expected, formatDescription(input))
}

func TestFormatRecommendation(t *testing.T) {
	input := "One\nTwo Three\nDescription Four\n  Five  "
	expected := "One\r\nTwo Three\r\nDescription Four\r\nFive"
	assert.Equal(t, expected, formatRecommendation(input))
}

func TestGetAccountID(t *testing.T) {
	m := &a.MockSTSClient{}
	output := getAccountID(m)
	assert.Equal(t, output, "012345678901")
}

func TestStringInSlice(t *testing.T) {
	assert.True(t, StringInSlice("test2", []string{"test1", "test2", "test3"}))
	assert.False(t, StringInSlice("test4", []string{"test1", "test2", "test3"}))
}

func TestPtrToInt64(t *testing.T) {
	assert.Equal(t, int64(123), *ptrToInt64(123))
}

func TestPtrToFloat(t *testing.T) {
	assert.Equal(t, 0.123, *ptrToFloat64(0.123))
}

func TestPtrToTime(t *testing.T) {
	now := time.Now()
	assert.Equal(t, now, *ptrToTime(now))
}

func TestPtrToBool(t *testing.T) {
	assert.True(t, *ptrToBool(true))
	assert.False(t, *ptrToBool(false))
}

func TestGetAccountAlias(t *testing.T) {
	m := &a.MockIAMClient{}
	output := getAccountAlias(m)
	assert.Equal(t, output, "testOne")
}

package airtest

import (
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/aws/aws-sdk-go/service/inspector/inspectoriface"
)

type MockSTSClient struct {
	stsiface.STSAPI
}

type MockIAMClient struct {
	iamiface.IAMAPI
}

func (m *MockSTSClient) GetCallerIdentity(in *sts.GetCallerIdentityInput) (out *sts.GetCallerIdentityOutput, err error) {
	return &sts.GetCallerIdentityOutput{
		Arn:     ptrToStr("test"),
		Account: ptrToStr("012345678901"),
	}, nil
}

func (m *MockIAMClient) ListAccountAliases(input *iam.ListAccountAliasesInput) (*iam.ListAccountAliasesOutput, error) {
	result := iam.ListAccountAliasesOutput{
		AccountAliases: []*string{
			ptrToStr("testOne"),
			ptrToStr("testTwo"),
		},
	}
	return &result, nil
}

type MockInspectorClient1 struct {
	inspectoriface.InspectorAPI
}

func (m *MockInspectorClient1) ListAssessmentRuns(input *inspector.ListAssessmentRunsInput) (*inspector.ListAssessmentRunsOutput, error) {
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

func (m *MockInspectorClient1) ListAssessmentTargets(input *inspector.ListAssessmentTargetsInput) (*inspector.ListAssessmentTargetsOutput, error) {
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

func (m *MockInspectorClient1) ListFindings(input *inspector.ListFindingsInput) (*inspector.ListFindingsOutput, error) {
	targetArns := []*string{
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0a"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0b"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0c"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0d"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0e"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0f"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0g"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0h"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0i"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0j"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-Jh4FVVe4/run/0-yWLFz8g4/finding/0-mb9Htk0k"),
	}
	return &inspector.ListFindingsOutput{
		FindingArns: targetArns,
	}, nil
}

func (m *MockInspectorClient1) ListAssessmentTemplates(input *inspector.ListAssessmentTemplatesInput) (*inspector.ListAssessmentTemplatesOutput, error) {
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

func (m *MockInspectorClient1) DescribeAssessmentRuns(input *inspector.DescribeAssessmentRunsInput) (*inspector.DescribeAssessmentRunsOutput, error) {
	assessmentRuns := []*inspector.AssessmentRun{
		{
			Arn:                   ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-RSL0ljsq/template/0-i0h82PKJ/run/0-2j38BEoa"),
			AssessmentTemplateArn: ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-4Jb8g2la"),
			CompletedAt: ptrToTime(time.Date(
				2019, 03, 23, 21, 34, 58, 651387237, time.UTC)),
			CreatedAt: ptrToTime(time.Date(
				2019, 03, 23, 20, 34, 58, 651387237, time.UTC)),
			DataCollected:     ptrToBool(true),
			DurationInSeconds: ptrToInt64(900),
			FindingCounts: map[string]*int64{
				"High":          ptrToInt64(117),
				"Informational": ptrToInt64(8),
				"Low":           ptrToInt64(0),
				"Medium":        ptrToInt64(9),
			},
			Name: ptrToStr("ao5o96ty-03b6-d11d-8344-755af84b8024_2ssd9bd2-ca9f-0e5c-76c1-21ac71cd347a"),
			Notifications: []*inspector.AssessmentRunNotification{
				{
					Date: ptrToTime(time.Date(
						2019, 03, 23, 21, 34, 58, 651387237, time.UTC)),
					Error:                ptrToBool(false),
					Event:                ptrToStr("ASSESSMENT_RUN_COMPLETED"),
					SnsPublishStatusCode: ptrToStr("SUCCESS"),
					SnsTopicArn:          ptrToStr("arn:aws:sns:eu-west-2:012345678901:Inspector-Scans"),
				},
			},
			RulesPackageArns: []*string{
				ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-sJBhCr0a"),
				ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-sJBhCr0b"),
				ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-sJBhCr0c"),
				ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-sJBhCr0d"),
				ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-sJBhCr0e"),
			},
			StartedAt: ptrToTime(time.Date(
				2019, 03, 23, 21, 34, 58, 651387237, time.UTC)),
			State: ptrToStr("COMPLETED"),
			StateChangedAt: ptrToTime(time.Date(
				2019, 03, 23, 21, 34, 58, 651387237, time.UTC)),
		},
	}
	return &inspector.DescribeAssessmentRunsOutput{
		AssessmentRuns: assessmentRuns,
	}, nil
}

func (m *MockInspectorClient1) ListRulesPackages(input *inspector.ListRulesPackagesInput) (*inspector.ListRulesPackagesOutput, error) {
	runArns := []*string{
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xa"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xb"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xc"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xd"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xe"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xf"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xg"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xh"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xi"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xj"),
		ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-SPzU33xk"),
	}
	return &inspector.ListRulesPackagesOutput{
		RulesPackageArns: runArns,
	}, nil
}

func (m *MockInspectorClient1) DescribeRulesPackages(input *inspector.DescribeRulesPackagesInput) (*inspector.DescribeRulesPackagesOutput, error) {
	rulesPackages := []*inspector.RulesPackage{
		{
			Arn:         ptrToStr("arn:aws:inspector:eu-west-1:357557129151:rulespackage/0-SPzU33xe"),
			Description: ptrToStr("These rules analyze the reachability of your instances over the network. Attacks can exploit your instances over the network by accessing services that are listening on open ports."),
			Name:        ptrToStr("Network Reachability"),
			Provider:    ptrToStr("Amazon Web Services, Inc."),
			Version:     ptrToStr("1.1"),
		},
		{
			Arn:         ptrToStr("arn:aws:inspector:eu-west-1:357557129151:rulespackage/0-SnojL3Z6"),
			Description: ptrToStr("The rules in this package help determine whether your systems are configured securely."),
			Name:        ptrToStr("Security Best Practices"),
			Provider:    ptrToStr("Amazon Web Services, Inc."),
			Version:     ptrToStr("1.0"),
		},
		{
			Arn:         ptrToStr("arn:aws:inspector:eu-west-1:357557129151:rulespackage/0-ubA5XvBh"),
			Description: ptrToStr("The rules in this package help verify whether the EC2 instances in your application are exposed to Common Vulnerabilities and Exposures (CVEs)."),
			Name:        ptrToStr("Common Vulnerabilities and Exposures"),
			Provider:    ptrToStr("Amazon Web Services, Inc."),
			Version:     ptrToStr("1.1"),
		},
		{
			Arn:         ptrToStr("arn:aws:inspector:eu-west-1:357557129151:rulespackage/0-lLmwe1zd"),
			Description: ptrToStr("These rules analyze the behavior of your instances during an assessment run, and provide guidance on how to make your instances more secure."),
			Name:        ptrToStr("Runtime Behavior Analysis"),
			Provider:    ptrToStr("Amazon Web Services, Inc."),
			Version:     ptrToStr("1.0"),
		},
		{
			Arn:         ptrToStr("arn:aws:inspector:eu-west-1:357557129151:rulespackage/0-sJBhCr0F"),
			Description: ptrToStr("The CIS Security Benchmarks program provides well-defined, un-biased and consensus-based industry best practices to help organizations assess and improve their security."),
			Name:        ptrToStr("CIS Operating System Security Configuration Benchmarks"),
			Provider:    ptrToStr("Amazon Web Services, Inc."),
			Version:     ptrToStr("1.0"),
		},
	}
	return &inspector.DescribeRulesPackagesOutput{
		RulesPackages: rulesPackages,
	}, nil
}

func (m *MockInspectorClient1) DescribeFindings(input *inspector.DescribeFindingsInput) (*inspector.DescribeFindingsOutput, error) {
	findings := []*inspector.Finding{
		{
			Arn: ptrToStr(""),
			AssetAttributes: &inspector.AssetAttributes{
				AmiId:             ptrToStr(""),
				SchemaVersion:     ptrToInt64(1),
				AgentId:           ptrToStr("a"),
				AutoScalingGroup:  ptrToStr("a"),
				Tags:              []*inspector.Tag{},
				Hostname:          ptrToStr("a"),
				Ipv4Addresses:     []*string{},
				NetworkInterfaces: []*inspector.NetworkInterface{},
			},
			SchemaVersion: ptrToInt64(1),
			ServiceAttributes: &inspector.ServiceAttributes{
				RulesPackageArn: ptrToStr("TBC"),
			},
			Description: ptrToStr("TBC"),
			CreatedAt: ptrToTime(time.Date(
				2019, 03, 23, 21, 34, 58, 651387237, time.UTC)),
			Title:          ptrToStr("TBC"),
			Severity:       ptrToStr("TBC"),
			AssetType:      ptrToStr("TBC"),
			Id:             ptrToStr("TBC"),
			Recommendation: ptrToStr("TBC"),
			Service:        ptrToStr("TBC"),
			UpdatedAt: ptrToTime(time.Date(
				2019, 03, 23, 21, 34, 58, 651387237, time.UTC)),
			IndicatorOfCompromise: ptrToBool(true),
			NumericSeverity:       ptrToFloat64(1.1),
			Confidence:            ptrToInt64(1),
			Attributes:            []*inspector.Attribute{},
			UserAttributes:        []*inspector.Attribute{},
		},
	}
	return &inspector.DescribeFindingsOutput{
		Findings: findings,
	}, nil
}

func (m *MockInspectorClient1) DescribeAssessmentTemplates(input *inspector.DescribeAssessmentTemplatesInput) (*inspector.DescribeAssessmentTemplatesOutput, error) {
	assessmentTemplates := []*inspector.AssessmentTemplate{
		{
			Arn:  ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-4Jb8g2li"),
			Name: ptrToStr("Template"),
			RulesPackageArns: []*string{
				ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-sJBhCr0a"),
				ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-sJBhCr0b"),
				ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-sJBhCr0c"),
				ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-sJBhCr0d"),
				ptrToStr("arn:aws:inspector:eu-west-2:012345678901:rulespackage/0-sJBhCr0e"),
			},
			CreatedAt: ptrToTime(time.Date(
				2019, 03, 23, 21, 34, 58, 651387237, time.UTC)),
			AssessmentRunCount:        ptrToInt64(6),
			AssessmentTargetArn:       ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n"),
			DurationInSeconds:         ptrToInt64(900),
			LastAssessmentRunArn:      ptrToStr("arn:aws:inspector:eu-west-2:012345678901:target/0-Ivm80T1n/template/0-4Jb8g2li/run/0-SkhjLkhv"),
			UserAttributesForFindings: []*inspector.Attribute{},
		},
	}
	return &inspector.DescribeAssessmentTemplatesOutput{
		AssessmentTemplates: assessmentTemplates,
		FailedItems:         nil,
	}, nil

}

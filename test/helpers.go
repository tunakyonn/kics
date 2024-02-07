package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Checkmarx/kics/pkg/model"
	"github.com/spf13/cobra"
)

const (
	// ValidUUIDRegex is a constant representing a regular expression rule to validate UUID string
	ValidUUIDRegex = `(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}$`
	positive       = "positive.tf"
	positiveYaml   = "test\\fixtures\\test_critical_custom_queries\\amazon_mq_broker_encryption_disabled\\test\\positive1.yaml"
)

type execute func() error

// CaptureOutput changes default stdout to intercept into a buffer, converts it to string and returns it
func CaptureOutput(funcToExec execute) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := funcToExec()

	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		if _, errs := io.Copy(&buf, r); errs != nil {
			return
		}
		outC <- buf.String()
	}()

	if errs := w.Close(); errs != nil {
		return "", errs
	}
	os.Stdout = old
	out := <-outC

	return out, err
}

// CaptureCommandOutput set cobra command args, if necessary, then capture the output
func CaptureCommandOutput(cmd *cobra.Command, args []string) (string, error) {
	if len(args) > 0 {
		cmd.SetArgs(args)
	}

	return CaptureOutput(cmd.Execute)
}

// ChangeCurrentDir gets current working directory and changes to its parent until finds the desired directory
// or fail
func ChangeCurrentDir(desiredDir string) error {
	for currentDir, err := os.Getwd(); GetCurrentDirName(currentDir) != desiredDir; currentDir, err = os.Getwd() {
		if err == nil {
			if err = os.Chdir(".."); err != nil {
				fmt.Print(formatCurrentDirError(err))
				return errors.New(formatCurrentDirError(err))
			}
		} else {
			return errors.New(formatCurrentDirError(err))
		}
	}
	return nil
}

func formatCurrentDirError(err error) string {
	return fmt.Sprintf("change path error = %v", err)
}

// GetCurrentDirName returns current working directory
func GetCurrentDirName(path string) string {
	dirs := strings.Split(path, string(os.PathSeparator))
	if dirs[len(dirs)-1] == "" && len(dirs) > 1 {
		return dirs[len(dirs)-2]
	}
	return dirs[len(dirs)-1]
}

// StringifyStruct stringify struct for pretty print
func StringifyStruct(v interface{}) (string, error) {
	jsonValue, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonValue), nil
}

// MapToStringSlice extract slice of keys from a map[string]string
func MapToStringSlice(stringKeyMap map[string]string) []string {
	keys := make([]string, len(stringKeyMap))

	i := 0
	for k := range stringKeyMap {
		keys[i] = k
		i++
	}
	return keys
}

var queryHigh = model.QueryResult{
	QueryName:                   "ALB protocol is HTTP",
	QueryID:                     "de7f5e83-da88-4046-871f-ea18504b1d43",
	Description:                 "ALB protocol is HTTP Description",
	DescriptionID:               "504b1d43",
	CISDescriptionIDFormatted:   "testCISID",
	CISDescriptionTitle:         "testCISTitle",
	CISDescriptionTextFormatted: "testCISDescription",
	Severity:                    model.SeverityHigh,
	Files: []model.VulnerableFile{
		{
			FileName:         positive,
			Line:             25,
			IssueType:        "MissingAttribute",
			SearchKey:        "aws_alb_listener[front_end].default_action.redirect",
			KeyExpectedValue: "'default_action.redirect.protocol' is equal 'HTTPS'",
			KeyActualValue:   "'default_action.redirect.protocol' is missing",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
		{
			FileName:         positive,
			Line:             19,
			IssueType:        "IncorrectValue",
			SearchKey:        "aws_alb_listener[front_end].default_action.redirect",
			KeyExpectedValue: "'default_action.redirect.protocol' is equal 'HTTPS'",
			KeyActualValue:   "'default_action.redirect.protocol' is equal 'HTTP'",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
	},
}

var queryMedium = model.QueryResult{
	QueryName:     "AmazonMQ Broker Encryption Disabled",
	Description:   "AmazonMQ Broker should have Encryption Options defined",
	QueryID:       "3db3f534-e3a3-487f-88c7-0a9fbf64b702",
	CloudProvider: "AWS",
	Severity:      model.SeverityMedium,
	Files: []model.VulnerableFile{
		{
			FileName:         positive,
			Line:             1,
			IssueType:        "MissingAttribute",
			SimilarityID:     "6b76f7a507e200bb2c73468ec9649b099da96a4efa0f49a3bdc88e12476d8ee7",
			SearchKey:        "resource.aws_mq_broker[positive1]",
			KeyExpectedValue: "resource.aws_mq_broker[positive1].encryption_options is defined",
			KeyActualValue:   "resource.aws_mq_broker[positive1].encryption_options is not defined",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
	},
}

var queryMedium2 = model.QueryResult{
	QueryName: "GuardDuty Detector Disabled",
	QueryID:   "704dadd3-54fc-48ac-b6a0-02f170011473",
	Severity:  model.SeverityMedium,
	Files: []model.VulnerableFile{
		{
			FileName:         filepath.Join("assets", "queries", "terraform", "aws", "guardduty_detector_disabled", "test", "positive.tf"),
			Line:             2,
			IssueType:        "IncorrectValue",
			SearchKey:        "aws_guardduty_detector[positive1].enable",
			KeyExpectedValue: "GuardDuty Detector should be Enabled",
			KeyActualValue:   "GuardDuty Detector is not Enabled",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
	},
	Platform:    "Terraform",
	Description: "Make sure that Amazon GuardDuty is Enabled",
}

var queryInfo = model.QueryResult{
	QueryName: "Resource Not Using Tags",
	QueryID:   "e38a8e0a-b88b-4902-b3fe-b0fcb17d5c10",
	Severity:  model.SeverityInfo,
	Files: []model.VulnerableFile{
		{
			FileName:         filepath.Join("assets", "queries", "terraform", "aws", "guardduty_detector_disabled", "test", "negative.tf"),
			Line:             1,
			IssueType:        "MissingAttribute",
			SearchKey:        "aws_guardduty_detector[{{negative1}}]",
			KeyExpectedValue: "aws_guardduty_detector[{{negative1}}].tags is defined and not null",
			KeyActualValue:   "aws_guardduty_detector[{{negative1}}].tags is undefined or null",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
		{
			FileName:         filepath.Join("assets", "queries", "terraform", "aws", "guardduty_detector_disabled", "test", "positive.tf"),
			Line:             1,
			IssueType:        "MissingAttribute",
			SearchKey:        "aws_guardduty_detector[{{positive1}}]",
			KeyExpectedValue: "aws_guardduty_detector[{{positive1}}].tags is defined and not null",
			KeyActualValue:   "aws_guardduty_detector[{{positive1}}].tags is undefined or null",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
	},
	Platform:    "Terraform",
	Description: "AWS services resource tags are an essential part of managing components",
}

var queryHighExperimental = model.QueryResult{
	QueryName:                   "ALB protocol is HTTP",
	QueryID:                     "de7f5e83-da88-4046-871f-ea18504b1d43",
	Description:                 "ALB protocol is HTTP Description",
	DescriptionID:               "504b1d43",
	CISDescriptionIDFormatted:   "testCISID",
	CISDescriptionTitle:         "testCISTitle",
	CISDescriptionTextFormatted: "testCISDescription",
	Severity:                    model.SeverityHigh,
	Experimental:                true,
	Files: []model.VulnerableFile{
		{
			FileName:         positive,
			Line:             25,
			IssueType:        "MissingAttribute",
			SearchKey:        "aws_alb_listener[front_end].default_action.redirect",
			KeyExpectedValue: "'default_action.redirect.protocol' is equal 'HTTPS'",
			KeyActualValue:   "'default_action.redirect.protocol' is missing",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
		{
			FileName:         positive,
			Line:             19,
			IssueType:        "IncorrectValue",
			SearchKey:        "aws_alb_listener[front_end].default_action.redirect",
			KeyExpectedValue: "'default_action.redirect.protocol' is equal 'HTTPS'",
			KeyActualValue:   "'default_action.redirect.protocol' is equal 'HTTP'",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
	},
}

var queryHighCWE = model.QueryResult{
	QueryName:                   "AMI Not Encrypted",
	QueryID:                     "97707503-a22c-4cd7-b7c0-f088fa7cf830",
	Description:                 "AWS AMI Encryption is not enabled",
	DescriptionID:               "a4342f0",
	CISDescriptionIDFormatted:   "testCISID",
	CISDescriptionTitle:         "testCISTitle",
	CISDescriptionTextFormatted: "testCISDescription",
	Severity:                    model.SeverityHigh,
	Files: []model.VulnerableFile{
		{
			FileName:         positive,
			Line:             30,
			IssueType:        "MissingAttribute",
			SearchKey:        "aws_alb_listener[front_end].default_action.redirect",
			KeyExpectedValue: "'default_action.redirect.protocol' is equal 'HTTPS'",
			KeyActualValue:   "'default_action.redirect.protocol' is missing",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
	},
	CWE: "22",
}

var queryCritical = model.QueryResult{
	QueryName:                   "AmazonMQ Broker Encryption Disabled",
	QueryID:                     "316278b3-87ac-444c-8f8f-a733a28da609",
	Description:                 "AmazonMQ Broker should have Encryption Options defined",
	DescriptionID:               "c5d562d9",
	CISDescriptionIDFormatted:   "testCISID",
	CISDescriptionTitle:         "testCISTitle",
	CISDescriptionTextFormatted: "testCISDescription",
	CloudProvider:               "AWS",
	Severity:                    model.SeverityCritical,
	Files: []model.VulnerableFile{
		{
			FileName:         positiveYaml,
			Line:             6,
			IssueType:        "MissingAttribute",
			SearchKey:        "aws_alb_listener[front_end].default_action.redirect",
			KeyExpectedValue: "'default_action.redirect.protocol' is equal 'HTTPS'",
			KeyActualValue:   "'default_action.redirect.protocol' is missing",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
	},
}

var SummaryMockCritical = model.Summary{
	Counters: model.Counters{
		ScannedFiles:           2,
		ParsedFiles:            2,
		FailedToScanFiles:      0,
		TotalQueries:           1,
		FailedToExecuteQueries: 0,
	},
	Queries: []model.QueryResult{
		queryCritical,
	},
	SeveritySummary: model.SeveritySummary{
		ScanID: "console",
		SeverityCounters: map[model.Severity]int{
			model.SeverityInfo:     0,
			model.SeverityLow:      0,
			model.SeverityMedium:   0,
			model.SeverityHigh:     0,
			model.SeverityCritical: 1,
		},
		TotalCounter: 1,
	},
	ScannedPaths: []string{
		"./",
	},
}

// SummaryMockCriticalFullPath test with full path to test file
var SummaryMockCriticalFullPath = model.Summary{
	Counters: model.Counters{
		ScannedFiles:           2,
		ParsedFiles:            2,
		FailedToScanFiles:      0,
		TotalQueries:           1,
		FailedToExecuteQueries: 0,
	},
	Queries: []model.QueryResult{
		queryCritical,
	},
	SeveritySummary: model.SeveritySummary{
		ScanID: "console",
		SeverityCounters: map[model.Severity]int{
			model.SeverityInfo:     0,
			model.SeverityLow:      0,
			model.SeverityMedium:   0,
			model.SeverityHigh:     0,
			model.SeverityCritical: 1,
		},
		TotalCounter: 1,
	},
	ScannedPaths: []string{
		".\\test\\fixtures\\test_critical_custom_queries\\amazon_mq_broker_encryption_disabled\\test\\",
	},
}

var queryCriticalASFF = model.QueryResult{
	QueryName:     "AmazonMQ Broker Encryption Disabled",
	QueryID:       "316278b3-87ac-444c-8f8f-a733a28da609",
	Description:   "AmazonMQ Broker should have Encryption Options defined",
	DescriptionID: "c5d562d9",
	CloudProvider: "AWS",
	Severity:      model.SeverityCritical,
	Files: []model.VulnerableFile{
		{
			FileName:         positiveYaml,
			Line:             6,
			IssueType:        "MissingAttribute",
			SearchKey:        "aws_alb_listener[front_end].default_action.redirect",
			KeyExpectedValue: "'default_action.redirect.protocol' is equal 'HTTPS'",
			KeyActualValue:   "'default_action.redirect.protocol' is missing",
			Value:            nil,
			VulnLines:        &[]model.CodeLine{},
		},
	},
}

var SummaryMockCriticalFullPathASFF = model.Summary{
	Counters: model.Counters{
		ScannedFiles:           2,
		ParsedFiles:            2,
		FailedToScanFiles:      0,
		TotalQueries:           1,
		FailedToExecuteQueries: 0,
	},
	Queries: []model.QueryResult{
		queryCriticalASFF,
	},
	SeveritySummary: model.SeveritySummary{
		ScanID: "console",
		SeverityCounters: map[model.Severity]int{
			model.SeverityInfo:     0,
			model.SeverityLow:      0,
			model.SeverityMedium:   0,
			model.SeverityHigh:     0,
			model.SeverityCritical: 1,
		},
		TotalCounter: 1,
	},
	ScannedPaths: []string{
		".\\test\\fixtures\\test_critical_custom_queries\\amazon_mq_broker_encryption_disabled\\test\\",
	},
}

// SummaryMock a summary to be used without running kics scan
var SummaryMock = model.Summary{
	Counters: model.Counters{
		ScannedFiles:           1,
		ParsedFiles:            1,
		FailedToScanFiles:      0,
		TotalQueries:           1,
		FailedToExecuteQueries: 0,
	},
	Queries: []model.QueryResult{
		queryHigh,
	},
	SeveritySummary: model.SeveritySummary{
		ScanID: "console",
		SeverityCounters: map[model.Severity]int{
			model.SeverityInfo:     0,
			model.SeverityLow:      0,
			model.SeverityMedium:   0,
			model.SeverityHigh:     2,
			model.SeverityCritical: 0,
		},
		TotalCounter: 2,
	},
	ScannedPaths: []string{
		"./",
	},
}

// ComplexSummaryMock a summary with more results to be used without running kics scan
var ComplexSummaryMock = model.Summary{
	Counters: model.Counters{
		ScannedFiles:           2,
		ParsedFiles:            2,
		FailedToScanFiles:      0,
		TotalQueries:           3,
		FailedToExecuteQueries: 0,
	},
	Queries: []model.QueryResult{
		queryHigh,
		queryMedium,
		queryHighCWE,
	},
	SeveritySummary: model.SeveritySummary{
		ScanID: "console",
		SeverityCounters: map[model.Severity]int{
			model.SeverityInfo:     0,
			model.SeverityLow:      0,
			model.SeverityMedium:   1,
			model.SeverityHigh:     2,
			model.SeverityCritical: 0,
		},
		TotalCounter: 5,
	},
	LatestVersion: model.Version{
		Latest: true,
	},
}

var ComplexSummaryMockWithExperimental = model.Summary{
	Counters: model.Counters{
		ScannedFiles:           2,
		ParsedFiles:            2,
		FailedToScanFiles:      0,
		TotalQueries:           2,
		FailedToExecuteQueries: 0,
	},
	Queries: []model.QueryResult{
		queryHighExperimental,
		queryMedium,
	},
	SeveritySummary: model.SeveritySummary{
		ScanID: "console",
		SeverityCounters: map[model.Severity]int{
			model.SeverityInfo:     0,
			model.SeverityLow:      0,
			model.SeverityMedium:   1,
			model.SeverityHigh:     2,
			model.SeverityCritical: 0,
		},
		TotalCounter: 3,
	},
	LatestVersion: model.Version{
		Latest: true,
	},
}

// ExampleSummaryMock a summary with specific results to CycloneDX report tests
var ExampleSummaryMock = model.Summary{
	Counters: model.Counters{
		ScannedFiles:           2,
		ParsedFiles:            2,
		FailedToScanFiles:      0,
		TotalQueries:           2,
		FailedToExecuteQueries: 0,
	},
	Queries: []model.QueryResult{
		queryInfo,
		queryMedium2,
	},
	SeveritySummary: model.SeveritySummary{
		ScanID: "console",
		SeverityCounters: map[model.Severity]int{
			model.SeverityInfo:     2,
			model.SeverityLow:      0,
			model.SeverityMedium:   1,
			model.SeverityHigh:     0,
			model.SeverityCritical: 0,
		},
		TotalCounter: 3,
	},
	ScannedPaths: []string{
		"./",
	},
}

// SimpleSummaryMock a summary with specific results to ASFF report tests
var SimpleSummaryMock = model.Summary{
	Counters: model.Counters{
		ScannedFiles:           1,
		ParsedFiles:            1,
		FailedToScanFiles:      0,
		TotalQueries:           1,
		FailedToExecuteQueries: 0,
	},
	Queries: []model.QueryResult{
		queryMedium,
	},
	SeveritySummary: model.SeveritySummary{
		ScanID: "console",
		SeverityCounters: map[model.Severity]int{
			model.SeverityInfo:     0,
			model.SeverityLow:      0,
			model.SeverityMedium:   1,
			model.SeverityHigh:     0,
			model.SeverityCritical: 0,
		},
		TotalCounter: 1,
	},
	ScannedPaths: []string{
		"./",
	},
}

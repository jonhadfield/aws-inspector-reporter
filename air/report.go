package air

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func generateAccountRegionXLSXData(accountResults accountResults) (data []dataRow) {

	for _, regionResult := range accountResults.regionResults {
		if len(regionResult.regionTemplateResults) == 0 {
			continue
		}

		for _, r := range regionResult.regionTemplateResults {
			for _, run := range r.runs {
				for _, f := range run.findings {
					var dr dataRow
					dr.region = regionResult.region
					dr.severity = strings.ToUpper(*f.Severity)
					dr.findingTitle = formatTitle(*f.Title)
					dr.instanceName = getInstanceName(f)
					dr.instanceID = *f.AssetAttributes.AgentId
					dr.createdAt = *f.CreatedAt
					dr.templateName = r.templateName
					dr.packageArn = *f.ServiceAttributes.RulesPackageArn
					dr.packageName = f.rulePackageName
					if f.AssetAttributes.AmiId != nil {
						dr.amiID = *f.AssetAttributes.AmiId
					}
					dr.template = r.templateArn
					dr.comment = f.comment
					dr.description = formatDescription(*f.Description)
					dr.recommendation = formatRecommendation(*f.Recommendation)
					if f.AssetAttributes.AutoScalingGroup != nil {
						dr.asgName = *f.AssetAttributes.AutoScalingGroup
					} else {
						dr.asgName = "-"
					}
					data = append(data, dr)
				}
			}
		}
	}

	sevLookup := map[string]int{}
	sevLookup["INFORMATIONAL"] = 1
	sevLookup["LOW"] = 2
	sevLookup["MEDIUM"] = 3
	sevLookup["HIGH"] = 4
	sort.Slice(data, func(i, j int) bool {
		return sevLookup[data[i].severity] > sevLookup[data[j].severity]
	})

	return data
}

type dataRow struct {
	createdAt      time.Time
	template       string
	region         string
	templateName   string
	packageArn     string
	packageName    string
	severity       string
	findingTitle   string
	instanceID     string
	instanceName   string
	amiID          string
	asgName        string
	description    string
	recommendation string
	comment        string
}

func generateSpreadsheet(accountsResults accountsResults, outputDir string) (string, error) {
	xlsx := excelize.NewFile()

	var headerStyle, highResultStyle, mediumResultStyle, lowResultStyle, infoResultStyle, ignoredResultStyle, defaultCenteredStyle int
	headerStyle, _ = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#000066"],"pattern":1},"font":{"bold":true,"italic":false,"family":"Calibri","size":14,"color":"#f2f2f2"},"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":true,"vertical":"","wrap_text":false}}`)
	highResultStyle, _ = xlsx.NewStyle(`{"font":{"bold":true,"italic":false,"family":"Calibri","size":12,"color":"#cc0000"},"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":true,"vertical":"","wrap_text":false}}`)
	mediumResultStyle, _ = xlsx.NewStyle(`{"font":{"bold":true,"italic":false,"family":"Calibri","size":12,"color":"#cc6600"},"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":true,"vertical":"","wrap_text":false}}`)
	lowResultStyle, _ = xlsx.NewStyle(`{"font":{"bold":true,"italic":false,"family":"Calibri","size":12,"color":"#003399"},"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":true,"vertical":"","wrap_text":false}}`)
	infoResultStyle, _ = xlsx.NewStyle(`{"font":{"bold":true,"italic":false,"family":"Calibri","size":12,"color":"#000000"},"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":true,"vertical":"","wrap_text":false}}`)
	ignoredResultStyle, _ = xlsx.NewStyle(`{"font":{"bold":true,"italic":false,"family":"Calibri","size":12,"color":"#000000"},"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":true,"vertical":"","wrap_text":false}}`)
	defaultCenteredStyle, _ = xlsx.NewStyle(`{"font":{"bold":false,"italic":false,"family":"Calibri","size":12,"color":"#000000"},"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":true,"vertical":"","wrap_text":true}}`)
	var firstSheet bool
	for _, accountResults := range accountsResults {
		accountSpreadsheetData := generateAccountRegionXLSXData(accountResults)
		if len(accountSpreadsheetData) == 0 {
			continue
		}
		sheetName := accountResults.accountAlias
		if !firstSheet {
			firstSheet = true
			xlsx.SetSheetName(xlsx.GetSheetName(1), sheetName)
		} else {
			_ = xlsx.NewSheet(accountResults.accountAlias)
		}

		_ = xlsx.SetCellValue(sheetName, "A1", "SEVERITY")
		_ = xlsx.SetCellValue(sheetName, "B1", "REGION")
		_ = xlsx.SetCellValue(sheetName, "C1", "TEMPLATE")
		_ = xlsx.SetCellValue(sheetName, "D1", "DATE")
		_ = xlsx.SetCellValue(sheetName, "E1", "INSTANCE ID")
		_ = xlsx.SetCellValue(sheetName, "F1", "INSTANCE NAME")
		_ = xlsx.SetCellValue(sheetName, "G1", "ASG")
		_ = xlsx.SetCellValue(sheetName, "H1", "RULES PACKAGE")
		_ = xlsx.SetCellValue(sheetName, "I1", "TITLE")
		_ = xlsx.SetCellValue(sheetName, "J1", "DESCRIPTION")
		_ = xlsx.SetCellValue(sheetName, "K1", "RECOMMENDATION")
		_ = xlsx.SetCellStyle(sheetName, "A1", "K1", headerStyle)
		_ = xlsx.SetColWidth(sheetName, "A", "A", 15)
		_ = xlsx.SetColWidth(sheetName, "B", "B", 13.5)
		_ = xlsx.SetColWidth(sheetName, "C", "C", 26)
		_ = xlsx.SetColWidth(sheetName, "D", "D", 22.5)
		_ = xlsx.SetColWidth(sheetName, "E", "E", 19)
		_ = xlsx.SetColWidth(sheetName, "F", "F", 24)
		_ = xlsx.SetColWidth(sheetName, "G", "G", 20)
		_ = xlsx.SetColWidth(sheetName, "H", "H", 44)
		_ = xlsx.SetColWidth(sheetName, "I", "I", 60)
		_ = xlsx.SetColWidth(sheetName, "J", "J", 70)
		_ = xlsx.SetColWidth(sheetName, "K", "K", 150)
		var lastRow string
		for i, dataRow := range accountSpreadsheetData {
			rowNum := i + 2
			strRowNum := strconv.Itoa(rowNum)
			resultCell := "A" + strRowNum
			regionCell := "B" + strRowNum
			templateCell := "C" + strRowNum
			dateCell := "D" + strRowNum
			instanceIDCell := "E" + strRowNum
			instanceNameCell := "F" + strRowNum
			instanceASGCell := "G" + strRowNum
			rulesPackageCell := "H" + strRowNum
			findingTitleCell := "I" + strRowNum
			descriptionCell := "J" + strRowNum
			recommendationCell := "K" + strRowNum
			_ = xlsx.SetCellValue(sheetName, resultCell, dataRow.severity)
			switch dataRow.severity {
			case "HIGH":
				_ = xlsx.SetCellStyle(sheetName, "A"+strRowNum, "A"+strRowNum, highResultStyle)
			case "MEDIUM":
				_ = xlsx.SetCellStyle(sheetName, "A"+strRowNum, "A"+strRowNum, mediumResultStyle)
			case "LOW":
				_ = xlsx.SetCellStyle(sheetName, "A"+strRowNum, "A"+strRowNum, lowResultStyle)
			case "INFORMATIONAL":
				_ = xlsx.SetCellStyle(sheetName, "A"+strRowNum, "A"+strRowNum, infoResultStyle)
			case "IGNORED":
				_ = xlsx.SetCellStyle(sheetName, "A"+strRowNum, "A"+strRowNum, ignoredResultStyle)
			}
			if dataRow.comment != "" {
				comment := fmt.Sprintf("{\"author\":\"%s\",\"text\":\" %s\"}", "-", dataRow.comment)
				_ = xlsx.AddComment(sheetName, "A"+strRowNum, comment)
			}
			_ = xlsx.SetCellValue(sheetName, regionCell, dataRow.region)
			_ = xlsx.SetCellValue(sheetName, templateCell, dataRow.templateName)
			_ = xlsx.SetCellValue(sheetName, dateCell, dataRow.createdAt.Format(time.ANSIC))
			_ = xlsx.SetCellValue(sheetName, instanceIDCell, dataRow.instanceID)
			instComment := fmt.Sprintf("{\"author\":\"%s\",\"text\":\" %s\"}", "AMI:", dataRow.amiID)
			_ = xlsx.AddComment(sheetName, instanceIDCell, instComment)
			_ = xlsx.SetCellValue(sheetName, instanceNameCell, dataRow.instanceName)
			_ = xlsx.SetCellValue(sheetName, rulesPackageCell, dataRow.packageName)
			_ = xlsx.SetCellValue(sheetName, instanceASGCell, dataRow.asgName)
			_ = xlsx.SetCellValue(sheetName, findingTitleCell, dataRow.findingTitle)
			_ = xlsx.SetCellValue(sheetName, descriptionCell, dataRow.description)
			_ = xlsx.SetCellValue(sheetName, recommendationCell, dataRow.recommendation)
			_ = xlsx.SetCellStyle(sheetName, "B"+strRowNum, "B"+strRowNum, defaultCenteredStyle)
			_ = xlsx.SetCellStyle(sheetName, "E"+strRowNum, "G"+strRowNum, defaultCenteredStyle)
			lastRow = strRowNum
		}
		_ = xlsx.AutoFilter(sheetName, "A1", "H"+lastRow, "")
	}

	timeStamp := time.Now().UTC().Format("20060102150405")
	var pathPrefix string
	if outputDir != "" {
		pathPrefix = outputDir
		if !strings.HasSuffix(outputDir, string(filepath.Separator)) {
			pathPrefix = outputDir + string(filepath.Separator)
		}
	}
	path := fmt.Sprintf("%sinspector_report_%s.xlsx", pathPrefix, timeStamp)
	err := xlsx.SaveAs(path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	absPath, _ := filepath.Abs(path)
	fmt.Println("report written to:", absPath)
	return absPath, err
}

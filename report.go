package air

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	yaml "gopkg.in/yaml.v2"
)

func loadReportConfig(reportFilePath string, debug bool) (reportConfig Report) {
	var err error
	if _, err = os.Stat(reportFilePath); err == nil {
		_, err = os.Open(reportFilePath)
		if err != nil && debug {
			fmt.Println(err)
		}
		var reportFileContent []byte
		reportFileContent, err = ioutil.ReadFile(reportFilePath)
		if err != nil && debug {
			fmt.Println(err)
		}
		err = yaml.Unmarshal(reportFileContent, &reportConfig)
		if err != nil && debug {
			fmt.Println(err)
		}
	} else if debug {
		fmt.Println(err)
	}
	return
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

func generateSpreadsheet(accountsResults accountsResults) (string, error) {
	// create spreadsheet
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

		xlsx.SetCellValue(sheetName, "A1", "SEVERITY")
		xlsx.SetCellValue(sheetName, "B1", "REGION")
		xlsx.SetCellValue(sheetName, "C1", "TEMPLATE")
		xlsx.SetCellValue(sheetName, "D1", "DATE")
		xlsx.SetCellValue(sheetName, "E1", "INSTANCE ID")
		xlsx.SetCellValue(sheetName, "F1", "INSTANCE NAME")
		xlsx.SetCellValue(sheetName, "G1", "ASG")
		xlsx.SetCellValue(sheetName, "H1", "RULES PACKAGE")
		xlsx.SetCellValue(sheetName, "I1", "TITLE")
		xlsx.SetCellValue(sheetName, "J1", "DESCRIPTION")
		xlsx.SetCellValue(sheetName, "K1", "RECOMMENDATION")

		xlsx.SetCellStyle(sheetName, "A1", "K1", headerStyle)
		xlsx.SetColWidth(sheetName, "A", "A", 15)
		xlsx.SetColWidth(sheetName, "B", "B", 13.5)
		xlsx.SetColWidth(sheetName, "C", "C", 26)
		xlsx.SetColWidth(sheetName, "D", "D", 22.5)
		xlsx.SetColWidth(sheetName, "E", "E", 19)
		xlsx.SetColWidth(sheetName, "F", "F", 24)
		xlsx.SetColWidth(sheetName, "G", "G", 20)
		xlsx.SetColWidth(sheetName, "H", "H", 44)
		xlsx.SetColWidth(sheetName, "I", "I", 60)
		xlsx.SetColWidth(sheetName, "J", "J", 70)
		xlsx.SetColWidth(sheetName, "K", "K", 150)
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
			xlsx.SetCellValue(sheetName, resultCell, dataRow.severity)
			switch dataRow.severity {
			case "HIGH":
				xlsx.SetCellStyle(sheetName, "A"+strRowNum, "A"+strRowNum, highResultStyle)
			case "MEDIUM":
				xlsx.SetCellStyle(sheetName, "A"+strRowNum, "A"+strRowNum, mediumResultStyle)
			case "LOW":
				xlsx.SetCellStyle(sheetName, "A"+strRowNum, "A"+strRowNum, lowResultStyle)
			case "INFORMATIONAL":
				xlsx.SetCellStyle(sheetName, "A"+strRowNum, "A"+strRowNum, infoResultStyle)
			case "IGNORED":
				xlsx.SetCellStyle(sheetName, "A"+strRowNum, "A"+strRowNum, ignoredResultStyle)
				if dataRow.comment != "" {
					comment := fmt.Sprintf("{\"author\":\"%s\",\"text\":\" %s\"}", "-", dataRow.comment)
					_ = xlsx.AddComment(sheetName, "A"+strRowNum, comment)
				}

			}
			xlsx.SetCellValue(sheetName, regionCell, dataRow.region)
			xlsx.SetCellValue(sheetName, templateCell, dataRow.templateName)
			xlsx.SetCellValue(sheetName, dateCell, dataRow.createdAt.Format(time.ANSIC))
			xlsx.SetCellValue(sheetName, instanceIDCell, dataRow.instanceID)
			instComment := fmt.Sprintf("{\"author\":\"%s\",\"text\":\" %s\"}", "AMI:", dataRow.amiID)
			_ = xlsx.AddComment(sheetName, instanceIDCell, instComment)
			xlsx.SetCellValue(sheetName, instanceNameCell, dataRow.instanceName)
			xlsx.SetCellValue(sheetName, rulesPackageCell, dataRow.packageName)
			xlsx.SetCellValue(sheetName, instanceASGCell, dataRow.asgName)
			xlsx.SetCellValue(sheetName, findingTitleCell, dataRow.findingTitle)
			xlsx.SetCellValue(sheetName, descriptionCell, dataRow.description)
			xlsx.SetCellValue(sheetName, recommendationCell, dataRow.recommendation)
			xlsx.SetCellStyle(sheetName, "B"+strRowNum, "B"+strRowNum, defaultCenteredStyle)
			xlsx.SetCellStyle(sheetName, "E"+strRowNum, "G"+strRowNum, defaultCenteredStyle)
			_ = xlsx.AutoFilter(sheetName, "A1", "H"+strRowNum, "")

		}
	}

	timeStamp := time.Now().UTC().Format("20060102150405")
	path := fmt.Sprintf("inspector_report_%s.xlsx", timeStamp)
	err := xlsx.SaveAs(path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	absPath, _ := filepath.Abs(path)
	fmt.Println("report written to:", absPath)
	return absPath, err
}

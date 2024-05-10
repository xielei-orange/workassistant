package common

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"regexp"
)

// ExcelRow 定义一个结构体，表示 Excel 数据
type ExcelRow struct {
	Header map[string]string // 表头和对应的单元格数据
}

// 默认require列的表头ID
// 传入表头及所需字段，返回所需表头的INT数组
func buildHeader(header []string, fields []string) []int {
	var data []int
	for i, v := range header {
		if contains(fields, v) {
			data = append(data, i)
		}
	}
	return data
}

func ReadExcel(filePath string, fields []string) ([]ExcelRow, error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}

	var data []ExcelRow
	for index := 0; index < file.SheetCount; index++ { // 修正索引范围
		sheetName := file.GetSheetName(index)
		rows, err := file.GetRows(sheetName)
		if err != nil {
			return nil, fmt.Errorf("failed to get rows from sheet %s: %v", sheetName, err)
		}

		if len(rows) == 0 {
			log.Printf("%s 表数据为空, 忽略\n", sheetName)
			continue // 如果表数据为空，继续检查下一个
		}
		header := rows[0]
		for _, v := range rows[1:] {
			resp, err := areElementsPresentAndNotEmpty(v, header, fields)
			if err != nil {
				continue
			}
			data = append(data, resp)

		}
	}
	return data, nil
}

// 检查是否在指定数组中
func isInArray(value string, array []string) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

func checkImageTag(value string) bool {
	imageNamespace := `^(nexus|family|belray|dicos|goods_selection|kong)`
	imageTag := `v\d.\d.\d{2}$`
	re, err := regexp.Compile(imageNamespace + "_prod_" + imageTag)
	if err != nil {
		panic(err)
	}
	return re.MatchString(value)
}
func verifyParameters(field string, value string) bool {
	switch field {
	case "命名空间":
		return isInArray(value, validNamespace)
	case "服务":
		return isInArray(value, validName)
	case "版本号":
		return checkImageTag(value)
	default:
		return false
	}
}

func contains(slice []string, elem string) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

// 检查数组中的特定索引是否存在且不为空
func areElementsPresentAndNotEmpty(row []string, header []string, fields []string) (ExcelRow, error) {
	indexes := buildHeader(header, fields)
	rowData := ExcelRow{Header: make(map[string]string)}
	for k, idx := range indexes {
		// 检查索引是否超出范围
		if idx >= len(row) {
			return ExcelRow{}, fmt.Errorf("index %d is out of range", idx)
		}
		// 检查元素是否为空以及是否符合规范
		if row[idx] == "" {
			return ExcelRow{}, fmt.Errorf("element at index %d is empty", idx)
		}
		// 检查元素是否为空以及是否符合规范
		if !verifyParameters(header[idx], row[idx]) {
			return ExcelRow{}, fmt.Errorf("parameter invalid")
		}
		rowData.Header[fields[k]] = row[idx]
	}
	return rowData, nil // 如果所有检查都通过，返回 true
}

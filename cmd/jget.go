package cmd

import (
	"fmt"
	"jdb/app"
	"jdb/base/data"
	"jdb/base/excel"
	"jdb/utils"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tealeg/xlsx"
)

var JGetCmd = &JGetCommand{}

var headColor = "FFFFFF"
var linkColor = "FF0000FF"
var bfColor = "6E60B0"

func init() {
	JGetCmd.JGetCmd = &cobra.Command{
		Use:              "get",
		Aliases:          []string{"g"},
		Short:            "获取数据库信息",
		Example:          "jdb get --host 127.0.0.1 --port 3306 -u root -p pass scheme",
		PreRun:           ComCmd.checkAndInit,
		Run:              JGetCmd.JGet,
		TraverseChildren: true,
	}

	// JGetCmd.JGetCmd.Flags().StringP("t", "t", "table", "获取哪部分的数据：scheme--库 table--表")
	ComCmd.addClient(JGetCmd.JGetCmd.Flags())
}

type JGetCommand struct {
	JGetCmd *cobra.Command
}

// JGet 获取数据
func (jGet *JGetCommand) JGet(cmd *cobra.Command, args []string) {

	for _, scheme := range app.Instance().JGetConf.Scheme {
		fmt.Println(fmt.Sprintf("开始处理数据库【%s】数据", scheme))
		jGet.getDB(scheme)
		fmt.Println(fmt.Sprintf("数据库【%s】数据处理完成", scheme))
	}

}

// getDB 获取一个数据库的数据
func (jGet *JGetCommand) getDB(scheme string) {

	// 获取全部表数据
	tables := (&data.Tables{}).GetByScheme(scheme)

	now := time.Now()
	fileName := fmt.Sprintf("%s%s.xlsx", scheme, now.Format("2006-01-02"))

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("表清单")
	if err != nil {
		fmt.Println("add sheet error: ", err)
		os.Exit(-1)
	}
	header := []string{"表名", "中文名", "字符集"}
	length := []float64{30, 30, 30}

	var row *xlsx.Row

	excel.AddHeaderSingle(sheet, header, length, headColor, bfColor)

	in := make(map[string]string)
	for _, table := range tables {

		tableSheet := table.Comment
		if table.Comment == "" {
			tableSheet = table.Name
		}

		// 去重，分表后数据有一致的
		if s := in[tableSheet]; s != "" {
			if s == table.Name[:len(table.Name)-2] || s[:len(s)-2] == table.Name[:len(table.Name)-2] || strings.Contains(s, table.Name) {
				fmt.Println(fmt.Sprintf("已存在表 [%s] [%s] 的同名表，忽略后续操作", tableSheet, table.Name))
				continue
			}

			tableSheet = fmt.Sprintf("表名异常_%s", tableSheet)
			in[tableSheet] = table.Name
		}

		in[tableSheet] = table.Name

		tableSheet = strings.Replace(tableSheet, "（", "_", -1)
		tableSheet = strings.Replace(tableSheet, "）", "", -1)
		tableSheet = strings.Replace(tableSheet, "(", "", -1)
		tableSheet = strings.Replace(tableSheet, ")", "", -1)
		tableSheet = strings.Replace(tableSheet, "-", "_", -1)

		// 添加表数据
		row = sheet.AddRow()
		row.SetHeightCM(1)
		excel.AddCommonCell(row, table.Name, "")

		excel.AddLinkCell(row, tableSheet, tableSheet, linkColor)

		excel.AddCommonCell(row, table.Collation, "")

		sheet, err := file.AddSheet(tableSheet)
		if err != nil {
			fmt.Println(fmt.Sprintf("添加表结构出错：%+v， 忽略此数据", err))
			row.Cells[0].GetStyle().Fill = xlsx.Fill{PatternType: xlsx.Solid_Cell_Fill, FgColor: xlsx.RGB_Dark_Red}
			continue
		}
		fmt.Println(fmt.Sprintf("开始处理数据库【%s】表【%s】数据", scheme, table.Name))
		jGet.getTable(sheet, scheme, table.Name)
		fmt.Println(fmt.Sprintf("数据库【%s】表【%s】数据处理完成", scheme, table.Name))

	}

	path := app.Instance().JGetConf.Out
	if path == "" {
		path, err = os.Getwd()
		fmt.Println(fmt.Sprintf("获取默认路径: %s", path))
		if err != nil {
			fmt.Println("获取默认目录失败")
			path = "."
		}
	}

	save := fmt.Sprintf("%s/%s", path, fileName)
	fmt.Println(fmt.Sprintf("\n保存数据库【%s】数据到【%s】", scheme, save))
	if err := file.Save(save); err != nil {
		fmt.Println(err)
	}
}

// getTable 获取表信息
func (jGet *JGetCommand) getTable(sheet *xlsx.Sheet, scheme, table string) {

	columns := (&data.Column{}).GetBySchemeAndTable(scheme, table)

	var row *xlsx.Row

	// 第一行
	row = sheet.AddRow()
	row.SetHeightCM(1)
	excel.AddBgCell(row, "表中文名", headColor, bfColor)
	excel.AddHMergeCell(row, sheet.Name, "", 6)
	excel.AddLinkCell(row, "表清单", "返回目录", linkColor)

	// 第二行
	row = sheet.AddRow()
	row.SetHeightCM(1)
	excel.AddBgCell(row, "表名称", headColor, bfColor)

	excel.AddHMergeCell(row, table, "", 7)

	// 第三行
	row = sheet.AddRow()
	row.SetHeightCM(1)
	excel.AddBgCell(row, "描述", headColor, bfColor)

	excel.AddHMergeCell(row, "", "", 7)

	header := []string{"中文名", "字段名", "字段类型", "长度", "能否为空", "默认值", "主键/索引", "extra", "描述"}
	length := []float64{20, 20, 10, 10, 10, 20, 10, 15, 30}

	excel.AddHeaderSingle(sheet, header, length, headColor, bfColor)

	// 组装字段数据
	var data [][]string
	for _, column := range columns {
		cType := ""
		l := ""
		types := strings.Split(column.Type, "(")
		if len(types) > 0 {
			cType = types[0]
		}
		if len(types) == 2 {
			l = strings.Replace(types[1], ")", "", 1)
		}

		pri := ""
		if column.Key == "PRI" {
			pri = "主键"
		} else if column.Key == "MUL" {
			pri = "普通索引"
		} else if column.Key == "UNI" {
			pri = "唯一索引"
		}

		data = append(data, []string{utils.IfNull(&column.Comment, &column.Comment, &column.Name), column.Name, cType, l, column.Nullable, column.Default, pri, column.Extra, ""})
	}

	excel.AddData(sheet, data)
}

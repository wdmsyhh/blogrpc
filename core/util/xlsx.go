package util

import (
	"errors"
	"os"

	"github.com/xuri/excelize/v2"
)

type xlsxIndexRange struct {
	verticalStart   int
	verticalEnd     int
	horizontalStart int
	horizontalEnd   int
}

type XlsxWriter struct {
	file                        *excelize.File
	sheets                      []*XlsxSheet
	styleId                     int
	style                       *excelize.Style
	defaultSheetHasBeenObtained bool
}

type XlsxSheet struct {
	writer          *XlsxWriter
	name            string
	rows            []*XlsxRow
	currentRowIndex int
	styleId         int
	style           *excelize.Style
	streamWriter    *excelize.StreamWriter
	mergeCellRanges [][]string
}

type XlsxRow struct {
	sheet                  *XlsxSheet
	indexRange             xlsxIndexRange
	firstCell              *XlsxCell
	currentHorizontalIndex int
	styleId                int
	style                  *excelize.Style
}

type XlsxCell struct {
	row              *XlsxRow
	value            interface{}
	preCells         []*XlsxCell
	nextCells        []*XlsxCell
	indexRange       *xlsxIndexRange
	cellName         string
	nextCellRowCount int
	styleId          int
	style            *excelize.Style
}

// 新建 xlsx 写入
//
// 使用示例：
//
//	writer := core_util.NewXlsxWriter()
//	sheet, err := writer.SetDefaultSheetName("订单维度")
//	if err != nil {
//	    panic(err)
//	}
//	err = sheet.SetStyle(&excelize.Style{
//	    Alignment: &excelize.Alignment{
//	        Vertical: "center",
//	    },
//	})
//	if err != nil {
//	    panic(err)
//	}
//
//	err = sheet.SetHeaders(
//	    "导购姓名",
//	    "导购手机号",
//	    "所属门店",
//	    "绑定分销员数",
//	    "绑定客户数",
//	    "推广订单数",
//	    "销售额（含运费）",
//	    "推广分销员",
//	    "分销员手机号",
//	    "周期内是否换绑",
//	    "原导购姓名/手机号",
//	    "与原导购解绑时间",
//	    "与现导购绑定时间",
//	    "订单编号",
//	    "下单时间",
//	    "订单状态",
//	    "订单实付金额（含运费和退款金额）",
//	    "运费",
//	    "退款状态",
//	    "退款金额",
//	)
//	if err != nil {
//	    panic(err)
//	}
//
//	row := sheet.AddRow()
//	cell := row.AddCells("name1", "12341234123", "门店A、门店B", 32, 209, 871, "￥9,871.00")
//	subCell := cell.AddCellsToNext("name2", "12312341234", "是", "大熊/13476577861", "2022/4/20 10:47:21", "2022/4/20 11:03:09")
//	subCell.AddCellsToNext("M2022042000000009", "2022/4/20 10:47:21", "待发货", "￥320.00", "￥6.00", "", "￥0.00")
//	subCell = cell.AddCellsToNext("name3", "12312341231", "否", "", "", "")
//	subCell.AddCellsToNext("M2022042000000002", "2022/4/20 10:47:21", "交易关闭", "￥189.00", "￥0.00", "已退款", "￥189.00")
//	subCell.AddCellsToNext("M2022041800000003", "2022/4/20 11:47:21", "交易完成", "￥220.00", "￥6.00", "部分退款", "￥32.00")
//	err = row.Flush()
//	if err != nil {
//	    panic(err)
//	}
//
//	row = sheet.AddRow()
//	cell = row.AddCells("name2", "12341234124", "门店A、门店B", 32, 209, 871, "￥9,871.00")
//	subCell = cell.AddCellsToNext("name4", "12312341232", "是", "大熊/13476577861", "2022/4/20 10:47:21", "2022/4/20 11:03:09")
//	subCell.AddCellsToNext("M2022042000000009", "2022/4/20 10:47:21", "待发货", "￥320.00", "￥6.00", "", "￥0.00")
//	subCell = cell.AddCellsToNext("name5", "12312341233", "否", "", "", "")
//	subCell.AddCellsToNext("M2022042000000002", "2022/4/20 10:47:21", "交易关闭", "￥189.00", "￥0.00", "已退款", "￥189.00")
//	subCell.AddCellsToNext("M2022041800000003", "2022/4/20 11:47:21", "交易完成", "￥220.00", "￥6.00", "部分退款", "￥32.00")
//	err = row.Flush()
//	if err != nil {
//	    panic(err)
//	}
//	err = sheet.Flush()
//	if err != nil {
//	    panic(err)
//	}
//	fp, err := os.OpenFile("./test.xlsx", os.O_RDWR|os.O_CREATE, 0755)
//	if err != nil {
//	    panic(err)
//	}
//	writer.WriteToFile(fp)
func NewXlsxWriter() *XlsxWriter {
	writer := &XlsxWriter{
		file: excelize.NewFile(),
		sheets: []*XlsxSheet{
			{
				name: "Sheet1",
			},
		},
	}
	writer.sheets[0].writer = writer
	return writer
}

// 设置表头
func (xs *XlsxSheet) SetHeaders(headers ...interface{}) error {
	if xs.currentRowIndex > 0 {
		return errors.New("setting headers is not allowed after adding rows")
	}
	row := xs.AddRow()
	row.AddCells(headers...)
	return row.Flush()
}

// 写入到指定文件
func (xw *XlsxWriter) WriteToFile(file *os.File) error {
	return xw.file.Write(file)
}

// 设置默认 Sheet 名称，默认为 Sheet1
func (xw *XlsxWriter) setDefaultSheetName(name string) (*XlsxSheet, error) {
	if len(xw.sheets) == 0 {
		sheet := &XlsxSheet{
			writer: xw,
			name:   name,
		}
		xw.sheets = append(xw.sheets, sheet)
		return sheet, nil
	}
	sheet := xw.sheets[0]
	if sheet.streamWriter != nil {
		return nil, errors.New("cannot set default sheet name after stream writer inited")
	}
	xw.file.SetSheetName(sheet.name, name)
	sheet.name = name
	return sheet, nil
}

// 获取默认 Sheet，Sheet 名称默认为 Sheet1，如果传入 name 则更新名称
func (xw *XlsxWriter) GetDefaultSheet(name string) *XlsxSheet {
	xw.defaultSheetHasBeenObtained = true
	if name != "" {
		sheet, _ := xw.setDefaultSheetName(name)
		return sheet
	}
	return xw.sheets[0]
}

// 新建 Sheet，首次新建直接返回默认 sheet
func (xw *XlsxWriter) NewSheet(name string) *XlsxSheet {
	if name == "" {
		return nil
	}
	if !xw.defaultSheetHasBeenObtained {
		xw.defaultSheetHasBeenObtained = true
		sheet, _ := xw.setDefaultSheetName(name)
		return sheet
	}
	for _, sheet := range xw.sheets {
		if name == sheet.name {
			return sheet
		}
	}
	newSheet := &XlsxSheet{
		name:   name,
		writer: xw,
	}
	xw.sheets = append(xw.sheets, newSheet)
	xw.file.NewSheet(name)
	return newSheet
}

// 设置表格样式
func (xw *XlsxWriter) SetStyle(style *excelize.Style) error {
	styleId, err := xw.file.NewStyle(style)
	if err != nil {
		return err
	}
	xw.styleId = styleId
	xw.style = style
	return nil
}

// 添加新的行
func (xs *XlsxSheet) AddRow() *XlsxRow {
	xs.currentRowIndex += 1
	row := &XlsxRow{
		sheet: xs,
		indexRange: xlsxIndexRange{
			verticalStart: xs.currentRowIndex,
			verticalEnd:   xs.currentRowIndex,
		},
	}
	xs.rows = append(xs.rows, row)
	return row
}

// 设置表格样式
func (xs *XlsxSheet) SetStyle(style *excelize.Style) error {
	styleId, err := xs.writer.file.NewStyle(style)
	if err != nil {
		return err
	}
	xs.styleId = styleId
	xs.style = style
	return nil
}

// 将 Sheet 中的数据刷新，刷新后不可继续写入
func (xs *XlsxSheet) Flush() error {
	return xs.streamWriter.Flush()
}

// 在当前行最后一列后横向追加单元格，返回追加后的最后一个单元格
// 如该行不存在单元格且未传入值，则会新建该行的第一个单元格并返回
func (xr *XlsxRow) AddCells(values ...interface{}) *XlsxCell {
	valueLen := len(values)
	if xr.firstCell == nil {
		xr.currentHorizontalIndex += 1
		xr.firstCell = &XlsxCell{
			row:   xr,
			value: "",
			indexRange: &xlsxIndexRange{
				verticalStart:   xr.indexRange.verticalStart,
				verticalEnd:     xr.indexRange.verticalEnd,
				horizontalStart: xr.currentHorizontalIndex,
				horizontalEnd:   xr.currentHorizontalIndex,
			},
		}
	}
	if valueLen == 0 {
		return xr.GetLastCell()
	}
	xr.firstCell.value = values[0]
	return xr.GetLastCell().AddCellsToNext(values[1:]...)
}

// 设置行的样式
func (xr *XlsxRow) SetStyle(style *excelize.Style) error {
	styleId, err := xr.sheet.writer.file.NewStyle(style)
	if err != nil {
		return err
	}
	xr.styleId = styleId
	xr.style = style
	return nil
}

// 获取一行的最后一列。
// 对于某一列存在多个单元格也只会定位到第一个单元格。
func (xr *XlsxRow) GetLastCell() *XlsxCell {
	if xr.firstCell == nil {
		return nil
	}
	cell := xr.firstCell
	for len(cell.nextCells) > 0 {
		cell = cell.nextCells[0]
	}
	return cell
}

// 刷新行，刷新后不可写入，每一行编辑结束后需调用此参数，否则新的行可能会出现问题
func (xr *XlsxRow) Flush() error {
	err := xr.sheet.initStreamWriter()
	if err != nil {
		return err
	}
	err = xr.firstCell.flush()
	if err != nil {
		return err
	}
	return xr.sheet.mergeCells()
}

// 在单元格后一列添加一个单元格，如果在一个单元格后多次增加单元格即为一对多。
func (xc *XlsxCell) AddCellToNext(value interface{}) *XlsxCell {
	xc.addNextCellRowCount(1)
	cell := &XlsxCell{
		row:      xc.row,
		value:    value,
		preCells: []*XlsxCell{xc},
		indexRange: &xlsxIndexRange{
			verticalStart:   xc.indexRange.verticalEnd,
			verticalEnd:     xc.indexRange.verticalEnd,
			horizontalStart: xc.indexRange.horizontalStart + 1,
			horizontalEnd:   xc.indexRange.horizontalStart + 1,
		},
	}
	xc.nextCells = append(xc.nextCells, cell)
	return cell
}

// 在单元格后多列增加多个单元格，返回添加后的最后一个 cell
func (xc *XlsxCell) AddCellsToNext(values ...interface{}) *XlsxCell {
	valueCount := len(values)
	if valueCount == 0 {
		return xc
	}
	newCell := xc.AddCellToNext(values[0])
	if valueCount-1 > 0 {
		return newCell.AddCellsToNext(values[1:valueCount]...)
	}
	return newCell
}

// 设置单元格样式
func (xc *XlsxCell) SetStyle(style *excelize.Style) error {
	styleId, err := xc.row.sheet.writer.file.NewStyle(style)
	if err != nil {
		return err
	}
	xc.styleId = styleId
	xc.style = style
	return nil
}

// 初始化写入流
func (xs *XlsxSheet) initStreamWriter() error {
	if xs.streamWriter != nil {
		return nil
	}
	var err error
	xs.streamWriter, err = xs.writer.file.NewStreamWriter(xs.name)
	return err
}

// 合并需要合并的单元格。
// 对行执行 flush 方法后会记录需要合并的单元格坐标，此方法对记录的坐标执行合并
func (xs *XlsxSheet) mergeCells() error {
	cellRange := xs.popMergeCellRange()
	for cellRange != nil {
		err := xs.streamWriter.MergeCell(cellRange[0], cellRange[1])
		if err != nil {
			return err
		}
		cellRange = xs.popMergeCellRange()
	}
	return nil
}

// 弹出一个需要合并的单元格坐标。
func (xs *XlsxSheet) popMergeCellRange() []string {
	length := len(xs.mergeCellRanges)
	if length == 0 {
		return nil
	}
	cellRange := xs.mergeCellRanges[length-1]
	if length == 1 {
		xs.mergeCellRanges = [][]string{}
	} else {
		xs.mergeCellRanges = xs.mergeCellRanges[:length-1]
	}
	return cellRange
}

// 增加下一列单元格计数，并更新前面的单元格的垂直坐标范围。
func (xc *XlsxCell) addNextCellRowCount(inc int) {
	if inc == 0 {
		return
	}
	// 新增下一列子单元格前已经存在子单元格了，则需要增加垂直坐标值。
	if xc.nextCellRowCount > 0 {
		xc.addIndexRangeVerticalEnd(inc)
	}
	xc.nextCellRowCount += inc
}

func (xc *XlsxCell) addIndexRangeVerticalEnd(inc int) {
	xc.indexRange.verticalEnd += inc
	if xc.indexRange.verticalEnd == 1 {
		return
	}
	preCellCount := len(xc.preCells)
	if preCellCount == 0 {
		xc.row.sheet.currentRowIndex += inc
		return
	}
	// 当前 preCells 只支持单个单元格。
	xc.preCells[preCellCount-1].addIndexRangeVerticalEnd(inc)
}

// 获取单元格样式ID。
// 优先级：单元格 > 行 > 表
func (xc *XlsxCell) getStyleId() int {
	if xc.styleId != 0 {
		return xc.styleId
	}
	if xc.row.styleId != 0 {
		return xc.row.styleId
	}
	if xc.row.sheet.styleId != 0 {
		return xc.row.sheet.styleId
	}
	return xc.row.sheet.writer.styleId
}

// 刷新单元格。
// 添加单元格后会执行此方法，更新当前单元格名称和垂直坐标范围，并记录需要合并的坐标。
func (xc *XlsxCell) flush() error {
	cellStart, err := excelize.CoordinatesToCellName(xc.indexRange.horizontalStart, xc.indexRange.verticalStart)
	if err != nil {
		return err
	}
	xc.cellName = cellStart
	if err := xc.writeStream(); err != nil {
		return err
	}
	// 垂直坐标发生变化
	if xc.indexRange.verticalEnd > xc.indexRange.verticalStart {
		cellEnd, err := excelize.CoordinatesToCellName(xc.indexRange.horizontalEnd, xc.indexRange.verticalEnd)
		if err != nil {
			return err
		}
		xc.row.sheet.mergeCellRanges = append(xc.row.sheet.mergeCellRanges, []string{cellStart, cellEnd})
	}
	return xc.flushNext()
}

// 刷新当前单元格的下一列单元格
func (xc *XlsxCell) flushNext() error {
	for _, cell := range xc.nextCells {
		if err := cell.flush(); err != nil {
			return err
		}
	}
	return nil
}

// 将单元格写入流
func (xc *XlsxCell) writeStream() error {
	err := xc.row.sheet.initStreamWriter()
	if err != nil {
		return err
	}
	return xc.row.sheet.streamWriter.SetRow(xc.cellName, []interface{}{excelize.Cell{Value: xc.value, StyleID: xc.getStyleId()}})
}

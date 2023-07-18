package library

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ModelFilterColumn struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	ColumnName string
	Selectable bool
	Filterable bool
}

type ModelFilterAggExpr struct {
	Column  ModelFilterColumn
	AggFunc string
}

type ModelFilterCondItem struct {
	Name  string
	Value []string
}

type ModelFilterCond struct {
	Column    ModelFilterColumn
	Condition []ModelFilterCondItem
	Operation string
}

type ModelFilterSort struct {
	Column    ModelFilterColumn
	Direction string
}

var (
	ErrInvalidColumnName        = errors.New("invalid column name")
	ErrInvalidTableName         = errors.New("invalid table name")
	ErrInvalidConditionName     = errors.New("invalid condition name")
	ErrInvalidSortDirection     = errors.New("invalid sort direction")
	ErrInvalidAggregateFunction = errors.New("invalid aggregate function")
	ErrInvalidFilterConjunction = errors.New("invalid filter conjunction function")
	ErrNoColumnSelected         = errors.New("no column selected")
	ErrInvalidPage              = errors.New("invalid page number")
	ErrInvalidPageSize          = errors.New("invalid page size number")
	ErrInvalidUnionMode         = errors.New("union mode is not valid")
	ErrInvalidLinkTable         = errors.New("invalid link table")
)

var (
	VALID_CONDITION_NAMES = map[string]bool{
		"equal":                 true,
		"not_equal":             true,
		"less_than":             true,
		"less_than_or_equal":    true,
		"greater_than":          true,
		"greater_than_or_equal": true,
		"between":               true,
		"not_between":           true,
		"in":                    true,
		"not_in":                true,
		"contains":              true,
		"not_contain":           true,
		"starts_with":           true,
		"not_starts_with":       true,
		"ends_with":             true,
		"not_ends_with":         true,
		"empty":                 true,
		"not_empty":             true,
	}

	VALID_AGGREGATE_FUNCTIONS = map[string]bool{
		"count":      true,
		"sum":        true,
		"avg":        true,
		"min":        true,
		"max":        true,
		"round":      true,
		"string_agg": true,
		"array_agg":  true,
	}
)

type ModelFilter struct {
	UnionMode struct {
		All   bool          `json:"all"`
		Items []ModelFilter `json:"items"`
	}
	Columns         []ModelFilterColumn  `json:"columns"`
	GroupBy         []ModelFilterColumn  `json:"groupBy"`
	Aggregate       []ModelFilterAggExpr `json:"aggregate"`
	Filter          []ModelFilterCond
	Sort            []ModelFilterSort          `json:"sort"`
	StringCondition string                     `json:"stringCondition"`
	Page            int                        `json:"page"`
	PageSize        int                        `json:"pageSize"`
	Distinct        bool                       `json:"distinct"`
	Search          string                     `json:"search"`
	SearchColumns   []ModelFilterColumn        `json:"searchColumns"`
	LinkTable       []ModelFilterLinkTableItem `json:"linkTable"`

	AllColumns    []ModelFilterColumn
	TenantId      int
	ParsingErrors []error
}

type ModelFilterLinkTableItem struct {
	Column1             ModelFilterColumn
	Column2             ModelFilterColumn
	Operator            string
	Mask                string
	Table               string
	TargetTableTenantId int
}

func getBoolFromAnyType(data interface{}, defaultValue bool) bool {
	if data == nil {
		return false
	}
	boolData, ok := data.(bool)
	if ok {
		return boolData
	}

	dataStr := ToString(data)
	if dataStr == "false" {
		return false
	} else {
		return defaultValue
	}
}

func getMapFromAnyType(data interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}
	mapData, ok := data.(map[string]interface{})
	if ok {
		return mapData
	}

	arrData, ok := data.([]interface{})
	if ok {
		rsl := make(map[string]interface{}, len(arrData))
		for k, v := range arrData {
			rsl[strconv.Itoa(k)] = v
		}
		return rsl
	}
	return nil
}

func NewModelFilter(allColumns []ModelFilterColumn) ModelFilter {
	return ModelFilter{
		AllColumns: allColumns,
	}
}

func (aggExpr *ModelFilterAggExpr) ToString() string {
	return fmt.Sprintf("%s(%s)", aggExpr.AggFunc, aggExpr.Column.ColumnName)
}

func (aggExpr *ModelFilterSort) ToString() string {
	return fmt.Sprintf("%s %s", aggExpr.Column.ColumnName, aggExpr.Direction)
}

// ParseToModelFilter sẽ parse dữ liệu từ map[string]interface{} thành ModelFilter
//
//	 allColumns là danh sách tất cả các cột có thể được chọn
//	 data là dữ liệu được parse
//	 tenantId : nếu tenantId là 1 giá trị bất kỳ nhỏ hơn 0 thì sẽ không áp dụng lọc theo tenant cho các bảng
//		Nếu tenantId >= 0 thì sẽ áp dụng lọc theo tenant cho tất các bảng : bảng chính và các bảng trong LinkTable
func ParseToModelFilter(allColumns []ModelFilterColumn, data map[string]interface{}, tenantId int) ModelFilter {
	var rsl = NewModelFilter(allColumns)
	rsl.parseColumns(data)
	if len(rsl.Columns) == 0 {
		panic(ErrNoColumnSelected.Error())
	}
	rsl.parseGroupBy(data)
	rsl.parseAggregate(data)
	rsl.parseFilter(data)
	rsl.parseSort(data)
	rsl.parseStringCondition(data)
	rsl.parsePage(data)
	rsl.parsePageSize(data)
	rsl.parseDistinct(data)
	rsl.parseLinkTable(data)
	rsl.parseUnionMode(data, tenantId)
	rsl.parseSearch(data)
	rsl.parseSearchColumns(data)

	rsl.setTenantIdIfNeeded(tenantId)
	return rsl
}

func (mf *ModelFilter) getColumn(name string) ModelFilterColumn {
	var rsl ModelFilterColumn
	for _, col := range mf.AllColumns {
		if col.Name == name || col.ColumnName == name {
			rsl = col
			break
		}
	}
	return rsl
}

func (mf *ModelFilter) setTenantIdIfNeeded(tenantId int) {
	if tenantId < 0 {
		mf.TenantId = -1
		for idx := range mf.LinkTable {
			mf.LinkTable[idx].TargetTableTenantId = -1
		}
		return
	}

	mf.TenantId = tenantId
	for idx := range mf.LinkTable {
		mf.LinkTable[idx].TargetTableTenantId = tenantId
	}
	tenantColumn := mf.getColumn("tenant_id_")

	// nếu không có cột tenant_id thì thêm cột tenant_id vào all columns
	if tenantColumn.ColumnName == "" {
		tenantColumn = ModelFilterColumn{
			Name:       "tenant_id",
			ColumnName: "tenant_id_",
			Title:      "tenant id",
			Type:       "number",
			Selectable: true,
			Filterable: true,
		}
		mf.AllColumns = append(mf.AllColumns, tenantColumn)
	}

	// Thêm điều kiện lọc theo tenant
	mf.Filter = append(mf.Filter, ModelFilterCond{
		Column: tenantColumn,
		Condition: []ModelFilterCondItem{
			{
				Name:  "equal",
				Value: []string{strconv.Itoa(tenantId)},
			},
		},
		Operation: "AND",
	})
}

func (mf *ModelFilter) addParsingError(err error) {
	// loop through errors and add them to the parsing errors
	mf.ParsingErrors = append(mf.ParsingErrors, err)
}

func (mf *ModelFilter) setAllColumnAsSelect() {
	// get selectable columns
	for _, col := range mf.AllColumns {
		if col.Selectable {
			mf.Columns = append(mf.Columns, col)
		}
	}
}

func (mf *ModelFilter) parseColumns(data map[string]interface{}) {
	mf.Columns = make([]ModelFilterColumn, 0)
	rawColumns := getMapFromAnyType(data["columns"])
	if rawColumns == nil {
		mf.setAllColumnAsSelect()
		return
	}

	for _, columnName := range rawColumns {
		column := mf.getColumn(ToString(columnName))
		if !column.Selectable {
			continue
		}
		if column.Name == "" {
			mf.addParsingError(fmt.Errorf(`%w: item "%s" in "columns" field`, ErrInvalidColumnName, columnName))
		}
		mf.Columns = append(mf.Columns, column)
	}

	if len(mf.Columns) == 0 {
		mf.setAllColumnAsSelect()

	}
}

func (mf *ModelFilter) parseGroupBy(data map[string]interface{}) {
	mf.GroupBy = make([]ModelFilterColumn, 0)
	rawGroupBy := getMapFromAnyType(data["groupBy"])
	if rawGroupBy == nil {
		return
	}

	for _, rawGroupByItem := range rawGroupBy {
		column := mf.getColumn(ToString(rawGroupByItem))
		if column.Name == "" {
			mf.addParsingError(fmt.Errorf(`%w: item "%s" in "groupBy" field`, ErrInvalidColumnName, ToString(rawGroupByItem)))
		}
		mf.GroupBy = append(mf.GroupBy, column)
	}
}

func (mf *ModelFilter) parseAggregate(data map[string]interface{}) {
	mf.Aggregate = make([]ModelFilterAggExpr, 0)
	rawAggregate := getMapFromAnyType(data["aggregate"])
	if rawAggregate == nil {
		return
	}

	for _, rawAggregateItem := range rawAggregate {
		rawAggregateItemMap := rawAggregateItem.(map[string]interface{})
		columnName := ToString(rawAggregateItemMap["column"])
		column := mf.getColumn(columnName)
		if column.Name == "" {
			mf.addParsingError(fmt.Errorf(`%w: item "%s" in "aggregate" field`, ErrInvalidColumnName, columnName))
		}
		aggFunc := ToString(rawAggregateItemMap["func"])
		if !VALID_AGGREGATE_FUNCTIONS[aggFunc] {
			mf.addParsingError(fmt.Errorf(`%w: function "%s" in "aggregate" field`, ErrInvalidAggregateFunction, aggFunc))
		}

		mf.Aggregate = append(mf.Aggregate, ModelFilterAggExpr{
			Column:  column,
			AggFunc: aggFunc,
		})
	}
}

func (mf *ModelFilter) validateCondName(op string) bool {
	_, ok := VALID_CONDITION_NAMES[op]
	return ok
}

func (mf *ModelFilter) parseFilter(data map[string]interface{}) {
	mf.Filter = make([]ModelFilterCond, 0)
	rawFilter := getMapFromAnyType(data["filter"])
	if rawFilter == nil {
		return
	}

	cond := ModelFilterCond{}
	for _, rawFilterItem := range rawFilter {
		rawFilterItemMap := rawFilterItem.(map[string]interface{})
		columnName := ToString(rawFilterItemMap["column"])
		column := mf.getColumn(columnName)
		if column.Name == "" {
			mf.addParsingError(fmt.Errorf(`%w: item "%s" in "filter" field`, ErrInvalidColumnName, columnName))
		}

		if !column.Filterable {
			continue
		}

		operation := ToString(rawFilterItemMap["operation"])
		operation = strings.ToUpper(operation)
		if operation != "AND" && operation != "OR" {
			mf.addParsingError(fmt.Errorf(`%w: operation "%s" in "filter" field`, ErrInvalidFilterConjunction, operation))
		}

		cond.Operation = operation
		cond.Column = column
		cond.Condition = make([]ModelFilterCondItem, 0)
		rawCondition := getMapFromAnyType(rawFilterItemMap["conditions"])
		if rawCondition == nil {
			continue
		}

		for _, rawConditionItem := range rawCondition {
			rawConditionItemMap := rawConditionItem.(map[string]interface{})
			rawValue := rawConditionItemMap["value"]
			operator := ToString(rawConditionItemMap["name"])
			if !mf.validateCondName(operator) {
				mf.addParsingError(fmt.Errorf(`%w: condition "%s" in "filter" field`, ErrInvalidConditionName, operator))
			}

			value := make([]string, 0)
			if operator == "in" || operator == "not_in" {
				rawValueSlice := getMapFromAnyType(rawValue)
				if rawValueSlice == nil {
					continue
				}
				for _, rawValueItem := range rawValueSlice {
					value = append(value, ToString(rawValueItem))
				}
			} else {
				value = append(value, ToString(rawValue))
			}

			cond.Condition = append(cond.Condition, ModelFilterCondItem{
				Value: value,
				Name:  operator,
			})
		}
		mf.Filter = append(mf.Filter, cond)
	}
}

func (mf *ModelFilter) parseSort(data map[string]interface{}) {
	mf.Sort = make([]ModelFilterSort, 0)
	if rawSort, ok := data["sort"]; ok {
		rawSortSlice := getMapFromAnyType(rawSort)
		if rawSortSlice == nil {
			return
		}
		for _, rawSortItem := range rawSortSlice {
			rawSortItemMap := rawSortItem.(map[string]interface{})
			columnName := ToString(rawSortItemMap["column"])
			column := mf.getColumn(columnName)
			if column.Name == "" {
				mf.addParsingError(fmt.Errorf(`%w: item "%s" in "sort" field`, ErrInvalidColumnName, columnName))
			}
			rawDirection := rawSortItemMap["type"]
			if rawDirection == nil {
				mf.addParsingError(fmt.Errorf(`%w: order not found field`, ErrInvalidSortDirection))
				continue
			}
			direction := ToString(rawDirection)
			if direction != "asc" && direction != "desc" {
				mf.addParsingError(fmt.Errorf(`%w: direction "%s" in "sort" field`, ErrInvalidSortDirection, direction))
			}
			mf.Sort = append(mf.Sort, ModelFilterSort{
				Column:    column,
				Direction: direction,
			})
			mf.AddColumnSeclect(column)
		}
	}
}

func (mf *ModelFilter) AddColumnSeclect(column ModelFilterColumn) {
	colExisted := false
	for _, col := range mf.Columns {
		if col.Name == column.Name {
			colExisted = true
			break
		}
	}
	if !colExisted {
		mf.Columns = append(mf.Columns, column)
	}
}

func (mf *ModelFilter) AddFilterCond(cond ModelFilterCond) {
	var oldCond *ModelFilterCond
	for _, c := range mf.Filter {
		if c.Column.Name == cond.Column.Name {
			oldCond = &c
			break
		}
	}

	if oldCond == nil {
		mf.Filter = append(mf.Filter, cond)
		return
	} else {
		oldCond.Condition = append(oldCond.Condition, cond.Condition...)
		return
	}
}

func (mf *ModelFilter) parseStringCondition(data map[string]interface{}) {
	mf.StringCondition = ""
	if rawStringCondition, ok := data["stringCondition"]; ok {
		mf.StringCondition = ToString(rawStringCondition)
	}
}

func (mf *ModelFilter) parsePage(data map[string]interface{}) {
	mf.Page = 1
	if rawPage, ok := data["page"]; ok {
		str := ToString(rawPage)
		// string to int
		page, err := strconv.Atoi(str)
		if err != nil || page < 1 {
			mf.addParsingError(fmt.Errorf(`%w: page "%s" in "page" field`, ErrInvalidPage, str))
		} else {
			mf.Page = page
		}
	}
}

func (mf *ModelFilter) parsePageSize(data map[string]interface{}) {
	mf.PageSize = 50
	if rawPageSize, ok := data["pageSize"]; ok {
		str := ToString(rawPageSize)
		// string to int
		pageSize, err := strconv.Atoi(str)
		if err != nil || pageSize < 0 {
			mf.addParsingError(fmt.Errorf(`%w: pageSize "%s" in "pageSize" field`, ErrInvalidPageSize, str))
		} else {
			mf.PageSize = pageSize
		}
	}
}

func (mf *ModelFilter) parseDistinct(data map[string]interface{}) {
	mf.Distinct = getBoolFromAnyType(data["distinct"], false)
}

func (mf *ModelFilter) parseSearch(data map[string]interface{}) {
	mf.Search = ""
	if rawSearch, ok := data["search"]; ok {
		mf.Search = ToString(rawSearch)
	}
}

func (mf *ModelFilter) parseSearchColumns(data map[string]interface{}) {
	mf.SearchColumns = make([]ModelFilterColumn, 0)
	if mf.Search == "" {
		return
	}

	// Lấy các cột được chỉ định trong searchColumns
	if rawSearchColumns, ok := data["searchColumns"]; ok {
		rawSearchColumnsSlice := ToString(rawSearchColumns)
		rawColumnName := strings.Split(rawSearchColumnsSlice, ",")
		for _, rawColumnNameItem := range rawColumnName {
			column := mf.getColumn(rawColumnNameItem)
			if column.Name == "" {
				mf.addParsingError(fmt.Errorf(`%w: item "%s" in "searchColumns" field`, ErrInvalidColumnName, rawColumnNameItem))
				continue
			}
			mf.SearchColumns = append(mf.SearchColumns, column)
		}
	}

	// Nếu không có searchColumns thì lấy tất cả các cột làm searchColumns
	if len(mf.SearchColumns) == 0 {
		mf.SearchColumns = mf.AllColumns
	}
}

func CheckValidTableName(tbName string) bool {
	rsl, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, tbName)
	return rsl
}

func (mf *ModelFilter) parseLinkTable(data map[string]interface{}) {
	mf.LinkTable = make([]ModelFilterLinkTableItem, 0)
	rawLinkTableSlice := getMapFromAnyType(data["linkTable"])

	if rawLinkTableSlice != nil {
		for _, rawLinkTableItem := range rawLinkTableSlice {
			rawLinkTableItemMap := rawLinkTableItem.(map[string]interface{})
			columnName1 := ToString(rawLinkTableItemMap["column1"])
			column1 := mf.getColumn(columnName1)
			if column1.Name == "" {
				mf.addParsingError(fmt.Errorf(`%w: item "%s" in "linkTable" field`, ErrInvalidColumnName, columnName1))
			}

			columnName2 := ToString(rawLinkTableItemMap["column2"])
			column2 := ModelFilterColumn{
				Name:       columnName2,
				Title:      columnName2,
				Type:       "text",
				ColumnName: columnName2,
				Selectable: true,
				Filterable: true,
			}

			mask, ok := rawLinkTableItemMap["mask"].(string)
			if !ok || mask == "" {
				mf.addParsingError(fmt.Errorf(`%w: item "%s" in "linkTable" field`, ErrInvalidColumnName, columnName2))
			}

			table, ok := rawLinkTableItemMap["table"].(string)
			if !ok || table == "" || CheckValidTableName(table) {
				mf.addParsingError(fmt.Errorf(`%w: table name "%s" in "linkTable" field`, ErrInvalidTableName, columnName2))
			}

			mf.LinkTable = append(mf.LinkTable, ModelFilterLinkTableItem{
				Column1:  column1,
				Column2:  column2,
				Operator: "=",
				Mask:     mask,
				Table:    table,
			})
		}
	} else {
		mf.addParsingError(fmt.Errorf(`%w: "linkTable" field is not object or array`, ErrInvalidLinkTable))
	}
}

func (mf *ModelFilter) parseUnionMode(data map[string]interface{}, tenantId int) {
	mf.UnionMode.Items = make([]ModelFilter, 0)
	mf.UnionMode.All = true

	if rawUnionMode, ok := data["unionMode"]; ok {
		rawUnionModeSlice := getMapFromAnyType(rawUnionMode)
		if rawUnionModeSlice == nil {
			mf.ParsingErrors = append(mf.ParsingErrors, fmt.Errorf(`%w: "unionMode" field`, ErrInvalidUnionMode))
			return
		}
		mf.UnionMode.All = getBoolFromAnyType(rawUnionModeSlice["all"], true)

		if items := getMapFromAnyType(rawUnionModeSlice["items"]); items != nil {
			for _, rawUnionModeItem := range items {
				rawUnionModeItemMap := rawUnionModeItem.(map[string]interface{})
				newMf := ParseToModelFilter(mf.AllColumns, rawUnionModeItemMap, tenantId)
				for _, e := range newMf.ParsingErrors {
					mf.addParsingError(e)
				}
				mf.UnionMode.Items = append(mf.UnionMode.Items, newMf)
			}
		} else {
			mf.ParsingErrors = append(mf.ParsingErrors, fmt.Errorf(`%w: "items" field is not object or array`, ErrInvalidUnionMode))
		}
	}
}

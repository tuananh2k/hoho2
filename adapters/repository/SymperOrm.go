package repository

import (
	"fmt"
	"hoho-framework-v2/library"
	"reflect"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"
)

type SymperOrm struct {
	pgdb *pg.DB
}
type anynomousModel struct{}

func isStruct(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Struct
}

func (sorm *SymperOrm) ModelWithTenant(model interface{}) *orm.Query {
	var value reflect.Value
	if isStruct(model) {
		value = reflect.ValueOf(model)
	} else {
		value = reflect.ValueOf(model).Elem()
	}
	field := value.FieldByName("TenantId")
	if !field.IsValid() {
		panic("[HOHO FRAMEWORK ERROR]: param in 'ModelWithTenant' dose not have 'TenantId' field")
	}
	tenantId := field.Int()
	return sorm.pgdb.Model(model).Where("tenant_id_ = ?", tenantId)
}

func (sorm *SymperOrm) CloseDB() error {
	return sorm.pgdb.Close()
}

func NewSymperOrm(pgdb *pg.DB) *SymperOrm {
	return &SymperOrm{
		pgdb: pgdb,
	}
}

func (sorm *SymperOrm) ModelFilterWithModel(model interface{}, mdf library.ModelFilter) *orm.Query {
	return sorm.ModelFilterWithOrmInst(sorm.pgdb.Model(model), mdf)
}

func (sorm *SymperOrm) ModelFilterWithOrmInst(ormInst *orm.Query, mdf library.ModelFilter) *orm.Query {
	if len(mdf.UnionMode.Items) > 0 {
		ormInst = sorm.SetUnionModelFilter(ormInst, mdf)
		return ormInst
	}

	sorm.SetFromByModelFilter(ormInst, mdf)
	sorm.SetWhereByModelFilter(ormInst, mdf)
	sorm.SetGroupByByModelFilter(ormInst, mdf)
	sorm.SetSelectByModelFilter(ormInst, mdf)
	sorm.SetSortByModelFilter(ormInst, mdf)
	sorm.SetLimitByByModelFilter(ormInst, mdf)

	return ormInst
}

func (sorm *SymperOrm) ModelFilterWithTbName(tbName string, mdf library.ModelFilter) *orm.Query {
	obj := new(anynomousModel)
	var query *orm.Query
	if library.CheckValidTableName(tbName) {
		query = sorm.pgdb.Model(obj)
	} else {
		query = sorm.pgdb.Model(obj).TableExpr(fmt.Sprintf(" (%s) AS tb1 ", tbName))
		tbName = "tb1"
	}
	tb := query.TableModel().Table()
	tb.SQLNameForSelects = types.Safe(tbName)
	return sorm.ModelFilterWithOrmInst(query, mdf)
}

func ToPlainQuery(query *orm.Query) string {
	str, err := query.AppendQuery(orm.NewFormatter(), nil)
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (sorm *SymperOrm) SetSelectByModelFilter(ormInst *orm.Query, mdf library.ModelFilter) {
	if len(mdf.Aggregate) > 0 {
		for _, aggItem := range mdf.Aggregate {
			aggItem.Column.ColumnName = getColumnNameWithPrefixIfNeeded(aggItem.Column, mdf)
			ormInst.ColumnExpr(aggItem.ToString())
		}
	} else {
		for _, col := range mdf.Columns {
			col.ColumnName = getColumnNameWithPrefixIfNeeded(col, mdf)
			if library.CheckValidTableName(col.ColumnName) {
				ormInst.Column(col.ColumnName)
			} else {
				ormInst.ColumnExpr(col.ColumnName)
			}
		}
		if mdf.Distinct {
			ormInst.Distinct()
		}
	}

	if len(mdf.GroupBy) > 0 {
		for _, groupItem := range mdf.GroupBy {
			groupItem.ColumnName = getColumnNameWithPrefixIfNeeded(groupItem, mdf)
			ormInst.ColumnExpr(groupItem.ColumnName)
		}
	}
}

func (sorm *SymperOrm) SetFromByModelFilter(ormInst *orm.Query, mdf library.ModelFilter) {
	var joinStr string
	var hadColumn bool

	if len(mdf.LinkTable) > 0 {
		for idx, linkTbItem := range mdf.LinkTable {
			hadColumn = false
			// Thêm các section join vào from
			joinStr = fmt.Sprintf(`LEFT JOIN %s AS tb%d ON tb1.%s = tb%d.%s`, linkTbItem.Table, idx+2, linkTbItem.Column1.ColumnName, idx+2, linkTbItem.Column2.ColumnName)
			if linkTbItem.TargetTableTenantId >= 0 {
				joinStr = fmt.Sprintf("%s AND tb%d.tenant_id_ = %d", joinStr, idx+2, linkTbItem.TargetTableTenantId)
			}
			ormInst.Join(joinStr)
			for colIdx := range mdf.Columns {
				col := &mdf.Columns[colIdx]
				if col.ColumnName == linkTbItem.Column1.ColumnName {
					col.ColumnName = fmt.Sprintf(`tb%d.%s AS %s`, idx+2, linkTbItem.Mask, linkTbItem.Column1.ColumnName)
					hadColumn = true
				}
			}

			if !hadColumn {
				newColumn := library.ModelFilterColumn{
					ColumnName: fmt.Sprintf(`tb%d.%s AS %s`, idx+2, linkTbItem.Mask, linkTbItem.Column1.ColumnName),
					Type:       "text",
					Selectable: true,
					Filterable: true,
				}
				mdf.Columns = append(mdf.Columns, newColumn)
			}
		}

		ormInst.TableModel().Table().Alias = "tb1"
	}
}

func (sorm *SymperOrm) SetWhereByModelFilter(ormInst *orm.Query, mdf library.ModelFilter) {
	condList := make([]string, 0)
	for _, filterConds := range mdf.Filter {
		filterConds.Column.ColumnName = getColumnNameWithPrefixIfNeeded(filterConds.Column, mdf)
		if condItemStr := sorm.getColumnCondition(filterConds); condItemStr != "" {
			condList = append(condList, condItemStr)
		}
	}

	if mdf.Search != "" && len(mdf.SearchColumns) > 0 {
		ormInst.WhereOrGroup(func(q *pg.Query) (*pg.Query, error) {
			for _, searchCol := range mdf.SearchColumns {
				if searchCol.Type != "text" {
					continue
				}
				q.WhereOr(fmt.Sprintf(`%s ILIKE '%%%s%%'`, searchCol.ColumnName, escapeString(mdf.Search)))
			}
			return q, nil
		})
	}

	if mdf.StringCondition != "" {
		ormInst.Where(mdf.StringCondition)
	}

	if len(condList) > 0 {
		ormInst.Where(strings.Join(condList, " AND "))
	}
}

func trimQuoteIfNeeded(dataType string, value string) string {
	if dataType != "number" {
		value = strings.Trim(value, "'")
	}
	return value
}

func (sorm *SymperOrm) getColumnCondition(cond library.ModelFilterCond) string {
	colName := cond.Column.ColumnName

	condItems := make([]string, 0)
	value := ""
	itemCondStr := ""

	for _, item := range cond.Condition {
		if len(item.Value) == 0 {
			continue
		}

		if item.Name == "in" || item.Name == "not_in" {
			// Phép toán theo loại : có thể nhận nhiều giá trị
			valueArr := make([]string, 0)
			for _, vitem := range item.Value {
				if vitem == "" && (cond.Column.Type == "date" || cond.Column.Type == "datetime" || cond.Column.Type == "number") {
					continue
				}
				vitem = escapeString(vitem)
				if cond.Column.Type != "number" {
					vitem = fmt.Sprintf(`'%s'`, vitem)
				}
				valueArr = append(valueArr, vitem)
			}
			if len(valueArr) == 0 {
				continue
			}
			value = fmt.Sprintf(`(%s)`, strings.Join(valueArr, ","))

			if item.Name == "in" {
				itemCondStr = fmt.Sprintf(`%s IN %s`, colName, value)

			} else if item.Name == "not_in" {
				itemCondStr = fmt.Sprintf(`%s NOT IN %s`, colName, value)
			}

		} else if item.Name == "empty" || item.Name == "not_empty" {
			// Phép toán theo loại : không nhận giá trị nào
			if item.Name == "empty" {
				if cond.Column.Type == "text" {
					itemCondStr = fmt.Sprintf(`%s = '' OR %s IS NULL`, colName, colName)
				} else {
					itemCondStr = fmt.Sprintf(`%s IS NULL`, colName)
				}

			} else if item.Name == "not_empty" {
				if cond.Column.Type == "text" {
					itemCondStr = fmt.Sprintf(`%s != '' OR %s IS NOT NULL`, colName, colName)
				} else {
					itemCondStr = fmt.Sprintf(`%s IS NOT NULL`, colName)
				}
			}

		} else {
			value = escapeString(item.Value[0])
			// Phép toán theo loại : chỉ nhận 1 giá trị
			if value == "" {
				continue
			}

			if cond.Column.Type != "number" {
				value = fmt.Sprintf(`'%s'`, value)
			}

			if item.Name == "equal" {
				itemCondStr = fmt.Sprintf(`%s = %s`, colName, value)

			} else if item.Name == "not_equal" {
				itemCondStr = fmt.Sprintf(`%s != %s`, colName, value)

			} else if item.Name == "less_than" {
				itemCondStr = fmt.Sprintf(`%s < %s`, colName, value)

			} else if item.Name == "less_than_or_equal" {
				itemCondStr = fmt.Sprintf(`%s <= %s`, colName, value)

			} else if item.Name == "greater_than" {
				itemCondStr = fmt.Sprintf(`%s > %s`, colName, value)

			} else if item.Name == "greater_than_or_equal" {
				itemCondStr = fmt.Sprintf(`%s >= %s`, colName, value)

			} else if item.Name == "contains" {
				itemCondStr = fmt.Sprintf(`%s ILIKE '%%%s%%'`, colName, trimQuoteIfNeeded(cond.Column.Type, value))

			} else if item.Name == "not_contains" {
				itemCondStr = fmt.Sprintf(`%s NOT ILIKE '%%%s%%'`, colName, trimQuoteIfNeeded(cond.Column.Type, value))

			} else if item.Name == "starts_with" {
				itemCondStr = fmt.Sprintf(`%s ILIKE '%s%%'`, colName, trimQuoteIfNeeded(cond.Column.Type, value))

			} else if item.Name == "not_starts_with" {
				itemCondStr = fmt.Sprintf(`%s NOT ILIKE '%s%%'`, colName, trimQuoteIfNeeded(cond.Column.Type, value))

			} else if item.Name == "ends_with" {
				itemCondStr = fmt.Sprintf(`%s ILIKE '%%%s'`, colName, trimQuoteIfNeeded(cond.Column.Type, value))

			} else if item.Name == "not_ends_with" {
				itemCondStr = fmt.Sprintf(`%s NOT ILIKE '%%%s'`, colName, trimQuoteIfNeeded(cond.Column.Type, value))
			}
		}

		condItems = append(condItems, fmt.Sprintf(`(%s)`, itemCondStr))
	}
	return strings.Join(condItems, cond.Operation)
}

func (sorm *SymperOrm) SetSortByModelFilter(ormInst *orm.Query, mdf library.ModelFilter) {
	for _, sortItem := range mdf.Sort {
		sortItem.Column.ColumnName = getColumnNameWithPrefixIfNeeded(sortItem.Column, mdf)
		ormInst.OrderExpr(sortItem.ToString())
	}
}

func (sorm *SymperOrm) SetGroupByByModelFilter(ormInst *orm.Query, mdf library.ModelFilter) {
	for _, groupItem := range mdf.GroupBy {
		groupItem.ColumnName = getColumnNameWithPrefixIfNeeded(groupItem, mdf)
		ormInst.Group(groupItem.ColumnName)
	}
}

func (sorm *SymperOrm) SetLimitByByModelFilter(ormInst *orm.Query, mdf library.ModelFilter) {
	ormInst.Offset(mdf.PageSize * (mdf.Page - 1)).Limit(mdf.PageSize)
}

func (sorm *SymperOrm) SetUnionModelFilter(ormInst *orm.Query, mdf library.ModelFilter) *orm.Query {
	firstInst := sorm.ModelFilterWithOrmInst(ormInst.Clone(), mdf.UnionMode.Items[0])
	for i := 1; i < len(mdf.UnionMode.Items); i++ {
		ormItem := sorm.ModelFilterWithOrmInst(ormInst.Clone(), mdf.UnionMode.Items[i])
		if mdf.UnionMode.All {
			firstInst.UnionAll(ormItem)
		} else {
			firstInst.Union(ormItem)
		}
	}
	return firstInst
}

func getColumnNameWithPrefixIfNeeded(column library.ModelFilterColumn, mdf library.ModelFilter) string {
	if strings.Contains(column.ColumnName, " AS ") {
		return column.ColumnName
	}

	if len(mdf.LinkTable) > 0 {
		for idx, linkTbItem := range mdf.LinkTable {
			if column.ColumnName == linkTbItem.Column1.ColumnName {
				return fmt.Sprintf("tb%d.%s", idx+2, linkTbItem.Mask)
			}
		}
		return "tb1." + column.ColumnName
	} else {
		return column.ColumnName
	}
}

func escapeString(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func (sorm *SymperOrm) GetModelFilterCountQuery(tbName string, mdf library.ModelFilter) *orm.Query {
	if len(mdf.UnionMode.Items) > 0 || len(mdf.GroupBy) > 0 {
		return nil
	}

	obj := new(anynomousModel)
	var ormInst *orm.Query
	if library.CheckValidTableName(tbName) {
		ormInst = sorm.pgdb.Model(obj)
	} else {
		ormInst = sorm.pgdb.Model(obj).TableExpr(fmt.Sprintf(" (%s) AS tb1 ", tbName))
		tbName = "tb1"
	}
	tb := ormInst.TableModel().Table()
	tb.SQLNameForSelects = types.Safe(tbName)

	sorm.SetFromByModelFilter(ormInst, mdf)
	sorm.SetWhereByModelFilter(ormInst, mdf)
	ormInst.ColumnExpr("COUNT(*) AS total")

	return ormInst
}

func (sorm *SymperOrm) GetTableNameAndAllColumnOfObj(model interface{}) (string, []string) {
	table := sorm.pgdb.Model(model).Table()
	fieldsMap := table.TableModel().Table().FieldsMap
	tbName := table.TableModel().Table().SQLNameForSelects
	columns := make([]string, 0)
	for _, field := range fieldsMap {
		columns = append(columns, field.SQLName)
	}
	return strings.ReplaceAll(string(tbName), "\"", ""), columns
}

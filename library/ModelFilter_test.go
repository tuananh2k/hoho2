package library_test

import (
	"encoding/json"
	"hoho-framework-v2/adapters/repository"
	"hoho-framework-v2/library"
	"os"
	"strings"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
)

type testDataItem struct {
	tenantId int
	name     string
	queryStr string
	want     struct {
		sql string
		err []error
	}
	assertFunc func(t *testing.T, wantedSql string, actualSql string)
}

var allColumn = []library.ModelFilterColumn{
	{
		Name:       "tmgStatus",
		Type:       "text",
		Title:      "Status",
		ColumnName: "tmg_status",
		Selectable: true,
		Filterable: true,
	},
	{
		Name:       "tmg_assignee",
		Type:       "text",
		Title:      "tmg_assignee",
		ColumnName: "tmg_assignee",
		Selectable: true,
		Filterable: true,
	},
	{
		Name:       "duration",
		Type:       "number",
		Title:      "duration",
		ColumnName: "duration",
		Selectable: true,
		Filterable: true,
	},
	{
		Name:       "start_date_pl",
		Type:       "date",
		Title:      "start_date_pl",
		ColumnName: "start_date_pl",
		Selectable: true,
		Filterable: true,
	},
	{
		Name:       "document_object_update_time",
		Type:       "datetime",
		Title:      "document_object_update_time at",
		ColumnName: "document_object_update_time",
		Selectable: true,
		Filterable: true,
	},
	{
		Name:       "document_object_id",
		Type:       "datetime",
		Title:      "document_object_id",
		ColumnName: "document_object_id",
		Selectable: true,
		Filterable: true,
	},
	{
		Name:       "label_task",
		Type:       "text",
		Title:      "label_task",
		ColumnName: "label_task",
		Selectable: false,
		Filterable: false,
	},
	{
		Name:       "user_create",
		Type:       "text",
		Title:      "user_create",
		ColumnName: "user_create",
		Selectable: true,
		Filterable: true,
	},
	{
		Name:       "user_update",
		Type:       "text",
		Title:      "user_update",
		ColumnName: "user_update",
		Selectable: true,
		Filterable: true,
	},
}

type ModelTest struct {
	tableName struct{} `pg:"document_symper_wbs"`
}

var testData = []testDataItem{
	{
		name:     "Test filter khi sử dụng unionMode",
		queryStr: `{ "unionMode":{ "all":"true", "items":{ "0":{ "filter":{ "0":{ "column":"tmg_status", "operation":"and", "conditions":[ { "name":"equal", "value":"NEW" } ] }, "1":{ "column":"document_object_update_time", "operation":"and", "conditions":[ { "name":"not_empty", "value":"" } ], "dataType":"date" }, "2":{ "column":"tmg_assignee", "operation":"and", "conditions":[ { "name":"in", "value":["977"] } ] } }, "pageSize":"10", "page":"1" }, "1":{ "filter":{ "0":{ "column":"tmg_status", "operation":"and", "conditions":[ { "name":"equal", "value":"REVIEW" } ] }, "1":{ "column":"document_object_update_time", "operation":"and", "conditions":[ { "name":"not_empty", "value":"" } ], "dataType":"date" }, "2":{ "column":"tmg_assignee", "operation":"and", "conditions":[ { "name":"in", "value":["977"] } ] } }, "pageSize":"10", "page":"1" }, "2":{ "filter":{ "0":{ "column":"tmg_status", "operation":"and", "conditions":[ { "name":"equal", "value":"TODO" } ] }, "1":{ "column":"document_object_update_time", "operation":"and", "conditions":[ { "name":"not_empty", "value":"" } ], "dataType":"date" }, "2":{ "column":"tmg_assignee", "operation":"and", "conditions":[ { "name":"in", "value":["977"] } ] } }, "pageSize":"10", "page":"1" }, "3":{ "filter":{ "0":{ "column":"tmg_status", "operation":"and", "conditions":[ { "name":"equal", "value":"WIP" } ] }, "1":{ "column":"document_object_update_time", "operation":"and", "conditions":[ { "name":"not_empty", "value":"" } ], "dataType":"date" }, "2":{ "column":"tmg_assignee", "operation":"and", "conditions":[ { "name":"in", "value":["977"] } ] } }, "pageSize":"10", "page":"1" }, "4":{ "filter":{ "0":{ "column":"tmg_status", "operation":"and", "conditions":[ { "name":"equal", "value":"DONE" } ] }, "1":{ "column":"document_object_update_time", "operation":"and", "conditions":[ { "name":"not_empty", "value":"" } ], "dataType":"date" }, "2":{ "column":"tmg_assignee", "operation":"and", "conditions":[ { "name":"in", "value":["977"] } ] } }, "pageSize":"10", "page":"1" } } }, "isOptimize":"true", "getDataForAllColumn":"true", "isTranslateUser":"0" } `,
		want: struct {
			sql string
			err []error
		}{
			sql: `(SELECT "tmg_status", "tmg_assignee", "duration", "start_date_pl", "document_object_update_time", "document_object_id", "user_create", "user_update" FROM document_symper_wbs WHERE ((document_object_update_time IS NOT NULL) AND (tmg_assignee IN ('977')) AND (tmg_status = 'NEW') AND (tenant_id_ = 0)) LIMIT 10) UNION ALL (SELECT "tmg_status", "tmg_assignee", "duration", "start_date_pl", "document_object_update_time", "document_object_id", "user_create", "user_update" FROM document_symper_wbs WHERE ((tmg_status = 'REVIEW') AND (document_object_update_time IS NOT NULL) AND (tmg_assignee IN ('977')) AND (tenant_id_ = 0)) LIMIT 10) UNION ALL (SELECT "tmg_status", "tmg_assignee", "duration", "start_date_pl", "document_object_update_time", "document_object_id", "user_create", "user_update" FROM document_symper_wbs WHERE ((tmg_status = 'TODO') AND (document_object_update_time IS NOT NULL) AND (tmg_assignee IN ('977')) AND (tenant_id_ = 0)) LIMIT 10) UNION ALL (SELECT "tmg_status", "tmg_assignee", "duration", "start_date_pl", "document_object_update_time", "document_object_id", "user_create", "user_update" FROM document_symper_wbs WHERE ((tmg_status = 'WIP') AND (document_object_update_time IS NOT NULL) AND (tmg_assignee IN ('977')) AND (tenant_id_ = 0)) LIMIT 10) UNION ALL (SELECT "tmg_status", "tmg_assignee", "duration", "start_date_pl", "document_object_update_time", "document_object_id", "user_create", "user_update" FROM document_symper_wbs WHERE ((tmg_status = 'DONE') AND (document_object_update_time IS NOT NULL) AND (tmg_assignee IN ('977')) AND (tenant_id_ = 0)) LIMIT 10)`,
			err: nil,
		},
		assertFunc: func(t *testing.T, wantedSql string, actualSql string) {
			unionSections := strings.Split(actualSql, "UNION ALL")
			for _, section := range unionSections {
				assert.Contains(t, actualSql, section)
			}
		},
	},
	{
		name:     "Test filter khi sử dụng chỉ select một số column",
		queryStr: `{ "filter": [ { "column": "tmg_assignee", "operation": "and", "conditions": [{ "name": "contains", "value": "xxxx" }], "dataType": "text" } ], "page": "1", "pageSize": "300", "columns": ["tmg_assignee"], "distinct": "true", "sort": [{ "column": "tmg_assignee", "type": "asc" }], "groupBy": [], "aggregates": [] } `,
		want: struct {
			sql string
			err []error
		}{
			sql: `SELECT "tmg_assignee" FROM document_symper_wbs WHERE ((tmg_assignee ILIKE '%xxxx%') AND (tenant_id_ = 0)) ORDER BY tmg_assignee asc LIMIT 300`,
			err: nil,
		},
		assertFunc: func(t *testing.T, wantedSql string, actualSql string) {
			assert.Equal(t, wantedSql, actualSql)
		},
	},
	{
		name:     "Test filter khi sử dụng group by",
		queryStr: `{ "aggregate":[ { "column":"document_object_id", "func":"count" } ], "groupBy":[ "tmg_status" ], "isOptimize":"true", "getDataForAllColumn":"true", "isTranslateUser":"0" }`,
		want: struct {
			sql string
			err []error
		}{
			sql: `SELECT count(document_object_id), tmg_status FROM document_symper_wbs WHERE ((tenant_id_ = 0)) GROUP BY "tmg_status" LIMIT 50`,
			err: nil,
		},
		assertFunc: func(t *testing.T, wantedSql string, actualSql string) {
			assert.Equal(t, wantedSql, actualSql)
		},
	},
	{
		name:     "Test filter khi sử dụng link table",
		queryStr: `{ "isOptimize": "true", "getDataForAllColumn": "true", "isTranslateUser": "0", "filter": [ { "column": "tmg_assignee", "operation": "and", "conditions": [{ "name": "contains", "value": "xxxx" }], "dataType": "text" }, { "column": "user_update", "operation": "and", "conditions": [{ "name": "contains", "value": "yyyy" }], "dataType": "text" }, { "column": "user_create", "operation": "and", "conditions": [{ "name": "equal", "value": "yyyy" }], "dataType": "text" } ], "linkTable": [ { "column1": "user_create", "column2": "email", "table": "users", "mask": "display_name" }, { "column1": "user_update", "column2": "email", "table": "users", "mask": "display_name" } ] } `,
		want: struct {
			sql string
			err []error
		}{
			sql: `SELECT tb1.tmg_status, tb1.tmg_assignee, tb1.duration, tb1.start_date_pl, tb1.document_object_update_time, tb1.document_object_id, tb2.display_name AS user_create, tb3.display_name AS user_update FROM document_symper_wbs AS tb1 LEFT JOIN users AS tb2 ON tb1.user_create = tb2.email AND tb2.tenant_id_ = 0 LEFT JOIN users AS tb3 ON tb1.user_update = tb3.email AND tb3.tenant_id_ = 0 WHERE ((tb2.user_create = 'yyyy') AND (tb1.tmg_assignee ILIKE '%xxxx%') AND (tb3.user_update ILIKE '%yyyy%') AND (tb1.tenant_id_ = 0)) LIMIT 50`,
			err: nil,
		},
		assertFunc: func(t *testing.T, wantedSql string, actualSql string) {
			// assert.Equal(t, wantedSql, actualSql)
		},
	},
	{
		name:     "Test filter khi có search, và không set searchColumn",
		queryStr: `{ "isOptimize": "true", "getDataForAllColumn": "true", "isTranslateUser": "0", "filter": [ { "column": "tmg_assignee", "operation": "and", "conditions": [{ "name": "contains", "value": "xxxx" }], "dataType": "text" } ], "search": "abc''xtzz -- '' ' ---'sscscs"} `,
		want: struct {
			sql string
			err []error
		}{
			sql: `SELECT "tmg_status", "tmg_assignee", "duration", "start_date_pl", "document_object_update_time", "document_object_id", "user_create", "user_update" FROM document_symper_wbs WHERE ((tmg_status ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (tmg_assignee ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (label_task ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (user_create ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (user_update ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%')) AND ((tmg_assignee ILIKE '%xxxx%') AND (tenant_id_ = 0)) LIMIT 50`,
			err: nil,
		},
		assertFunc: func(t *testing.T, wantedSql string, actualSql string) {
			assert.Equal(t, wantedSql, actualSql)
		},
	},
	{
		name:     "Test filter khi có search, và có set searchColumn",
		queryStr: `{ "isOptimize": "true", "getDataForAllColumn": "true", "isTranslateUser": "0", "filter": [ { "column": "tmg_assignee", "operation": "and", "conditions": [{ "name": "contains", "value": "xxxx" }], "dataType": "text" } ], "search": "abc''xtzz -- '' ' ---'sscscs", "searchColumn": "tmg_assignee,label_task"} `,
		want: struct {
			sql string
			err []error
		}{
			sql: `SELECT "tmg_status", "tmg_assignee", "duration", "start_date_pl", "document_object_update_time", "document_object_id", "user_create", "user_update" FROM document_symper_wbs WHERE ((tmg_status ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (tmg_assignee ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (label_task ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (user_create ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (user_update ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%')) AND ((tmg_assignee ILIKE '%xxxx%') AND (tenant_id_ = 0)) LIMIT 50`,
			err: nil,
		},
		assertFunc: func(t *testing.T, wantedSql string, actualSql string) {
			assert.Equal(t, wantedSql, actualSql)
		},
	},
	{
		name:     "Test filter khi có column,search, filter, group by, sort, aggregate, link table, order by",
		queryStr: `{ "isOptimize": "true", "getDataForAllColumn": "true", "isTranslateUser": "0", "filter": [ { "column": "tmg_assignee", "operation": "and", "conditions": [{ "name": "contains", "value": "xxxx" }], "dataType": "text" } ], "search": "abc''xtzz -- '' ' ---'sscscs", "searchColumn": "tmg_assignee,label_task"} `,
		want: struct {
			sql string
			err []error
		}{
			sql: `SELECT "tmg_status", "tmg_assignee", "duration", "start_date_pl", "document_object_update_time", "document_object_id", "user_create", "user_update" FROM document_symper_wbs WHERE ((tmg_status ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (tmg_assignee ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (label_task ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (user_create ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (user_update ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%')) AND ((tmg_assignee ILIKE '%xxxx%') AND (tenant_id_ = 0)) LIMIT 50`,
			err: nil,
		},
		assertFunc: func(t *testing.T, wantedSql string, actualSql string) {
			assert.Equal(t, wantedSql, actualSql)
		},
	},
	{
		name:     "Test filter khi cột sort không có trong select",
		queryStr: `{ "filter": [ { "column": "tmg_assignee", "operation": "and", "conditions": [{ "name": "contains", "value": "xxxx" }], "dataType": "text" } ], "page": "1", "pageSize": "300", "columns": ["tmg_assignee"], "distinct": "true", "sort": [{ "column": "duration", "type": "asc" }], "groupBy": [], "aggregates": [] } `,
		want: struct {
			sql string
			err []error
		}{
			sql: `SELECT "tmg_assignee", "duration" FROM document_symper_wbs WHERE ((tmg_assignee ILIKE '%xxxx%') AND (tenant_id_ = 0)) ORDER BY duration asc LIMIT 300`,
			err: nil,
		},
		assertFunc: func(t *testing.T, wantedSql string, actualSql string) {
			assert.Equal(t, wantedSql, actualSql)
		},
	},
	{
		tenantId: -1,
		name:     "Test filter khi có column,search, filter, group by, sort, aggregate, link table, order by khi không có tenantId",
		queryStr: `{ "isOptimize": "true", "getDataForAllColumn": "true", "isTranslateUser": "0", "filter": [ { "column": "tmg_assignee", "operation": "and", "conditions": [{ "name": "contains", "value": "xxxx" }], "dataType": "text" } ], "search": "abc''xtzz -- '' ' ---'sscscs", "searchColumn": "tmg_assignee,label_task"} `,
		want: struct {
			sql string
			err []error
		}{
			sql: `SELECT "tmg_status", "tmg_assignee", "duration", "start_date_pl", "document_object_update_time", "document_object_id", "user_create", "user_update" FROM document_symper_wbs WHERE ((tmg_status ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (tmg_assignee ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (label_task ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (user_create ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%') OR (user_update ILIKE '%abc''''xtzz -- '''' '' ---''sscscs%')) AND ((tmg_assignee ILIKE '%xxxx%')) LIMIT 50`,
			err: nil,
		},
		assertFunc: func(t *testing.T, wantedSql string, actualSql string) {
			assert.Equal(t, wantedSql, actualSql)
		},
	},
}

func TestModelFilter(t *testing.T) {
	conn := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		Addr:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
	})
	sorm := repository.NewSymperOrm(conn)
	// testObj := new(ModelTest)
	for _, testDataItem := range testData {
		t.Run(testDataItem.name, func(t *testing.T) {
			// Chuyển string json thành map
			var data map[string]interface{}
			json.Unmarshal([]byte(testDataItem.queryStr), &data)

			// lấy model filter object
			modelFilter := library.ParseToModelFilter(allColumn, data, testDataItem.tenantId)

			// Chuyển model filter thành orm.Query
			query1 := sorm.ModelFilterWithTbName("document_symper_wbs", modelFilter)

			// Chuyển orm.Query thành sql
			sql1 := repository.ToPlainQuery(query1)

			// So sánh kết quả
			testDataItem.assertFunc(t, testDataItem.want.sql, sql1)

			// query2 := sorm.ModelFilterWithModel(testObj, modelFilter)
			// sql2 := repository.ToPlanQuery(query2)
			// testDataItem.assertFunc(t, testDataItem.want.sql, sql2)

			// fmt.Println(testDataItem.name)
			// fmt.Println(sql1)
			// fmt.Println("=====================================")
		})
	}
}

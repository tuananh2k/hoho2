package accesscontrol

import (
	"encoding/json"
	"errors"
	"fmt"
	"hoho-framework-v2/adapters/request"
	"hoho-framework-v2/library"
	"hoho-framework-v2/library/auth"
	"os"
	"strings"
)

var (
	ErrGetUserOperation      = errors.New("error when get user operation")
	ErrReplaceFormulaByValue = errors.New("error when replace formula by value")
)

type OperationItem struct {
	RoleIdentifier   string `json:"roleIdentifier"`
	ObjectIdentifier string `json:"objectIdentifier"`
	Action           string `json:"action"`
	Filter           string `json:"filter"`
	ActionPackId     string `json:"actionPackId"`
	ObjectType       string `json:"objectType"`
}

func CheckActionWithCurrentRole(authObject auth.AuthObject, objectIdentifier, action string) bool {
	if authObject.IsBa() {
		return true
	}
	roleIdentifier := authObject.GetUserRole()
	roleIdentifier = "orgchart:118:196c70e7-024e-486e-8f62-2181baa9c3a6"
	return CheckRoleActionRemote(authObject, roleIdentifier, objectIdentifier, action)
}

func CheckRoleActionRemote(authObject auth.AuthObject, roleIdentifier, objectIdentifier, action string) bool {
	listAction := GetRoleActionRemote(authObject, roleIdentifier, objectIdentifier)
	fmt.Println("listAction")
	fmt.Println(listAction)
	return library.StringInSlice(action, listAction)
}

func GetRoleActionRemote(authObject auth.AuthObject, roleIdentifier, objectIdentifier string) []string {
	listAction := make([]string, 0)
	dataRes, e := request.Make(os.Getenv("ACCESS_CONTROL_SERVICE") + "/roles/" + roleIdentifier + "/accesscontrol/" + objectIdentifier).SetHeaders(map[string]string{"Authorization": "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEzMTEiLCJmaXJzdE5hbWUiOiJUXHUwMGY5bmcxIiwibGFzdE5hbWUiOiJUclx1MWVhN24gTFx1MDBlYSIsInVzZXJOYW1lIjoidHVuZ3RsQGNlbnRvZ3JvdXAudm4iLCJkaXNwbGF5TmFtZSI6IlRyXHUxZWE3biBMXHUwMGVhIFRcdTAwZjluZzEiLCJlbWFpbCI6InR1bmd0bEBjZW50b2dyb3VwLnZuIiwidGVuYW50SWQiOiIxIiwicmVzZXRQc3dkSW5mbyI6bnVsbCwidHlwZSI6InVzZXIiLCJpcCI6IjU4LjE4Ny4yNTAuMTM4IiwidXNlckFnZW50IjoiTW96aWxsYVwvNS4wIChNYWNpbnRvc2g7IEludGVsIE1hYyBPUyBYIDEwXzE1XzcpIEFwcGxlV2ViS2l0XC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBDaHJvbWVcLzEwOS4wLjAuMCBTYWZhcmlcLzUzNy4zNiIsImRlbGVnYXRlZEJ5Ijp7ImlkIjoiMTkiLCJlbWFpbCI6ImhvYW5nbmRAY2VudG9ncm91cC52biIsIm5hbWUiOiJOZ3V5XHUxZWM1biBcdTAxMTBcdTAwZWNuaCBIb1x1MDBlMG5nIiwiZXhwIjoiMCJ9LCJ0ZW5hbnQiOnsiaWQiOiIxIn0sImlzX2Nsb3VkIjp0cnVlLCJ0ZW5hbnRfZG9tYWluIjoic3ltcGVyLnZuIiwiZXhwIjoiMCIsImlhdCI6MTY3NTE0MDE2Niwicm9sZSI6Im9yZ2NoYXJ0OjExODoxOTZjNzBlNy0wMjRlLTQ4NmUtOGY2Mi0yMTgxYmFhOWMzYTYifQ==.NTk5NTY5YWRjZWRmN2ZiOWI4ZjljYzYxNDhmZWQ2YmRmMGI1OTBhYTU0M2QzYTI1ZWQzOTYwODAyYzU1YmI0ZjViNjMwYTM2OTE3NzI4YmFkNDllYjU1NWI2NWQ1NWNmZWQ4NzBmMDQxYjYzYjE1ZjY1ODU1ZTAwZWIyOWUxZGZiNmZkYTRkNTU1YTljZGYyYTI0NzAxOTMyMmJjOGJmMjRiZWEzNzY3MDY3ZjRkM2YxMmZjOGU0Yzk3YmNjYjZhMjQ1YzIwNzMzMDgxZTQ1NDRlNTAxMTI2NTIwNWU5MjJhYWRhYWM0MmExNjQ5NDM2NGQxZGRkZjkxZjg4NjhmNTE4ZDkyNjYyZDcyMWUxMDcyYzc3OWU0MDc4NzM1M2IwMWYzNTk4NGE4YWMzODVlYmIwMzZkNmQ1OGExNWYwZTBmNDBmYTBmZWQ1OWQ3YWRiOGM5ZTgxNzRkOWJhN2IzNDJmNGQ1Y2MwNDQzN2ZkZjFkMzNiYjc4MGE1OGIzZTRhYjY4MjAxZTE2YjMxNDdmODY5ZmM2MjQ2Y2Y4NDgwNTc4OWNkMjFhYTQyNmVkMDMxZTNiMWEwODk2OGQ1YTJhNmU3ODBmY2MyNWE0Y2Q5MDNlYTdkZjE5Y2QxNjg1NDNmNWE1NjdkMTgxYjE1N2NjYzQ5MGM5ZGU1YjQxNTVkYTQ="}).Get()
	if e != nil || dataRes.Status != 200 {
		return listAction
	}
	dataJsonStr := dataRes.Data.(map[string]interface{})
	dataApi := dataJsonStr["data"].([]interface{})
	for _, v := range dataApi {
		vJson := v.(map[string]interface{})
		listAction = append(listAction, vJson["action"].(string))
	}
	return listAction
}

func GetFilterString(authObject auth.AuthObject, objectIdentifier, action string) (string, error) {
	if authObject.IsBa() {
		return "", nil
	}
	oprs, err := GetOperations(authObject, objectIdentifier, action)
	if err != nil {
		return "", err
	}

	if len(oprs) == 0 {
		return "", nil
	}

	filterItems := make([]string, 0)
	groupByActionPackId := make(map[string][]string)
	hasOperationWithoutFilter := false

	for _, opr := range oprs {
		if opr.Filter != "" {
			if _, ok := groupByActionPackId[opr.ActionPackId]; !ok {
				groupByActionPackId[opr.ActionPackId] = make([]string, 0)
			}
			groupByActionPackId[opr.ActionPackId] = append(groupByActionPackId[opr.ActionPackId], opr.Filter)
		} else {
			hasOperationWithoutFilter = true
			break
		}
	}

	// nếu có ít nhất 1 operation không có filter thì trả về chuỗi rỗng
	if hasOperationWithoutFilter {
		return "", nil
	}

	for _, item := range groupByActionPackId {
		// Các filter cùng action pack nối với nhau bằng AND
		filterItems = append(filterItems, "("+strings.Join(item, " AND ")+")")
	}
	// Các filter khác action pack nối với nhau bằng OR
	cond := strings.Join(filterItems, " OR ")
	// replace all new line character by space
	cond = strings.ReplaceAll(cond, "\n", " ")
	cond = strings.ReplaceAll(cond, "\r\n", " ")
	cond = strings.ReplaceAll(cond, "\t", " ")

	cond, err = replaceFormulaByValue(cond, authObject)
	if err != nil {
		return cond, err
	}
	return cond, nil
}

func replaceFormulaByValue(cond string, authObject auth.AuthObject) (string, error) {
	allFormula := library.GetSubStringByFunction(cond, "ref", '(', ')', true)
	mapFormula := make(map[string]string)

	for _, formula := range allFormula {
		mapFormula[formula] = formula
	}

	syqlUrl := os.Getenv("SYQL_SERVICE") + "/formulas/compileClientBulk"
	data := make(map[string]interface{}, 0)
	jsonData, _ := json.Marshal(mapFormula)

	data["formulas"] = string(jsonData)
	data["variables"] = "{}"

	dataRes, err := request.Make(syqlUrl).SetHeaders(map[string]string{"Authorization": authObject.GetToken()}).SetBody(data).Post()
	if err != nil {
		return cond, fmt.Errorf("%w : %v", ErrReplaceFormulaByValue, err)
	}

	if dataRes.Status != 200 {
		return cond, fmt.Errorf("%w : status code %d", ErrReplaceFormulaByValue, dataRes.Status)
	}

	dataJsonStr := dataRes.Data.(map[string]interface{})
	if dataJsonStr == nil {
		return cond, fmt.Errorf("%w : data is nil", ErrReplaceFormulaByValue)
	}

	for formula := range mapFormula {
		value := dataJsonStr[formula]
		if value != nil {
			cond = strings.ReplaceAll(cond, formula, library.ToString(value))
		}
	}

	return cond, nil
}

// GetOperations lấy danh sách các operation của user
//
//	authObject: auth object của user
//	objectIdentifier: object cần lấy operation
//	action: action cần lấy operation, nếu truyền vào giá trị rỗng thì sẽ lấy tất cả các operation ứng với object
func GetOperations(authObject auth.AuthObject, objectIdentifier, action string) ([]OperationItem, error) {
	roleIdentifier := authObject.GetUserRole()
	rsl := make([]OperationItem, 0)
	url := os.Getenv("ACCESS_CONTROL_SERVICE") + "/roles/" + roleIdentifier + "/accesscontrol/" + objectIdentifier

	dataRes, e := request.Make(url).SetHeaders(map[string]string{"Authorization": authObject.GetToken()}).Get()

	if e != nil {
		return rsl, fmt.Errorf("%w : %v", ErrGetUserOperation, e)
	}

	if dataRes.Status != 200 {
		return rsl, fmt.Errorf("%w : status code %d", ErrGetUserOperation, dataRes.Status)
	}

	dataJsonStr := dataRes.Data.(map[string]interface{})
	dataApi := dataJsonStr["data"].([]interface{})
	for _, v := range dataApi {
		vJson := v.(map[string]interface{})
		if action == "" || vJson["action"] == action {
			operationItem := OperationItem{}
			operationItem.Action = library.ToString(vJson["action"])
			operationItem.ObjectIdentifier = library.ToString(vJson["objectIdentifier"])
			operationItem.ActionPackId = library.ToString(vJson["actionPackId"])
			operationItem.Filter = library.ToString(vJson["filter"])
			operationItem.RoleIdentifier = library.ToString(vJson["roleIdentifier"])
			operationItem.ObjectType = library.ToString(vJson["objectType"])
			rsl = append(rsl, operationItem)
		}
	}
	return rsl, nil
}

func GetAllOperations(authObject auth.AuthObject) ([]OperationItem, error) {
	roleIdentifier := authObject.GetUserRole()
	url := os.Getenv("ACCESS_CONTROL_SERVICE") + "/roles/" + roleIdentifier + "/accesscontrol"

	rsl := make([]OperationItem, 0)
	dataRes, e := request.Make(url).SetHeaders(map[string]string{"Authorization": authObject.GetToken()}).Get()

	if e != nil {
		return rsl, fmt.Errorf("%w : %v", ErrGetUserOperation, e)
	}

	if dataRes.Status != 200 {
		return rsl, fmt.Errorf("%w : status code %d", ErrGetUserOperation, dataRes.Status)
	}

	dataJsonStr := dataRes.Data.(map[string]interface{})
	dataApi := dataJsonStr["data"].([]interface{})
	for _, v := range dataApi {
		vJson := v.(map[string]interface{})
		operationItem := OperationItem{}
		operationItem.Action = library.ToString(vJson["action"])
		operationItem.ObjectIdentifier = library.ToString(vJson["objectIdentifier"])
		operationItem.ActionPackId = library.ToString(vJson["actionPackId"])
		operationItem.Filter = library.ToString(vJson["filter"])
		operationItem.RoleIdentifier = library.ToString(vJson["roleIdentifier"])
		operationItem.ObjectType = library.ToString(vJson["objectType"])
		rsl = append(rsl, operationItem)
	}
	return rsl, nil
}

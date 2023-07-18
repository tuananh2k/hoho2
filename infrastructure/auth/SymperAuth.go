package auth

/*
create by: Hoangnd
create at: 2023-01-01
des: Xử lý dữ liệu từ token của user
*/

import (
	"io"
	"os"
	"strconv"
	"strings"

	sAuth "hoho-framework-v2/library/auth"

	"github.com/valyala/fastjson"
)

type baObject struct {
	id            string
	name          string
	email         string
	resetPswdInfo string
}
type tenant struct {
	id      int
	domain  string
	isCloud string
}

type authObject struct {
	id            string
	userName      string
	displayName   string
	email         string
	resetPswdInfo string
	accType       string
	userAgent     string
	ip            string
	role          string
	exp           int64
	iat           int64
	tenant        tenant
	baInfo        baObject
	token         string
}

func NewAuthObject(data []byte, jwt string) (sAuth.AuthObject, error) {
	var p fastjson.Parser
	v, err := p.Parse(string(data))
	if err != nil {
		return nil, err
	}
	typeObj := string(v.GetStringBytes("type"))
	baObj := &baObject{}
	tenantStr := v.Get("tenant")
	t := string(tenantStr.GetStringBytes("id"))
	tI, _ := strconv.Atoi(t)
	tenantObj := &tenant{
		id:      tI,
		domain:  string(tenantStr.GetStringBytes("domain")),
		isCloud: string(tenantStr.GetStringBytes("isCloud")),
	}
	if typeObj == "ba" {
		baObj = &baObject{
			id:            string(v.GetStringBytes("id")),
			name:          string(v.GetStringBytes("name")),
			email:         string(v.GetStringBytes("email")),
			resetPswdInfo: string(v.GetStringBytes("resetPswdInfo")),
		}
		userDelegate := v.Get("userDelegate")
		v = userDelegate
	}
	userId := string(v.GetStringBytes("id"))
	userUserName := string(v.GetStringBytes("userName"))
	userDisplayName := string(v.GetStringBytes("displayName"))
	userEmail := string(v.GetStringBytes("email"))
	userResetPswdInfo := string(v.GetStringBytes("resetPswdInfo"))
	userType := string(v.GetStringBytes("type"))
	userIp := string(v.GetStringBytes("ip"))
	userUserAgent := string(v.GetStringBytes("userAgent"))
	userRole := string(v.GetStringBytes("role"))
	return &authObject{
			baInfo:        *baObj,
			id:            userId,
			userName:      userUserName,
			displayName:   userDisplayName,
			email:         userEmail,
			resetPswdInfo: userResetPswdInfo,
			accType:       userType,
			ip:            userIp,
			role:          userRole,
			userAgent:     userUserAgent,
			tenant:        *tenantObj,
			token:         jwt,
		},
		err
}

func (au *authObject) GetUserId() string {
	return au.id
}
func (au *authObject) GetUserDisplayName() string {
	return au.displayName
}
func (au *authObject) GetUserTenantId() int {
	return au.tenant.id
}
func (au *authObject) GetUserUserName() string {
	return au.userName
}
func (au *authObject) GetUserEmail() string {
	return au.email
}
func (au *authObject) GetUserResetPswdInfo() string {
	return au.resetPswdInfo
}
func (au *authObject) GetUserAccType() string {
	return au.accType
}
func (au *authObject) GetUserUserAgent() string {
	return au.userAgent
}
func (au *authObject) GetUserIp() string {
	return au.ip
}
func (au *authObject) GetUserRole() string {
	return au.role
}
func (au *authObject) GetUserExp() int64 {
	return au.exp
}
func (au *authObject) GetUserIat() int64 {
	return au.iat
}
func (au *authObject) GetTenant() map[string]interface{} {
	return au.getTenant()
}
func (au *authObject) GetBaInfo() map[string]interface{} {
	return au.getBaInfo()
}
func (au *authObject) GetToken() string {
	if strings.Contains(au.token, "new_symper_authen_!") {
		return au.token
	}
	return au.token + "new_symper_authen_!"
}

func (au *authObject) GetBaId() string {
	return au.baInfo.id
}
func (au *authObject) IsBa() bool {
	return au.baInfo.id != ""
}
func (au *authObject) GetBaEmail() string {
	return au.baInfo.email
}
func (au *authObject) GetBaName() string {
	return au.baInfo.name
}

func (au *authObject) GetAll() interface{} {
	return map[string]interface{}{
		"id":            au.id,
		"userName":      au.userName,
		"displayName":   au.displayName,
		"email":         au.email,
		"resetPswdInfo": au.resetPswdInfo,
		"accType":       au.accType,
		"userAgent":     au.userAgent,
		"ip":            au.ip,
		"role":          au.role,
		"exp":           au.exp,
		"iat":           au.iat,
		"tenant":        map[string]interface{}{"id": au.tenant.id, "domain": au.tenant.domain, "isCloud": au.tenant.isCloud},
		"baInfo":        map[string]interface{}{"id": au.baInfo.id, "name": au.baInfo.name, "email": au.baInfo.email, "resetPswdInfo": au.baInfo.resetPswdInfo},
		"token":         au.token,
	}
}

func (au *authObject) getTenant() map[string]interface{} {
	return map[string]interface{}{"id": au.tenant.id, "domain": au.tenant.domain, "isCloud": au.tenant.isCloud}
}
func (au *authObject) getBaInfo() map[string]interface{} {
	return map[string]interface{}{"id": au.baInfo.id, "name": au.baInfo.name, "email": au.baInfo.email, "resetPswdInfo": au.baInfo.resetPswdInfo}
}

func GetPublicKey() []byte {
	file, err := os.Open("crypt/public.pem")
	if err != nil {
		return nil
	}
	defer file.Close()
	fileByte, _ := io.ReadAll(file)
	return fileByte
}

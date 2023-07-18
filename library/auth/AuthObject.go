package auth

type AuthObject interface {
	GetUserId() string
	GetUserDisplayName() string
	GetUserTenantId() int
	GetAll() interface{}
	GetUserUserName() string
	GetUserEmail() string
	GetUserResetPswdInfo() string
	GetUserAccType() string
	GetUserUserAgent() string
	GetUserIp() string
	GetUserRole() string
	GetUserExp() int64
	GetUserIat() int64
	GetTenant() map[string]interface{}
	GetBaInfo() map[string]interface{}
	GetToken() string
	GetBaId() string
	GetBaEmail() string
	GetBaName() string
	IsBa() bool
}

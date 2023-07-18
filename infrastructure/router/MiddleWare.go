package router

import (
	"encoding/json"
	"errors"
	"hoho-framework-v2/adapters/request"
	iAuth "hoho-framework-v2/infrastructure/auth"
	aAuth "hoho-framework-v2/library/auth"
	"os"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/*
Hàm lấy thông tin token và kiểm tra
- hợp lệ
- Xác thực
*/
func getMiddleWareConfig(authObject *aAuth.AuthObject) middleware.JWTConfig {
	return middleware.JWTConfig{
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				signingKey := iAuth.GetPublicKey()
				signingK, _ := jwt.ParseRSAPublicKeyFromPEM(signingKey)
				return signingK, nil
			}
			representativeToken := c.Request().Header.Get("S-Representative")
			if representativeToken != "" {
				res, eRes := request.Make(os.Getenv("ACCOUNT_SERVICE") + "/auth/representative").
					SetBody(map[string]interface{}{"hashKey": representativeToken}).
					SetHeaders(map[string]string{"Authorization": "Bearer " + auth}).Post()
				if eRes == nil {
					dataRes := res.Data.(map[string]interface{})
					auth = dataRes["data"].(string)
				}
			}
			var re = regexp.MustCompile(`[a-z0-9A-X]*::`)
			newAuth := re.ReplaceAllString(auth, "")
			newAuth = strings.Replace(newAuth, "new_symper_authen_!", "", -1)
			token, err := jwt.Parse(newAuth, keyFunc)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}
			claims, _ := json.Marshal(token.Claims)
			*authObject, _ = iAuth.NewAuthObject(claims, "Bearer "+token.Raw)
			return token, nil
		},
	}
}

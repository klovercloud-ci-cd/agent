package v1

import (
	"github.com/klovercloud-ci/api/common"
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/dependency"
	"github.com/labstack/echo/v4"
)

// handle user authentication and authorization here.
func AuthenticationAndAuthorizationHandler(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) (err error) {
		if config.EnableAuthentication {
			token := context.Request().Header.Get("token")
			if token ==""{
				return common.GenerateErrorResponse(context,"[ERROR]: Invalid token!","Failed to parse token!")
			}else{
				res,_:=dependency.GetJwtService().ValidateToken(token)
				if !res{
					return common.GenerateErrorResponse(context,"[ERROR]: Invalid token!","Please provide a valid token!")
				}
			}
			return handler(context)
		}else{
			return handler(context)
		}
	}
}

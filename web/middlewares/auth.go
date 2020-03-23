package middlewares

import (
	"fmt"
	"strings"
	"time"

	"flower/config"
	"flower/entities"
	"flower/model"
	"flower/repositories"
	"flower/utils"
	"flower/web/controllers"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	jwtInternal "gopkg.in/dgrijalva/jwt-go.v3"
)

const (
	PERMISSION_READ      = "r"
	PERMISSION_WRITE     = "w"
	PERMISSION_ALLOW_ALL = ""
	PERMISSION_DENY_ALL  = "deny"
)

var authMiddleware *jwt.GinJWTMiddleware
var authMiddlewareFunc gin.HandlerFunc
var LoginHandler func(c *gin.Context)
var RefreshHandler func(c *gin.Context)
var modulePermissionMap map[string]string

func InitAuth() {
	conf := config.GetConfig()
	authMiddleware = &jwt.GinJWTMiddleware{
		Realm:           "flower",
		Key:             []byte("4rd5rAPPz52dhvjlSK1lTT08dOVNxFKN0GeosDeps3cAUwlL5YzkGBDbD3sGvogd"),
		Authenticator:   authenticator,
		PayloadFunc:     payloadFunc,
		IdentityHandler: identityHandler,
		Authorizator:    authorizator,
	}
	if timeout, err := time.ParseDuration(conf.AuthTimeout); err == nil {
		authMiddleware.Timeout = timeout
	}
	if maxRefresh, err := time.ParseDuration(conf.AuthMaxRefresh); err == nil {
		authMiddleware.MaxRefresh = maxRefresh
	}
	authMiddlewareFunc = authMiddleware.MiddlewareFunc()
	LoginHandler = authMiddleware.LoginHandler
	RefreshHandler = authMiddleware.RefreshHandler
	modulePermissionMap = make(map[string]string)
	modulePermissionMap["get_user"] = PERMISSION_ALLOW_ALL
	modulePermissionMap["post_funcroles"] = PERMISSION_READ // search
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authMiddlewareFunc(ctx)
	}
}

func LogoutHandler(ctx *gin.Context) {
	user := controllers.GetCurrentUserFromCtx(ctx)
	if user == nil {
		controllers.Error404f(ctx, "User not found")
		return
	}
	dal := repositories.GetAccountDal()
	if err := dal.UpdateSecurityToken(user.UserName, utils.GetIDGenerate().GetID().String()); err != nil {
		controllers.Error500(ctx, err)
		return
	}
	controllers.NoContent(ctx)
}

// Callback function that should perform the authentication of the user based on login info.
// Must return user data as user identifier, it will be stored in Claim Array. Required.
// Check error (e) to determine the appropriate error message.
func authenticator(ctx *gin.Context) (interface{}, error) {
	req := &model.LoginReq{}
	if err := controllers.ShouldBindValidJSON(ctx, req); err != nil {
		return nil, errors.WithStack(err)
	}
	user, err := authenticate(req.UserName, req.Password)
	if err != nil {
		glog.Errorf("Failed validate userName password\n%+v", err)
		return nil, err
	}
	return user, nil
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if e, ok := data.(map[string]interface{}); ok {
		return jwt.MapClaims(e)
	}
	return jwt.MapClaims{}
}

// Set the identity handler function
func identityHandler(claims jwtInternal.MapClaims) interface{} {
	dal := repositories.GetAccountDal()
	userName, ok := claims["UserName"].(string)
	if !ok || userName == "" {
		glog.Warningf("User name not found in claims")
		return nil
	}
	st, ok := claims["SecurityToken"].(string)
	if !ok || st == "" {
		glog.Warningf("Security token not found in claims")
		return nil
	}
	user, err := dal.GetValidExtendByUserName(userName)
	if err != nil {
		glog.Warningf("Invalid user %v", userName)
		return nil
	}
	if user.SecurityToken != st {
		glog.Infof("Security token was updated")
		return nil
	}
	user.UserPasswd = "" // clear pwd
	return user
}

func authorizator(data interface{}, ctx *gin.Context) bool {
	if data == nil {
		return false
	}
	user, ok := data.(*entities.AccountExtend)
	if !ok {
		return false
	}
	if user.FuncRoleId != model.DEFAULT_FUNC_ROLE_ID_INT {
		// not admin role
		// fmt.Println(user.FuncRole.Content)
		parts := strings.Split(ctx.Request.URL.Path, "/")
		if len(parts) < 2 {
			glog.Warningf("Buggy url %v", ctx.Request.URL.Path)
			return false
		}
		module := parts[1]
		method := strings.ToLower(ctx.Request.Method)
		permission, exist := modulePermissionMap[fmt.Sprintf("%v_%v", method, module)]
		if !exist {
			if method == "get" {
				permission = PERMISSION_READ
			} else {
				permission = PERMISSION_WRITE
			}
		}
		if permission != "" {
			result := gjson.Get(user.FuncRole.Content, pluralToSingular(module))
			if !result.Exists() {
				return false
			}
			if !strings.Contains(result.String(), permission) {
				return false
			}
		}
	}
	ctx.Set(controllers.CTX_KEY_USER, user)
	return true
}

func authenticate(userName, pwd string) (map[string]interface{}, error) {
	dal := repositories.GetAccountDal()
	e, err := dal.GetValidByUserName(userName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if e == nil {
		return nil, errors.Errorf("Invalid userName / password")
	}
	if e.UserPasswd != model.EncryptPassword(pwd) {
		return nil, errors.Errorf("Invalid userName / password")
	}
	return map[string]interface{}{
		"UserName":      e.UserName,
		"SecurityToken": e.SecurityToken,
	}, nil
}

func pluralToSingular(text string) string {
	i := len(text) - 1
	if text[i] == 's' {
		return text[:i]
	}
	return text
}

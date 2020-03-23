package routers

import (
	"github.com/gin-gonic/gin"
	"flower/web/controllers"
	"flower/web/middlewares"
)

type Router struct {
	RouterEngine *gin.Engine
	//AuthMiddleware *jwt.GinJWTMiddleware
	Controller *controllers.Controller
}

func NewRouter(controller *controllers.Controller) *Router {
	engine := gin.Default()
	router := &Router{
		RouterEngine: engine,
		Controller:   controller,
	}
	return router
}
func (router *Router) AddBaseRouter() {
	router.RouterEngine.GET("/ping", router.Controller.Ping)
	//GetID
	router.RouterEngine.GET("/id", router.Controller.GetID)

	router.AddAuth()
	router.AddFuncRole()
	router.AddUser()
}

func (router *Router) AddAuth() {
	middlewares.InitAuth()
	router.RouterEngine.POST("/login", middlewares.LoginHandler)
	router.RouterEngine.POST("/logout", middlewares.Auth(), middlewares.LogoutHandler)
	router.RouterEngine.PUT("/auth/password", middlewares.Auth(), router.Controller.UpdatePassword)
}

func (router *Router) AddFuncRole() {
	router.RouterEngine.Group("/funcroles").
		Use(middlewares.Auth()).
		POST("", router.Controller.QueryFuncRoles).
		DELETE("", router.Controller.DeleteFuncRoles)
	router.RouterEngine.Group("/funcrole").
		Use(middlewares.Auth()).
		POST("", router.Controller.AddFuncRole).
		PUT("", router.Controller.ModifyFuncRole).
		DELETE("/:id", router.Controller.DeleteFuncRole)
}

func (router *Router) AddUser() {
	//
	router.RouterEngine.Group("/users").
		Use(middlewares.Auth()).
		DELETE("", router.Controller.DeleteUsers).
		POST("", router.Controller.QueryUsers)
	router.RouterEngine.Group("/user").
		Use(middlewares.Auth()).
		GET("", router.Controller.GetCurrentUser).
		POST("", router.Controller.AddUser).
		PUT("", router.Controller.ModifyUser).
		PUT("/:id/valid", router.Controller.ModifyUserValid).
		PUT("/:id/invalid", router.Controller.ModifyUserInvalid).
		DELETE("/:id", router.Controller.DeleteUser)
}

package main

import (
	"fmt"

	"flower/config"
	"flower/web/controllers"
	"flower/web/routers"
)

func main() {
	//flag.Parse()
	conf := config.GetConfig()
	fmt.Println(conf)
	router := routers.NewRouter(controllers.New(conf))
	router.AddBaseRouter()
	router.RouterEngine.Run(conf.HttpServiceAddr)
}

package main

import (
	"fmt"

	"github.com/panzhenyu12/flower/config"
	"github.com/panzhenyu12/flower/web/controllers"
	"github.com/panzhenyu12/flower/web/routers"
)

func main() {
	//flag.Parse()
	conf := config.GetConfig()
	fmt.Println(conf)
	router := routers.NewRouter(controllers.New(conf))
	router.AddBaseRouter()
	router.RouterEngine.Run(conf.HttpServiceAddr)
}

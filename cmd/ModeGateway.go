package cmd

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"message/internal/config"
	"message/internal/repo"
	"message/internal/router"
	"message/package/util_debug"
)

func MainModeGateway() {
	debugOn, e := config.GetDebugMode()
	if e != nil {
		log.Fatal("GetDebugMode: ", e)
	}

	if debugOn {
		addr, e := config.GetDebugPprofGateway()
		if e != nil {
			log.Fatal("config e: ", e)
		}

		log.Println("GetDebugPprofGateway addr: ", addr)
		go util_debug.InitPProf(addr)
	}

	v, _ := config.GetDependMongo()
	repo.InitMongo(v)

	redisDsn, e := config.GetDependRedisDsn()
	if e != nil {
		log.Fatal("GetDependRedisDsn: ", e)
	}
	log.Println("redisDsn: ", redisDsn)
	e = repo.InitRedis(context.Background(), redisDsn, 0, 0)
	if e != nil {
		log.Fatal("InitRedis: ", e)
	}

	r := gin.Default()
	v1 := r.Group("/api/v1/sessions")
	{
		v1.GET("/:session_id", router.FetchSession)
	}
	v2 := r.Group("/api/v2/sessions")
	{
		v2.GET("/:session_id", router.FetchSessionV2)
	}
	v1r := r.Group("/api/v1r/session")
	{
		//v1r.GET("/:session_id", router.GetSessionV2)
		v1r.GET("/:session_id", router.GetSessionV2)
	}

	val, _ := config.GetGwPortHttp()
	r.Run(val)
}

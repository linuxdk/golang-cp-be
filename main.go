package main

import (
  "github.com/gin-gonic/gin"
  "golang-cp-be/server"
  "golang-cp-be/config"
)

/*

/authorizations

*/

func init() {
  config.InitConfigurations()
}

func main() {

  r := gin.Default()
  //r.GET("/ping", controller.Ping)

  v1 := r.Group("v1")

  server.V1Routes(v1) //Added all routes

  r.Run() // listen and serve on 0.0.0.0:8080
}
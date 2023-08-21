package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/cmuh-cmax/cmax-backend/api"
)

func router() *gin.Engine{
  r := gin.Default()

  r.GET("/api/version", api.GetVersion)

  return r
}

func main(){
  fmt.Println("Going to start the server.")
  r := router()
  r.Run(":8000")
}

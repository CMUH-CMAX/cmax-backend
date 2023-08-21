package api

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func GetVersion(r *gin.Context){
  version := "v0.0.0"
  r.String(http.StatusOK, version)
}
package server

import "github.com/gin-gonic/gin"

type GinHandler = func(ctx *gin.Context)

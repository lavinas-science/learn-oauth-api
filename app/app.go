package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lavinas-science/learn-oauth-api/domain/access_token"
	"github.com/lavinas-science/learn-oauth-api/http"
	"github.com/lavinas-science/learn-oauth-api/repository/db"
	"github.com/lavinas-science/learn-oauth-api/repository/rest"
)

var (
	router = gin.Default()
)

func StartApp() {
	atHandler := http.NewHandler(access_token.NewService(db.NewRepository(), rest.NewRepository()))
	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)
	router.Run(":9090")
}

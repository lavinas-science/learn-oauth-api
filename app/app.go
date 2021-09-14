package app

import (
	"github.com/lavinas-science/learn-oauth-api/domain/access_token"
	"github.com/lavinas-science/learn-oauth-api/http"
	"github.com/lavinas-science/learn-oauth-api/repository/db"
	"github.com/gin-gonic/gin"
)


var (
	router = gin.Default()
)

func StartApp() {
	// service := access_token.NewService(db.NewRepository())
	atHandler := http.NewHandler(access_token.NewService(db.NewRepository()))
	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.Run(":8080")

}
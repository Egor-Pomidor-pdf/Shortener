package handler

import (
	"github.com/wb-go/wbf/ginext"
)

func NewRouter(notifyHandler * Handler) *ginext.Engine {
	router := ginext.New("release")
	router.Use(ginext.Logger())
	router.Use(ginext.Recovery())
	router.POST("/shorten", notifyHandler.postShorten)
	router.GET("/s/:short", notifyHandler.getRedirect)
	router.GET("/analytics/:short", notifyHandler.getAnalytics)
	return router
}

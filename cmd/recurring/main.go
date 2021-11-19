package main

import (
	"recurring-standard/pkg/recurring/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger()) // Logger middleware

	r.GET("/", controllers.Home)
	r.NoRoute(controllers.PageNotFound)

	// Create Sub Router for  customised API version
	authorized := r.Group("/api/v1")
	authorized.Use(controllers.IsAuthorized)

	// subscription
	authorized.POST("/subscription", controllers.SetSubscription)
	authorized.POST("/subscription/list", controllers.GetSubscriptionList)
	authorized.POST("/subscription/search", controllers.GetSubscriptionSearch)
	authorized.POST("/subscription/status", controllers.GetSubscriptionStatus)
	authorized.POST("/subscription/cancel", controllers.Unsubscribe)

	r.HandleMethodNotAllowed = true
	r.Run(":9090")
}

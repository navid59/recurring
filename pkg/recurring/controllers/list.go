/**
Get subscription list of a Merchant

Problems : -
To do: -
Note : -
*/
package controllers

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

type Merchant struct {
	MerchantSignature string `json:"MerchantSignature" binding:"required"`
}

func GetSubscriptionList(c *gin.Context) {

	var merchant Merchant
	if err := c.BindJSON(&merchant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request Errors",
		})
		return
	}

	// Verify Audience
	if _, ok := IsLicensed(c.Request.Header.Get("token"), merchant.MerchantSignature); !ok {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Request has conflict with your authion",
		})
		return
	}

	// Creates a client.
	ctx := context.Background()

	// Set your Google Cloud Platform project ID.
	projectID := "netopia-payments"

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	var members []Subscription
	query := datastore.NewQuery("subscribers").Namespace("recurring").
		Filter("MerchantSignature =", merchant.MerchantSignature)

	keys, err := client.GetAll(ctx, query, &members)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": http.StatusUnprocessableEntity, "message": "", "error": err})
		return
	}
	for i, key := range keys {
		members[i].Id = key.ID
	}

	c.JSON(http.StatusOK, gin.H{
		"members": members,
	})

}

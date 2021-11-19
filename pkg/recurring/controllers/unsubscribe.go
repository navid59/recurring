/**
Unsubscribe a member

Problems : -
To do: -
Note : -
*/
package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func Unsubscribe(c *gin.Context) {
	/* Validate Request Body & asssign to VAR member*/
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request Errors",
			"error":   err,
		})
		return
	}

	// Verify Audience
	if _, ok := IsLicensed(c.Request.Header.Get("token"), user.MerchantSignature); !ok {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Request has conflict with your authion",
		})
		return
	}

	// Creates a client.
	ctx := context.Background()
	projectID := "netopia-payments"
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	/* Get User Info */
	subscriberkey := datastore.IDKey("subscribers", user.UserId, nil)
	subscriberkey.Namespace = "recurring"
	entity := &Subscription{}
	if err := client.Get(ctx, subscriberkey, entity); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Data not found",
		})
		return
	}

	// Verify Audience & Ownership of data
	if ownership := entity.MerchantSignature; len(ownership) > 0 && ownership != user.MerchantSignature {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Request has conflict with your authion",
		})
		return
	}

	if entity.Status == false {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotModified,
			"message": "Member already is unsubscribed",
		})
		return
	}

	/* Unscubscribe the Member */
	entity.Status = false
	entity.Flags = FlagCancel
	tmpTime := time.Now()
	entity.UpdatedAt = tmpTime.String()

	if _, err := client.Put(ctx, subscriberkey, entity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Request failed!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Unsubscribed successfully!",
	})
}

/**
Looking for a subscription

Problems : The ID return ZERO
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
	"google.golang.org/api/iterator"
)

func GetSubscriptionSearch(c *gin.Context) {
	/* Validate Request Body & asssign to VAR member*/
	var member Subscriber
	if err := c.BindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request Errors",
			"error":   err,
		})
		return
	}

	// Verify Audience
	if _, ok := IsLicensed(c.Request.Header.Get("token"), member.MerchantSignature); !ok {
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

	/* Get Example of search */
	query := datastore.NewQuery("subscribers").Namespace("recurring").
		Filter("Plan =", member.Plan).
		Filter("Email =", member.UserEmail).
		Filter("MerchantSignature =", member.MerchantSignature)

	var entitis []Subscription
	it := client.Run(ctx, query)
	for {
		var entity Subscription
		_, err := it.Next(&entity)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next task: %v", err)
		} else {
			entitis = append(entitis, entity)
			// Hear need to get the KEY & Assign the KEY
			// fmt.Printf("Nume : %q, Prenume:  %q, Key:%q\n", entity.Name, entity.LastName, entity.Id)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"members": entitis,
	})
}

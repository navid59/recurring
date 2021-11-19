package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

type Subscription struct {
	// Key               *datastore.Key `datastore:"Key"`
	Id                int64         `datastore:"ID" json:"Id"`
	Name              string        `datastore:"Name" json:"Name" binding:"required"`
	LastName          string        `datastore:"LastName" json:"LastName" binding:"required"`
	Email             string        `datastore:"Email" json:"Email" binding:"required"`
	Adress            string        `datastore:"Adress" json:"Adress"`
	Tel               string        `datastore:"Tel" json:"Tel"`
	MerchantSignature string        `datastore:"MerchantSignature" json:"MerchantSignature" binding:"required"`
	Plan              string        `datastore:"Plan" json:"Plan"  binding:"required"`
	PlanDetails       PlanDetails   `datastore:"PlanDetails" json:"PlanDetails"  binding:"required"`
	StartDate         time.Time     `datastore:"StartDate" json:"StartDate" binding:"required"`
	EndDate           time.Time     `datastore:"EndDate" json:"EndDate"`
	PaymentConfig     PaymentConfig `datastore:"PaymentConfig" json:"PaymentConfig" binding:"required"`
	Status            bool
	Flags             string
	CreatedAt         time.Time
	UpdatedAt         string
}

type Subscriber struct {
	Plan              string `json:"plan" binding:"required"`
	UserEmail         string `json:"userEmail" binding:"required"`
	MerchantSignature string `json:"merchantSignature" binding:"required"`
}

type User struct {
	UserId            int64  `json:"userId" binding:"required"`
	MerchantSignature string `json:"merchantSignature" binding:"required"`
}

type PlanDetails struct {
	RecurrenceType string  `datastore:"RecurrenceType" json:"RecurrenceType" binding:"required"`
	RecurrenceDay  int8    `datastore:"RecurrenceDay" json:"RecurrenceDay" binding:"required"`
	PosID          string  `datastore:"PosID" json:"PosID"`
	Description    string  `datastore:"Description" json:"Description"`
	GracePeriod    int     `datastore:"GracePeriod" json:"GracePeriod"`
	Amount         float64 `datastore:"Amount" json:"Amount" binding:"required"`
	Currency       string  `datastore:"Currency" json:"Currency" binding:"required"`
}

type PaymentConfig struct {
	InstantPayment string     `datastore:"InstantPayment" json:"InstantPayment" binding:"required"`
	Instrument     Instrument `datastore:"Instrument" json:"Instrument" binding:"required"`
	ThreeDS2       ThreeDS2   `datastore:"ThreeDS2" json:"ThreeDS2"`
}

type SubscriptionDetails struct {
	PaymentStatus   bool
	PaymentMessage  string
	ValidateStatus  bool
	ValidateMessage string
}

type Instrument struct {
	Type       string `datastore:"Type" json:"Type"`
	Account    string `datastore:"Account" json:"Account"`
	ExpMonth   int    `datastore:"ExpMonth" json:"ExpMonth"`
	ExpYear    int    `datastore:"ExpYear" json:"ExpYear"`
	SecretCode string `datastore:"SecretCode" json:"SecretCode"`
	Token      string `datastore:"Token" json:"Token"`
}

type ThreeDS2 struct {
	BrowserUserAgent    string `datastore:"BROWSER_USER_AGENT" json:"BROWSER_USER_AGENT"`
	Os                  string `datastore:"OS" json:"OS"`
	OsVersion           string `datastore:"OS_VERSION" json:"OS_VERSION"`
	Mobile              string `datastore:"MOBILE" json:"MOBILE"`
	ScreenPoint         string `datastore:"SCREEN_POINT" json:"SCREEN_POINT"`
	ScreenPrint         string `datastore:"SCREEN_PRINT" json:"SCREEN_PRINT"`
	BrowserColorDepth   string `datastore:"BROWSER_COLOR_DEPTH" json:"BROWSER_COLOR_DEPTH"`
	BrowserScreenHeight string `datastore:"BROWSER_SCREEN_HEIGHT" json:"BROWSER_SCREEN_HEIGHT"`
	BrowserScreenWidth  string `datastore:"BROWSER_SCREEN_WIDTH" json:"BROWSER_SCREEN_WIDTH"`
	BrowserPlugins      string `datastore:"BROWSER_PLUGINS" json:"BROWSER_PLUGINS"`
	BrowserJavaEnabled  string `datastore:"BROWSER_JAVA_ENABLED" json:"BROWSER_JAVA_ENABLED"`
	BrowserLanguage     string `datastore:"BROWSER_LANGUAGE" json:"BROWSER_LANGUAGE"`
	BrowserTz           string `datastore:"BROWSER_TZ" json:"BROWSER_TZ"`
	BrowserTzOffset     string `datastore:"BROWSER_TZ_OFFSET" json:"BROWSER_TZ_OFFSET"`
	IPAddress           string `datastore:"IP_ADDRESS" json:"IP_ADDRESS"`
}

const FlagNew = "new"
const FlagValid = "validated"
const FlagActiv = "active"
const FlagCancel = "canceled"

func SetSubscription(c *gin.Context) {
	var newSubscription Subscription
	if err := c.BindJSON(&newSubscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request Errors",
			"error":   err,
		})
		return
	}

	if err := valid(newSubscription.Email); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Email validation failed!",
		})
		return
	}

	// Verify Audience
	if _, ok := IsLicensed(c.Request.Header.Get("token"), newSubscription.MerchantSignature); !ok {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Request has conflict with your authion",
		})
		return
	}

	/* Verify Plan Details */
	// Valid RecurrenceType
	if err := hasValidRecurrenceType(newSubscription.PlanDetails.RecurrenceType); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Wrong recurrence type",
		})
		return
	}

	// Valid RecurrenceDay
	if err := hasValidRecurrenceDay(newSubscription.PlanDetails.RecurrenceType, newSubscription.PlanDetails.RecurrenceDay); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Wrong recurrence Day",
		})
		return
	}

	// Valid Amount
	if err := hasValidAmount(newSubscription.PlanDetails.Amount); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Wrong amount",
		})
		return
	}

	// Valid Currency
	if err := hasValidCurrency(newSubscription.PlanDetails.Currency); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Wrong currency",
		})
		return
	}

	// Valid Card Exp Year
	if err := isValidYear(newSubscription.PaymentConfig.Instrument.ExpYear); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Expired Card",
		})
		return
	}

	// Valid Card Exp Month
	if err := isValidMonth(newSubscription.PaymentConfig.Instrument.ExpMonth); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Invalid Month",
		})
		return
	}

	// Valid Card number
	if err := isValidCard(newSubscription.PaymentConfig.Instrument.Account); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Invalid Card",
		})
		return
	}

	// Valid Card ccv2
	if err := isValidCCV2(newSubscription.PaymentConfig.Instrument.SecretCode); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Invalid ccv2",
		})
		return
	}

	// Validate Instant Payment Type
	if err := hasValidInstantPaymentType(newSubscription.PaymentConfig.InstantPayment); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Invalid instant payment type",
		})
		return
	}

	/**
	add the new subscriber
	*/
	ctx := context.Background()
	projectID := "netopia-payments"

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	newSubscriber := &newSubscription

	/* To make hash string like password, if there is case*/
	// strHash := sha256.Sum256([]byte(newSubscriber.Password))
	// newSubscriber.Password = base64.StdEncoding.EncodeToString(strHash[:])

	// set predefined data
	newSubscriber.Status = true
	newSubscriber.Flags = FlagNew
	newSubscriber.CreatedAt = time.Now()
	// newSubscriber.UpdatedAt = ""

	key := datastore.IncompleteKey("subscribers", nil)
	key.Namespace = "recurring"
	entityKey, err := client.Put(ctx, key, newSubscriber)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code":    "ERROR",
			"message": "Error durring ADD Data",
			"Error":   err.Error(),
		})
		return
	}

	// Define payment detaile
	subscriptionDetails := SubscriptionDetails{}
	// Payment Section
	switch needPayOrValidate := newSubscription.PaymentConfig.InstantPayment; needPayOrValidate {
	case "Yes":
		if err := PayCart(entityKey.ID, newSubscription, c); err == false {
			subscriptionDetails.PaymentStatus = false
			subscriptionDetails.PaymentMessage = "Payment is failed"
		} else {
			subscriptionDetails.PaymentStatus = true
			subscriptionDetails.PaymentMessage = "Payment successfully"
		}
		break
	case "ValidateCard":
		if err := PayOneRon(); err == false {
			subscriptionDetails.ValidateStatus = true
			subscriptionDetails.ValidateMessage = "Payment successfully"
		}
		break
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "subscriber registered successfully!",
		"data": gin.H{
			"subscriptionId": entityKey.ID,
			"details":        subscriptionDetails,
		},
	})
}

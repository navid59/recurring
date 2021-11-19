package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/datastore"
)

const STATUS_PAID = "STATUS_PAID"                       //0x03; // capturate (card) // status : 3
const STATUS_DECLINED = "STATUS_DECLINED"               //0x0c; // declined status	// status : 12
const STATUS_FRAUD = "STATUS_FRAUD"                     //0x0d; // fraud status		// status : 13
const STATUS_3D_PENDING_AUTH = "STATUS_3D_PENDING_AUTH" //0x0f; // 3D authorized 	// Status : 15

// const STATUS_SCHEDULED 								= 7;	//0x07; //scheduled status, specific to Model_Purchase_Sms_Online / Model_Purchase_Sms_Offline
// const STATUS_PROGRAMMED_RECURRENT_PAYMENT 			= 19;	//0x13; //specific to recurrent card purchases
// const STATUS_CANCELED_PROGRAMMED_RECURRENT_PAYMENT 	= 20;	//0x14; //specific to cancelled recurrent card purchases
// const STATUS_EXPIRED								= 23;	//0x17; //cancel a not payed purchase

const FlagApproved = "Approved"
const FlagDeclined = "Declined"
const FlagInvalidAccount = "Invalid Account"

type PaymentsLog struct {
	id             *datastore.Key
	SubscriptionID int64        `json:"subscriptionID"`
	PaymentNtpID   string       `json:"PaymentNtpID" binding:"required"`
	PaymentResult  StartRespons `json:"PaymentResult" binding:"required"`
	PaymentComment string       `json:"PaymentComment"`
	Status         string
	Flags          string
	CreatedAt      time.Time
}

func setPaymentArchive(SubscriptionID int64, sprResponse ActionResponse) bool {
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

	newPaymentLog := PaymentsLog{}

	// set predefined data
	newPaymentLog.SubscriptionID = SubscriptionID
	newPaymentLog.PaymentNtpID = sprResponse.Details.Payment.NtpID
	newPaymentLog.PaymentResult = sprResponse.Details
	newPaymentLog.PaymentComment = sprResponse.Message

	switch sprResponse.Details.Payment.Status {
	case 3:
		newPaymentLog.Status = STATUS_PAID
		break
	case 15:
		newPaymentLog.Status = STATUS_3D_PENDING_AUTH
		break
	case 12:
		newPaymentLog.Status = STATUS_DECLINED
		break
	case 13:
		newPaymentLog.Status = STATUS_FRAUD
		break
	default:
		newPaymentLog.Status = "NO_IDEA"
	}

	newPaymentLog.Flags = sprResponse.Details.Error.Message
	newPaymentLog.CreatedAt = time.Now()

	key := datastore.IncompleteKey("history", nil)
	key.Namespace = "recurring"
	entityKey, err := client.Put(ctx, key, &newPaymentLog)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Printf("The Log ID is : %v", entityKey)
	return true
}

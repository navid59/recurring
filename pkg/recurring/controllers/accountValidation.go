package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

const InstantPaymentYes = "Yes"                   // Exact Fix Date
const InstantPaymentNO = "No"                     // After X Day
const InstantPaymentValidateCard = "ValidateCard" // After X Day

const ErrorCodeNone3DS = "00"
const ErrorCode3DS = "100"
const StatusApprovedNone3DS = 3
const StatusApproved3DS = 15

func PayCart(SubscriptionID int64, newSubscription Subscription, c *gin.Context) bool {
	// 1- setPaymentRequest ->  DONE
	// 2- SendPaymentRequest -> DONE
	// 3- Check if Card has 3dS -> DONE
	// 4- if YES, do 3DS auth and pay return TREU / FALSE -> NOT DONE YET
	// 5- if NO, Just Pay and return TRUE / FALSE -> DONE
	// 6- Save data to Archive / History in any case -> DONE

	// 7- check if payment is with success
	//		  TRUE :-> SET for Recurring entity The TOKEN, NtpID
	//	 		FALSE :-> Not add entity , return False

	startJson := SetPaymentRequest(newSubscription)
	sprResult := SendPaymentRequest(startJson)

	if err := setPaymentArchive(SubscriptionID, sprResult); !err {
		log.Fatalf("Failed to register in Arhive")
	}

	if (sprResult.Details.Error.Code == ErrorCodeNone3DS) && (sprResult.Details.Payment.Status == StatusApprovedNone3DS) {
		return true
	} else if (sprResult.Details.Error.Code == ErrorCode3DS) && sprResult.Details.Payment.Status == StatusApproved3DS {
		TreeDSAuthResult := ThreeDSAuthorize(sprResult.Details)
		if err := setPaymentArchive(SubscriptionID, TreeDSAuthResult); !err {
			log.Fatalf("Failed to register in Arhive")
		}
	} else {
		fmt.Println("-------------------------------")
		fmt.Println(sprResult.Details.Error.Code)
		fmt.Println(sprResult.Details.Payment.Status)
		fmt.Println("-------------------------------")
	}
	return false
}

func PayOneRon() bool {
	return false
}

func hasValidInstantPaymentType(InstantPaymentType string) bool {
	validTypes := map[string]bool{InstantPaymentYes: true, InstantPaymentNO: true, InstantPaymentValidateCard: true}
	if validTypes[InstantPaymentType] {
		return true
	} else {
		return false
	}
}

package controllers

import (
	"net/mail"
	"time"
)

const RecurrenceTypeFix = "Fix"         // Exact Fix Date
const RecurrenceTypeDynamic = "Dynamic" // After X Day

const CurrencyRon = "RON"

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func hasValidRecurrenceType(recurrenceType string) bool {
	validTypes := map[string]bool{RecurrenceTypeFix: true, RecurrenceTypeDynamic: true}
	if validTypes[recurrenceType] {
		return true
	} else {
		return false
	}
}

func hasValidAmount(amount float64) bool {
	if amount <= 0 {
		return false
	}

	return true
}

func isValidMonth(month int) bool {
	if month <= 0 || month > 12 {
		return false
	}
	return true
}

func isValidYear(year int) bool {
	cYear := time.Now().Year()
	if year < cYear {
		return false
	}
	return true
}

func isValidCard(cardNo string) bool {
	if len(cardNo) != 16 {
		return false
	}
	return true
}

func isValidCCV2(ccv2 string) bool {
	if len(ccv2) != 3 {
		return false
	}
	return true
}

func hasValidCurrency(currency string) bool {
	if currency != CurrencyRon {
		return false
	}
	return true
}

func hasValidRecurrenceDay(recurrenceType string, recurrenceDay int8) bool {
	if recurrenceDay <= 0 {
		return false
	}

	if recurrenceType == RecurrenceTypeFix {
		if recurrenceDay > 31 {
			return false
		}
	}

	return true

}

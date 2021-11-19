package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Define the Constants
const RECURRING_POS_SIGNATURE = "LXTP-3WDM-WVXL-GC8B-Y5DA"

const PAYMENT_API_SANDBOX_URL = "https://secure.sandbox.netopia-payments.com/payment/card/start"
const PAYMENT_API_LIVE_URL = "https://secure.netopia-payments.com/payment/card/start"
const PAYMENT_API_KEY = "Uxf3OY--rDK3Qae8CiJJUlAcuRJFp7tzGY4M8KocQaCGyfEqUGhGskv0"
const CONFIG_EMAIL_TEMPLATE = "confirm"
const CONFIG_NOTIFY_URL = "http://35.204.43.65/demoV2/example/ipn.php"
const CONFIG_REDIRECT_URL = "http://35.204.43.65/demoV2/example/backUrl.php"
const CONFIG_LANGUAGE = "RO"
const PAYMENT_OPTIONS_INSTALLMENTS = 1
const PAYMENT_OPTIONS_BONUS = 0
const ORDER_INSTALLMENTS_SELECTED = 1

// The start JSON for Payment API
type Payload struct {
	Config  Config  `json:"config"`
	Payment Payment `json:"payment"`
	Order   Order   `json:"order"`
}
type Config struct {
	EmailTemplate string `json:"emailTemplate"`
	NotifyURL     string `json:"notifyUrl"`
	RedirectURL   string `json:"redirectUrl"`
	Language      string `json:"language"`
}
type Options struct {
	Installments int `json:"installments"`
	Bonus        int `json:"bonus"`
}

// type Instrument struct {
// 	Type       string `datastore:"Type" json:"Type"`
// 	Account    string `datastore:"Account" json:"Account"`
// 	ExpMonth   int    `datastore:"ExpMonth" json:"ExpMonth"`
// 	ExpYear    int    `datastore:"ExpYear" json:"ExpYear"`
// 	SecretCode string `datastore:"SecretCode" json:"SecretCode"`
// 	Token      string `datastore:"Token" json:"Token"`
// }

// type ThreeDS2 struct {
// 	BrowserUserAgent    string `datastore:"BROWSER_USER_AGENT" json:"BROWSER_USER_AGENT"`
// 	Os                  string `datastore:"OS" json:"OS"`
// 	OsVersion           string `datastore:"OS_VERSION" json:"OS_VERSION"`
// 	Mobile              string `datastore:"MOBILE" json:"MOBILE"`
// 	ScreenPoint         string `datastore:"SCREEN_POINT" json:"SCREEN_POINT"`
// 	ScreenPrint         string `datastore:"SCREEN_PRINT" json:"SCREEN_PRINT"`
// 	BrowserColorDepth   string `datastore:"BROWSER_COLOR_DEPTH" json:"BROWSER_COLOR_DEPTH"`
// 	BrowserScreenHeight string `datastore:"BROWSER_SCREEN_HEIGHT" json:"BROWSER_SCREEN_HEIGHT"`
// 	BrowserScreenWidth  string `datastore:"BROWSER_SCREEN_WIDTH" json:"BROWSER_SCREEN_WIDTH"`
// 	BrowserPlugins      string `datastore:"BROWSER_PLUGINS" json:"BROWSER_PLUGINS"`
// 	BrowserJavaEnabled  string `datastore:"BROWSER_JAVA_ENABLED" json:"BROWSER_JAVA_ENABLED"`
// 	BrowserLanguage     string `datastore:"BROWSER_LANGUAGE" json:"BROWSER_LANGUAGE"`
// 	BrowserTz           string `datastore:"BROWSER_TZ" json:"BROWSER_TZ"`
// 	BrowserTzOffset     string `datastore:"BROWSER_TZ_OFFSET" json:"BROWSER_TZ_OFFSET"`
// 	IPAddress           string `datastore:"IP_ADDRESS" json:"IP_ADDRESS"`
// }

type Payment struct {
	Options    Options    `json:"options"`
	Instrument Instrument `json:"instrument"`
	Data       ThreeDS2   `json:"data"`
}
type Billing struct {
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	City       string `json:"city"`
	Country    int    `json:"country"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Details    string `json:"details"`
}
type Shipping struct {
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	City       string `json:"city"`
	Country    int    `json:"country"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Details    string `json:"details"`
}
type Products struct {
	Name     string  `json:"name"`
	Code     string  `json:"code"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Vat      int     `json:"vat"`
}
type Installments struct {
	Selected  int   `json:"selected"`
	Available []int `json:"available"`
}
type Order struct {
	NtpID        string       `json:"ntpID"`
	PosSignature string       `json:"posSignature"`
	DateTime     string       `json:"dateTime"`
	Description  string       `json:"description"`
	OrderID      string       `json:"orderID"`
	Amount       float64      `json:"amount"`
	Currency     string       `json:"currency"`
	Billing      Billing      `json:"billing"`
	Shipping     Shipping     `json:"shipping"`
	Products     []Products   `json:"products"`
	Installments Installments `json:"installments"`
	Data         interface{}  `json:"data"`
}

type StartRespons struct {
	CustomerAction struct {
		AuthenticationToken string `json:"authenticationToken"`
		FormData            struct {
			BackUrl string `json:"backUrl"`
			PaReq   string `json:"paReq"`
		} `json:"formData"`
		Type string `json:"type"`
		Url  string `json:"url"`
	} `json:"customerAction"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Payment struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
		Data     struct {
			AuthCode      string `json:"AuthCode"`
			Bin           string `json:"BIN"`
			Issuer        string `json:"ISSUER"`
			IssuerCountry string `json:"ISSUER_COUNTRY"`
			Rrn           string `json:"RRN"`
		} `json:"data"`
		NtpID  string `json:"ntpID"`
		Status int    `json:"status"`
		Token  string `json:"token"`
	} `json:"payment"`
}

type ActionResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Details StartRespons `json:"details"`
}

type authRespons struct {
	Error   string
	Message string
}

func SetPaymentRequest(newSubscription Subscription) Payload {
	data := Payload{
		Config: Config{
			EmailTemplate: CONFIG_EMAIL_TEMPLATE,
			NotifyURL:     CONFIG_NOTIFY_URL,
			RedirectURL:   CONFIG_REDIRECT_URL,
			Language:      CONFIG_LANGUAGE,
		},
		Payment: Payment{
			Options: Options{
				Installments: PAYMENT_OPTIONS_INSTALLMENTS,
				Bonus:        PAYMENT_OPTIONS_BONUS,
			},
			Instrument: Instrument{
				Type:       newSubscription.PaymentConfig.Instrument.Type,
				Account:    newSubscription.PaymentConfig.Instrument.Account,
				ExpMonth:   newSubscription.PaymentConfig.Instrument.ExpMonth,
				ExpYear:    newSubscription.PaymentConfig.Instrument.ExpYear,
				SecretCode: newSubscription.PaymentConfig.Instrument.SecretCode,
				Token:      newSubscription.PaymentConfig.Instrument.Token,
			},
			Data: ThreeDS2{
				BrowserUserAgent:    newSubscription.PaymentConfig.ThreeDS2.BrowserUserAgent,
				Os:                  newSubscription.PaymentConfig.ThreeDS2.Os,
				OsVersion:           newSubscription.PaymentConfig.ThreeDS2.OsVersion,
				Mobile:              newSubscription.PaymentConfig.ThreeDS2.Mobile,
				ScreenPoint:         newSubscription.PaymentConfig.ThreeDS2.ScreenPoint,
				ScreenPrint:         newSubscription.PaymentConfig.ThreeDS2.ScreenPrint,
				BrowserColorDepth:   newSubscription.PaymentConfig.ThreeDS2.BrowserColorDepth,
				BrowserScreenHeight: newSubscription.PaymentConfig.ThreeDS2.BrowserScreenHeight,
				BrowserScreenWidth:  newSubscription.PaymentConfig.ThreeDS2.BrowserScreenWidth,
				BrowserPlugins:      newSubscription.PaymentConfig.ThreeDS2.BrowserPlugins,
				BrowserJavaEnabled:  newSubscription.PaymentConfig.ThreeDS2.BrowserJavaEnabled,
				BrowserLanguage:     newSubscription.PaymentConfig.ThreeDS2.BrowserLanguage,
				BrowserTz:           newSubscription.PaymentConfig.ThreeDS2.BrowserTz,
				BrowserTzOffset:     newSubscription.PaymentConfig.ThreeDS2.BrowserTzOffset,
				IPAddress:           newSubscription.PaymentConfig.ThreeDS2.IPAddress,
			},
		},
		Order: Order{
			NtpID:        "",
			PosSignature: RECURRING_POS_SIGNATURE,
			DateTime:     "2021-11-02T13:06:22+02:00",
			Description:  "DEMO API FROM GOLANG - V0",
			OrderID:      "RandomRecurringOrder_" + randomBase64String(10),
			Amount:       newSubscription.PlanDetails.Amount,
			Currency:     newSubscription.PlanDetails.Currency,
			Billing: Billing{
				Email:      newSubscription.Email,
				Phone:      newSubscription.Tel,
				FirstName:  newSubscription.Name,
				LastName:   newSubscription.LastName,
				City:       "STATIC.CITY",
				Country:    642,
				State:      "STATIC.STATE",
				PostalCode: "000000",
				Details:    "STATIC_BILLING_DETAILS",
			},
			Shipping: Shipping{
				Email:      newSubscription.Email,
				Phone:      newSubscription.Tel,
				FirstName:  newSubscription.Name,
				LastName:   newSubscription.LastName,
				City:       "STATIC.CITY",
				Country:    642,
				State:      "STATIC.STATE",
				PostalCode: "000000",
				Details:    "STATIC_SHIPING_DETAILS",
			},
			Products: []Products{
				{
					Name:     newSubscription.PlanDetails.Description,
					Code:     "MyProDuctCode",
					Category: "RECURRING",
					Price:    newSubscription.PlanDetails.Amount,
					Vat:      0,
				},
			},
			Installments: Installments{
				Selected:  ORDER_INSTALLMENTS_SELECTED,
				Available: []int{0},
			},
			Data: nil,
		},
	}
	return data
}

func SendPaymentRequest(startJson Payload) ActionResponse {
	actionResponse := ActionResponse{}
	payloadBytes, err := json.Marshal(startJson)
	if err != nil {
		actionResponse.Code = http.StatusNotAcceptable
		actionResponse.Message = "Payment request not acceptable"
		return actionResponse
	}

	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", PAYMENT_API_SANDBOX_URL, body)
	if err != nil {
		actionResponse.Code = http.StatusMisdirectedRequest
		actionResponse.Message = "Misdirected Request"
		return actionResponse
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", PAYMENT_API_KEY)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		actionResponse.Code = http.StatusPaymentRequired
		actionResponse.Message = "Payment failed"
	}

	if resp.StatusCode == http.StatusOK {
		r, _ := ioutil.ReadAll(resp.Body)
		startResp := StartRespons{}
		err = json.Unmarshal(r, &startResp)
		if err != nil {
			// Handel error
		}
		// fmt.Println(startResp.Error.Code)
		// fmt.Println(startResp.Error.Message)

		switch startResp.Error.Code {
		case "100":
			actionResponse.Code = http.StatusNotImplemented
			actionResponse.Message = "Payment Not complated - The card Has 3DS & is not implimented yet"
			actionResponse.Details = startResp
			break
		case "00":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "Payment request accepted"
			actionResponse.Details = startResp
			break
		case "56":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "duplicated Order ID"
			actionResponse.Details = startResp
			break
		case "99":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "There is another order with a different price"
			actionResponse.Details = startResp
			break
		case "19":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "Expire Card Error"
			actionResponse.Details = startResp
			break
		case "20":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "Founduri Error"
			actionResponse.Details = startResp
			break
		case "21":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "CVV Error"
			actionResponse.Details = startResp
			break
		case "22":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "CVV Error"
			actionResponse.Details = startResp
			break
		case "34":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "Card Tranzactie nepermisa Error"
			actionResponse.Details = startResp
			break
		default:
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "Payment API return code : " + startResp.Error.Code
			actionResponse.Details = startResp
		}
	} else {
		// If Payment Request START,  has another status except 200
		actionResponse.Code = resp.StatusCode
		actionResponse.Message = "Problem durring payment"
		actionResponse.Details = StartRespons{}
	}
	defer resp.Body.Close()
	return actionResponse

}

func Ipn(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "This is an test for IPN",
	})
	// return
}

func BackUrl(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusOK,
		"message": "This is an test for BackURL",
	})
	// return
}

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const PAYMENT_API_3DS_VERIFY_AUTH_SANDBOX_URL = "https://secure.sandbox.netopia-payments.com/payment/card/verify-auth"

type VerifyAuthPayload struct {
	AuthenticationToken string   `json:"authenticationToken"`
	NtpID               string   `json:"ntpID"`
	FormData            FormData `json:"formData"`
}

type FormData struct {
	PaRes string `json:"paRes"`
}

func SendRequestVerifyAuth(startRespons StartRespons, paRes string) ActionResponse {
	fmt.Println(paRes)
	fmt.Println(startRespons.CustomerAction.AuthenticationToken)
	fmt.Println(startRespons.Payment.NtpID)

	actionResponse := ActionResponse{}

	data := VerifyAuthPayload{
		AuthenticationToken: startRespons.CustomerAction.AuthenticationToken,
		NtpID:               startRespons.Payment.NtpID,
		FormData: FormData{
			PaRes: paRes,
		},
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(http.MethodPost, PAYMENT_API_3DS_VERIFY_AUTH_SANDBOX_URL, body)
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
		actionResponse.Message = "Verify Auth failed"
	}

	if resp.StatusCode == http.StatusOK {
		r, _ := ioutil.ReadAll(resp.Body)
		verifyAutoResp := StartRespons{}
		err = json.Unmarshal(r, &verifyAutoResp)
		if err != nil {
			// Handel error
		}
		// fmt.Println(startResp.Error.Code)
		// fmt.Println(startResp.Error.Message)

		switch verifyAutoResp.Error.Code {
		case "00":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "3DS - Payment request accepted"
			actionResponse.Details = verifyAutoResp
			break
		case "56":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "3DS - duplicated Order ID"
			actionResponse.Details = verifyAutoResp
			break
		case "99":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "3DS - There is another order with a different price"
			actionResponse.Details = verifyAutoResp
			break
		case "19":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "3DS - Expire Card Error"
			actionResponse.Details = verifyAutoResp
			break
		case "20":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "3DS - Founduri Error"
			actionResponse.Details = verifyAutoResp
			break
		case "21":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "3DS - CVV Error"
			actionResponse.Details = verifyAutoResp
			break
		case "22":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "3DS - CVV Error"
			actionResponse.Details = verifyAutoResp
			break
		case "34":
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "3DS - Card Tranzactie nepermisa Error"
			actionResponse.Details = verifyAutoResp
			break
		default:
			actionResponse.Code = http.StatusOK
			actionResponse.Message = "3DS - Payment API return code : " + verifyAutoResp.Error.Code
			actionResponse.Details = verifyAutoResp
		}
	} else {
		// If Payment Request START,  has another status except 200
		actionResponse.Code = resp.StatusCode
		actionResponse.Message = "Problem durring verify auth in 3DS"
		actionResponse.Details = StartRespons{}
	}

	defer resp.Body.Close()
	return actionResponse
}

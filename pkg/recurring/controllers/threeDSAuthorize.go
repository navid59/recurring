package controllers

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

func ThreeDSAuthorize(startRespons StartRespons) ActionResponse {
	authUrl := startRespons.CustomerAction.Url
	backUrl := startRespons.CustomerAction.FormData.BackUrl
	paReq := startRespons.CustomerAction.FormData.PaReq
	getPaRes, err := call(authUrl, backUrl, paReq)
	if err != nil {
		// Handel Error
	}
	VerifyAuthResult := SendRequestVerifyAuth(startRespons, getPaRes)
	return VerifyAuthResult
}

func call(authUrl string, backUrl string, paReq string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormField("paReq")
	if err != nil {
	}
	_, err = io.Copy(fw, strings.NewReader(paReq))
	if err != nil {
		return "", err
	}
	fw, err = writer.CreateFormField("backUrl")
	if err != nil {
	}
	_, err = io.Copy(fw, strings.NewReader(backUrl))
	if err != nil {
		return "", err
	}

	// Close multipart writer.
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, authUrl, bytes.NewReader(body.Bytes()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, _ := client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}
	content, _ := ioutil.ReadAll(rsp.Body)

	paResStr := getPaRes(string(content))
	return paResStr, nil
}

func getPaRes(str string) string {
	firstPos := strings.Index(str, `input type="hidden" name="paRes" value="`)
	paResPos := firstPos + 40
	paResStr := str[paResPos : paResPos+60]
	return paResStr
}

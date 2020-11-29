package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	post()
}

func post() {

	type Payload struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		RedirectURI  string `json:"redirect_uri"`
	}

	data := Payload{
		GrantType:    "authorization_code",
		ClientID:     "hPYbH5ghUKk0VTMPnkvwGDSm4B4EGApC",
		ClientSecret: "ag20pruBTFmntG8s2DnUrlI7ne6vDVvCIfgLsahiVxnoef-GgxONu7vHZhMjem_0",
		Code:         "9w4LZo91HX4nE78j",
		RedirectURI:  "http://localhost:5000/callback",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error na payloadBytes")
		return
	}
	body := bytes.NewReader(payloadBytes)

	//req, err := http.NewRequest("POST", "https://auth.atlassian.com/oauth/token", body)
	req, err := http.NewRequest("POST", "http://localhost:5000/callback", body)
	if err != nil {
		fmt.Println("error na req")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error na resp \n", err)
		return
	}
	defer resp.Body.Close()
}

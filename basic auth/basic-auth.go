package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/*
type fields struct {
	Assignee
	Reporter
	Priority
	Sprint
	FixVersions
	Description
	ShortDescription
	Status
	Comments
}
*/

type assignee struct {
	emailAddress string
	displayName  string
}

func main() {
	basicAuth()
}

func basicAuth() {
	var m1 map[string]string
	var m2 map[string]map[string]string
	var m3 map[string]map[string]map[string]string
	var m4 map[string]map[string]map[string]map[string]string
	var m5 map[string]map[string]map[string]map[string]map[string]string

	var username string = "w.demko@gmail.com"
	var passwd string = "gjKuO2RAgrFq2gMnqYE5E04C"
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://projectx-tadw.atlassian.net/rest/api/3/issue/A2A-17", nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)

	json.Unmarshal([]byte(s), &m1)
	json.Unmarshal([]byte(s), &m2)
	json.Unmarshal([]byte(s), &m3)
	json.Unmarshal([]byte(s), &m4)
	json.Unmarshal([]byte(s), &m5)
	fmt.Println("Assignee: ", m3["fields"]["assignee"]["displayName"])
	fmt.Println("Reporter: ", m3["fields"]["reporter"]["displayName"])
	fmt.Println("Priority: ", m3["fields"]["priority"]["id"])
	fmt.Println("Sprint: ", m3["fields"]["customfield_10020"]["name"])
	fmt.Println("FixVersions: ", m3["fields"]["fixVersions"]["name"])
	fmt.Println("Description: ", m5["fields"]["description"]["content"]["content"]["text"])
	fmt.Println("ShortDescription: ", m2["fields"]["summary"])
	fmt.Println("Status: ", m3["fields"]["status"]["name"])
	//fmt.Println("Comments: ",
	//return s
}

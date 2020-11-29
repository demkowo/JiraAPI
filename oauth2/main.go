package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

//errorRepo in case of curl error
type errorRepo struct {
	F1 string `json:"error"`
	F2 string `json:"error_description"`
}

// flag for new token creation
var ret bool = false

// JiraID is a variable to store task nr
var jid string = ""

//to prevent repetetive actions
var end string = ""
var e = &end

// ReqToken is what we get from Jira
type ReqToken struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

// AccToken is a struct that represents a single post
type AccToken struct {
	Act string `json:"access_token"`
	Scp string `json:"scope"`
	Exp int    `json:"expires_in"`
	Tty string `json:"token_type"`
}

// JiraTask data fetched from Jira to be presented on page
type JiraTask struct {
	Assignee         string
	Reporter         string
	Priority         string
	Sprint           string
	FixVersions      string
	Description      string
	ShortDescription string
	Status           string
	Comments         string
}

// credentials
var cid string = "hPYbH5ghUKk0VTMPnkvwGDSm4B4EGApC"
var sec string = "ag20pruBTFmntG8s2DnUrlI7ne6vDVvCIfgLsahiVxnoef-GgxONu7vHZhMjem_0"
var code string = "FTNU1kuZl5OhkTw2"

//piinter to template
var tpl *template.Template

func init() {
	log.Println("wywolano init")
	log.Println("wgrywanie templatow")
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/getJiraID", getJiraID).Methods("POST")
	r.HandleFunc("/getJiraID", redirect).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(redirect)
	http.ListenAndServe(":5000", r)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	log.Println("wywolano redirect()")
	fmt.Println("status:", *e)
	if *e != "" {
		log.Println("ale nic sie nie zadzialo")
		return
	}
	log.Println("przekierowanie do main z redirect")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("wywolano index")
	log.Println("wgrano template dla index")
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func getJiraID(w http.ResponseWriter, r *http.Request) {
	log.Println("wywolano getJiraID")

	id := &jid
	*id = r.FormValue("taskID")

	rt := &ret

	fmt.Println("jid: ", jid)
	log.Println("przekierowanie do index z getJiraID")
	if *rt != false {
		requestToken()
	} else {
		requestData()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func requestToken() {
	log.Println("wywolano curl - request token")

	//data for curl script
	text := "#!/bin/bash \n\n curl --request POST --url 'https://auth.atlassian.com/oauth/token' --header 'Content-Type: application/json' --data '{\"grant_type\": \"authorization_code\",\"client_id\": \"hPYbH5ghUKk0VTMPnkvwGDSm4B4EGApC\",\"client_secret\": \"ag20pruBTFmntG8s2DnUrlI7ne6vDVvCIfgLsahiVxnoef-GgxONu7vHZhMjem_0\",\"code\": \"" + code + "\",\"redirect_uri\": \"https://localhost\"}'"

	//creating curl script
	fmt.Println("\n\ntext: ", text, "\n\n ")
	err := ioutil.WriteFile("command.sh", []byte(text), 0755)

	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}

	//executing curl script
	out, err := exec.Command("/bin/sh", "./command.sh").Output()
	if err != nil {
		log.Fatal("Error with bash: ", err)
		return
	}
	s := string(out)
	fmt.Println("s: ", s)
	data := AccToken{}
	de := errorRepo{}
	json.Unmarshal([]byte(s), &data)
	dataAct := string(data.Act)
	fmt.Println("data: ", dataAct)

	//verification if aythorization code is still valid
	if len(dataAct) < 20 {
		log.Println("Error: not enough data", len(dataAct), &de)
		json.Unmarshal([]byte(s), &de)
		log.Println("Error: ", &de)
		return
	}

	//creating curl script for data request
	text = "#!/bin/bash \n\ncurl --request GET --url https://projectx-tadw.atlassian.net/rest/api/2/issue/" + jid + " --header 'Authorization: Bearer " + dataAct + "' --header 'Accept: application/json'"
	//text := "#!/bin/bash \n\ncurl --request GET --url https://projectx-tadw.atlassian.net/rest/api/2/issue/" + jid + " --header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik16bERNemsxTVRoRlFVRTJRa0ZGT0VGRk9URkJOREJDTVRRek5EZzJSRVpDT1VKRFJrVXdNZyJ9.eyJodHRwczovL2F0bGFzc2lhbi5jb20vb2F1dGhDbGllbnRJZCI6ImhQWWJINWdoVUtrMFZUTVBua3Z3R0RTbTRCNEVHQXBDIiwiaHR0cHM6Ly9hdGxhc3NpYW4uY29tL2VtYWlsRG9tYWluIjoiZ21haWwuY29tIiwiaHR0cHM6Ly9hdGxhc3NpYW4uY29tL3N5c3RlbUFjY291bnRJZCI6IjVmYmQ1ZTI4ZjJkZjZjMDA3NjhhMzRlZSIsImh0dHBzOi8vYXRsYXNzaWFuLmNvbS9zeXN0ZW1BY2NvdW50RW1haWxEb21haW4iOiJjb25uZWN0LmF0bGFzc2lhbi5jb20iLCJodHRwczovL2F0bGFzc2lhbi5jb20vdmVyaWZpZWQiOnRydWUsImh0dHBzOi8vYXRsYXNzaWFuLmNvbS9maXJzdFBhcnR5IjpmYWxzZSwiaHR0cHM6Ly9hdGxhc3NpYW4uY29tLzNsbyI6dHJ1ZSwiaXNzIjoiaHR0cHM6Ly9hdGxhc3NpYW4tYWNjb3VudC1wcm9kLnB1czIuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfDVmNzc4ZTFmNGQwOWY3MDA3NmZiNzdhZCIsImF1ZCI6ImFwaS5hdGxhc3NpYW4uY29tIiwiaWF0IjoxNjA2NjEwMTU0LCJleHAiOjE2MDY2MTM3NTQsImF6cCI6ImhQWWJINWdoVUtrMFZUTVBua3Z3R0RTbTRCNEVHQXBDIiwic2NvcGUiOiJyZWFkOmppcmEtd29yayByZWFkOmppcmEtdXNlciJ9.h7SyZ3erOqPd10okB3CYrH0b8Fx2m25J_J-h5WBVfZ84x8lunXQlOIad7ba4XYTRpBkMyldOZ8WcGGSHaNhXgT80K13EWgoEfUeD0vtAzMs8Kk1IXo1yax-uRef9G4nMFGPLAsadYWhPCCVV_YKD5aKJ7Xyx1nFNvLe51Pb8lVJJMj-jnxpJOGAa2VpBXD57Z80iC0pWt-Rnpf4hykYrEYVEa1B7VZRkd9Uu79JlsfkEF7eBlGVw2EXRbAIgxg5nK0uNoEwlboZ5BXSpJx3SLDrwrCpvemKQ1HJA_oAhN2JivbqoJKhYHR_NCHxIqA95WMChLo7hH7Ohu69Hf0wLDQ' --header 'Accept: application/json'"
	fmt.Println("\n\ntext: ", text, "\n\n ")

	err = ioutil.WriteFile("command2.sh", []byte(text), 0755)

	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}

	//requestData()
}

func requestData() {
	var m1 map[string]string
	var m2 map[string]map[string]string
	var m3 map[string]map[string]map[string]string
	var m4 map[string]map[string]map[string]map[string]string
	var m5 map[string]map[string]map[string]map[string]map[string]string

	log.Println("wywolano curl - request data")

	out, err := exec.Command("/bin/sh", "./command2.sh").Output()
	if err != nil {
		log.Fatal("wyklada sie na bash: ", err)
		return
	}
	s := string(out)
	fmt.Println("s: ", s)
	if s == `{"errorMessages":["Issue does not exist or you do not have permission to see it."],"errors":{}}` {
		fmt.Println("mozna wywolac pobranie tokenu")
		requestToken()

		out, err := exec.Command("/bin/sh", "./command2.sh").Output()
		if err != nil {
			log.Fatal("wyklada sie na bash: ", err)
			return
		}
		s := string(out)
		fmt.Println("s: ", s)
	}
	json.Unmarshal([]byte(s), &m1)
	json.Unmarshal([]byte(s), &m2)
	json.Unmarshal([]byte(s), &m3)
	json.Unmarshal([]byte(s), &m4)
	json.Unmarshal([]byte(s), &m5)
	fmt.Println("")
	fmt.Println("Assignee: ", m3["fields"]["assignee"]["displayName"])
	fmt.Println("Reporter: ", m3["fields"]["reporter"]["displayName"])
	fmt.Println("Priority: ", m3["fields"]["priority"]["id"])
	fmt.Println("Sprint: ", m3["fields"]["customfield_10020"]["name"])
	fmt.Println("FixVersions: ", m3["fields"]["fixVersions"]["name"])
	fmt.Println("Description: ", m5["fields"]["description"]["content"]["content"]["text"])
	fmt.Println("ShortDescription: ", m2["fields"]["summary"])
	fmt.Println("Status: ", m3["fields"]["status"]["name"])
	fmt.Println("")
	*e = "done"
	log.Println("===========  zakonczono ===========")
}

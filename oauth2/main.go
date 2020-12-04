package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

// Values represents values
type Values struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Code         string `json:"code"`
	Jid          string
}

var val = Values{
	GrantType: "refresh_token",

	AccessToken:  "",
	Scope:        "manage:jira-configuration manage:jira-project write:jira-work read:jira-work read:jira-user manage:jira-data-provider offline_access",
	ExpiresIn:    3600,
	TokenType:    "Bearer",
	Jid:          "",
	ClientID:     "[Client Id From Jira Configuration]",
	ClientSecret: "[Client Id Secret Jira Configuration]",
	Code:         "[Code received after Jira app authorizatopn]",
	RefreshToken: "[Refresh Token received after Jira app authorization]",
}

/*
//ReqRefTok reperesents struct of data required for refresh token request
type ReqRefTok struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

var rrt = ReqRefTok{
	GrantType:    "refresh_token",
	ClientID:     cid,
	ClientSecret: sec,
	RefreshToken: reft,
}

//GetRefTok reperesents struct of data received after refresh token request
type GetRefTok struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

var grt GetRefTok
*/
//var posts []GetRefTok = []GetRefTok{}

//errorRepo in case of curl error
type errorRepo struct {
	F1 string `json:"error"`
	F2 string `json:"error_description"`
}

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

var tpl *template.Template // pointer to templates
var end bool = false       // flag to make sure that there are no actions when program ends
var tf bool = false        // flag to check if token is valid
var tok string = ""        // Bearer token
//var jid string = ""        // JiraID is a variable to store task nr
var er string = "" // error info for logs

func init() {
	//log.Println("wywolano init")
	//log.Println("wgrywanie templatow")
	start()
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	r := mux.NewRouter()
	//r.NotFoundHandler = http.HandlerFunc(redirect)
	r.HandleFunc("/", index)
	r.HandleFunc("/getJiraID", getJiraID).Methods("POST")
	//r.HandleFunc("/", getRefreshToken).Methods("POST")
	//r.HandleFunc("/getJiraID", redirect).Methods("GET")
	http.ListenAndServe(":5000", r)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	log.Println("wywolano redirect()")

	e := &end

	fmt.Println("status completed = ", *e)
	if *e == true {
		log.Println("No redirect actions, status 'completed' = true")
		return
	}
	log.Println("przekierowanie do main z redirect")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("index()")
	tpl.ExecuteTemplate(w, "index.gohtml", nil)

	refreshToken()
	/*
		// -- ZADANIE NA POZNIEJ -- upewnic sie czy potrzeba kolejkowanie

		//getJiraID starts on demand
		t := &tok
		if tok == "" {
			fmt.Println("token = ", *t)
			// -- ZADANIE NA POZNIEJ -- upewnic sie czy potrzeba dodac czas dla token zeby sie nie nadpisywaly
			prepareRequestToken()
			executeRequestToken()
			if tok == "" {
				log.Println("token: ", tok)
				refreshToken()
			}
			prepareRequestData()
			executeRequestData()
		} else {
			fmt.Println("token = ", *t)
			prepareRequestData()
			executeRequestData()
		}
	*/
}

func getJiraID(w http.ResponseWriter, r *http.Request) {
	log.Println("wywolano getJiraID")

	/*
		id := &jid
		*id = r.FormValue("taskID")
	*/
	val.Jid = r.FormValue("taskID")
	fmt.Println("jid = ", val.Jid)

	log.Println("wywolano http.Redirect z getJiraID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func prepareRequestToken() {
	log.Println("wywolano curl - request token")

	//data for curl script
	text := "#!/bin/bash \n\ncurl --request POST --url 'https://auth.atlassian.com/oauth/token' --header 'Content-Type: application/json' --data '{\"grant_type\": \"authorization_code\",\"client_id\": \"hPYbH5ghUKk0VTMPnkvwGDSm4B4EGApC\",\"client_secret\": \"ag20pruBTFmntG8s2DnUrlI7ne6vDVvCIfgLsahiVxnoef-GgxONu7vHZhMjem_0\",\"code\": \"" + val.Code + "\",\"redirect_uri\": \"https://localhost\"}'"

	//creating curl script
	fmt.Println("\n\ntext: ", text, "\n\n ")
	err := ioutil.WriteFile("command.sh", []byte(text), 0755)

	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
}

func executeRequestToken() {
	log.Println("wywolano execute request token")

	e := &er
	//executing curl script
	out, err := exec.Command("/bin/sh", "./command.sh").Output()
	if err != nil {
		log.Fatal("Error with bash command.sh: ", err)
		return
	}

	s := string(out)
	*e = s

	// when everything is fine
	data := AccToken{}
	json.Unmarshal([]byte(s), &data)

	//authorization varification (if user code is valid)
	if len(*e) < 80 {
		log.Println("ERROR exReqT, can't get authorization token", er)
	} else {
		fmt.Println("token:", er)
	}

	t := &tok
	*t = data.Act
	return
}

func refreshToken() {
	log.Println("refreshToken()")

	// data to bytes
	dby, err := json.Marshal(val)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	body := bytes.NewReader(dby)

	// send refresh token request
	req, err := http.NewRequest("POST", "https://auth.atlassian.com/oauth/token", body)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// bytes to getResponseToken struct
	v := &val
	err = json.Unmarshal(out, v)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(v.TokenType)
	fmt.Println(v.RefreshToken)
	fmt.Println(v.ClientSecret)
	fmt.Println(v.ClientID)
	fmt.Println(v.ExpiresIn)
	fmt.Println(v.AccessToken)

	defer resp.Body.Close()
}

func requestData() {

}

func prepareRequestData() {
	log.Println("wywolano prepareRequestData")

	t := &tok
	//creating curl script for data request
	text := "#!/bin/bash \n\ncurl --request GET --url https://projectx-tadw.atlassian.net/rest/api/2/issue/" + val.Jid + " --header 'Authorization: Bearer " + *t + "' --header 'Accept: application/json'"
	//text := "#!/bin/bash \n\ncurl --request GET --url https://projectx-tadw.atlassian.net/rest/api/2/issue/" + jid + " --header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik16bERNemsxTVRoRlFVRTJRa0ZGT0VGRk9URkJOREJDTVRRek5EZzJSRVpDT1VKRFJrVXdNZyJ9.eyJodHRwczovL2F0bGFzc2lhbi5jb20vb2F1dGhDbGllbnRJZCI6ImhQWWJINWdoVUtrMFZUTVBua3Z3R0RTbTRCNEVHQXBDIiwiaHR0cHM6Ly9hdGxhc3NpYW4uY29tL2VtYWlsRG9tYWluIjoiZ21haWwuY29tIiwiaHR0cHM6Ly9hdGxhc3NpYW4uY29tL3N5c3RlbUFjY291bnRJZCI6IjVmYmQ1ZTI4ZjJkZjZjMDA3NjhhMzRlZSIsImh0dHBzOi8vYXRsYXNzaWFuLmNvbS9zeXN0ZW1BY2NvdW50RW1haWxEb21haW4iOiJjb25uZWN0LmF0bGFzc2lhbi5jb20iLCJodHRwczovL2F0bGFzc2lhbi5jb20vdmVyaWZpZWQiOnRydWUsImh0dHBzOi8vYXRsYXNzaWFuLmNvbS9maXJzdFBhcnR5IjpmYWxzZSwiaHR0cHM6Ly9hdGxhc3NpYW4uY29tLzNsbyI6dHJ1ZSwiaXNzIjoiaHR0cHM6Ly9hdGxhc3NpYW4tYWNjb3VudC1wcm9kLnB1czIuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfDVmNzc4ZTFmNGQwOWY3MDA3NmZiNzdhZCIsImF1ZCI6ImFwaS5hdGxhc3NpYW4uY29tIiwiaWF0IjoxNjA2NjEwMTU0LCJleHAiOjE2MDY2MTM3NTQsImF6cCI6ImhQWWJINWdoVUtrMFZUTVBua3Z3R0RTbTRCNEVHQXBDIiwic2NvcGUiOiJyZWFkOmppcmEtd29yayByZWFkOmppcmEtdXNlciJ9.h7SyZ3erOqPd10okB3CYrH0b8Fx2m25J_J-h5WBVfZ84x8lunXQlOIad7ba4XYTRpBkMyldOZ8WcGGSHaNhXgT80K13EWgoEfUeD0vtAzMs8Kk1IXo1yax-uRef9G4nMFGPLAsadYWhPCCVV_YKD5aKJ7Xyx1nFNvLe51Pb8lVJJMj-jnxpJOGAa2VpBXD57Z80iC0pWt-Rnpf4hykYrEYVEa1B7VZRkd9Uu79JlsfkEF7eBlGVw2EXRbAIgxg5nK0uNoEwlboZ5BXSpJx3SLDrwrCpvemKQ1HJA_oAhN2JivbqoJKhYHR_NCHxIqA95WMChLo7hH7Ohu69Hf0wLDQ' --header 'Accept: application/json'"
	fmt.Println("\n\ntext: ", text, "\n\n ")

	err := ioutil.WriteFile("command2.sh", []byte(text), 0755)

	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}

	return
}

func executeRequestData() {
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
		//////////////wyslac info ze trzeba pobrac token

		out, err := exec.Command("/bin/sh", "./command2.sh").Output()
		if err != nil {
			log.Fatal("wyklada sie na bash: ", err)
			return
		}
		s := string(out)
		fmt.Println("s: ", s)
	}
	var m1 map[string]string
	var m2 map[string]map[string]string
	var m3 map[string]map[string]map[string]string
	var m4 map[string]map[string]map[string]map[string]string
	var m5 map[string]map[string]map[string]map[string]map[string]string

	json.Unmarshal([]byte(s), &m1)
	json.Unmarshal([]byte(s), &m2)
	json.Unmarshal([]byte(s), &m3)
	json.Unmarshal([]byte(s), &m4)
	json.Unmarshal([]byte(s), &m5)
	fmt.Println("")
	fmt.Println("Issue: ", m1["key"])
	fmt.Println("Assignee: ", m3["fields"]["assignee"]["displayName"])
	fmt.Println("Reporter: ", m3["fields"]["reporter"]["displayName"])
	fmt.Println("Priority: ", m3["fields"]["priority"]["id"])
	fmt.Println("Sprint: ", m3["fields"]["customfield_10020"]["name"])
	fmt.Println("FixVersions: ", m3["fields"]["fixVersions"]["name"])
	fmt.Println("Description: ", m5["fields"]["description"]["content"]["content"]["text"])
	fmt.Println("ShortDescription: ", m2["fields"]["summary"])
	fmt.Println("Status: ", m3["fields"]["status"]["name"])
	fmt.Println("")
	e := &end
	*e = true
	log.Println("===========  zakonczono ===========")
	return
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//ReqRefTok reperesents struct of data required for refresh token request
type ReqRefTok struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

var rrt = ReqRefTok{
	GrantType:    "refresh_token",
	ClientID:     "[Client Id From Jira Configuration]",
	ClientSecret: "[Client Id Secret Jira Configuration]",
	RefreshToken: "[Refresh Token received after Jira app authorization]",
}

//GetRefTok reperesents struct of data received after refresh token request
type GetRefTok struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	F1          string `json:"error"`
	F2          string `json:"error_description"`
}

// GetDataResponse is a struct with Jira ticket data
type GetDataResponse struct {
	Expand string `json:"expand"`
	ID     string `json:"id"`
	Self   string `json:"self"`
	Key    string `json:"key"`
	Fields struct {
		Statuscategorychangedate string `json:"statuscategorychangedate"`
		Issuetype                struct {
			Self        string `json:"self"`
			ID          string `json:"id"`
			Description string `json:"description"`
			IconURL     string `json:"iconUrl"`
			Name        string `json:"name"`
			Subtask     bool   `json:"subtask"`
			AvatarID    int    `json:"avatarId"`
		} `json:"issuetype"`
		Timespent interface{} `json:"timespent"`
		Project   struct {
			Self           string `json:"self"`
			ID             string `json:"id"`
			Key            string `json:"key"`
			Name           string `json:"name"`
			ProjectTypeKey string `json:"projectTypeKey"`
			Simplified     bool   `json:"simplified"`
			AvatarUrls     struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
		} `json:"project"`
		FixVersions []struct {
			Self        string `json:"self"`
			ID          string `json:"id"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Archived    bool   `json:"archived"`
			Released    bool   `json:"released"`
		} `json:"fixVersions"`
		Aggregatetimespent interface{} `json:"aggregatetimespent"`
		Resolution         interface{} `json:"resolution"`
		Resolutiondate     interface{} `json:"resolutiondate"`
		Workratio          int         `json:"workratio"`
		LastViewed         string      `json:"lastViewed"`
		Issuerestriction   struct {
			Issuerestrictions struct {
			} `json:"issuerestrictions"`
			ShouldDisplay bool `json:"shouldDisplay"`
		} `json:"issuerestriction"`
		Watches struct {
			Self       string `json:"self"`
			WatchCount int    `json:"watchCount"`
			IsWatching bool   `json:"isWatching"`
		} `json:"watches"`
		Created          string `json:"created"`
		Customfield10020 []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			State   string `json:"state"`
			BoardID int    `json:"boardId"`
		} `json:"customfield_10020"`
		Customfield10021 interface{} `json:"customfield_10021"`
		Customfield10022 interface{} `json:"customfield_10022"`
		Customfield10023 interface{} `json:"customfield_10023"`
		Priority         struct {
			Self    string `json:"self"`
			IconURL string `json:"iconUrl"`
			Name    string `json:"name"`
			ID      string `json:"id"`
		} `json:"priority"`
		Customfield10024 string      `json:"customfield_10024"`
		Customfield10025 interface{} `json:"customfield_10025"`
		Labels           []string    `json:"labels"`
		Customfield10016 interface{} `json:"customfield_10016"`
		Customfield10017 interface{} `json:"customfield_10017"`
		Customfield10018 struct {
			HasEpicLinkFieldDependency bool `json:"hasEpicLinkFieldDependency"`
			ShowField                  bool `json:"showField"`
			NonEditableReason          struct {
				Reason  string `json:"reason"`
				Message string `json:"message"`
			} `json:"nonEditableReason"`
		} `json:"customfield_10018"`
		Customfield10019              string        `json:"customfield_10019"`
		Timeestimate                  interface{}   `json:"timeestimate"`
		Aggregatetimeoriginalestimate interface{}   `json:"aggregatetimeoriginalestimate"`
		Versions                      []interface{} `json:"versions"`
		Issuelinks                    []interface{} `json:"issuelinks"`
		Assignee                      struct {
			Self         string `json:"self"`
			AccountID    string `json:"accountId"`
			EmailAddress string `json:"emailAddress"`
			AvatarUrls   struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			TimeZone    string `json:"timeZone"`
			AccountType string `json:"accountType"`
		} `json:"assignee"`
		Updated string `json:"updated"`
		Status  struct {
			Self           string `json:"self"`
			Description    string `json:"description"`
			IconURL        string `json:"iconUrl"`
			Name           string `json:"name"`
			ID             string `json:"id"`
			StatusCategory struct {
				Self      string `json:"self"`
				ID        int    `json:"id"`
				Key       string `json:"key"`
				ColorName string `json:"colorName"`
				Name      string `json:"name"`
			} `json:"statusCategory"`
		} `json:"status"`
		Components           []interface{} `json:"components"`
		Timeoriginalestimate interface{}   `json:"timeoriginalestimate"`
		Description          string        `json:"description"`
		Customfield10010     interface{}   `json:"customfield_10010"`
		Customfield10014     interface{}   `json:"customfield_10014"`
		Customfield10015     interface{}   `json:"customfield_10015"`
		Timetracking         struct {
		} `json:"timetracking"`
		Customfield10005      interface{}   `json:"customfield_10005"`
		Customfield10006      interface{}   `json:"customfield_10006"`
		Security              interface{}   `json:"security"`
		Customfield10007      interface{}   `json:"customfield_10007"`
		Customfield10008      interface{}   `json:"customfield_10008"`
		Customfield10009      interface{}   `json:"customfield_10009"`
		Attachment            []interface{} `json:"attachment"`
		Aggregatetimeestimate interface{}   `json:"aggregatetimeestimate"`
		Summary               string        `json:"summary"`
		Creator               struct {
			Self       string `json:"self"`
			AccountID  string `json:"accountId"`
			AvatarUrls struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			TimeZone    string `json:"timeZone"`
			AccountType string `json:"accountType"`
		} `json:"creator"`
		Subtasks []interface{} `json:"subtasks"`
		Reporter struct {
			Self       string `json:"self"`
			AccountID  string `json:"accountId"`
			AvatarUrls struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active      bool   `json:"active"`
			TimeZone    string `json:"timeZone"`
			AccountType string `json:"accountType"`
		} `json:"reporter"`
		Aggregateprogress struct {
			Progress int `json:"progress"`
			Total    int `json:"total"`
		} `json:"aggregateprogress"`
		Customfield10000 string      `json:"customfield_10000"`
		Customfield10001 interface{} `json:"customfield_10001"`
		Customfield10002 interface{} `json:"customfield_10002"`
		Customfield10003 interface{} `json:"customfield_10003"`
		Customfield10004 interface{} `json:"customfield_10004"`
		Environment      interface{} `json:"environment"`
		Duedate          interface{} `json:"duedate"`
		Progress         struct {
			Progress int `json:"progress"`
			Total    int `json:"total"`
		} `json:"progress"`
		Comment struct {
			Comments []struct {
				Self   string `json:"self"`
				ID     string `json:"id"`
				Author struct {
					Self         string `json:"self"`
					AccountID    string `json:"accountId"`
					EmailAddress string `json:"emailAddress"`
					AvatarUrls   struct {
						Four8X48  string `json:"48x48"`
						Two4X24   string `json:"24x24"`
						One6X16   string `json:"16x16"`
						Three2X32 string `json:"32x32"`
					} `json:"avatarUrls"`
					DisplayName string `json:"displayName"`
					Active      bool   `json:"active"`
					TimeZone    string `json:"timeZone"`
					AccountType string `json:"accountType"`
				} `json:"author"`
				Body         string `json:"body"`
				UpdateAuthor struct {
					Self         string `json:"self"`
					AccountID    string `json:"accountId"`
					EmailAddress string `json:"emailAddress"`
					AvatarUrls   struct {
						Four8X48  string `json:"48x48"`
						Two4X24   string `json:"24x24"`
						One6X16   string `json:"16x16"`
						Three2X32 string `json:"32x32"`
					} `json:"avatarUrls"`
					DisplayName string `json:"displayName"`
					Active      bool   `json:"active"`
					TimeZone    string `json:"timeZone"`
					AccountType string `json:"accountType"`
				} `json:"updateAuthor"`
				Created   string `json:"created"`
				Updated   string `json:"updated"`
				JsdPublic bool   `json:"jsdPublic"`
			} `json:"comments"`
			MaxResults int `json:"maxResults"`
			Total      int `json:"total"`
			StartAt    int `json:"startAt"`
		} `json:"comment"`
		Votes struct {
			Self     string `json:"self"`
			Votes    int    `json:"votes"`
			HasVoted bool   `json:"hasVoted"`
		} `json:"votes"`
		Worklog struct {
			StartAt    int           `json:"startAt"`
			MaxResults int           `json:"maxResults"`
			Total      int           `json:"total"`
			Worklogs   []interface{} `json:"worklogs"`
		} `json:"worklog"`
	} `json:"fields"`
}

var grt GetRefTok
var code string = "[Code received after Jira app authorizatopn]"
var jid string = ""        // JiraID is a variable to store task nr
var tpl *template.Template // pointer to templates
var end bool = false       // flag to make sure that there are no actions when program ends
var gdr GetDataResponse

/*
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

var tf bool = false        // flag to check if token is valid
var tok string = ""        // Bearer token
var er string = ""         // error info for logs
*/
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
	requestData()
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

	id := &jid
	*id = r.FormValue("taskID")

	fmt.Println("jid = ", *id)

	log.Println("wywolano http.Redirect z getJiraID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
func prepareRequestToken() {
	log.Println("wywolano curl - request token")

	//data for curl script
	text := "#!/bin/bash \n\ncurl --request POST --url 'https://auth.atlassian.com/oauth/token' --header 'Content-Type: application/json' --data '{\"grant_type\": \"authorization_code\",\"client_id\": \"hPYbH5ghUKk0VTMPnkvwGDSm4B4EGApC\",\"client_secret\": \"ag20pruBTFmntG8s2DnUrlI7ne6vDVvCIfgLsahiVxnoef-GgxONu7vHZhMjem_0\",\"code\": \"" + code + "\",\"redirect_uri\": \"https://localhost\"}'"

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
*/
func refreshToken() {
	log.Println("refreshToken()")

	// data to bytes
	dby, err := json.Marshal(rrt)

	if err != nil {
		fmt.Println("Error 1: ", err)
	}

	body := bytes.NewReader(dby)

	// send refresh token request
	req, err := http.NewRequest("POST", "https://auth.atlassian.com/oauth/token", body)
	if err != nil {
		fmt.Println("Error 2: ", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error 3: ", err)
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error 4: ", err)
	}

	// bytes to getResponseToken struct
	v := &grt
	err = json.Unmarshal(out, v)
	if err != nil {
		fmt.Println("Error 5: ", err)
	}

	fmt.Println(v)
	fmt.Println(v.F1)
	fmt.Println(v.F2)

	defer resp.Body.Close()
}

func requestData() {
	log.Println("requestData()")

	j := &jid
	link := "https://projectx-tadw.atlassian.net/rest/api/2/issue/" + *j
	fmt.Println(link)
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println("Error 6: ", err)
	}

	v := &grt
	bea := "Bearer " + v.AccessToken

	req.Header.Set("Authorization", bea)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error 7: ", err)
	}

	out, err := ioutil.ReadAll(resp.Body)
	va := &gdr
	err = json.Unmarshal(out, va)

	if err != nil {
		fmt.Println("Error 8: ", err)
	}

	fmt.Println(va.Fields.Assignee.DisplayName)
	fmt.Println(va.Fields.Assignee.EmailAddress)
	defer resp.Body.Close()
}

/*
func prepareRequestData() {
	log.Println("wywolano prepareRequestData")

	t := &tok
	//creating curl script for data request
	text := "#!/bin/bash \n\ncurl --request GET --url https://projectx-tadw.atlassian.net/rest/api/2/issue/" + jid + " --header 'Authorization: Bearer " + *t + "' --header 'Accept: application/json'"
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
*/

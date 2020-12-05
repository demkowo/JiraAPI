package main

import (
	"html/template"
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

func init() {
	start()
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/getJiraID", getJiraID).Methods("POST")
	http.ListenAndServe(":5000", r)
}

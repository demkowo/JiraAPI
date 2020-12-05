package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

	refreshToken(w, r)
	requestData(w, r)

}

func getJiraID(w http.ResponseWriter, r *http.Request) {
	log.Println("wywolano getJiraID")

	id := &jid
	*id = r.FormValue("taskID")

	fmt.Println("jid = ", *id)

	log.Println("wywolano http.Redirect z getJiraID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	log.Println("refreshToken()")

	// data to bytes
	dby, err := json.Marshal(rrt)

	if err != nil {
		fmt.Println("Error 1: ", err)
	}

	body := bytes.NewReader(dby)

	// send request
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

func requestData(w http.ResponseWriter, r *http.Request) {
	log.Println("requestData()")

	j := &jid
	link := "https://projectx-tadw.atlassian.net/rest/api/2/issue/" + *j
	fmt.Println(link)

	// send request
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

	//va := &gdr

	tpl.ExecuteTemplate(io.Writter, "index.gohtml", va)
}

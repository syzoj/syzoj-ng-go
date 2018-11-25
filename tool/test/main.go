package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type GroupCreateRequest struct {
	GroupName string `json:"group_name"`
}
type ProblemsetCreateRequest struct {
	GroupName      string `json:"group_name"`
	ProblemsetName string `json:"problemset_name"`
	ProblemsetType string `json:"problemset_type"`
}
type ProblemCreateRequest struct {
	GroupName      string `json:"group_name"`
	ProblemsetName string `json:"problemset_name"`
	ProblemName    string `json:"problem_name"`
	ProblemType    string `json:"problem_type"`
}

func Test(client *http.Client, method string, ep string, data interface{}) {
	fmt.Printf("Doing %s to %s with data %+v\n", method, ep, data)
	var body io.Reader
	if data != nil {
		bodyContent, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(bodyContent)
	}
	req, err := http.NewRequest(method, "http://127.0.0.1:5900"+ep, body)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Printf("Result: %+v\n", result)
	fmt.Printf("Headers: %+v\n", resp.Header)
	u, _ := url.Parse("http://127.0.0.1:5900")
	client.Jar.SetCookies(u, resp.Cookies())
}

func main() {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}
	Test(client, "POST", "/api/auth/register", RegisterRequest{"aaa", "B"})
	Test(client, "POST", "/api/auth/login", LoginRequest{"aaa", "B"})
	u, _ := url.Parse("http://127.0.0.1:5900/")
	fmt.Println(cookieJar.Cookies(u))
	Test(client, "POST", "/api/group/create", GroupCreateRequest{"aaa"})
	Test(client, "POST", "/api/group/problemset/create", ProblemsetCreateRequest{"aaa", "bbb", "standard"})
	Test(client, "POST", "/api/group/problemset/problem/create", ProblemCreateRequest{"aaa", "bbb", "1", "standard"})
}

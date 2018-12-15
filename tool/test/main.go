package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type CreateProblemRequest struct {
	Statement ProblemStatement `json:"statement"`
}
type ProblemStatement struct {
	Title     string `json:"title"`
	Statement string `json:"statement"`
}
type ProblemsetCreateRequest struct{}
type ProblemsetAddProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	Name         string    `json:"name"`
	ProblemId    uuid.UUID `json:"problem_id"`
}
type ProblemsetListProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
}
type ProblemsetViewProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	Name         string    `json;"name"`
}
type ProblemsetSubmitRequest struct {
	ProblemsetId uuid.UUID                     `json:"problemset_id"`
	ProblemName  string                        `json:"problem_name"`
	Type         string                        `json:"type"`
	Traditional  *TraditionalSubmissionRequest `json:"traditional"`
}
type TraditionalSubmissionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code'`
}

func Test(client *http.Client, method string, ep string, data interface{}) (result interface{}) {
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
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Printf("Result: %+v\n", result)
	fmt.Printf("Headers: %+v\n", resp.Header)
	u, _ := url.Parse("http://127.0.0.1:5900")
	client.Jar.SetCookies(u, resp.Cookies())
	return
}

func main() {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}
	Test(client, "POST", "/api/auth/register", RegisterRequest{"aaa", "B"})
	Test(client, "POST", "/api/auth/logout", nil)
	Test(client, "POST", "/api/auth/login", LoginRequest{"aaa", "B"})
	u, _ := url.Parse("http://127.0.0.1:5900/")
	fmt.Println(cookieJar.Cookies(u))

	resp1 := Test(client, "POST", "/api/problem/create", CreateProblemRequest{
		Statement: ProblemStatement{
			Title:     "title",
			Statement: "statement",
		},
	})
	problemId := uuid.MustParse(resp1.(map[string]interface{})["data"].(map[string]interface{})["problem_id"].(string))

	resp2 := Test(client, "POST", "/api/problemset/create", ProblemsetCreateRequest{})
	problemsetId := uuid.MustParse(resp2.(map[string]interface{})["data"].(map[string]interface{})["problemset_id"].(string))
	Test(client, "POST", "/api/problemset/add", ProblemsetAddProblemRequest{
		ProblemsetId: problemsetId,
		ProblemId:    problemId,
		Name:         "1",
	})
	Test(client, "GET", "/api/problemset/list", ProblemsetListProblemRequest{
		ProblemsetId: problemsetId,
	})
	Test(client, "GET", "/api/problemset/view", ProblemsetViewProblemRequest{
		ProblemsetId: problemsetId,
		Name:         "1",
	})
	Test(client, "POST", "/api/problemset/submit", ProblemsetSubmitRequest{
		ProblemsetId: problemsetId,
		ProblemName:  "1",
		Type:         "traditional",
		Traditional: &TraditionalSubmissionRequest{
			Language: "cpp",
			Code: `
#include<cstdio>
int main(){
    freopen("problem.in","r",stdin);
    freopen("problem.out","w",stdout);
    int a,b;
    scanf("%d%d",&a,&b);
    printf("%d\n",a+b);
}
`,
		},
	})
}

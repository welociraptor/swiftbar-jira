package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/welociraptor/swiftbar-jira/config"
	"github.com/welociraptor/swiftbar-jira/jira"
)

var (
	//go:embed output.tmpl
	files embed.FS
)

func main() {
	configuration := config.Load()

	responses := ExecuteQueries(configuration)

	tmpl, err := template.New("output.tmpl").ParseFS(files, "output.tmpl")
	if err != nil {
		panic(err)
	}

	tmpl.Execute(os.Stdout, responses)
}

type Response struct {
	Header string
	Issues []struct {
		Key    string
		Fields struct {
			Summary string
		}
	}
}

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func ExecuteQueries(c *config.Configuration) []Response {
	responseCh := make(chan Response, len(c.Queries))

	responses := make([]Response, len(c.Queries))

	for _, query := range c.Queries {
		go executeQuery(query, c.JiraUrl, c.JiraToken, responseCh)
	}

	for i, _ := range c.Queries {
		responses[i] = <-responseCh
	}

	return responses
}

func executeQuery(q jira.Query, jiraUrl, jiraToken string, ch chan<- Response) {
	req, _ := http.NewRequest(http.MethodPost, jiraUrl, q.GetBuffer())

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", jiraToken))
	req.Header.Add("Content-Type", "application/json")

	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer closeBody(resp.Body)

	resb, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jiraResp := Response{}
	jiraResp.Header = q.Header

	err = json.Unmarshal(resb, &jiraResp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ch <- jiraResp
}

package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
)

var (
	//go:embed output.tmpl
	files     embed.FS
	JiraUrl   string
	JiraToken string
)

func main() {
	if JiraUrl == "" {
		JiraUrl = os.Getenv("JIRA_URL")
	}

	if JiraToken == "" {
		JiraToken = os.Getenv("JIRA_TOKEN")
	}

	req, _ := http.NewRequest(http.MethodPost, JiraUrl, Query())

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", JiraToken))
	req.Header.Add("Content-Type", "application/json")

	c := &http.Client{}

	jiraResp, err := DoRequest(c, req)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("output.tmpl").ParseFS(files, "output.tmpl")
	if err != nil {
		panic(err)
	}

	tmpl.Execute(os.Stdout, jiraResp)
}

type Response struct {
	Issues []struct {
		Key    string
		Fields struct {
			Summary string
		}
	}
}

func DoRequest(c *http.Client, r *http.Request) (*Response, error) {
	resp, err := c.Do(r)

	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		return nil, err
	}

	resb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jiraResp Response

	err = json.Unmarshal(resb, &jiraResp)

	return &jiraResp, err
}

func Query() *bytes.Buffer {
	query := map[string]interface{}{
		"jql": "labels = SRE ORDER BY created DESC",
		/*		"startAt":    0,
				"maxResults": 25, */
		"fields": []string{
			"summary", "status",
		},
	}

	b, err := json.Marshal(query)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}

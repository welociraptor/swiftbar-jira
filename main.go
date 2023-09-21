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

	jiraResp, err := DoRequest(req)
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

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func DoRequest(r *http.Request) (*Response, error) {
	c := &http.Client{}

	resp, err := c.Do(r)
	defer closeBody(resp.Body)
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

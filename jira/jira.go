package jira

import (
	"bytes"
	"encoding/json"
)

type Query struct {
	Header     string   `json:"-"`
	JQL        string   `json:"jql"`
	StartAt    int      `json:"startAt"`
	MaxResults int      `json:"maxResults"`
	Fields     []string `json:"fields"`
}

func (q *Query) GetBuffer() *bytes.Buffer {
	q.StartAt = 0
	q.MaxResults = 10

	b, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(b)
}

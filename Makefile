jira.5m.bin:
	go build -ldflags="-s -w -X 'main.JiraToken=${JIRA_TOKEN}' -X 'main.JiraUrl=${JIRA_URL}'" -o jira.5m.bin

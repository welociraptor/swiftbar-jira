package config

import (
	"fmt"
	"github.com/welociraptor/swiftbar-jira/jira"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Configuration struct {
	JiraUrl   string `yaml:"jiraUrl"`
	JiraToken string `yaml:"jiraToken"`
	Queries   []jira.Query
}

func Load() *Configuration {
	path := os.Getenv("SWIFTBAR_PLUGIN_PATH")
	if path != "" {
		path = path[:strings.LastIndex(path, "/")] + "/swiftbar-jira.yaml"
	} else {
		path = "swiftbar-jira.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configuration := Configuration{}

	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &configuration
}

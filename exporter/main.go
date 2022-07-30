package main

// https://developer.atlassian.com/cloud/jira/platform/rest/v3/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Username    string `mapstructure:"JIRA_USER"`
	Token       string `mapstructure:"JIRA_TOKEN"`
	JiraBaseURL string `mapstructure:"JIRA_BASE_URL"`
}

type rawJiraIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Labels   []string `json:"labels"`
		Assignee struct {
			Name string `json:"DisplayName"`
		} `json:"assignee"`
		Components []struct {
			Name string `json:"name"`
		} `json:"components"`
		Reporter struct {
			Name string `json:"DisplayName"`
		}
		IssueType struct {
			Name string `json:"name"`
		} `json:"issuetype"`
		Team struct {
			Name string `json:"value"`
		} `json:"customfield_10032"`
		Sprints  []Sprint `json:"customfield_10019"`
		Priority struct {
			Value string `json:"name"`
		} `json:"priority"`
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
		StoryPoint float32 `json:"customfield_10025"`
		Subtasks   []struct {
			Key string `json:"key"`
		} `json:"subtasks"`
		Parent struct {
			Key string `json:"key"`
		}
		CreatedDate    string `json:"created"`
		ResolutionDate string `json:"resolutiondate"`
	} `json:"fields"`
}

type Sprint struct {
	Id        int    `json:"id"`
	State     string `json:"state"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Name      string `json:"name"`
}

type Issue struct {
	Key            string
	Labels         []string
	Components     []string
	Reporter       string
	Assignee       string
	Type           string
	Team           string
	Status         string
	StoryPoint     float32
	Sprints        []Sprint
	SubTasks       []*Issue
	Parent         *Issue
	CreatedDate    string
	ResolutionDate string
}

func newIssue(raw *rawJiraIssue) *Issue {
	issue := Issue{
		Key:            raw.Key,
		Labels:         raw.Fields.Labels,
		Reporter:       raw.Fields.Reporter.Name,
		Assignee:       raw.Fields.Assignee.Name,
		Type:           raw.Fields.IssueType.Name,
		Team:           raw.Fields.Team.Name,
		Status:         raw.Fields.Status.Name,
		StoryPoint:     raw.Fields.StoryPoint,
		Sprints:        raw.Fields.Sprints,
		CreatedDate:    raw.Fields.CreatedDate,
		ResolutionDate: raw.Fields.ResolutionDate,
	}

	for _, value := range raw.Fields.Components {
		issue.Components = append(issue.Components, value.Name)
	}

	return &issue
}

type JiraClient struct {
	config Config
}

func NewJiraClient() *JiraClient {
	viper.BindEnv("JIRA_USER")
	viper.BindEnv("JIRA_TOKEN")
	viper.BindEnv("JIRA_BASE_URL")
	var config Config
	viper.Unmarshal(&config)

	return &JiraClient{
		config: config,
	}
}

func (jira *JiraClient) getRawIssue(key string) (*rawJiraIssue, error) {

	url := jira.config.JiraBaseURL + key

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("Internal error - could not create new request %v", err)
	}

	request.SetBasicAuth(jira.config.Username, jira.config.Token)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("Internal Error - could not execute request: %v", err)
	}

	decoder := json.NewDecoder(response.Body)
	raw := rawJiraIssue{}
	decoder.Decode(&raw)

	return &raw, nil
}

func (jira *JiraClient) GetIssue(key string) (*Issue, error) {
	raw, err := jira.getRawIssue(key)

	if err != nil {
		return nil, err
	}

	issue := newIssue(raw)

	for _, subtask := range raw.Fields.Subtasks {
		rawSubtask, err := jira.getRawIssue(subtask.Key)

		if err != nil {
			return nil, fmt.Errorf("Could not fetch subtask, %v", err)
		}

		issue.SubTasks = append(issue.SubTasks, newIssue(rawSubtask))
	}

	rawParentKey := raw.Fields.Parent.Key
	rawParent, err := jira.getRawIssue(rawParentKey)

	if err != nil {
		return nil, fmt.Errorf("Could not fetch subtask, %v", err)
	}

	issue.Parent = newIssue(rawParent)

	return issue, nil
}

func main() {

	if len(os.Args) == 1 {
		log.Fatal("Error - please specify ticket ID (usage: ./main ticketID)")
	}

	ticketID := os.Args[1]

	jira := NewJiraClient()
	issue, err := jira.GetIssue(ticketID)

	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%v", issue.SubTasks)

	fmt.Printf("%v", issue)
}

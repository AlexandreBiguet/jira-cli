# Jira-cli

Note: Super custom, will probably not work on your side

## Build

`go build exporter/main.go`

## Usage

Expects 3 env variables to be defined:

- `JIRA_USER` : your username
- `JIRA_TOKEN`: the api you created [somehow](https://id.atlassian.com/manage-profile/security/api-tokens)
- `JIRA_BASE_URL`: Something like `https://your-company-name.atlassian.net/rest/api/3/issue/` (TODO: shouldn't contain resource path)

To execute:

`./main <ticketID>`

# Jira-cli

## Don't use that tool

This 'cli' will probably not work on your side as jira lets its users define any custom fields. As such, the mapping between jira custom field and their meaning is probably different between you and what is hardcoded here.

Next time our custom fields change, I'll add an intermediate representation to make this "configurable".

## Build

`go build exporter/main.go`

## Usage

Expects 3 env variables to be defined:

- `JIRA_USER` : your username
- `JIRA_TOKEN`: the api token you created [somehow](https://id.atlassian.com/manage-profile/security/api-tokens)
- `JIRA_BASE_URL`: something like `https://your-company-name.atlassian.net/rest/api/3/issue/` (TODO: shouldn't contain resource path)

To execute:

`./main <ticketID>`

package github

import "github.com/google/go-github/github"

const (
	EventIssues      = "GitHub:IssuesEvent"
	EventPullRequest = "GitHub:PullRequestEvent"
	EventStatus      = "GitHub:StatusEvent"
	EventPing        = "GitHub:PingEvent"
)

type IssuesEvent github.IssuesEvent

func (e IssuesEvent) EventName() string {
	return EventIssues
}

type PullRequestEvent github.PullRequestEvent

func (e PullRequestEvent) EventName() string {
	return EventPullRequest
}

type StatusEvent github.StatusEvent

func (e StatusEvent) EventName() string {
	return EventStatus
}

type PingEvent github.PingEvent

func (e PingEvent) EventName() string {
	return EventPing
}

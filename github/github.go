package github

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strconv"
	"time"
)

type GitHub struct {
	client *github.Client
}

func New(accessToken string) *GitHub {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	return &GitHub{
		client,
	}
}

func (g *GitHub) GenerateReport() (string, error) {
	issues, err := g.pickup()
	if err != nil {
		return "", err
	}
	report := g.dailyReport(issues)
	return report, nil
}

func (g *GitHub) pickup() ([]*github.Issue, error) {
	year, month, day := time.Now().Date()
	loc, _ := time.LoadLocation("UTC")
	startTime := time.Date(year, month, day, 1, 0, 0, 0, loc)
	listOptions := &github.IssueListOptions{
		Filter:    "created",
		State:     "all",
		Labels:    []string{},
		Sort:      "updated",
		Direction: "asc",
		Since:     startTime,
	}
	issues, _, err := g.client.Issues.List(true, listOptions)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (g *GitHub) dailyReport(issues []*github.Issue) string {
	report := ""
	for _, issue := range issues {
		report = report + "[" + *issue.Repository.FullName + "] " + "[#" + strconv.Itoa(*issue.Number) + " " + *issue.Title + "](" + *issue.HTMLURL + ")\n\n"
	}
	return report
}

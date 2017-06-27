package beater

import (
	"github.com/elastic/beats/libbeat/common"
	"time"
)

type Event struct {
	ReadTime          time.Time
	DocumentType      string `json:"document_type"`
	FullName          string `json:"github_repository_full_name"`
	Owner             string `json:"github_repository_owner"`
	Name              string `json:"github_repository_name"`
	StargazersCount   int    `json:"github_repository_stargazers_count"`
	PullRequestsCount int    `json:"github_repository_pull_requests_count"`
	OpenIssuesCount   int    `json:"github_repository_open_issues_count"`
	ForksCount        int    `json:"github_repository_forks_count"`
	ReleasesCount     int    `json:"github_repository_releases_count"`
}

func (h *Event) ToMapStr() common.MapStr {

	event := common.MapStr{
		"@timestamp":                            common.Time(h.ReadTime),
		"type":                                  h.DocumentType,
		"document_type":                         h.DocumentType,
		"github_repository_full_name":           h.FullName,
		"github_repository_owner":               h.Owner,
		"github_repository_name":                h.Name,
		"github_repository_stargazers_count":    h.StargazersCount,
		"github_repository_pull_requests_count": h.PullRequestsCount,
		"github_repository_open_issues_count":   h.OpenIssuesCount,
		"github_repository_forks_count":         h.ForksCount,
		"github_repository_releases_count":      h.ReleasesCount,
	}

	return event
}

package beater

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
)

type Event struct {
	ReadTime          time.Time
	DocumentType      string `json:"document_type"`
	FullName          string `json:"github.repository.full_name"`
	Owner             string `json:"github.repository.owner"`
	Name              string `json:"github.repository.name"`
	StargazersCount   int    `json:"github.repository.stargazers_count"`
	PullRequestsCount int    `json:"github.repository.pull_requests_count"`
	OpenIssuesCount   int    `json:"github.repository.open_issues_count"`
	ForksCount        int    `json:"github.repository.forks_count"`
	ReleasesCount     int    `json:"github.repository.releases_count"`
}

func (h *Event) ToMapStr() common.MapStr {

	event := common.MapStr{
		"@timestamp":                            common.Time(h.ReadTime),
		"type":                                  h.DocumentType,
		"document_type":                         h.DocumentType,
		"github.repository.fullname":            h.FullName,
		"github.repository.owner":               h.Owner,
		"github.repository.name":                h.Name,
		"github.repository.stargazers_count":    h.StargazersCount,
		"github.repository.pull_requests_count": h.PullRequestsCount,
		"github.repository.open_issues_count":   h.OpenIssuesCount,
		"github.repository.forks_count":         h.ForksCount,
		"github.repository.releases_count":      h.ReleasesCount,
	}

	return event
}

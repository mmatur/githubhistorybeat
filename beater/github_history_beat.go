package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/mmatur/githubhistorybeat/config"
	"sort"
	"time"
)

type GithubHistoryBeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &GithubHistoryBeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *GithubHistoryBeat) Run(b *beat.Beat) error {
	logp.Info("githubhistorybeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()

	for _, cr := range bt.config.Repositories {
		times := make(map[time.Time]*Event)

		repo := NewGithubRepository(bt, cr.Owner, cr.Name)

		date := repo.GetRoundedCreateAt(cr.TimeInterval)
		for {
			times[date] = &Event{
				ReadTime:          date,
				DocumentType:      cr.DocumentType,
				FullName:          fmt.Sprintf("%s/%s", repo.Owner, repo.Name),
				Owner:             repo.Owner,
				Name:              repo.Name,
				StargazersCount:   0,
				OpenIssuesCount:   0,
				PullRequestsCount: 0,
				ForksCount:        0,
				ReleasesCount:     0,
			}
			date = date.Add(cr.TimeInterval)
			if date.After(time.Now().Round(cr.TimeInterval)) {
				break
			}
		}
		repo.FetchForks(times, cr.TimeInterval)
		repo.FetchStargazers(times, cr.TimeInterval)
		repo.FetchReleases(times, cr.TimeInterval)
		repo.FetchIssues(times, cr.TimeInterval)
		repo.FetchPullRequest(times, cr.TimeInterval)

		times_sorted := []*Event{}

		for _, v := range times {
			times_sorted = append(times_sorted, v)
		}

		sort.Slice(times_sorted, func(i, j int) bool {
			return times_sorted[i].ReadTime.Before(times_sorted[j].ReadTime)
		})

		var lastEvent *Event
		for _, t := range times_sorted {
			event := t
			if lastEvent != nil {
				event.ForksCount += lastEvent.ForksCount
				event.StargazersCount += lastEvent.StargazersCount
				event.ReleasesCount += lastEvent.ReleasesCount
				event.OpenIssuesCount += lastEvent.OpenIssuesCount
				event.PullRequestsCount += lastEvent.PullRequestsCount
			}
			lastEvent = event
			bt.client.PublishEvent(event.ToMapStr())
			logp.Info("Event sent")
		}
	}

	time.Sleep(30 * time.Second)
	return nil
}

func (bt *GithubHistoryBeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

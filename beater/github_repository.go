package beater

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/elastic/beats/libbeat/logp"
)

type GithubRepository struct {
	Owner        string
	Name         string
	CreatedAt    time.Time
	Client       *github.Client
	Token        string
	Context      context.Context
	Stargazers   []*github.Stargazer
	Forks        []*github.Repository
	Issues       []*github.Issue
	PullRequests []*github.PullRequest
	Releases     []*github.RepositoryRelease
}

func NewGithubRepository(bt *GithubHistoryBeat, owner string, repository string) *GithubRepository {
	logp.Info("Token Used : %s", bt.config.Token)
	client := CreateGithubClient(bt.config.Token)
	gr := &GithubRepository{
		Owner:        owner,
		Name:         repository,
		CreatedAt:    time.Now(),
		Client:       client,
		Token:        bt.config.Token,
		Context:      context.Background(),
		Stargazers:   []*github.Stargazer{},
		Forks:        []*github.Repository{},
		Issues:       []*github.Issue{},
		PullRequests: []*github.PullRequest{},
		Releases:     []*github.RepositoryRelease{},
	}
	gr.FetchRepositoryInfo()
	return gr
}

func CreateGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func (gr *GithubRepository) GetRoundedCreateAt(timeInterval time.Duration) time.Time {
	return gr.CreatedAt.Round(timeInterval)
}

func (gr *GithubRepository) FetchRepositoryInfo() error {
	logp.Info("Start fetching repository info for %s/%s", gr.Owner, gr.Name)
	repository, _, err := gr.Client.Repositories.Get(gr.Context, gr.Owner, gr.Name)
	if err != nil {
		return err
	}
	gr.CreatedAt = repository.CreatedAt.Time
	return nil
}

func (gr *GithubRepository) FetchStargazers(times map[time.Time]*Event, timeInterval time.Duration) ([]*github.Stargazer, error) {
	logp.Info("Start fetching stargazers for %s/%s", gr.Owner, gr.Name)
	opt := &github.ListOptions{PerPage: 100}
	for {
		stargazers, response, err := gr.Client.Activity.ListStargazers(gr.Context, gr.Owner, gr.Name, opt)
		logp.Info("Remaining github http call %d/%d", response.Rate.Remaining, response.Rate.Limit)
		if err != nil {
			return gr.Stargazers, err
		}
		gr.Stargazers = append(gr.Stargazers, stargazers...)
		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}
	for _, r := range gr.Stargazers {
		rounded := r.StarredAt.Time.Round(timeInterval)
		v := times[rounded]
		logp.Info("rouded:%s", rounded)
		v.StargazersCount++
	}
	return gr.Stargazers, nil
}

func (gr *GithubRepository) FetchForks(times map[time.Time]*Event, timeInterval time.Duration) ([]*github.Repository, error) {
	logp.Info("Start fetching forks for %s/%s", gr.Owner, gr.Name)
	opt := &github.RepositoryListForksOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		forks, response, err := gr.Client.Repositories.ListForks(gr.Context, gr.Owner, gr.Name, opt)
		logp.Info("Remaining github http call %d/%d", response.Rate.Remaining, response.Rate.Limit)
		if err != nil {
			return gr.Forks, err
		}
		gr.Forks = append(gr.Forks, forks...)
		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}
	for _, f := range gr.Forks {
		rounded := f.CreatedAt.Round(timeInterval)
		v := times[rounded]
		v.ForksCount++
	}

	return gr.Forks, nil
}

func (gr *GithubRepository) FetchIssues(times map[time.Time]*Event, timeInterval time.Duration) ([]*github.Issue, error) {
	logp.Info("Start fetching issues for %s/%s", gr.Owner, gr.Name)
	opt := &github.IssueListByRepoOptions{
		State:       "all",
		Direction:   "asc",
		Sort:        "created",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		issues, response, err := gr.Client.Issues.ListByRepo(gr.Context, gr.Owner, gr.Name, opt)
		logp.Info("Remaining github http call %d/%d", response.Rate.Remaining, response.Rate.Limit)
		if err != nil {
			return gr.Issues, err
		}
		gr.Issues = append(gr.Issues, issues...)
		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}
	for _, i := range gr.Issues {
		rounded := i.CreatedAt.Round(timeInterval)
		if i.PullRequestLinks == nil {
			v := times[rounded]
			v.OpenIssuesCount++
			if i.ClosedAt != nil {
				rounded := i.ClosedAt.Round(timeInterval)
				v := times[rounded]
				v.OpenIssuesCount--
			}
		}
	}
	return gr.Issues, nil
}

func (gr *GithubRepository) FetchPullRequest(times map[time.Time]*Event, timeInterval time.Duration) ([]*github.PullRequest, error) {
	logp.Info("Start fetching PR for %s/%s", gr.Owner, gr.Name)
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		pullRequests, response, err := gr.Client.PullRequests.List(gr.Context, gr.Owner, gr.Name, opt)
		logp.Info("Remaining github http call %d/%d", response.Rate.Remaining, response.Rate.Limit)
		if err != nil {
			return gr.PullRequests, err
		}
		gr.PullRequests = append(gr.PullRequests, pullRequests...)
		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}
	for _, p := range gr.PullRequests {
		rounded := p.CreatedAt.Round(timeInterval)
		v := times[rounded]
		v.PullRequestsCount++

	}
	return gr.PullRequests, nil
}

func (gr *GithubRepository) FetchReleases(times map[time.Time]*Event, timeInterval time.Duration) ([]*github.RepositoryRelease, error) {
	logp.Info("Start fetching releases for %s/%s", gr.Owner, gr.Name)
	opt := &github.ListOptions{PerPage: 100}
	for {
		releases, response, err := gr.Client.Repositories.ListReleases(gr.Context, gr.Owner, gr.Name, opt)
		logp.Info("Remaining github http call %d/%d", response.Rate.Remaining, response.Rate.Limit)
		if err != nil {
			return gr.Releases, err
		}
		gr.Releases = append(gr.Releases, releases...)
		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}
	for _, r := range gr.Releases {
		rounded := r.CreatedAt.Time.Round(timeInterval)
		v := times[rounded]
		v.ReleasesCount++
	}
	return gr.Releases, nil
}

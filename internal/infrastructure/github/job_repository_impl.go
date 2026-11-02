package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	domainrepo "github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
	"github.com/cli/go-gh/v2/pkg/api"
)

// maxConcurrentRequests limits the number of parallel GitHub API requests
const maxConcurrentRequests = 10

// perPage is the page size used for paginated GitHub API requests
const perPage = 100

// orgRepoListTTL is how long the organization's repository list is cached.
// The list rarely changes, so there is no need to fetch it on every refresh.
const orgRepoListTTL = 5 * time.Minute

// JobRepositoryImpl implements the JobRepository interface using GitHub API
type JobRepositoryImpl struct {
	restClient *api.RESTClient

	cacheMu       sync.Mutex
	cachedOrg     string
	cachedRepos   []repoRef
	cachedReposAt time.Time
}

// NewJobRepository creates a new instance of JobRepositoryImpl
func NewJobRepository() (domainrepo.JobRepository, error) {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w\nPlease run 'gh auth login' to authenticate with GitHub", err)
	}

	return &JobRepositoryImpl{
		restClient: restClient,
	}, nil
}

// FetchActiveJobs retrieves all active jobs for a repository or organization
func (j *JobRepositoryImpl) FetchActiveJobs(ctx context.Context, owner, repo, org string) ([]*entity.Job, error) {
	targets, skipErrors, err := j.resolveTargets(ctx, owner, repo, org)
	if err != nil {
		return nil, err
	}

	runs, err := j.fetchActiveRuns(ctx, targets, skipErrors)
	if err != nil {
		return nil, err
	}

	// Fetch jobs for each run in parallel, keeping the original order
	jobsPerRun := make([][]*entity.Job, len(runs))
	forEachConcurrent(len(runs), func(i int) {
		jobs, err := j.getJobsForRun(ctx, runs[i].repo, runs[i].run)
		if err != nil {
			return // Skip this run if we can't get jobs
		}
		jobsPerRun[i] = jobs
	})

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var allJobs []*entity.Job
	for _, jobs := range jobsPerRun {
		allJobs = append(allJobs, jobs...)
	}
	return allJobs, nil
}

// resolveTargets decides which repositories to scan for workflow runs.
// skipErrors is true when scanning a whole organization, where some repositories may be unreadable.
func (j *JobRepositoryImpl) resolveTargets(ctx context.Context, owner, repo, org string) (targets []repoRef, skipErrors bool, err error) {
	if org == "" {
		return []repoRef{{owner: owner, name: repo}}, false, nil
	}

	// GitHub has no org-level endpoint for workflow runs, so each repository is queried individually.
	targets, err = j.orgRepositories(ctx, org)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch repositories for organization %s: %w", org, err)
	}
	return targets, true, nil
}

// orgRepositories returns the organization's repositories, using a cached list while it is fresh
func (j *JobRepositoryImpl) orgRepositories(ctx context.Context, org string) ([]repoRef, error) {
	j.cacheMu.Lock()
	defer j.cacheMu.Unlock()

	if j.cachedOrg == org && time.Since(j.cachedReposAt) < orgRepoListTTL {
		return j.cachedRepos, nil
	}

	repos, err := j.fetchOrgRepositories(ctx, org)
	if err != nil {
		return nil, err
	}

	refs := make([]repoRef, 0, len(repos))
	for _, r := range repos {
		refs = append(refs, repoRef{owner: r.Owner.Login, name: r.Name})
	}

	j.cachedOrg = org
	j.cachedRepos = refs
	j.cachedReposAt = time.Now()
	return refs, nil
}

// repoRef identifies a repository
type repoRef struct {
	owner string
	name  string
}

// runRef is a workflow run together with the repository it belongs to
type runRef struct {
	repo repoRef
	run  workflowRun
}

// activeRunStatuses are the workflow run statuses that can have active jobs
var activeRunStatuses = []string{"in_progress", "queued"}

// fetchActiveRuns fetches in_progress and queued runs for the given repositories in parallel.
// When skipErrors is true, repositories that fail (e.g. Actions disabled or no permission) are ignored.
func (j *JobRepositoryImpl) fetchActiveRuns(ctx context.Context, targets []repoRef, skipErrors bool) ([]runRef, error) {
	n := len(targets) * len(activeRunStatuses)
	runsPerReq := make([][]workflowRun, n)
	errs := make([]error, n)

	forEachConcurrent(n, func(i int) {
		target := targets[i/len(activeRunStatuses)]
		status := activeRunStatuses[i%len(activeRunStatuses)]
		runs, err := j.fetchWorkflowRuns(ctx, target.owner, target.name, status)
		if err != nil {
			errs[i] = fmt.Errorf("failed to fetch %s runs for %s/%s: %w", status, target.owner, target.name, err)
			return
		}
		runsPerReq[i] = runs
	})

	var result []runRef
	for i, runs := range runsPerReq {
		if errs[i] != nil {
			if skipErrors {
				continue
			}
			return nil, errs[i]
		}
		for _, run := range runs {
			result = append(result, runRef{repo: targets[i/len(activeRunStatuses)], run: run})
		}
	}
	return result, nil
}

// forEachConcurrent calls fn for each index in [0, n) with bounded concurrency and waits for all to finish
func forEachConcurrent(n int, fn func(i int)) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentRequests)
	for i := 0; i < n; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			fn(i)
		}()
	}
	wg.Wait()
}

// fetchOrgRepositories lists all non-archived repositories in an organization
func (j *JobRepositoryImpl) fetchOrgRepositories(ctx context.Context, org string) ([]repositoryResponse, error) {
	var result []repositoryResponse

	for page := 1; ; page++ {
		path := fmt.Sprintf("orgs/%s/repos?type=all&per_page=%d&page=%d", org, perPage, page)
		var repos []repositoryResponse
		if err := j.getJSON(ctx, path, &repos); err != nil {
			return nil, err
		}

		for _, r := range repos {
			// Archived repositories are read-only and can't run workflows
			if !r.Archived {
				result = append(result, r)
			}
		}

		if len(repos) < perPage {
			break
		}
	}

	return result, nil
}

// fetchWorkflowRuns fetches workflow runs with a specific status, following pagination
func (j *JobRepositoryImpl) fetchWorkflowRuns(ctx context.Context, owner, repo, status string) ([]workflowRun, error) {
	var allRuns []workflowRun

	for page := 1; ; page++ {
		path := fmt.Sprintf("repos/%s/%s/actions/runs?status=%s&per_page=%d&page=%d", owner, repo, status, perPage, page)
		var runs workflowRunsResponse
		if err := j.getJSON(ctx, path, &runs); err != nil {
			return nil, fmt.Errorf("failed to request workflow runs: %w", err)
		}

		allRuns = append(allRuns, runs.WorkflowRuns...)

		// If we got fewer items than per_page, we've reached the last page
		if len(runs.WorkflowRuns) < perPage {
			break
		}
	}

	return allRuns, nil
}

// getJobsForRun fetches and converts active jobs for a specific workflow run
func (j *JobRepositoryImpl) getJobsForRun(ctx context.Context, repo repoRef, run workflowRun) ([]*entity.Job, error) {
	path := fmt.Sprintf("repos/%s/%s/actions/runs/%d/jobs", repo.owner, repo.name, run.ID)
	var jobs jobsResponse
	if err := j.getJSON(ctx, path, &jobs); err != nil {
		return nil, fmt.Errorf("failed to request jobs: %w", err)
	}

	var result []*entity.Job
	for _, job := range jobs.Jobs {
		if j.isActiveJob(job.Status) {
			result = append(result, &entity.Job{
				ID:           job.ID,
				RunID:        job.RunID,
				Name:         job.Name,
				Status:       job.Status,
				RunnerID:     job.RunnerID,
				RunnerName:   job.RunnerName,
				StartedAt:    job.StartedAt,
				WorkflowName: run.Name,
				Repository:   run.Repository.FullName,
				HtmlUrl:      job.HtmlUrl,
			})
		}
	}

	return result, nil
}

// getJSON sends a GET request and decodes the JSON response into v
func (j *JobRepositoryImpl) getJSON(ctx context.Context, path string, v any) error {
	response, err := j.restClient.RequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if err := json.NewDecoder(response.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

// isActiveJob checks if a job status is considered active (in_progress or queued)
func (j *JobRepositoryImpl) isActiveJob(status string) bool {
	return status == "in_progress" || status == "queued"
}

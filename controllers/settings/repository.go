package settings

import (
	"context"
	"fmt"
	"strings"

	settingsv1beta1 "github.com/Hunter-Thompson/github-operator/apis/settings/v1beta1"
	"github.com/Hunter-Thompson/github-operator/pkg/gh"
	"github.com/go-logr/logr"
	"github.com/google/go-github/v52/github"
)

func (r *RepositoryReconciler) RepoExists(ctx context.Context, repo *settingsv1beta1.Repository, reqLogger logr.Logger) (bool, error) {
	ghClient := gh.Login(ctx)

	_, _, err := ghClient.Repositories.Get(ctx, repo.Spec.Organization, repo.Name)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found []") {
			return false, nil
		}
		return false, fmt.Errorf("failed to get repo: %w", err)
	}

	return true, nil
}

func (r *RepositoryReconciler) CreateRepo(ctx context.Context, repo *settingsv1beta1.Repository, reqLogger logr.Logger) error {
	ghClient := gh.Login(ctx)

	create := &github.Repository{
		Name:                &repo.Name,
		Description:         repo.GetDescription(),
		Homepage:            repo.GetHomepage(),
		Topics:              repo.GetTopics(),
		HasIssues:           repo.GetHasIssues(),
		HasPages:            repo.GetHasPages(),
		HasProjects:         repo.GetHasProjects(),
		HasWiki:             repo.GetHasWiki(),
		Private:             repo.GetPrivate(),
		AllowSquashMerge:    repo.GetAllowSquashMerge(),
		AllowMergeCommit:    repo.GetAllowMergeCommit(),
		AllowRebaseMerge:    repo.GetAllowRebaseMerge(),
		AllowAutoMerge:      repo.GetAllowAutoMerge(),
		DeleteBranchOnMerge: repo.GetDeleteBranchOnMerge(),
		IsTemplate:          repo.GetIsTemplate(),
	}

	if repo.Spec.SquashMergeCommitTitle != nil {
		create.SquashMergeCommitTitle = repo.Spec.SquashMergeCommitTitle
	}

	if repo.Spec.SquashMergeCommitMessage != nil {
		create.SquashMergeCommitMessage = repo.Spec.SquashMergeCommitMessage
	}

	if repo.Spec.MergeCommitTitle != nil {
		create.MergeCommitTitle = repo.Spec.MergeCommitTitle
	}

	if repo.Spec.MergeCommitMessage != nil {
		create.MergeCommitMessage = repo.Spec.MergeCommitMessage
	}

	_, _, err := ghClient.Repositories.Create(ctx, repo.Spec.Organization, create)
	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}

	reqLogger.Info("created repository")
	return nil
}

func (r *RepositoryReconciler) EditRepoSettings(ctx context.Context, repo *settingsv1beta1.Repository, reqLogger logr.Logger) error {
	ghClient := gh.Login(ctx)

	edit := &github.Repository{
		Name:                &repo.Name,
		Description:         repo.GetDescription(),
		Homepage:            repo.GetHomepage(),
		Topics:              repo.GetTopics(),
		HasIssues:           repo.GetHasIssues(),
		HasPages:            repo.GetHasPages(),
		HasProjects:         repo.GetHasProjects(),
		HasWiki:             repo.GetHasWiki(),
		Private:             repo.GetPrivate(),
		AllowSquashMerge:    repo.GetAllowSquashMerge(),
		AllowMergeCommit:    repo.GetAllowMergeCommit(),
		AllowRebaseMerge:    repo.GetAllowRebaseMerge(),
		AllowAutoMerge:      repo.GetAllowAutoMerge(),
		DeleteBranchOnMerge: repo.GetDeleteBranchOnMerge(),
		IsTemplate:          repo.GetIsTemplate(),
	}

	if repo.Spec.SquashMergeCommitTitle != nil {
		edit.SquashMergeCommitTitle = repo.Spec.SquashMergeCommitTitle
	}

	if repo.Spec.SquashMergeCommitMessage != nil {
		edit.SquashMergeCommitMessage = repo.Spec.SquashMergeCommitMessage
	}

	if repo.Spec.MergeCommitTitle != nil {
		edit.MergeCommitTitle = repo.Spec.MergeCommitTitle
	}

	if repo.Spec.MergeCommitMessage != nil {
		edit.MergeCommitMessage = repo.Spec.MergeCommitMessage
	}

	_, _, err := ghClient.Repositories.Edit(ctx, repo.Spec.Organization, repo.GetName(), edit)
	if err != nil {
		return fmt.Errorf("failed to edit repo: %w", err)
	}

	reqLogger.Info("edited repository")
	return nil
}

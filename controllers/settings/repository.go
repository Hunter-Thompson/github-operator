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

func (r *RepositoryReconciler) EditRepoCollaboraters(ctx context.Context, repo *settingsv1beta1.Repository, reqLogger logr.Logger) error {
	ghClient := gh.Login(ctx)

	if repo.Spec.RepositoryCollaborators == nil {
		return nil
	}

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allMembers []*github.User
	for {
		users, resp, err := ghClient.Organizations.ListMembers(ctx, repo.Spec.Organization, opt)
		if err != nil {
			return err
		}
		allMembers = append(allMembers, users...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, adminUser := range repo.Spec.RepositoryCollaborators.AdminPermission {
		err := addAdminPerm(ctx, repo, ghClient, adminUser, allMembers, reqLogger)
		if err != nil {
			return fmt.Errorf("failed to addAdminPerm: %w", err)
		}
	}

	for _, pullUser := range repo.Spec.RepositoryCollaborators.PullPermission {
		err := addPullPerm(ctx, repo, ghClient, pullUser, allMembers, reqLogger)
		if err != nil {
			return fmt.Errorf("failed to addPullPerm: %w", err)
		}
	}

	for _, pushUser := range repo.Spec.RepositoryCollaborators.PushPermission {
		err := addPushPerm(ctx, repo, ghClient, pushUser, allMembers, reqLogger)
		if err != nil {
			return fmt.Errorf("failed to addPushPerm: %w", err)
		}
	}

	return nil
}

func addAdminPerm(ctx context.Context, repo *settingsv1beta1.Repository, ghClient *github.Client, adminUser string, allMembers []*github.User, reqLogger logr.Logger) error {
	permLevel, _, err := ghClient.Repositories.GetPermissionLevel(ctx, repo.Spec.Organization, repo.GetName(), adminUser)
	if err != nil {
		if strings.Contains(err.Error(), "is not a user []") {
			reqLogger.Error(err, "not a user")
			return nil
		}
		return fmt.Errorf("failed to get perm level for repo: %w", err)
	}

	if permLevel.GetPermission() == "admin" {
		reqLogger.Info(adminUser + " already has admin permission")
		return nil
	}

	add := false

	for _, member := range allMembers {
		// only add if user is actually a part of the org
		if member.GetLogin() == adminUser {
			add = true
			break
		}
	}

	if add {
		_, _, err = ghClient.Repositories.AddCollaborator(ctx, repo.Spec.Organization, repo.GetName(), adminUser, &github.RepositoryAddCollaboratorOptions{
			Permission: "admin",
		})
		if err != nil {
			return fmt.Errorf("failed to add admin collaborator for repo: %w", err)
		}
		reqLogger.Info("gave " + adminUser + " admin permission")
	} else {
		reqLogger.Info(adminUser + " is not a part of the " + repo.Spec.Organization + " org")
	}

	return nil
}

func addPullPerm(ctx context.Context, repo *settingsv1beta1.Repository, ghClient *github.Client, pullUser string, allMembers []*github.User, reqLogger logr.Logger) error {
	permLevel, _, err := ghClient.Repositories.GetPermissionLevel(ctx, repo.Spec.Organization, repo.GetName(), pullUser)
	if err != nil {
		if strings.Contains(err.Error(), "is not a user []") {
			reqLogger.Error(err, "not a user")
			return nil
		}
		return fmt.Errorf("failed to get perm level for repo: %w", err)
	}

	if permLevel.GetPermission() == "read" || permLevel.GetPermission() == "admin" || permLevel.GetPermission() == "write" {
		reqLogger.Info(pullUser + " already has read permission")
		return nil
	}

	add := false

	for _, member := range allMembers {
		// only add if user is actually a part of the org
		if member.GetLogin() == pullUser {
			add = true
			break
		}
	}

	if add {
		_, _, err = ghClient.Repositories.AddCollaborator(ctx, repo.Spec.Organization, repo.GetName(), pullUser, &github.RepositoryAddCollaboratorOptions{
			Permission: "pull",
		})
		if err != nil {
			return fmt.Errorf("failed to add pull collaborator for repo: %w", err)
		}
		reqLogger.Info("gave " + pullUser + " pull permission")
	} else {
		reqLogger.Info(pullUser + " is not a part of the " + repo.Spec.Organization + " org")
	}
	return nil
}

func addPushPerm(ctx context.Context, repo *settingsv1beta1.Repository, ghClient *github.Client, pushUser string, allMembers []*github.User, reqLogger logr.Logger) error {
	permLevel, _, err := ghClient.Repositories.GetPermissionLevel(ctx, repo.Spec.Organization, repo.GetName(), pushUser)
	if err != nil {
		if strings.Contains(err.Error(), "is not a user []") {
			reqLogger.Error(err, "not a user")
			return nil
		}
		return fmt.Errorf("failed to get perm level for repo: %w", err)
	}

	if permLevel.GetPermission() == "admin" || permLevel.GetPermission() == "write" {
		reqLogger.Info(pushUser + " already has push permission")
		return nil
	}

	add := false

	for _, member := range allMembers {
		// only add if user is actually a part of the org
		if member.GetLogin() == pushUser {
			add = true
			break
		}
	}

	if add {
		_, _, err = ghClient.Repositories.AddCollaborator(ctx, repo.Spec.Organization, repo.GetName(), pushUser, &github.RepositoryAddCollaboratorOptions{
			Permission: "push",
		})
		if err != nil {
			return fmt.Errorf("failed to add push collaborator for repo: %w", err)
		}
		reqLogger.Info("gave " + pushUser + " push permission")
	} else {
		reqLogger.Info(pushUser + " is not a part of the " + repo.Spec.Organization + " org")
	}

	return nil
}

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

func (r *GlobalRepositoryReconciler) EditRepoSettings(ctx context.Context, gr *settingsv1beta1.GlobalRepository, repoName string, reqLogger logr.Logger) error {
	ghClient := gh.Login(ctx)

	edit := &github.Repository{
		Name:                &repoName,
		Description:         gr.GetDescription(),
		Homepage:            gr.GetHomepage(),
		Topics:              gr.GetTopics(),
		HasIssues:           gr.GetHasIssues(),
		HasPages:            gr.GetHasPages(),
		HasProjects:         gr.GetHasProjects(),
		HasWiki:             gr.GetHasWiki(),
		Private:             gr.GetPrivate(),
		AllowSquashMerge:    gr.GetAllowSquashMerge(),
		AllowMergeCommit:    gr.GetAllowMergeCommit(),
		AllowRebaseMerge:    gr.GetAllowRebaseMerge(),
		AllowAutoMerge:      gr.GetAllowAutoMerge(),
		DeleteBranchOnMerge: gr.GetDeleteBranchOnMerge(),
		IsTemplate:          gr.GetIsTemplate(),
	}

	if gr.Spec.SquashMergeCommitTitle != nil {
		edit.SquashMergeCommitTitle = gr.Spec.SquashMergeCommitTitle
	}

	if gr.Spec.SquashMergeCommitMessage != nil {
		edit.SquashMergeCommitMessage = gr.Spec.SquashMergeCommitMessage
	}

	if gr.Spec.MergeCommitTitle != nil {
		edit.MergeCommitTitle = gr.Spec.MergeCommitTitle
	}

	if gr.Spec.MergeCommitMessage != nil {
		edit.MergeCommitMessage = gr.Spec.MergeCommitMessage
	}

	_, _, err := ghClient.Repositories.Edit(ctx, gr.Spec.Organization, repoName, edit)
	if err != nil {
		return fmt.Errorf("failed to edit repo: %w", err)
	}

	reqLogger.Info("edited repository")
	return nil
}

//     pull - team members can pull, but not push to or administer this repository
//     push - team members can pull and push, but not administer this repository
//     admin - team members can pull, push and administer this repository
//     maintain - team members can manage the repository without access to sensitive or destructive actions.
//     triage - team members can proactively manage issues and pull requests without write access.
func (r *GlobalRepositoryReconciler) EditRepoTeams(ctx context.Context, gr *settingsv1beta1.GlobalRepository, repoName string, reqLogger logr.Logger) error {
	if gr.Spec.RepositoryTeams == nil {
		return nil
	}

	ghClient := gh.Login(ctx)
	opt := &github.ListOptions{
		PerPage: 10,
	}

	allTeams := []*github.Team{}
	for {
		repoTeam, resp, err := ghClient.Repositories.ListTeams(ctx, gr.Spec.Organization, repoName, opt)
		if err != nil {
			return err
		}
		allTeams = append(allTeams, repoTeam...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, adminTeam := range gr.Spec.RepositoryTeams.AdminPermission {
		err := addTeamToRepoPerm(ctx, gr, repoName, ghClient, adminTeam, allTeams, "admin", reqLogger)
		if err != nil {
			return fmt.Errorf("failed to add admin perm for %s to %s", adminTeam, repoName)
		}
	}

	for _, maintainTeam := range gr.Spec.RepositoryTeams.MaintainPermission {
		err := addTeamToRepoPerm(ctx, gr, repoName, ghClient, maintainTeam, allTeams, "maintain", reqLogger)
		if err != nil {
			return fmt.Errorf("failed to add maintain perm for %s to %s", maintainTeam, repoName)
		}
	}

	for _, pushTeam := range gr.Spec.RepositoryTeams.PushPermission {
		err := addTeamToRepoPerm(ctx, gr, repoName, ghClient, pushTeam, allTeams, "push", reqLogger)
		if err != nil {
			return fmt.Errorf("failed to add push perm for %s to %s", pushTeam, repoName)
		}
	}

	for _, pullTeam := range gr.Spec.RepositoryTeams.PullPermission {
		err := addTeamToRepoPerm(ctx, gr, repoName, ghClient, pullTeam, allTeams, "pull", reqLogger)
		if err != nil {
			return fmt.Errorf("failed to add pull perm for %s to %s", pullTeam, repoName)
		}
	}

	for _, triageTeam := range gr.Spec.RepositoryTeams.TriagePermission {
		err := addTeamToRepoPerm(ctx, gr, repoName, ghClient, triageTeam, allTeams, "triage", reqLogger)
		if err != nil {
			return fmt.Errorf("failed to add triage perm for %s to %s", triageTeam, repoName)
		}
	}

	return nil
}

func addTeamToRepoPerm(ctx context.Context, gr *settingsv1beta1.GlobalRepository, repoName string, ghClient *github.Client, team string, allTeams []*github.Team, perm string, reqLogger logr.Logger) error {
	for _, repoTeam := range allTeams {
		if repoTeam.GetName() == team {
			if checkTeamPerm(repoTeam, perm) {
				l := fmt.Sprintf("%s already has %s perm", team, perm)
				reqLogger.Info(l)
				return nil
			}
		}
	}

	_, err := ghClient.Teams.AddTeamRepoBySlug(ctx, gr.Spec.Organization, team, gr.Spec.Organization, repoName, &github.TeamAddTeamRepoOptions{
		Permission: perm,
	})
	if err != nil {
		return fmt.Errorf("failed to add perm for team for repo: %w", err)
	}
	reqLogger.Info("gave " + team + " " + perm + " permission")
	return nil
}

func checkTeamPerm(team *github.Team, perm string) bool {
	return team.GetPermission() == perm
}

func (r *GlobalRepositoryReconciler) EditRepoCollaboraters(ctx context.Context, gr *settingsv1beta1.GlobalRepository, repoName string, reqLogger logr.Logger) error {
	ghClient := gh.Login(ctx)

	if gr.Spec.RepositoryCollaborators == nil {
		return nil
	}

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allMembers []*github.User
	for {
		users, resp, err := ghClient.Organizations.ListMembers(ctx, gr.Spec.Organization, opt)
		if err != nil {
			return err
		}
		allMembers = append(allMembers, users...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, adminUser := range gr.Spec.RepositoryCollaborators.AdminPermission {
		err := addGlobalAdminPerm(ctx, gr, repoName, ghClient, adminUser, allMembers, reqLogger)
		if err != nil {
			return fmt.Errorf("failed to addAdminPerm: %w", err)
		}
	}

	for _, pullUser := range gr.Spec.RepositoryCollaborators.PullPermission {
		err := addGlobalPullPerm(ctx, gr, repoName, ghClient, pullUser, allMembers, reqLogger)
		if err != nil {
			return fmt.Errorf("failed to addPullPerm: %w", err)
		}
	}

	for _, pushUser := range gr.Spec.RepositoryCollaborators.PushPermission {
		err := addGlobalPushPerm(ctx, gr, repoName, ghClient, pushUser, allMembers, reqLogger)
		if err != nil {
			return fmt.Errorf("failed to addPushPerm: %w", err)
		}
	}

	return nil
}

func addGlobalAdminPerm(ctx context.Context, gr *settingsv1beta1.GlobalRepository, repoName string, ghClient *github.Client, adminUser string, allMembers []*github.User, reqLogger logr.Logger) error {
	permLevel, _, err := ghClient.Repositories.GetPermissionLevel(ctx, gr.Spec.Organization, repoName, adminUser)
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
		_, _, err = ghClient.Repositories.AddCollaborator(ctx, gr.Spec.Organization, repoName, adminUser, &github.RepositoryAddCollaboratorOptions{
			Permission: "admin",
		})
		if err != nil {
			return fmt.Errorf("failed to add admin collaborator for repo: %w", err)
		}
		reqLogger.Info("gave " + adminUser + " admin permission")
	} else {
		reqLogger.Info(adminUser + " is not a part of the " + gr.Spec.Organization + " org")
	}

	return nil
}

func addGlobalPullPerm(ctx context.Context, gr *settingsv1beta1.GlobalRepository, repoName string, ghClient *github.Client, pullUser string, allMembers []*github.User, reqLogger logr.Logger) error {
	permLevel, _, err := ghClient.Repositories.GetPermissionLevel(ctx, gr.Spec.Organization, repoName, pullUser)
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
		_, _, err = ghClient.Repositories.AddCollaborator(ctx, gr.Spec.Organization, repoName, pullUser, &github.RepositoryAddCollaboratorOptions{
			Permission: "pull",
		})
		if err != nil {
			return fmt.Errorf("failed to add pull collaborator for repo: %w", err)
		}
		reqLogger.Info("gave " + pullUser + " pull permission")
	} else {
		reqLogger.Info(pullUser + " is not a part of the " + gr.Spec.Organization + " org")
	}
	return nil
}

func addGlobalPushPerm(ctx context.Context, gr *settingsv1beta1.GlobalRepository, repoName string, ghClient *github.Client, pushUser string, allMembers []*github.User, reqLogger logr.Logger) error {
	permLevel, _, err := ghClient.Repositories.GetPermissionLevel(ctx, gr.Spec.Organization, repoName, pushUser)
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
		_, _, err = ghClient.Repositories.AddCollaborator(ctx, gr.Spec.Organization, repoName, pushUser, &github.RepositoryAddCollaboratorOptions{
			Permission: "push",
		})
		if err != nil {
			return fmt.Errorf("failed to add push collaborator for repo: %w", err)
		}
		reqLogger.Info("gave " + pushUser + " push permission")
	} else {
		reqLogger.Info(pushUser + " is not a part of the " + gr.Spec.Organization + " org")
	}

	return nil
}

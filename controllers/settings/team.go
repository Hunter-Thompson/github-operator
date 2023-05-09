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

func (r *TeamReconciler) TeamExists(ctx context.Context, team *settingsv1beta1.Team, reqLogger logr.Logger) (bool, error) {
	ghClient := gh.Login(ctx)

	_, _, err := ghClient.Teams.GetTeamBySlug(ctx, team.Spec.Organization, team.Name)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found []") {
			return false, nil
		}
		return false, fmt.Errorf("failed to get team: %w", err)
	}

	return true, nil
}

func (r *TeamReconciler) CreateTeam(ctx context.Context, team *settingsv1beta1.Team, reqLogger logr.Logger) error {
	ghClient := gh.Login(ctx)

	create := github.NewTeam{
		Name:        team.GetName(),
		Description: team.GetDescription(),
		Maintainers: team.GetMaintainers(),
		RepoNames:   team.GetRepoNames(),
	}

	if team.Spec.Privacy != nil {
		create.Privacy = team.Spec.Privacy
	}

	if team.Spec.ParentTeamID != nil {
		create.ParentTeamID = team.Spec.ParentTeamID
	}

	_, _, err := ghClient.Teams.CreateTeam(ctx, team.Spec.Organization, create)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	reqLogger.Info("created team")
	return nil
}

func (r *TeamReconciler) EditTeam(ctx context.Context, team *settingsv1beta1.Team, reqLogger logr.Logger) error {
	ghClient := gh.Login(ctx)

	edit := github.NewTeam{
		Name:        team.GetName(),
		Description: team.GetDescription(),
		RepoNames:   team.GetRepoNames(),
	}

	if team.Spec.Privacy != nil {
		edit.Privacy = team.Spec.Privacy
	}

	if team.Spec.ParentTeamID != nil {
		edit.ParentTeamID = team.Spec.ParentTeamID
	}

	_, _, err := ghClient.Teams.EditTeamBySlug(ctx, team.Spec.Organization, team.GetName(), edit, false)
	if err != nil {
		return fmt.Errorf("failed to edit team: %w", err)
	}

	reqLogger.Info("edited team")
	return nil
}

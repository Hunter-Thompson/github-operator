package settings

import (
	"context"
	"fmt"

	"github.com/Hunter-Thompson/github-operator/pkg/gh"
	"github.com/go-logr/logr"
	"github.com/google/go-github/v52/github"

	settingsv1beta1 "github.com/Hunter-Thompson/github-operator/apis/settings/v1beta1"
)

func (r *InviteUserReconciler) UserExists(ctx context.Context, iv *settingsv1beta1.InviteUser, userName string, reqLogger logr.Logger) (bool, error) {
	ghClient := gh.Login(ctx)

	member, resp, err := ghClient.Organizations.GetOrgMembership(ctx, userName, iv.Spec.Organization)
	if err != nil {
		if resp.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}

	if member.GetState() == "active" {
		return true, nil
	}

	return false, nil
}

func (r *InviteUserReconciler) UserInviteExists(ctx context.Context, iv *settingsv1beta1.InviteUser, userName string, userEmail string, reqLogger logr.Logger) (bool, error) {
	ghClient := gh.Login(ctx)

	var allPendingOrgInv []*github.Invitation
	opt := &github.ListOptions{
		PerPage: 10,
	}

	for {
		inv, resp, err := ghClient.Organizations.ListPendingOrgInvitations(ctx, iv.Spec.Organization, &github.ListOptions{})
		if err != nil {
			return false, err
		}
		allPendingOrgInv = append(allPendingOrgInv, inv...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, inv := range allPendingOrgInv {
		if inv.GetEmail() == userEmail {
			return true, nil
		}
	}

	l := fmt.Sprintf("%s | %s is not yet invited to org", userName, userEmail)
	reqLogger.Info(l)

	return false, nil
}

func (r *InviteUserReconciler) InviteUser(ctx context.Context, iv *settingsv1beta1.InviteUser, userEmail string, reqLogger logr.Logger) error {
	ghClient := gh.Login(ctx)

	_, _, err := ghClient.Organizations.CreateOrgInvitation(ctx, iv.Spec.Organization, &github.CreateOrgInvitationOptions{
		Email:  github.String(userEmail),
		Role:   github.String("direct_member"),
		TeamID: []int64{},
	})
	if err != nil {
		return err
	}

	return nil
}

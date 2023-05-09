/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package settings

import (
	"context"
	"time"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/Hunter-Thompson/github-operator/apis/settings/v1beta1"
	settingsv1beta1 "github.com/Hunter-Thompson/github-operator/apis/settings/v1beta1"
)

// TeamReconciler reconciles a Team object
type TeamReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=settings.github.com,resources=teams,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=settings.github.com,resources=teams/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=settings.github.com,resources=teams/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Team object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *TeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)

	team := &v1beta1.Team{}
	err := r.Client.Get(ctx, req.NamespacedName, team)
	if err != nil && k8sErrors.IsNotFound(err) {
		// Request object not found, could have been deleted after reconcile
		// request. Owned objects are automatically garbage collected.
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	exists, err := r.TeamExists(ctx, team, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	if !exists {
		err := r.CreateTeam(ctx, team, reqLogger)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	if exists {
		err := r.EditTeam(ctx, team, reqLogger)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&settingsv1beta1.Team{}).
		WithOptions(controller.Options{
			RateLimiter: workqueue.NewItemExponentialFailureRateLimiter(time.Minute*1, time.Minute*5),
		}).
		Complete(r)
}

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

package controllers

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	globalsv1beta2 "github.com/jnnkrdb/configrdb/api/v1beta2"
)

// GlobalSecretReconciler reconciles a GlobalSecret object
type GlobalSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=globals.jnnkrdb.de,resources=globalsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=globals.jnnkrdb.de,resources=globalsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=globals.jnnkrdb.de,resources=globalsecrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GlobalSecret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *GlobalSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var _log = log.FromContext(ctx).WithName("GlobalSecret")

	gs := &globalsv1beta2.GlobalSecret{}

	// parse the ctrl.Request into a globalsecret
	if err := r.Get(ctx, req.NamespacedName, gs); err != nil {

		// if the error is an "NotFound" error, then the globalsecret probably was deleted
		// returning no error
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		// if the error is something else, print the globalsecret and the error
		_log.Error(err, "error reconciling globalsecret", *gs)
		return ctrl.Result{}, err
	}

	_log.Info("reconciling globalsecret", fmt.Sprintf("[%s/%s]", gs.Namespace, gs.Name))

	// ---------------------------------------------------------------------------------------- remove all secrets, if the globalsecret is marked to be deleted
	// check, if the globalsecret is marked to be deleted
	if gs.GetDeletionTimestamp() != nil {

		// check, wether the globalconfig has the required finalizer or not
		if controllerutil.ContainsFinalizer(gs, globalsv1beta2.FinalizerGlobal) {

			_log.Info("finalizing globalsecret", fmt.Sprintf("%s/%s", gs.Namespace, gs.Name))

			// receiving a list of secrets, which are connected to this specific globalsecret
			var secretList = &v1.SecretList{}
			if err := r.List(ctx, secretList, client.MatchingLabels{"createdByConfRDB": globalsv1beta2.GroupVersion.Version, "globalsecretuid": string(gs.UID)}); err != nil {
				_log.Error(err, "error receiving list of secrets", client.MatchingLabels{"createdByConfRDB": globalsv1beta2.GroupVersion.Version, "globalsecretuid": string(gs.UID)})
				return ctrl.Result{Requeue: true}, err
			}

			// removing all the secrets in the list
			for _, scrt := range secretList.Items {
				_log.Info("removing secret", fmt.Sprintf("Secret[%s/%s]", scrt.Namespace, scrt.Name))
				if err := r.Delete(ctx, &scrt, &client.DeleteOptions{}); err != nil {
					_log.Error(err, "error removing secret", fmt.Sprintf("Secret[%s/%s]", scrt.Namespace, scrt.Name))
					return ctrl.Result{Requeue: true}, err
				}
			}

			_log.Info("finished finalizing globalsecrets", fmt.Sprintf("%s/%s", gs.Namespace, gs.Name))

			// remove the finalizer from the globalconfig
			controllerutil.RemoveFinalizer(gs, globalsv1beta2.FinalizerGlobal)
			if err := r.Update(ctx, gs); err != nil {
				return ctrl.Result{Requeue: true}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// check, wether the globalconfig has the required finalizer or not
	// if not, then add the finalizer
	if controllerutil.ContainsFinalizer(gs, globalsv1beta2.FinalizerGlobal) {

		// add the desired finalizer and update the object
		controllerutil.AddFinalizer(gs, globalsv1beta2.FinalizerGlobal)
		if err := r.Update(ctx, gs); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
	}
	// ---------------------------------------------------------------------------------------- start processing the globalsecret

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GlobalSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&globalsv1beta2.GlobalSecret{}).
		Complete(r)
}

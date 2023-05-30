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
	"encoding/base64"
	"fmt"
	"reflect"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch

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
	var _log = log.FromContext(ctx).WithName(fmt.Sprintf("GlobalSecret [%s]", req.NamespacedName))
	_log.Info("start reconciling")

	// ---------------------------------------------------------------------------------------- get the current globalconfig from the reconcile request
	// create caching object
	gs := &globalsv1beta2.GlobalSecret{}

	// parse the ctrl.Request into a globalsecret
	if err := r.Get(ctx, req.NamespacedName, gs); err != nil {

		_log.Error(err, "error reconciling globalsecret")

		// if the error is an "NotFound" error, then the globalsecret probably was deleted
		// returning no error
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		// if the error is something else, return the error
		return ctrl.Result{}, err
	}

	// ---------------------------------------------------------------------------------------- add neccessary finalizer, if not added
	// check, wether the globalsecret has the required finalizer or not
	// if not, then add the finalizer
	if controllerutil.ContainsFinalizer(gs, globalsv1beta2.FinalizerGlobal) {
		_log.Info("appending finalizer")

		// add the desired finalizer and update the object
		controllerutil.AddFinalizer(gs, globalsv1beta2.FinalizerGlobal)
		if err := r.Update(ctx, gs); err != nil {
			_log.Error(err, "error adding finalizer")
			return ctrl.Result{Requeue: true}, err
		}
	}

	// ---------------------------------------------------------------------------------------- receiving a list of secrets, which are connected to this specific globalsecret
	// start the finalizing routine
	_log.Info("receiving a list of secrets, which are connected to this specific globalsecret")

	var secretList = &v1.SecretList{}
	if err := r.List(ctx, secretList, globalsv1beta2.MatchingLables(gs.UID)); err != nil {
		_log.Error(err, "error receiving list of secrets", globalsv1beta2.MatchingLables(gs.UID))
		return ctrl.Result{Requeue: true}, err
	}

	// ---------------------------------------------------------------------------------------- remove all secrets, if the globalsecret is marked to be deleted
	// check, if the globalsecret is marked to be deleted
	if gs.GetDeletionTimestamp() != nil {
		// check, wether the globalconfig has the required finalizer or not
		if controllerutil.ContainsFinalizer(gs, globalsv1beta2.FinalizerGlobal) {
			// start the finalizing routine
			_log.Info("finalizing globalsecret")

			// removing all the secrets in the list
			for _, scrt := range secretList.Items {

				_log.Info("removing secret", "Secret", fmt.Sprintf("[%s/%s]", scrt.Namespace, scrt.Name))
				if err := r.Delete(ctx, &scrt, &client.DeleteOptions{}); err != nil {

					_log.Error(err, "error removing secret", "Secret", fmt.Sprintf("[%s/%s]", scrt.Namespace, scrt.Name))
					return ctrl.Result{Requeue: true}, err
				}
			}

			_log.Info("finished finalizing globalsecrets")

			// remove the finalizer from the globalsecret
			controllerutil.RemoveFinalizer(gs, globalsv1beta2.FinalizerGlobal)
			if err := r.Update(ctx, gs); err != nil {
				return ctrl.Result{Requeue: true}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// ---------------------------------------------------------------------------------------- start processing the globalsecret
	_log.Info("calculating the namespaces")
	var matches, avoids []v1.Namespace
	var err error
	var scrt = &v1.Secret{}

	// calculate the neccessary namespaces
	if matches, avoids, err = gs.Spec.Namespaces.CalculateNamespaces(_log, ctx, r.Client); err != nil {
		_log.Error(err, "error calculating the namespaces")
		return ctrl.Result{Requeue: true}, err
	}

	// remove existing secrets from the avoids
	_log.Info("removing already existing secrets in namespaces to avoid")
	for i := range avoids {
		nsLog := _log.WithValues("current Secret", fmt.Sprintf("[%s/%s]", avoids[i].Name, gs.Name))

		if err = r.Get(ctx, types.NamespacedName{Namespace: avoids[i].Name, Name: gs.Name}, scrt, &client.GetOptions{}); err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			nsLog.Error(err, "error receiving secretdata")
			return ctrl.Result{Requeue: true}, err
		}
		if err = r.Delete(ctx, scrt, &client.DeleteOptions{}); err != nil {
			nsLog.Error(err, "error removing secret")
			return ctrl.Result{Requeue: true}, err
		}
	}

	// create or update the secrets from the matching namespaces
	for i := range matches {
		nsLog := _log.WithValues("current Secret", fmt.Sprintf("[%s/%s]", matches[i].Name, gs.Name))

		if err = r.Get(ctx, types.NamespacedName{Namespace: matches[i].Name, Name: gs.Name}, scrt, &client.GetOptions{}); err != nil && !errors.IsNotFound(err) {
			nsLog.Error(err, "error requesting secretdata")
			return ctrl.Result{Requeue: true}, err
		}

		// if the configmap does not exist, then create a new configmap
		if errors.IsNotFound(err) {

			nsLog.Info("creating configmap")
			// create the actual object
			scrt.Name = gs.Name
			scrt.Namespace = matches[i].Name
			// cm.Annotations = globalsv1beta2.Annotations(gs.ResourceVersion)
			data := make(map[string]string)
			for k, v := range gs.Spec.Data {
				if unenc, err := base64.StdEncoding.DecodeString(v); err != nil {
					nsLog.Error(err, "error converting base64 data into secret data bytes")
					return ctrl.Result{Requeue: true}, err
				} else {
					data[k] = string(unenc)
				}
			}
			scrt.StringData = data
			scrt.Type = v1.SecretType(gs.Spec.Type)
			scrt.Immutable = func() *bool { b := true; return &b }()
			scrt.Labels = globalsv1beta2.MatchingLables(gs.UID)
			if err = r.Create(ctx, scrt, &client.CreateOptions{}); err != nil {
				nsLog.Error(err, "error creating new secret")
				return ctrl.Result{Requeue: true}, err
			}
			continue
		}

		nsLog.Info("updating secret")
		// since all secrets where created with the immutable=true flag, we can not simple update them,
		// we have to delete the secret and then create the new secret
		if !reflect.DeepEqual(scrt.Data, gs.Spec.Data) {
			if err = r.Delete(ctx, scrt, &client.DeleteOptions{}); err != nil {
				nsLog.Error(err, "error creating new secret")
				return ctrl.Result{Requeue: true}, err
			}

			// provide data
			scrt.Name = gs.Name
			scrt.Namespace = matches[i].Name
			// cm.Annotations = globalsv1beta2.Annotations(gs.ResourceVersion)
			data := make(map[string]string)
			for k, v := range gs.Spec.Data {
				if unenc, err := base64.StdEncoding.DecodeString(v); err != nil {
					nsLog.Error(err, "error converting base64 data into secret data bytes")
					return ctrl.Result{Requeue: true}, err
				} else {
					data[k] = string(unenc)
				}
			}
			scrt.StringData = data
			scrt.Immutable = func() *bool { b := true; return &b }()
			scrt.Labels = globalsv1beta2.MatchingLables(gs.UID)

			// recreate the configmap
			if err = r.Create(ctx, scrt, &client.CreateOptions{}); err != nil {
				nsLog.Error(err, "error creating new secret")
				return ctrl.Result{Requeue: true}, err
			}
		}
	}

	return ctrl.Result{
		RequeueAfter: (3 * time.Minute),
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GlobalSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&globalsv1beta2.GlobalSecret{}).
		Complete(r)
}

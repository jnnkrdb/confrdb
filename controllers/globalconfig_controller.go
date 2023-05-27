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

	"github.com/go-logr/logr"
	globalsv1beta2 "github.com/jnnkrdb/configrdb/api/v1beta2"
)

// GlobalConfigReconciler reconciles a GlobalConfig object
type GlobalConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=globals.jnnkrdb.de,resources=globalconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=globals.jnnkrdb.de,resources=globalconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=globals.jnnkrdb.de,resources=globalconfigs/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GlobalConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *GlobalConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var _log = log.FromContext(ctx).WithName("GlobalConfig")

	// create caching object
	gc := &globalsv1beta2.GlobalConfig{}

	// parse the ctrl.Request into a globalconfig
	if err := r.Get(ctx, req.NamespacedName, gc); err != nil {

		// if the error is an "NotFound" error, then the globalconfig probably was deleted
		// returning no error
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		// if the error is something else, print the globalconfig and the error
		_log.Error(err, "error reconciling globalconfig", *gc)
		return ctrl.Result{}, err
	}

	_log.Info("reconciling globalconfig", fmt.Sprintf("[%s/%s]", gc.Namespace, gc.Name))

	// ---------------------------------------------------------------------------------------- remove all configmaps, if the globalconfig is marked to be deleted
	// check, if the globalconfig is marked to be deleted
	if gc.GetDeletionTimestamp() != nil {

		// check, wether the globalconfig has the required finalizer or not
		if controllerutil.ContainsFinalizer(gc, globalsv1beta2.FinalizerGlobal) {

			// start the finalizing routine
			if err := r.finalize(ctx, _log, gc); err != nil {
				return ctrl.Result{Requeue: true}, err
			}

			// remove the finalizer from the globalconfig
			controllerutil.RemoveFinalizer(gc, globalsv1beta2.FinalizerGlobal)
			if err := r.Update(ctx, gc); err != nil {
				return ctrl.Result{Requeue: true}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// check, wether the globalconfig has the required finalizer or not
	// if not, then add the finalizer
	if controllerutil.ContainsFinalizer(gc, globalsv1beta2.FinalizerGlobal) {

		// add the desired finalizer and update the object
		controllerutil.AddFinalizer(gc, globalsv1beta2.FinalizerGlobal)
		if err := r.Update(ctx, gc); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
	}

	// ---------------------------------------------------------------------------------------- check the existing and non existing configmaps

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GlobalConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&globalsv1beta2.GlobalConfig{}).
		Complete(r)
}

// remove all objects if the finalize function gets called
func (r *GlobalConfigReconciler) finalize(ctx context.Context, reqLogger logr.Logger, gc *globalsv1beta2.GlobalConfig) error {

	reqLogger.Info("finalizing globalconfig", fmt.Sprintf("%s/%s", gc.Namespace, gc.Name))

	// receiving a list of configmaps, which are connected to this specific globalconfig
	var configMapList = &v1.ConfigMapList{}
	if err := r.List(ctx, configMapList, client.MatchingLabels{"createdByConfRDB": globalsv1beta2.GroupVersion.Version, "globalconfiguid": string(gc.UID)}); err != nil {
		reqLogger.Error(err, "error receiving list of configmaps", fmt.Sprintf("label [%s]=%s", "createdByConfRDB", globalsv1beta2.GroupVersion.Version), fmt.Sprintf("label [%s]=%s", "globalconfiguuid", string(gc.UID)))
		return err
	}

	// removing all the configmaps in the list
	for _, cm := range configMapList.Items {
		reqLogger.Info("removing configmap", fmt.Sprintf("ConfigMap[%s/%s]", cm.Namespace, cm.Name))
		if err := r.Delete(ctx, &cm, &client.DeleteOptions{}); err != nil {
			reqLogger.Error(err, "error removing configmap", fmt.Sprintf("ConfigMap[%s/%s]", cm.Namespace, cm.Name))
			return err
		}
	}

	reqLogger.Info("finished finalizing globalconfigs", fmt.Sprintf("%s/%s", gc.Namespace, gc.Name))

	return nil
}

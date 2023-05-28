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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

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
	var _log = log.FromContext(ctx).WithName(fmt.Sprintf("GlobalConfig [%s]", req.NamespacedName))
	_log.Info("start reconciling")

	// ---------------------------------------------------------------------------------------- get the current globalconfig from the reconcile request
	// create caching object
	gc := &globalsv1beta2.GlobalConfig{}

	// parse the ctrl.Request into a globalconfig
	if err := r.Get(ctx, req.NamespacedName, gc); err != nil {

		_log.Error(err, "error reconciling globalconfig")

		// if the error is an "NotFound" error, then the globalconfig probably was deleted
		// returning no error
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		// if the error is something else, print the globalconfig and the error
		return ctrl.Result{}, err
	}

	// ---------------------------------------------------------------------------------------- receiving a list of configmaps, which are connected to this specific globalconfig
	// start the finalizing routine
	_log.Info("receiving a list of configmaps, which are connected to this specific globalconfig")

	var configMapList = &v1.ConfigMapList{}
	if err := r.List(ctx, configMapList, globalsv1beta2.MatchingLables(gc.UID)); err != nil {
		_log.Error(err, "error receiving list of configmaps", globalsv1beta2.MatchingLables(gc.UID))
		return ctrl.Result{Requeue: true}, err
	}

	// ---------------------------------------------------------------------------------------- remove all configmaps, if the globalconfig is marked to be deleted
	// check, if the globalconfig is marked to be deleted
	if gc.GetDeletionTimestamp() != nil {
		// check, wether the globalconfig has the required finalizer or not
		if controllerutil.ContainsFinalizer(gc, globalsv1beta2.FinalizerGlobal) {
			// start the finalizing routine
			_log.Info("finalizing globalconfig", fmt.Sprintf("%s/%s", gc.Namespace, gc.Name))

			// removing all the configmaps in the list
			for _, cm := range configMapList.Items {
				_log.Info("removing configmap")
				if err := r.Delete(ctx, &cm, &client.DeleteOptions{}); err != nil {
					_log.Error(err, "error removing configmap", fmt.Sprintf("ConfigMap[%s/%s]", cm.Namespace, cm.Name))
					return ctrl.Result{Requeue: true}, err
				}
			}

			_log.Info("finished finalizing globalconfig", fmt.Sprintf("%s/%s", gc.Namespace, gc.Name))

			// remove the finalizer from the globalconfig
			controllerutil.RemoveFinalizer(gc, globalsv1beta2.FinalizerGlobal)
			if err := r.Update(ctx, gc); err != nil {
				_log.Error(err, "error updating finalizer")
				return ctrl.Result{Requeue: true}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// ---------------------------------------------------------------------------------------- add neccessary finalizer, if not added
	// check, wether the globalconfig has the required finalizer or not
	// if not, then add the finalizer
	if controllerutil.ContainsFinalizer(gc, globalsv1beta2.FinalizerGlobal) {
		_log.Info("appending finalizer")

		// add the desired finalizer and update the object
		controllerutil.AddFinalizer(gc, globalsv1beta2.FinalizerGlobal)
		if err := r.Update(ctx, gc); err != nil {
			_log.Error(err, "error adding finalizer")
			return ctrl.Result{Requeue: true}, err
		}
	}

	// ---------------------------------------------------------------------------------------- start processing the globalconfig
	_log.Info("calculating the namespaces")
	var matches, avoids []v1.Namespace
	var err error
	var cm = &v1.ConfigMap{}

	if matches, avoids, err = gc.Spec.Namespaces.CalculateNamespaces(_log, ctx, r.Client); err != nil {
		_log.Error(err, "error calculating the namespaces")
		return ctrl.Result{Requeue: true}, err
	}

	// remove existing configmaps from the avoids
	_log.Info("removing already existing configmap in namespaces to avoid")
	for i := range avoids {
		nsLog := _log.WithValues(fmt.Sprintf("current ConfigMap[%s/%s]", avoids[i].Name, gc.Name))

		if err = r.Get(ctx, types.NamespacedName{Namespace: avoids[i].Name, Name: gc.Name}, cm, &client.GetOptions{}); err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			nsLog.Error(err, "error receiving configmapdata")
			return ctrl.Result{Requeue: true}, err
		}
		if err = r.Delete(ctx, cm, &client.DeleteOptions{}); err != nil {
			nsLog.Error(err, "error removing configmap")
			return ctrl.Result{Requeue: true}, err
		}
	}

	// create or update the configmaps from the matching namespaces
	for i := range matches {
		nsLog := _log.WithValues(fmt.Sprintf("current ConfigMap[%s/%s]", matches[i].Name, gc.Name))

		if err = r.Get(ctx, types.NamespacedName{Namespace: matches[i].Name, Name: gc.Name}, cm, &client.GetOptions{}); err != nil && errors.IsNotFound(err) {
			nsLog.Error(err, "error requesting configmapdata")
			return ctrl.Result{Requeue: true}, err
		}

		// if the configmap does not exist, then create a new configmap
		if errors.IsNotFound(err) {
			// create the actual object
			cm.Annotations = globalsv1beta2.Annotations(gc.ResourceVersion)
			cm.Data = gc.Spec.Data
			cm.Immutable = func() *bool { b := true; return &b }()
			cm.Labels = globalsv1beta2.MatchingLables(gc.UID)
			if err := r.Create(ctx, cm, &client.CreateOptions{}); err != nil {
				nsLog.Error(err, "error creating new configmap")
				return ctrl.Result{Requeue: true}, err
			}
			continue
		}

		// since all configmaps where created with the immutable=true flag, we can not simple update them,
		// we have to delete the configmap and then create the new configmap
		if value, ok := cm.Annotations[globalsv1beta2.RVAnnotation]; ok {
			nsLog.Error(fmt.Errorf("annotation [%s] does not exist on this object", globalsv1beta2.RVAnnotation), "error comparing the annotations")
			return ctrl.Result{Requeue: true}, err
		} else {
			if value != gc.ResourceVersion {
				if err = r.Delete(ctx, cm, &client.DeleteOptions{}); err != nil {
					nsLog.Error(err, "error creating new configmap")
					return ctrl.Result{Requeue: true}, err
				}

				// provide data
				cm.Name = gc.Name
				cm.Namespace = matches[i].Name
				cm.Annotations = globalsv1beta2.Annotations(gc.ResourceVersion)
				cm.Data = gc.Spec.Data
				cm.Immutable = func() *bool { b := true; return &b }()
				cm.Labels = globalsv1beta2.MatchingLables(gc.UID)

				// recreate the configmap
				if err = r.Create(ctx, cm, &client.CreateOptions{}); err != nil {
					nsLog.Error(err, "error creating new configmap")
					return ctrl.Result{Requeue: true}, err
				}
			}
		}
	}

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GlobalConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&globalsv1beta2.GlobalConfig{}).
		Complete(r)
}

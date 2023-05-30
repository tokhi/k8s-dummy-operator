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

	"github.com/pingcap/errors"
	dummyv1alpha1 "github.com/tokhi/k8s-dummy-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var logger = log.Log.WithName("controllers").WithName("Dummy")

// DummyReconciler reconciles a Dummy object
type DummyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dummy.interview.com,resources=dummies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dummy.interview.com,resources=dummies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dummy.interview.com,resources=dummies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Dummy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DummyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log := logger.WithValues("dummy", req.NamespacedName)

	log.Info("Reconciling Dummy", "name", req.NamespacedName)

	// fetch the Dummy instance
	instance := &dummyv1alpha1.Dummy{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	log.Info("Message:", "message", instance.Spec.Message)

	instance.Status.SpecEcho = instance.Spec.Message

	// status update
	err = r.Status().Update(ctx, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("SpecEcho:", "specEcho", instance.Status.SpecEcho)

	// Check if the Pod already exist
	pod := &corev1.Pod{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, pod)

	if err != nil && errors.IsNotFound(err) {
		pod := newDumyPodWithNginx(instance)
		log.Info("Creating Pod for Dummy", "name", instance.Name, "namespace", instance.Namespace)
		if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
			return reconcile.Result{}, err
		}

		instance.Status.PodStatus = "Pending"
		err = r.Status().Update(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Info("PodStatus:", "podStatus", instance.Status.PodStatus)

		// create the pod
		err = r.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		instance.Status.PodStatus = "Running"
		err = r.Status().Update(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Info("PodStatus:", "podStatus", instance.Status.PodStatus)

	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		log.Info("Pod is already running")
	}

	return ctrl.Result{}, nil
}

func newDumyPodWithNginx(dummy *dummyv1alpha1.Dummy) *corev1.Pod {
	labels := map[string]string{
		"app": dummy.Name,
	}
	containers := []corev1.Container{
		{
			Name:  "nginx",
			Image: "nginx",
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dummy.Name,
			Namespace: dummy.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: containers,
		},
	}
	return pod
}

func (r *DummyReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&dummyv1alpha1.Dummy{}).
		Complete(r)
}

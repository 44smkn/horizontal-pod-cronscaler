/*
Copyright 2021.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crlog "sigs.k8s.io/controller-runtime/pkg/log"

	autoscalingv1beta1 "github.com/44smkn/horizontal-pod-cronscaler/api/v1beta1"
	"github.com/44smkn/horizontal-pod-cronscaler/pkg/scaling"
)

// HorizontalPodCronscalerReconciler reconciles a HorizontalPodCronscaler object
type HorizontalPodCronscalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=autoscaling.44smkn.github.io,resources=horizontalpodcronscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autoscaling.44smkn.github.io,resources=horizontalpodcronscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=autoscaling.44smkn.github.io,resources=horizontalpodcronscalers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments;replicaset,verbs=get;list;watch;update;patch

func (r *HorizontalPodCronscalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := crlog.FromContext(ctx)

	hpc := &autoscalingv1beta1.HorizontalPodCronscaler{}
	if err := r.Get(ctx, req.NamespacedName, hpc); err != nil {
		log.Error(err, "unable to fetch HorizontalPodCronscaler")
		return ctrl.Result{}, err
	}

	if !hpc.DeletionTimestamp.IsZero() {
		scaling.RemoveSchedule(req.NamespacedName.String())
		log.Info("remove schedule")
		return ctrl.Result{}, nil
	}

	scaling.AddSchedule(ctx, hpc, r.Client, log)
	log.Info("added schedule", "schedule", hpc.Spec.Schedule)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HorizontalPodCronscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1beta1.HorizontalPodCronscaler{}).
		Complete(r)
}

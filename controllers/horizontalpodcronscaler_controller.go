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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	cronv3 "github.com/robfig/cron/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crlog "sigs.k8s.io/controller-runtime/pkg/log"

	autoscalingv1beta1 "github.com/44smkn/horizontal-pod-cronscaler/api/v1beta1"
)

type namespacedName = string

// HorizontalPodCronscalerReconciler reconciles a HorizontalPodCronscaler object
type HorizontalPodCronscalerReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	cron        *cronv3.Cron
	cronMapping map[namespacedName]cronv3.EntryID
	ch          chan *Schedule
}

func NewHorizontalPodCronscalerReconciler(client client.Client, scheme *runtime.Scheme) *HorizontalPodCronscalerReconciler {
	cron := cronv3.New()
	go cron.Start()
	return &HorizontalPodCronscalerReconciler{
		Client:      client,
		Scheme:      scheme,
		cron:        cron,
		cronMapping: make(map[string]cronv3.EntryID),
		ch:          make(chan *Schedule),
	}
}

//+kubebuilder:rbac:groups=autoscaling.44smkn.github.io,resources=horizontalpodcronscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autoscaling.44smkn.github.io,resources=horizontalpodcronscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=autoscaling.44smkn.github.io,resources=horizontalpodcronscalers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments;replicaset,verbs=get;list;watch;update;patch

func (r *HorizontalPodCronscalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := crlog.FromContext(ctx)

	log.Info("start reconcile")
	hpc := &autoscalingv1beta1.HorizontalPodCronscaler{}
	if err := r.Get(ctx, req.NamespacedName, hpc); err != nil {
		log.Error(err, "unable to fetch HorizontalPodCronscaler")
		return ctrl.Result{}, err
	}

	name := fmt.Sprintf("%s/%s", hpc.Namespace, hpc.Name)
	sched := &Schedule{
		name:           name,
		cronExpression: hpc.Spec.Schedule,
		shouldDelete:   !hpc.DeletionTimestamp.IsZero(),
		fn: func() {
			log.Info(fmt.Sprintf("scale %s replicas to %v", name, hpc.Spec.Replicas))

			patch := &unstructured.Unstructured{}
			gvk := schema.FromAPIVersionAndKind(hpc.Spec.ScaleTargetRef.APIVersion, hpc.Spec.ScaleTargetRef.Kind)
			patch.SetGroupVersionKind(gvk)
			patch.SetNamespace(hpc.Namespace)
			patch.SetName(hpc.Spec.ScaleTargetRef.Name)

			patch.UnstructuredContent()["spec"] = map[string]interface{}{
				"replicas": hpc.Spec.Replicas,
			}
			err := r.Client.Patch(ctx, patch, client.Apply, &client.PatchOptions{
				FieldManager: "horizontal-pod-cronscaler",
				Force:        pointer.BoolPtr(true),
			})
			if err != nil {
				log.Error(err, "failed to patch")
			}

			hpc.Status = autoscalingv1beta1.HorizontalPodCronscalerStatus{
				LastTargetReplicas: &hpc.Spec.Replicas,
				LastSchedule:       metav1.NewTime(time.Now()),
			}

			if err != nil {
				log.Error(err, "failed to addFunc")
			}
		},
	}
	log.Info("before send message to channel")
	r.ch <- sched
	log.Info("added schedule", "schedule", hpc.Spec.Schedule)

	return ctrl.Result{}, nil
}

type Schedule struct {
	name           namespacedName
	cronExpression string
	fn             func()
	shouldDelete   bool
}

func (r *HorizontalPodCronscalerReconciler) putSchedule(log logr.Logger) {
	for p := range r.ch {
		if id, ok := r.cronMapping[p.name]; ok {
			r.cron.Remove(id)
		}
		if !p.shouldDelete {
			id, err := r.cron.AddFunc(p.cronExpression, p.fn)
			if err != nil {
				log.Error(err, "failed to add func")
			}
			r.cronMapping[p.name] = id
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *HorizontalPodCronscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	go r.putSchedule(crlog.Log)
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1beta1.HorizontalPodCronscaler{}).
		Complete(r)
}

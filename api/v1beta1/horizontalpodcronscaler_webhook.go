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

package v1beta1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var horizontalpodcronscalerlog = logf.Log.WithName("horizontalpodcronscaler-resource")

func (r *HorizontalPodCronscaler) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-autoscaling-44smkn-github-io-v1beta1-horizontalpodcronscaler,mutating=true,failurePolicy=fail,sideEffects=None,groups=autoscaling.44smkn.github.io,resources=horizontalpodcronscalers,verbs=create;update,versions=v1beta1,name=mhorizontalpodcronscaler.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &HorizontalPodCronscaler{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *HorizontalPodCronscaler) Default() {
	horizontalpodcronscalerlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	r.Spec.Reference = fmt.Sprintf("%s/%s", r.Spec.ScaleTargetRef.Kind, r.Spec.ScaleTargetRef.Name)
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-autoscaling-44smkn-github-io-v1beta1-horizontalpodcronscaler,mutating=false,failurePolicy=fail,sideEffects=None,groups=autoscaling.44smkn.github.io,resources=horizontalpodcronscalers,verbs=create;update,versions=v1beta1,name=vhorizontalpodcronscaler.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &HorizontalPodCronscaler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *HorizontalPodCronscaler) ValidateCreate() error {
	horizontalpodcronscalerlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *HorizontalPodCronscaler) ValidateUpdate(old runtime.Object) error {
	horizontalpodcronscalerlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *HorizontalPodCronscaler) ValidateDelete() error {
	horizontalpodcronscalerlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

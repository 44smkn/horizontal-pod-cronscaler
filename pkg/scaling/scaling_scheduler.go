package scaling

import (
	"context"
	"fmt"
	"sync"
	"time"

	autoscalingv1beta1 "github.com/44smkn/horizontal-pod-cronscaler/api/v1beta1"
	"github.com/go-logr/logr"
	"github.com/robfig/cron/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NamespacedName = string

var (
	entries     map[NamespacedName]cron.EntryID = make(map[string]cron.EntryID)
	mu          sync.Mutex
	scalingCron *cron.Cron = cron.New()
	once        sync.Once
)

func AddSchedule(ctx context.Context, hpc *autoscalingv1beta1.HorizontalPodCronscaler, cli client.Client, log logr.Logger) {
	namespacedName := fmt.Sprintf("%s/%s", hpc.Namespace, hpc.Name)
	if id, ok := entries[namespacedName]; ok {
		scalingCron.Remove(id)
	}

	id, err := scalingCron.AddFunc(hpc.Spec.Schedule, func() {
		log.Info(fmt.Sprintf("scale %s replicas to %v", namespacedName, hpc.Spec.Replicas))

		patch := &unstructured.Unstructured{}
		patch.SetGroupVersionKind(hpc.GroupVersionKind())
		patch.SetNamespace(hpc.Namespace)
		patch.SetName(hpc.Name)

		patch.UnstructuredContent()["spec"] = map[string]interface{}{
			"replicas": hpc.Spec.Replicas,
		}
		err := cli.Patch(ctx, patch, client.Apply, &client.PatchOptions{
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
	})
	if err != nil {
		log.Error(err, "failed to addFunc")
	}
	mu.Lock()
	entries[namespacedName] = id
	mu.Unlock()

	once.Do(func() {
		log.Info("start cron process...")
		go scalingCron.Start()
	})

	log.Info(scalingCron.Location().String())
}

func RemoveSchedule(name NamespacedName) {
	mu.Lock()
	delete(entries, name)
	mu.Unlock()
}

func Entries() []cron.Entry {
	return scalingCron.Entries()
}

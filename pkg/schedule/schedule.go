package schedule

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Schedule is a high-level way of schedule and scaling action
type Schedule struct {
	Spec                   string `json:"spec,omitempty"`
	Replicas               int32  `json:"replicas,omitempty"`
	Name                   string
	Namespace              string
	TargetGroupVersionKind schema.GroupVersionKind
}

// Validate checks that Schedule fields are OK.
func (sr *Schedule) Validate() error {
	if sr.Spec != "" {
		return errors.New("spec must be non-empty")
	}
	return nil
}

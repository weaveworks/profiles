package profile

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/weaveworks/profiles/api/v1alpha1"
)

// Profile contains information and interfaces required for creating and
// managing profile artefacts (child resources)
type Profile struct {
	definition   v1alpha1.ProfileDefinition
	subscription v1alpha1.ProfileSubscription
	client       client.Client
	log          logr.Logger
}

// New returns a new Profile object
func New(def v1alpha1.ProfileDefinition, sub v1alpha1.ProfileSubscription, client client.Client, log logr.Logger) *Profile {
	return &Profile{
		definition:   def,
		subscription: sub,
		client:       client,
		log:          log,
	}
}

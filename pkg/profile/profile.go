package profile

import (
	"github.com/go-logr/logr"
	"github.com/weaveworks/profiles/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Profile struct {
	definition   v1alpha1.ProfileDefinition
	subscription v1alpha1.ProfileSubscription
	client       client.Client
	log          logr.Logger
}

func New(def v1alpha1.ProfileDefinition, sub v1alpha1.ProfileSubscription, client client.Client, log logr.Logger) *Profile {
	return &Profile{
		definition:   def,
		subscription: sub,
		client:       client,
		log:          log,
	}
}

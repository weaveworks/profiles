package profile

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
)

// Profile contains information and interfaces required for creating and
// managing profile artefacts (child resources)
type Profile struct {
	definition   profilesv1.ProfileDefinition
	subscription profilesv1.ProfileSubscription
	client       client.Client
	log          logr.Logger
	ctx          context.Context
}

// New returns a new Profile object
func New(ctx context.Context, def profilesv1.ProfileDefinition, sub profilesv1.ProfileSubscription, client client.Client, log logr.Logger) *Profile {
	return &Profile{
		definition:   def,
		subscription: sub,
		client:       client,
		log:          log,
		ctx:          ctx,
	}
}

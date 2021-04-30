package profile

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/git"
)

// Profile contains information and interfaces required for creating and
// managing profile artefacts (child resources)
type Profile struct {
	definition profilesv1.ProfileDefinition
	instance   profilesv1.ProfileInstance
	client     client.Client
	log        logr.Logger
	ctx        context.Context
}

// ProfileGetter is a func that can fetch a profile definition
type ProfileGetter func(repoURL, branch string, log logr.Logger) (profilesv1.ProfileDefinition, error)

var getProfileDefinition = git.GetProfileDefinition

// New returns a new Profile object
func New(ctx context.Context, def profilesv1.ProfileDefinition, sub profilesv1.ProfileInstance, client client.Client, log logr.Logger) *Profile {
	return &Profile{
		definition: def,
		instance:   sub,
		client:     client,
		log:        log,
		ctx:        ctx,
	}
}

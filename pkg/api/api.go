package api

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/protos"
)

//go:generate counterfeiter -o fakes/fake_catalog.go . Catalog
// Catalog is an interface for the Catalog
type Catalog interface {
	// Get will return a specific profile from the catalog
	Get(sourceName, profileName string) *profilesv1.ProfileCatalogEntry
	// GetWithVersion will return a specific profile from the catalog
	GetWithVersion(logger logr.Logger, sourceName, profileName, version string) *profilesv1.ProfileCatalogEntry
	// ProfilesGreaterThanVersion returns all profiles which are of a greater version for a given profile with a version.
	ProfilesGreaterThanVersion(logger logr.Logger, sourceName, profileName, version string) []profilesv1.ProfileCatalogEntry
	// Search will return a list of profiles which match query
	Search(query string) []profilesv1.ProfileCatalogEntry
	// Search will return a list of all profiles
	SearchAll() []profilesv1.ProfileCatalogEntry
}

type CatalogAPI interface {
	protos.ProfilesServiceServer
}

type ProfilesCatalogService struct {
	profileCatalog Catalog
	logger         logr.Logger
}

func (p *ProfilesCatalogService) Get(ctx context.Context, request *protos.GetRequest) (*protos.GetResponse, error) {
	sourceName := request.GetSourceName()
	profileName := request.GetProfileName()
	logger := p.logger.WithValues("func", "Get", "catalog", sourceName, "profile", profileName)
	if sourceName == "" || profileName == "" {
		errMsg := fmt.Errorf("missing query param: sourceName: %q, profileName: %q", sourceName, profileName)
		logger.Error(errMsg, "profile and/or catalog not set")
		return nil, status.Errorf(codes.InvalidArgument, errMsg.Error())
	}
	result := p.profileCatalog.Get(sourceName, profileName)
	if result == nil {
		return nil, status.Errorf(codes.NotFound, "profile not found")
	}
	logger.Info("profile found", "profile", result)
	return &protos.GetResponse{
		Item: protos.TransformCatalogEntry(result),
	}, nil
}

func (p *ProfilesCatalogService) GetWithVersion(ctx context.Context, request *protos.GetWithVersionRequest) (*protos.GetWithVersionResponse, error) {
	sourceName := request.GetSourceName()
	profileName := request.GetProfileName()
	version := request.GetVersion()
	logger := p.logger.WithValues("func", "GetWithVersion", "catalog", sourceName, "profile", profileName, "version", version)
	if sourceName == "" || profileName == "" || version == "" {
		errMsg := fmt.Errorf("missing query param: sourceName: %q, profileName: %q, version: %q", sourceName, profileName, version)
		logger.Error(errMsg, "catalog, profile and/or version not set")
		return nil, status.Errorf(codes.InvalidArgument, errMsg.Error())
	}
	result := p.profileCatalog.GetWithVersion(logger, sourceName, profileName, version)
	if result == nil {
		return nil, status.Errorf(codes.NotFound, "profile not found")
	}
	logger.Info("profile found", "profile", result)
	return &protos.GetWithVersionResponse{
		Item: protos.TransformCatalogEntry(result),
	}, nil
}

func (p *ProfilesCatalogService) ProfilesGreaterThanVersion(ctx context.Context, request *protos.ProfilesGreaterThanVersionRequest) (*protos.ProfilesGreaterThanVersionResponse, error) {
	sourceName := request.GetSourceName()
	profileName := request.GetProfileName()
	version := request.GetVersion()
	logger := p.logger.WithValues("func", "ProfilesGreaterThanVersion", "catalog", sourceName, "profile", profileName, "version", version)
	if sourceName == "" || profileName == "" || version == "" {
		errMsg := fmt.Errorf("missing query param: sourceName: %q, profileName: %q, version: %q", sourceName, profileName, version)
		logger.Error(errMsg, "catalog, profile and/or version not set")
		return nil, status.Errorf(codes.InvalidArgument, errMsg.Error())
	}
	result := p.profileCatalog.ProfilesGreaterThanVersion(logger, sourceName, profileName, version)
	if result == nil {
		return nil, status.Errorf(codes.NotFound, "profile not found")
	}
	logger.Info("profile found", "profile", result)
	return &protos.ProfilesGreaterThanVersionResponse{
		Items: protos.TransformCatalogEntryList(result),
	}, nil
}

func (p *ProfilesCatalogService) Search(ctx context.Context, request *protos.SearchRequest) (*protos.SearchResponse, error) {
	query := request.GetName()
	logger := p.logger.WithValues("func", "Search", "name", query)
	var result []profilesv1.ProfileCatalogEntry
	if query == "" {
		logger.Info("Searching for all available profiles")
		result = p.profileCatalog.SearchAll()
	} else {
		logger.Info("Searching for profiles matching name", "name", query)
		result = p.profileCatalog.Search(query)
	}

	logger.Info("found profiles", "profiles", result)
	return &protos.SearchResponse{
		Items: protos.TransformCatalogEntryList(result),
	}, nil
}

var _ protos.ProfilesServiceServer = &ProfilesCatalogService{}

// NewCatalog .
func NewCatalog(profileCatalog Catalog, logger logr.Logger) *ProfilesCatalogService {
	return &ProfilesCatalogService{
		profileCatalog: profileCatalog,
		logger:         logger,
	}
}

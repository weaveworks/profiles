package api

import (
	"context"

	"github.com/go-logr/logr"

	profiles "github.com/weaveworks/profiles/pkg/protos"
)

type CatalogGRPC interface {
	profiles.ProfilesServiceServer
}

type ProfilesCatalogService struct {
	profileCatalog Catalog
	logger         logr.Logger
}

func (p *ProfilesCatalogService) Get(ctx context.Context, request *profiles.GetRequest) (*profiles.GetResponse, error) {
	panic("implement me")
}

func (p *ProfilesCatalogService) GetWithVersion(ctx context.Context, request *profiles.GetWithVersionRequest) (*profiles.GetWithVersionResponse, error) {
	panic("implement me")
}

func (p *ProfilesCatalogService) ProfilesGreaterThanVersion(ctx context.Context, request *profiles.ProfilesGreaterThanVersionRequest) (*profiles.ProfilesGreaterThanVersionResponse, error) {
	panic("implement me")
}

func (p *ProfilesCatalogService) Search(ctx context.Context, request *profiles.SearchRequest) (*profiles.SearchResponse, error) {
	panic("implement me")
}

var _ profiles.ProfilesServiceServer = &ProfilesCatalogService{}

// NewCatalog .
func NewCatalog(profileCatalog Catalog, logger logr.Logger) *ProfilesCatalogService {
	return &ProfilesCatalogService{
		profileCatalog: profileCatalog,
		logger:         logger,
	}
}

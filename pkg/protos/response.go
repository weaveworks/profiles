package protos

import profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"

// GRPCProfileCatalogEntry defines a return type for the grpc-gateway based catalog entry item.
type GRPCProfileCatalogEntry struct {
	Item profilesv1.ProfileCatalogEntry `json:"item"`
}

// GRPCProfileCatalogEntryList defines a return type for the grpc-gateway based catalog entry list.
type GRPCProfileCatalogEntryList struct {
	Items []profilesv1.ProfileCatalogEntry `json:"items"`
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.1
// source: profiles.proto

package protos

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Empty parameters.
type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{0}
}

// GetRequest defines parameters for the Get endpoint.
type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the catalog
	SourceName string `protobuf:"bytes,1,opt,name=source_name,json=sourceName,proto3" json:"source_name,omitempty"`
	// Name of the profile
	ProfileName string `protobuf:"bytes,2,opt,name=profile_name,json=profileName,proto3" json:"profile_name,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{1}
}

func (x *GetRequest) GetSourceName() string {
	if x != nil {
		return x.SourceName
	}
	return ""
}

func (x *GetRequest) GetProfileName() string {
	if x != nil {
		return x.ProfileName
	}
	return ""
}

// GetResponse defines response parameters for Get endpoint.
type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Item *ProfileCatalogEntry `protobuf:"bytes,1,opt,name=item,proto3" json:"item,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{2}
}

func (x *GetResponse) GetItem() *ProfileCatalogEntry {
	if x != nil {
		return x.Item
	}
	return nil
}

// ProfileDescription defines details about a given profile.
type ProfileCatalogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Defines the branch or tag to use
	Tag string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	// Name of the catalog the profile is listed in
	CatalogSource string `protobuf:"bytes,2,opt,name=catalog_source,json=catalogSource,proto3" json:"catalog_source,omitempty"`
	// The full URL path to the profile.yaml
	Url string `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	// The fields below are inlined from ProfileDescription.
	// Name of the profile
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the profile
	Description string `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	// The maintainer of the profile
	Maintainer string `protobuf:"bytes,6,opt,name=maintainer,proto3" json:"maintainer,omitempty"`
	// Any prerequisites that should be met for this profile to be installable
	Prerequisites []string `protobuf:"bytes,7,rep,name=prerequisites,proto3" json:"prerequisites,omitempty"`
}

func (x *ProfileCatalogEntry) Reset() {
	*x = ProfileCatalogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileCatalogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileCatalogEntry) ProtoMessage() {}

func (x *ProfileCatalogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileCatalogEntry.ProtoReflect.Descriptor instead.
func (*ProfileCatalogEntry) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{3}
}

func (x *ProfileCatalogEntry) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *ProfileCatalogEntry) GetCatalogSource() string {
	if x != nil {
		return x.CatalogSource
	}
	return ""
}

func (x *ProfileCatalogEntry) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *ProfileCatalogEntry) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProfileCatalogEntry) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ProfileCatalogEntry) GetMaintainer() string {
	if x != nil {
		return x.Maintainer
	}
	return ""
}

func (x *ProfileCatalogEntry) GetPrerequisites() []string {
	if x != nil {
		return x.Prerequisites
	}
	return nil
}

// GetWithVersionRequest defines request parameters for GetWithVersion endpoint.
type GetWithVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the catalog
	SourceName string `protobuf:"bytes,1,opt,name=source_name,json=sourceName,proto3" json:"source_name,omitempty"`
	// Name of the profile
	ProfileName string `protobuf:"bytes,2,opt,name=profile_name,json=profileName,proto3" json:"profile_name,omitempty"`
	// Version of the profile
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetWithVersionRequest) Reset() {
	*x = GetWithVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWithVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWithVersionRequest) ProtoMessage() {}

func (x *GetWithVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWithVersionRequest.ProtoReflect.Descriptor instead.
func (*GetWithVersionRequest) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{4}
}

func (x *GetWithVersionRequest) GetSourceName() string {
	if x != nil {
		return x.SourceName
	}
	return ""
}

func (x *GetWithVersionRequest) GetProfileName() string {
	if x != nil {
		return x.ProfileName
	}
	return ""
}

func (x *GetWithVersionRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// GetWithVersionResponse defines response parameters for GetWithVersion endpoint.
type GetWithVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Item *ProfileCatalogEntry `protobuf:"bytes,1,opt,name=item,proto3" json:"item,omitempty"`
}

func (x *GetWithVersionResponse) Reset() {
	*x = GetWithVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWithVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWithVersionResponse) ProtoMessage() {}

func (x *GetWithVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWithVersionResponse.ProtoReflect.Descriptor instead.
func (*GetWithVersionResponse) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{5}
}

func (x *GetWithVersionResponse) GetItem() *ProfileCatalogEntry {
	if x != nil {
		return x.Item
	}
	return nil
}

// ProfilesGreaterThanVersionRequest defines request parameters for ProfilesGreaterThanVersion endpoint.
type ProfilesGreaterThanVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the catalog
	SourceName string `protobuf:"bytes,1,opt,name=source_name,json=sourceName,proto3" json:"source_name,omitempty"`
	// Name of the profile
	ProfileName string `protobuf:"bytes,2,opt,name=profile_name,json=profileName,proto3" json:"profile_name,omitempty"`
	// Version of the profile
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *ProfilesGreaterThanVersionRequest) Reset() {
	*x = ProfilesGreaterThanVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfilesGreaterThanVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfilesGreaterThanVersionRequest) ProtoMessage() {}

func (x *ProfilesGreaterThanVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfilesGreaterThanVersionRequest.ProtoReflect.Descriptor instead.
func (*ProfilesGreaterThanVersionRequest) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{6}
}

func (x *ProfilesGreaterThanVersionRequest) GetSourceName() string {
	if x != nil {
		return x.SourceName
	}
	return ""
}

func (x *ProfilesGreaterThanVersionRequest) GetProfileName() string {
	if x != nil {
		return x.ProfileName
	}
	return ""
}

func (x *ProfilesGreaterThanVersionRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// ProfilesGreaterThanVersionResponse defines response parameters for ProfilesGreaterThanVersion endpoint.
type ProfilesGreaterThanVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*ProfileCatalogEntry `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *ProfilesGreaterThanVersionResponse) Reset() {
	*x = ProfilesGreaterThanVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfilesGreaterThanVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfilesGreaterThanVersionResponse) ProtoMessage() {}

func (x *ProfilesGreaterThanVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfilesGreaterThanVersionResponse.ProtoReflect.Descriptor instead.
func (*ProfilesGreaterThanVersionResponse) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{7}
}

func (x *ProfilesGreaterThanVersionResponse) GetItems() []*ProfileCatalogEntry {
	if x != nil {
		return x.Items
	}
	return nil
}

// SearchRequest defines request parameters for Search endpoint.
type SearchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Defines a name to search for that is included in a profile's name
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *SearchRequest) Reset() {
	*x = SearchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchRequest) ProtoMessage() {}

func (x *SearchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchRequest.ProtoReflect.Descriptor instead.
func (*SearchRequest) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{8}
}

func (x *SearchRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// SearchResponse defines response parameters for Search endpoint.
type SearchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*ProfileCatalogEntry `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *SearchResponse) Reset() {
	*x = SearchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchResponse) ProtoMessage() {}

func (x *SearchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchResponse.ProtoReflect.Descriptor instead.
func (*SearchResponse) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{9}
}

func (x *SearchResponse) GetItems() []*ProfileCatalogEntry {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_profiles_proto protoreflect.FileDescriptor

var file_profiles_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x17, 0x77, 0x65, 0x61, 0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x50, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x22, 0x4f, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x40, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2c, 0x2e, 0x77, 0x65, 0x61, 0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x69,
	0x74, 0x65, 0x6d, 0x22, 0xdc, 0x01, 0x0a, 0x13, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x43,
	0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x74,
	0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x25, 0x0a,
	0x0e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x53, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a,
	0x6d, 0x61, 0x69, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x0d,
	0x70, 0x72, 0x65, 0x72, 0x65, 0x71, 0x75, 0x69, 0x73, 0x69, 0x74, 0x65, 0x73, 0x18, 0x07, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0d, 0x70, 0x72, 0x65, 0x72, 0x65, 0x71, 0x75, 0x69, 0x73, 0x69, 0x74,
	0x65, 0x73, 0x22, 0x75, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x57, 0x69, 0x74, 0x68, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x5a, 0x0a, 0x16, 0x47, 0x65, 0x74,
	0x57, 0x69, 0x74, 0x68, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2c, 0x2e, 0x77, 0x65, 0x61, 0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x04, 0x69, 0x74, 0x65, 0x6d, 0x22, 0x81, 0x01, 0x0a, 0x21, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x47, 0x72, 0x65, 0x61, 0x74, 0x65, 0x72, 0x54, 0x68, 0x61, 0x6e, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x68, 0x0a, 0x22, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x47, 0x72, 0x65, 0x61, 0x74, 0x65, 0x72, 0x54, 0x68, 0x61, 0x6e,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x42, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c,
	0x2e, 0x77, 0x65, 0x61, 0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x22, 0x23, 0x0a, 0x0d, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x54, 0x0a, 0x0e, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x77, 0x65, 0x61, 0x76,
	0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c,
	0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x32, 0xa0,
	0x05, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x83, 0x01, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x23, 0x2e, 0x77, 0x65, 0x61,
	0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x24, 0x2e, 0x77, 0x65, 0x61, 0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x31, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2b, 0x12, 0x29, 0x2f,
	0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x7b, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x2f, 0x7b, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x12, 0xae, 0x01, 0x0a, 0x0e, 0x47, 0x65, 0x74,
	0x57, 0x69, 0x74, 0x68, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2e, 0x2e, 0x77, 0x65,
	0x61, 0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x69, 0x74, 0x68, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x77, 0x65,
	0x61, 0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x69, 0x74, 0x68, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3b, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x35, 0x12, 0x33, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x2f, 0x7b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d,
	0x2f, 0x7b, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x2f,
	0x7b, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x7d, 0x12, 0xe4, 0x01, 0x0a, 0x1a, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x47, 0x72, 0x65, 0x61, 0x74, 0x65, 0x72, 0x54, 0x68, 0x61,
	0x6e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x2e, 0x77, 0x65, 0x61, 0x76, 0x65,
	0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x47, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x72, 0x54, 0x68, 0x61, 0x6e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x3b, 0x2e, 0x77, 0x65, 0x61, 0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72,
	0x6b, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x47, 0x72, 0x65, 0x61, 0x74, 0x65, 0x72, 0x54, 0x68,
	0x61, 0x6e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x4d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x47, 0x12, 0x45, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x7b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x2f, 0x7b, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x7d, 0x2f, 0x7b, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x7d, 0x2f, 0x61,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73,
	0x12, 0x6f, 0x0a, 0x06, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x26, 0x2e, 0x77, 0x65, 0x61,
	0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x27, 0x2e, 0x77, 0x65, 0x61, 0x76, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x0e, 0x12, 0x0c, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x77, 0x65, 0x61, 0x76, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_profiles_proto_rawDescOnce sync.Once
	file_profiles_proto_rawDescData = file_profiles_proto_rawDesc
)

func file_profiles_proto_rawDescGZIP() []byte {
	file_profiles_proto_rawDescOnce.Do(func() {
		file_profiles_proto_rawDescData = protoimpl.X.CompressGZIP(file_profiles_proto_rawDescData)
	})
	return file_profiles_proto_rawDescData
}

var file_profiles_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_profiles_proto_goTypes = []interface{}{
	(*Empty)(nil),                              // 0: weave.works.profiles.v1.Empty
	(*GetRequest)(nil),                         // 1: weave.works.profiles.v1.GetRequest
	(*GetResponse)(nil),                        // 2: weave.works.profiles.v1.GetResponse
	(*ProfileCatalogEntry)(nil),                // 3: weave.works.profiles.v1.ProfileCatalogEntry
	(*GetWithVersionRequest)(nil),              // 4: weave.works.profiles.v1.GetWithVersionRequest
	(*GetWithVersionResponse)(nil),             // 5: weave.works.profiles.v1.GetWithVersionResponse
	(*ProfilesGreaterThanVersionRequest)(nil),  // 6: weave.works.profiles.v1.ProfilesGreaterThanVersionRequest
	(*ProfilesGreaterThanVersionResponse)(nil), // 7: weave.works.profiles.v1.ProfilesGreaterThanVersionResponse
	(*SearchRequest)(nil),                      // 8: weave.works.profiles.v1.SearchRequest
	(*SearchResponse)(nil),                     // 9: weave.works.profiles.v1.SearchResponse
}
var file_profiles_proto_depIdxs = []int32{
	3, // 0: weave.works.profiles.v1.GetResponse.item:type_name -> weave.works.profiles.v1.ProfileCatalogEntry
	3, // 1: weave.works.profiles.v1.GetWithVersionResponse.item:type_name -> weave.works.profiles.v1.ProfileCatalogEntry
	3, // 2: weave.works.profiles.v1.ProfilesGreaterThanVersionResponse.items:type_name -> weave.works.profiles.v1.ProfileCatalogEntry
	3, // 3: weave.works.profiles.v1.SearchResponse.items:type_name -> weave.works.profiles.v1.ProfileCatalogEntry
	1, // 4: weave.works.profiles.v1.ProfilesService.Get:input_type -> weave.works.profiles.v1.GetRequest
	4, // 5: weave.works.profiles.v1.ProfilesService.GetWithVersion:input_type -> weave.works.profiles.v1.GetWithVersionRequest
	6, // 6: weave.works.profiles.v1.ProfilesService.ProfilesGreaterThanVersion:input_type -> weave.works.profiles.v1.ProfilesGreaterThanVersionRequest
	8, // 7: weave.works.profiles.v1.ProfilesService.Search:input_type -> weave.works.profiles.v1.SearchRequest
	2, // 8: weave.works.profiles.v1.ProfilesService.Get:output_type -> weave.works.profiles.v1.GetResponse
	5, // 9: weave.works.profiles.v1.ProfilesService.GetWithVersion:output_type -> weave.works.profiles.v1.GetWithVersionResponse
	7, // 10: weave.works.profiles.v1.ProfilesService.ProfilesGreaterThanVersion:output_type -> weave.works.profiles.v1.ProfilesGreaterThanVersionResponse
	9, // 11: weave.works.profiles.v1.ProfilesService.Search:output_type -> weave.works.profiles.v1.SearchResponse
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_profiles_proto_init() }
func file_profiles_proto_init() {
	if File_profiles_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_profiles_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profiles_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profiles_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profiles_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileCatalogEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profiles_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWithVersionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profiles_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWithVersionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profiles_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfilesGreaterThanVersionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profiles_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfilesGreaterThanVersionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profiles_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profiles_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_profiles_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_profiles_proto_goTypes,
		DependencyIndexes: file_profiles_proto_depIdxs,
		MessageInfos:      file_profiles_proto_msgTypes,
	}.Build()
	File_profiles_proto = out.File
	file_profiles_proto_rawDesc = nil
	file_profiles_proto_goTypes = nil
	file_profiles_proto_depIdxs = nil
}

package api_test

import (
	"context"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/api"
	catfakes "github.com/weaveworks/profiles/pkg/api/fakes"
	"github.com/weaveworks/profiles/pkg/protos"
)

var _ = Describe("API", func() {
	var (
		catalogAPI  api.CatalogAPI
		fakeCatalog *catfakes.FakeCatalog
	)

	BeforeEach(func() {
		fakeCatalog = new(catfakes.FakeCatalog)
		catalogAPI = api.NewCatalogAPI(fakeCatalog, logr.Discard())
	})

	Context("Get", func() {
		When("a matching profile exists", func() {
			BeforeEach(func() {
				fakeCatalog.GetReturns(&profilesv1.ProfileCatalogEntry{
					ProfileDescription: profilesv1.ProfileDescription{
						Name:        "nginx-1",
						Description: "nginx 1",
					},
					CatalogSource: "foo",
				})
			})

			It("returns the matching profiles from the catalog", func() {
				result, err := catalogAPI.Get(context.Background(), &protos.GetRequest{
					ProfileName: "nginx-1",
					SourceName:  "foo",
				})
				Expect(err).NotTo(HaveOccurred())
				expected := &protos.GetResponse{
					Item: &protos.ProfileCatalogEntry{
						CatalogSource: "foo",
						Name:          "nginx-1",
						Description:   "nginx 1",
					},
				}
				Expect(result).To(Equal(expected))
			})
		})
		When("there is no matching profile", func() {
			It("return a not found error", func() {
				result, err := catalogAPI.Get(context.Background(), &protos.GetRequest{
					ProfileName: "invalid",
					SourceName:  "invalid",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("profile not found"))
				Expect(grpcErr.Code()).To(Equal(codes.NotFound))
				Expect(result).To(BeNil())
			})
		})
		When("source name empty", func() {
			It("returns a proper error", func() {
				result, err := catalogAPI.Get(context.Background(), &protos.GetRequest{
					ProfileName: "foo",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("missing query param: sourceName: \"\", profileName: \"foo\""))
				Expect(grpcErr.Code()).To(Equal(codes.InvalidArgument))
				Expect(result).To(BeNil())
			})
		})
		When("profile name empty", func() {
			It("returns a proper error", func() {
				result, err := catalogAPI.Get(context.Background(), &protos.GetRequest{
					SourceName: "foo",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("missing query param: sourceName: \"foo\", profileName: \"\""))
				Expect(grpcErr.Code()).To(Equal(codes.InvalidArgument))
				Expect(result).To(BeNil())
			})
		})
	})
	Context("GetWithVersion", func() {
		When("a matching profile exists", func() {
			BeforeEach(func() {
				fakeCatalog.GetWithVersionReturns(&profilesv1.ProfileCatalogEntry{
					ProfileDescription: profilesv1.ProfileDescription{
						Name:        "nginx-1",
						Description: "nginx 1",
					},
					CatalogSource: "foo",
					Tag:           "v0.0.1",
				})
			})

			It("returns the matching profiles from the catalog", func() {
				result, err := catalogAPI.GetWithVersion(context.Background(), &protos.GetWithVersionRequest{
					ProfileName: "nginx-1",
					SourceName:  "foo",
					Version:     "v0.0.1",
				})
				Expect(err).NotTo(HaveOccurred())
				expected := &protos.GetWithVersionResponse{
					Item: &protos.ProfileCatalogEntry{
						CatalogSource: "foo",
						Name:          "nginx-1",
						Description:   "nginx 1",
						Tag:           "v0.0.1",
					},
				}
				Expect(result).To(Equal(expected))
			})
		})
		When("there is no matching profile", func() {
			It("return a not found error", func() {
				result, err := catalogAPI.GetWithVersion(context.Background(), &protos.GetWithVersionRequest{
					ProfileName: "invalid",
					SourceName:  "invalid",
					Version:     "invalid",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("profile not found"))
				Expect(grpcErr.Code()).To(Equal(codes.NotFound))
				Expect(result).To(BeNil())
			})
		})
		When("source name empty", func() {
			It("returns a proper error", func() {
				result, err := catalogAPI.GetWithVersion(context.Background(), &protos.GetWithVersionRequest{
					ProfileName: "foo",
					Version:     "whatever",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("missing query param: sourceName: \"\", profileName: \"foo\", version: \"whatever\""))
				Expect(grpcErr.Code()).To(Equal(codes.InvalidArgument))
				Expect(result).To(BeNil())
			})
		})
		When("profile name empty", func() {
			It("returns a proper error", func() {
				result, err := catalogAPI.GetWithVersion(context.Background(), &protos.GetWithVersionRequest{
					SourceName: "foo",
					Version:    "whatever",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("missing query param: sourceName: \"foo\", profileName: \"\", version: \"whatever\""))
				Expect(grpcErr.Code()).To(Equal(codes.InvalidArgument))
				Expect(result).To(BeNil())
			})
		})
		When("version name empty", func() {
			It("returns a proper error", func() {
				result, err := catalogAPI.GetWithVersion(context.Background(), &protos.GetWithVersionRequest{
					SourceName:  "foo",
					ProfileName: "bar",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("missing query param: sourceName: \"foo\", profileName: \"bar\", version: \"\""))
				Expect(grpcErr.Code()).To(Equal(codes.InvalidArgument))
				Expect(result).To(BeNil())
			})
		})
	})

	Context("Search", func() {
		When("a query matches some profiles", func() {
			BeforeEach(func() {
				fakeCatalog.SearchReturns([]profilesv1.ProfileCatalogEntry{
					{
						ProfileDescription: profilesv1.ProfileDescription{
							Name:        "nginx-1",
							Description: "nginx 1",
						},
						CatalogSource: "foo",
					},
				})
			})
			It("returns valid results", func() {
				result, err := catalogAPI.Search(context.Background(), &protos.SearchRequest{
					Name: "nginx",
				})
				Expect(err).NotTo(HaveOccurred())
				expected := &protos.SearchResponse{
					Items: []*protos.ProfileCatalogEntry{
						{
							CatalogSource: "foo",
							Name:          "nginx-1",
							Description:   "nginx 1",
						},
					},
				}
				Expect(result).To(Equal(expected))
			})
		})
		When("no query is provided", func() {
			BeforeEach(func() {
				fakeCatalog.SearchAllReturns([]profilesv1.ProfileCatalogEntry{
					{
						ProfileDescription: profilesv1.ProfileDescription{
							Name:        "nginx-1",
							Description: "nginx 1",
						},
						CatalogSource: "foo",
					},
					{
						ProfileDescription: profilesv1.ProfileDescription{
							Name:        "redis-1",
							Description: "redis 1",
						},
						CatalogSource: "foo",
					},
				})
			})
			It("returns all profiles", func() {
				result, err := catalogAPI.Search(context.Background(), &protos.SearchRequest{})
				Expect(err).NotTo(HaveOccurred())
				expected := &protos.SearchResponse{
					Items: []*protos.ProfileCatalogEntry{
						{
							CatalogSource: "foo",
							Name:          "nginx-1",
							Description:   "nginx 1",
						},
						{
							CatalogSource: "foo",
							Name:          "redis-1",
							Description:   "redis 1",
						},
					},
				}
				Expect(result).To(Equal(expected))
			})
		})
		When("there are no profiles", func() {
			It("returns an empty response", func() {
				result, err := catalogAPI.Search(context.Background(), &protos.SearchRequest{})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Items).To(BeEmpty())
			})
		})
	})

	Context("ProfilesGreaterThanVersion", func() {
		When("there are higher versions for a profile available", func() {
			BeforeEach(func() {
				fakeCatalog.ProfilesGreaterThanVersionReturns([]profilesv1.ProfileCatalogEntry{{
					ProfileDescription: profilesv1.ProfileDescription{
						Name:        "nginx-1",
						Description: "nginx 1",
					},
					CatalogSource: "foo",
					Tag:           "v0.0.2",
				}})
			})

			It("returns the profile with the higher version from the catalog", func() {
				result, err := catalogAPI.ProfilesGreaterThanVersion(context.Background(), &protos.ProfilesGreaterThanVersionRequest{
					ProfileName: "nginx-1",
					SourceName:  "foo",
					Version:     "v0.0.1",
				})
				Expect(err).NotTo(HaveOccurred())
				expected := &protos.ProfilesGreaterThanVersionResponse{
					Items: []*protos.ProfileCatalogEntry{{
						CatalogSource: "foo",
						Name:          "nginx-1",
						Description:   "nginx 1",
						Tag:           "v0.0.2",
					}},
				}
				Expect(result).To(Equal(expected))
			})
		})
		When("there is no matching profile", func() {
			It("return a not found error", func() {
				result, err := catalogAPI.ProfilesGreaterThanVersion(context.Background(), &protos.ProfilesGreaterThanVersionRequest{
					ProfileName: "invalid",
					SourceName:  "invalid",
					Version:     "invalid",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("profile not found"))
				Expect(grpcErr.Code()).To(Equal(codes.NotFound))
				Expect(result).To(BeNil())
			})
		})
		When("source name empty", func() {
			It("returns a proper error", func() {
				result, err := catalogAPI.ProfilesGreaterThanVersion(context.Background(), &protos.ProfilesGreaterThanVersionRequest{
					ProfileName: "foo",
					Version:     "whatever",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("missing query param: sourceName: \"\", profileName: \"foo\", version: \"whatever\""))
				Expect(grpcErr.Code()).To(Equal(codes.InvalidArgument))
				Expect(result).To(BeNil())
			})
		})
		When("profile name empty", func() {
			It("returns a proper error", func() {
				result, err := catalogAPI.ProfilesGreaterThanVersion(context.Background(), &protos.ProfilesGreaterThanVersionRequest{
					SourceName: "foo",
					Version:    "whatever",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("missing query param: sourceName: \"foo\", profileName: \"\", version: \"whatever\""))
				Expect(grpcErr.Code()).To(Equal(codes.InvalidArgument))
				Expect(result).To(BeNil())
			})
		})
		When("version name empty", func() {
			It("returns a proper error", func() {
				result, err := catalogAPI.ProfilesGreaterThanVersion(context.Background(), &protos.ProfilesGreaterThanVersionRequest{
					SourceName:  "foo",
					ProfileName: "bar",
				})
				grpcErr, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcErr.Message()).To(Equal("missing query param: sourceName: \"foo\", profileName: \"bar\", version: \"\""))
				Expect(grpcErr.Code()).To(Equal(codes.InvalidArgument))
				Expect(result).To(BeNil())
			})
		})
	})
})

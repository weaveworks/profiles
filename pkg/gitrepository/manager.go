package gitrepository

import (
	"context"
	"fmt"
	"strings"
	"time"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//Manager is responsible for managing gitrepository resources
type Manager struct {
	kClient   Kubernetes
	namespace string
	ctx       context.Context
	timeout   time.Duration
	interval  time.Duration
}

//Instance contains a tag and path of profile.yaml
type Instance struct {
	Tag  string
	Path string
}

//go:generate counterfeiter -o fakes/fake_kubernetes.go . Kubernetes
//Kubernetes interface for itneracting with kubernetes
type Kubernetes interface {
	Get(ctx context.Context, key client.ObjectKey, obj client.Object) error
	Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error
	Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error
}

//NewManager returns a Manager struct
func NewManager(ctx context.Context, namespace string, kClient Kubernetes, timeout, interval time.Duration) *Manager {
	return &Manager{
		kClient:   kClient,
		namespace: namespace,
		ctx:       context.TODO(),
		timeout:   timeout,
		interval:  interval,
	}
}

//CreateAndWaitForResources creates the gitrepository resources and waits for them to be created
func (m *Manager) CreateAndWaitForResources(r profilesv1.Repository, instances []Instance) ([]*sourcev1.GitRepository, error) {
	var gitResources []*sourcev1.GitRepository
	for _, instance := range instances {
		gitRes := makeGitRepository(r, instance.Tag, instance.Path, m.namespace)
		err := m.kClient.Create(m.ctx, gitRes)
		if err != nil {
			return nil, fmt.Errorf("failed to create gitrepository: %w", err)
		}
		gitResources = append(gitResources, gitRes)
	}

	for _, gitRes := range gitResources {
		err := m.waitForURL(gitRes)
		if err != nil {
			return nil, err
		}
	}
	return gitResources, nil
}

func (m *Manager) waitForURL(gitRes *sourcev1.GitRepository) error {
	startTime := time.Now()
	for {
		err := m.kClient.Get(m.ctx, client.ObjectKeyFromObject(gitRes), gitRes)
		if err != nil {
			return fmt.Errorf("failed to get gitrepository: %w", err)
		}
		if gitRes.Status.URL != "" {
			return nil
		}

		if time.Since(startTime) > m.timeout {
			return fmt.Errorf("timed out waiting for %s/%s gitrepository.Status.URL to be populated", gitRes.Namespace, gitRes.Name)
		}
		time.Sleep(m.interval)
	}
}

func makeGitRepository(r profilesv1.Repository, tag, path, namespace string) *sourcev1.GitRepository {
	ignore := fmt.Sprintf(`# exclude all
/*
# include deploy dir
!/%s`, path)

	repo := &sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeGitRepoName(tag, r.URL),
			Namespace: namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       sourcev1.GitRepositoryKind,
			APIVersion: sourcev1.GroupVersion.String(),
		},
		Spec: sourcev1.GitRepositorySpec{
			URL: r.URL,
			Reference: &sourcev1.GitRepositoryRef{
				Tag: tag,
			},
			Ignore: &ignore,
		},
	}

	if r.SecretRef != nil {
		repo.Spec.SecretRef = r.SecretRef
	}
	return repo
}

func makeGitRepoName(tag, url string) string {
	urlParts := strings.Split(url, "/")
	repo := strings.TrimRight(urlParts[len(urlParts)-1], ".git")

	splitTag := strings.Split(tag, "/")
	if len(splitTag) == 2 {
		return fmt.Sprintf("%s-%s-%s", repo, splitTag[0], splitTag[1])
	}
	return fmt.Sprintf("%s-%s", repo, tag)
}

//DeleteResources deletes the gitrepository resources
func (m *Manager) DeleteResources(gitRepos []*sourcev1.GitRepository) error {
	for _, res := range gitRepos {
		err := m.kClient.Delete(m.ctx, res)
		if err != nil {
			return err
		}
	}
	return nil
}

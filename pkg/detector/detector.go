package detector

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/fluxcd/pkg/version"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-logr/logr"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"github.com/weaveworks/profiles/pkg/catalog"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// given a repo url
// discover tags
// create gitrepository resources
// fetch the tarballs and read the profile.yaml
// populate catalog
// ???
// profit

func Detect(catalogSource profilesv1.ProfileCatalogSource, catalogManager *catalog.Catalog, kClient client.Client, logger logr.Logger) error {
	logger.Info(fmt.Sprintf("scanning repo: %s", catalogSource.Spec.Repo))
	cmd := exec.Command("git", "ls-remote", "--tags", catalogSource.Spec.Repo)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	logger.Info(string(out))
	rawTags := strings.Split(string(out), "\n")
	// fmt.Println(tags)

	tags := getValidTags(rawTags)
	fmt.Println(tags)

	gitRepoResources := createAndWaitForGitRepos(tags, catalogSource.Spec.Repo, kClient)
	fmt.Println("finished creating and waiting for git repos")
	profiles := readProfileFromTarball(catalogSource, gitRepoResources)
	fmt.Println("finished fetching and unmarshalling tarballs")
	for _, profile := range profiles {
		fmt.Println("adding profile:")
		fmt.Printf("%v\n", profile)
	}
	catalogManager.Update(catalogSource.Name, profiles...)
	fmt.Println("finished")

	return nil
}

func readProfileFromTarball(catalogSource profilesv1.ProfileCatalogSource, gitRepoResources []*sourcev1.GitRepository) []profilesv1.ProfileDescription {
	var profileDescriptions []profilesv1.ProfileDescription
	for _, gitRepo := range gitRepoResources {
		// dirName := strings.Split(gitRepo.Spec.Reference.Tag, "/")[0]
		url := gitRepo.Status.URL

		// ONLY NEEDED IF RUNNING LOCALLY. You must open a port-forward yourself to forward connections to source controller
		// kubectl -n flux-system port-forward source-controller-<random-guids> 8888:9090
		// url = strings.Replace(url, "source-controller.flux-system.svc.cluster.local.", "localhost:8888", -1)

		fmt.Printf("url: %s\n", url)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		profileYaml := extractTarGz(resp.Body)
		decoder := yaml.NewDecoder(profileYaml)
		var profileDef profilesv1.ProfileDefinition
		err = decoder.Decode(&profileDef)
		if err != nil {
			panic(err)
		}
		fmt.Println("profile.yaml:")
		fmt.Printf("%v\n", profileDef)
		profileDescriptions = append(profileDescriptions, profilesv1.ProfileDescription{
			Name:          strings.Split(gitRepo.Spec.Reference.Tag, "/")[0],
			Description:   profileDef.Spec.Description,
			Version:       strings.Split(gitRepo.Spec.Reference.Tag, "/")[1],
			CatalogSource: catalogSource.Name,
			URL:           catalogSource.Spec.Repo,
		})
	}
	return profileDescriptions
}

// copied from https://stackoverflow.com/questions/57639648/how-to-decompress-tar-gz-file-in-go
func extractTarGz(gzipStream io.Reader) io.Reader {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			fmt.Println("dir found")
			continue
		case tar.TypeReg:
			fmt.Println("file found")
			// outFile, err := os.Create(header.Name)
			// if err != nil {
			// 	log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			// }
			// if _, err := io.Copy(outFile, tarReader); err != nil {
			// 	log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			// }
			// outFile.Close()

			return tarReader
			// out, err := ioutil.ReadAll(tarReader)
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Println(string(out))

		default:
			fmt.Println("not sure what was found")
			log.Fatalf(
				"ExtractTarGz: uknown type: %v in %s",
				header.Typeflag,
				header.Name)
		}

	}
	return nil
}

func createAndWaitForGitRepos(tags []string, repo string, kClient client.Client) []*sourcev1.GitRepository {
	var gitRepoResources []*sourcev1.GitRepository
	for _, tag := range tags {
		gitRepo := makeGitRepository(tag, repo)
		gitRepoResources = append(gitRepoResources, gitRepo)
		err := kClient.Create(context.TODO(), gitRepo)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("created resources")
	for _, gitRep := range gitRepoResources {
		fmt.Println(gitRep.Name)
	}
	done := false
	fmt.Println("checking status")
	for !done {
		done = waitForStatus(gitRepoResources, kClient)
		fmt.Println("waiting...")
		time.Sleep(time.Second)
	}
	return gitRepoResources
}

func waitForStatus(gitRepoResources []*sourcev1.GitRepository, kClient client.Client) bool {
	for _, gitRep := range gitRepoResources {
		err := kClient.Get(context.TODO(), types.NamespacedName{Name: gitRep.Name, Namespace: gitRep.Namespace}, gitRep)
		if err != nil {
			panic(err)
		}
		if gitRep.Status.URL == "" {
			return false
		}
	}
	return true
}

func makeGitRepository(tag, url string) *sourcev1.GitRepository {
	ref := &sourcev1.GitRepositoryRef{
		Tag: tag,
	}
	profileName := strings.Split(tag, "/")[0]
	ignore := fmt.Sprintf(`# exclude all
/*
# include deploy dir
!/%s/profile.yaml`, profileName)

	tag = strings.Replace(tag, "/", "-", -1)
	tag = strings.Replace(tag, ".", "-", -1)

	return &sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tag,
			Namespace: "profiles-system",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       sourcev1.GitRepositoryKind,
			APIVersion: sourcev1.GroupVersion.String(),
		},
		Spec: sourcev1.GitRepositorySpec{
			URL:       url,
			Reference: ref,
			Ignore:    &ignore,
		},
	}
}

func getValidTags(rawTags []string) []string {
	var tags []string
	for _, tag := range rawTags {
		bits := strings.Split(tag, "refs/tags/")
		unParsedTag := bits[len(bits)-1]

		splitTag := strings.Split(unParsedTag, "/")
		if len(splitTag) != 2 {
			continue
		}
		_, err := version.ParseVersion(splitTag[1])
		if err != nil {
			continue
		}
		tags = append(tags, unParsedTag)
	}
	return tags
}

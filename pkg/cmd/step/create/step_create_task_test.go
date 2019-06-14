package create_test

import (
	"fmt"
	"github.com/jenkins-x/jx/pkg/cmd/step/create"
	"github.com/jenkins-x/jx/pkg/cmd/testhelpers"
	"github.com/jenkins-x/jx/pkg/log"
	"github.com/jenkins-x/jx/pkg/prow"
	"github.com/jenkins-x/jx/pkg/tekton"

	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jenkins-x/jx/pkg/gits/mocks"
	"github.com/jenkins-x/jx/pkg/helm/mocks"
	"github.com/jenkins-x/jx/pkg/kube"
	"github.com/knative/pkg/kmp"
	"github.com/satori/go.uuid"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/ghodss/yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/jenkins-x/jx/pkg/config"
	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/jenkinsfile"
	"github.com/jenkins-x/jx/pkg/tekton/tekton_helpers_test"
	"github.com/jenkins-x/jx/pkg/tests"
	"github.com/stretchr/testify/assert"
	pipelineapi "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerateTektonCRDs(t *testing.T) {
	t.Parallel()

	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)

	testData := path.Join("test_data", "step_create_task")
	_, err := os.Stat(testData)
	assert.NoError(t, err)

	testVersionsDir := path.Join(testData, "stable_versions")
	packsDir := path.Join(testData, "packs")
	_, err = os.Stat(packsDir)
	assert.NoError(t, err)

	resolver := func(importFile *jenkinsfile.ImportFile) (string, error) {
		dirPath := []string{packsDir, "import_dir", importFile.Import}
		// lets handle cross platform paths in `importFile.File`
		path := append(dirPath, strings.Split(importFile.File, "/")...)
		return filepath.Join(path...), nil
	}

	cases := []struct {
		name           string
		language       string
		repoName       string
		organization   string
		branch         string
		kind           string
		expectingError bool
	}{
		{
			name:         "js_build_pack",
			language:     "javascript",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "build-pack",
			kind:         "release",
		},
		{
			name:         "maven_build_pack",
			language:     "maven",
			repoName:     "jx-demo-qs",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "from_yaml",
			language:     "none",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "really-long",
			kind:         "release",
		},
		{
			name:           "no_pipeline_config",
			language:       "none",
			repoName:       "anything",
			organization:   "anything",
			branch:         "anything",
			kind:           "release",
			expectingError: true,
		},
		{
			name:         "per_step_container_build_pack",
			language:     "apps",
			repoName:     "golang-qs-test",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "kaniko_entrypoint",
			language:     "none",
			repoName:     "jx",
			organization: "jenkins-x",
			branch:       "fix-kaniko-special-casing",
			kind:         "pullrequest",
		},
		{
			name:         "set-agent-container-with-agentless-build-pack",
			language:     "no-default-agent",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "no-default-agent",
			kind:         "release",
		},
		{
			name:         "override-agent-container-with-build-pack",
			language:     "override-default-agent",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "override-default-agent",
			kind:         "release",
		},
		{
			name:         "override-steps",
			language:     "maven",
			repoName:     "jx-demo-qs",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "override_block_step",
			language:     "apps",
			repoName:     "golang-qs-test",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "loop-in-buildpack-syntax",
			language:     "maven",
			repoName:     "jx-demo-qs",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "containeroptions-on-pipelineconfig",
			language:     "maven",
			repoName:     "jx-demo-qs",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "default-in-jenkins-x-yml",
			language:     "none",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "really-long",
			kind:         "release",
		},
		{
			name:         "default-in-buildpack",
			language:     "default-pipeline",
			repoName:     "golang-qs-test",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "add-env-to-default-in-buildpack",
			language:     "default-pipeline",
			repoName:     "golang-qs-test",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "override-default-in-buildpack",
			language:     "default-pipeline",
			repoName:     "golang-qs-test",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "override-default-in-jenkins-x-yml",
			language:     "none",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "really-long",
			kind:         "release",
		},
		{
			name:         "remove-stage",
			language:     "maven",
			repoName:     "jx-demo-qs",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:           "remove-pipeline",
			language:       "none",
			repoName:       "anything",
			organization:   "anything",
			branch:         "anything",
			kind:           "pullRequest",
			expectingError: true,
		},
		{
			name:         "remove-stage-from-jenkins-x-yml",
			language:     "none",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "really-long",
			kind:         "release",
		},
		{
			name:           "remove-pipeline-from-jenkins-x-yml",
			language:       "none",
			repoName:       "anything",
			organization:   "anything",
			branch:         "anything",
			kind:           "pullRequest",
			expectingError: true,
		},
		{
			name:         "replace-stage-steps",
			language:     "maven",
			repoName:     "jx-demo-qs",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "append-and-prepend-stage-steps",
			language:     "maven",
			repoName:     "jx-demo-qs",
			organization: "abayer",
			branch:       "master",
			kind:         "release",
		},
		{
			name:         "replace-stage-steps-in-jenkins-x-yml",
			language:     "none",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "really-long",
			kind:         "release",
		},
		{
			name:         "append-and-prepend-stage-steps-in-jenkins-x-yml",
			language:     "none",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "really-long",
			kind:         "release",
		},
		{
			name:         "correct-pipeline-stage-is-removed",
			language:     "none",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "really-long",
			kind:         "release",
		},
		{
			name:         "command-as-multiline-script",
			language:     "none",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "really-long",
			kind:         "release",
		},
		{
			name:         "pipeline-timeout",
			language:     "none",
			repoName:     "js-test-repo",
			organization: "abayer",
			branch:       "really-long",
			kind:         "release",
		},
	}

	k8sObjects := []runtime.Object{
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      kube.ConfigMapJenkinsDockerRegistry,
				Namespace: "jx",
			},
			Data: map[string]string{
				"docker.registry": "1.2.3.4:5000",
			},
		},
	}
	jxObjects := []runtime.Object{}
	repoOwnerUUID, err := uuid.NewV4()
	assert.NoError(t, err)
	repoOwner := repoOwnerUUID.String()
	repoNameUUID, err := uuid.NewV4()
	assert.NoError(t, err)
	repoName := repoNameUUID.String()
	fakeRepo := gits.NewFakeRepository(repoOwner, repoName)
	fakeGitProvider := gits.NewFakeProvider(fakeRepo)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			caseDir := path.Join(testData, tt.name)
			_, err = os.Stat(caseDir)
			assert.NoError(t, err)

			projectConfig, projectConfigFile, err := config.LoadProjectConfig(caseDir)
			if err != nil {
				t.Fatalf("Error loading %s/jenkins-x.yml: %s", caseDir, err)
			}

			createTask := &create.StepCreateTaskOptions{
				Pack:         tt.language,
				DryRun:       true,
				SourceName:   "source",
				PodTemplates: assertLoadPodTemplates(t),
				GitInfo: &gits.GitRepository{
					Host:         "github.com",
					Name:         tt.repoName,
					Organisation: tt.organization,
				},
				Branch:       tt.branch,
				PipelineKind: tt.kind,
				NoKaniko:     true,
				Trigger:      string(pipelineapi.PipelineTriggerTypeManual),
				StepOptions: opts.StepOptions{
					CommonOptions: &opts.CommonOptions{
						ServiceAccount: "tekton-bot",
					},
				},
				BuildNumber: "1",
				VersionResolver: &opts.VersionResolver{
					VersionsDir: testVersionsDir,
				},
				DefaultImage: "maven",
			}
			testhelpers.ConfigureTestOptionsWithResources(createTask.CommonOptions, k8sObjects, jxObjects, gits_test.NewMockGitter(), fakeGitProvider, helm_test.NewMockHelmer(), nil)

			crds, err := createTask.GenerateTektonCRDs(packsDir, projectConfig, projectConfigFile, resolver, "jx")
			if tt.expectingError {
				if err == nil {
					t.Fatalf("Expected an error generating CRDs")
				}
			} else {
				if err != nil {
					t.Fatalf("Error generating CRDs: %s", err)
				}

				taskList := &pipelineapi.TaskList{}
				for _, task := range crds.Tasks() {
					taskList.Items = append(taskList.Items, *task)
				}

				resourceList := &pipelineapi.PipelineResourceList{}
				for _, resource := range crds.Resources() {
					resourceList.Items = append(resourceList.Items, *resource)
				}

				if d := cmp.Diff(tekton_helpers_test.AssertLoadPipeline(t, caseDir), crds.Pipeline()); d != "" {
					t.Errorf("Generated Pipeline did not match expected: \n%s", d)
				}
				if d, _ := kmp.SafeDiff(tekton_helpers_test.AssertLoadTasks(t, caseDir), taskList, cmpopts.IgnoreFields(corev1.ResourceRequirements{}, "Requests")); d != "" {
					t.Errorf("Generated Tasks did not match expected: \n%s", d)
				}
				if d := cmp.Diff(tekton_helpers_test.AssertLoadPipelineResources(t, caseDir), resourceList); d != "" {
					t.Errorf("Generated PipelineResources did not match expected: %s", d)
				}

				if d := cmp.Diff(tekton_helpers_test.AssertLoadPipelineRun(t, caseDir), crds.PipelineRun()); d != "" {
					t.Errorf("Generated PipelineRun did not match expected: %s", d)
				}
				if d := cmp.Diff(tekton_helpers_test.AssertLoadPipelineStructure(t, caseDir), crds.Structure()); d != "" {
					t.Errorf("Generated PipelineStructure did not match expected: %s", d)
				}

				pa := tekton.GeneratePipelineActivity(createTask.BuildNumber, createTask.Branch, createTask.GitInfo, &prow.PullRefs{})

				expectedActivityKey := &kube.PromoteStepActivityKey{
					PipelineActivityKey: kube.PipelineActivityKey{
						Name:     fmt.Sprintf("%s-%s-%s-1", tt.organization, tt.repoName, tt.branch),
						Pipeline: fmt.Sprintf("%s/%s/%s", tt.organization, tt.repoName, tt.branch),
						Build:    "1",
						GitInfo:  createTask.GitInfo,
					},
				}
				if d := cmp.Diff(expectedActivityKey, pa); d != "" {
					t.Errorf("not match expected: %s", d)
				}
			}
		})
	}
}

func assertLoadPodTemplates(t *testing.T) map[string]*corev1.Pod {
	fileName := filepath.Join("test_data", "step_create_task", "PodTemplates.yml")
	if tests.AssertFileExists(t, fileName) {
		configMap := &corev1.ConfigMap{}
		data, err := ioutil.ReadFile(fileName)
		if assert.NoError(t, err, "Failed to load file %s", fileName) {
			err = yaml.Unmarshal(data, configMap)
			if assert.NoError(t, err, "Failed to unmarshall YAML file %s", fileName) {
				podTemplates := make(map[string]*corev1.Pod)
				for k, v := range configMap.Data {
					pod := &corev1.Pod{}
					if v != "" {
						err := yaml.Unmarshal([]byte(v), pod)
						if assert.NoError(t, err, "Failed to parse pod template") {
							podTemplates[k] = pod
						}
					}
				}
				return podTemplates
			}
		}
	}
	return nil
}

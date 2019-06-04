package cmd_test

import (
	"fmt"
	"github.com/jenkins-x/jx/pkg/jx/cmd/cmd_test_helpers"
	"testing"

	"github.com/jenkins-x/jx/pkg/kube"

	helm_test "github.com/jenkins-x/jx/pkg/helm/mocks"
	"github.com/jenkins-x/jx/pkg/prow/config"

	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/prow"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/stretchr/testify/assert"

	"github.com/jenkins-x/jx/pkg/jx/cmd"
	"github.com/jenkins-x/jx/pkg/jx/cmd/opts"
	resources_test "github.com/jenkins-x/jx/pkg/kube/resources/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	protectionRepoName = "test-repo"
	protectionOrgName  = "test-org"
	protectionContext  = "test-context"
)

func TestStartProtection(t *testing.T) {
	o := cmd.StartProtectionOptions{
		CommonOptions: &opts.CommonOptions{},
	}

	cmd_test_helpers.ConfigureTestOptionsWithResources(o.CommonOptions,
		[]runtime.Object{},
		[]runtime.Object{},
		&gits.GitFake{},
		nil,
		helm_test.NewMockHelmer(),
		resources_test.NewMockInstaller(),
	)

	kubeClient, ns, err := o.KubeClientAndDevNamespace()
	assert.NoError(t, err)
	// First configure a repo in prow
	repo := fmt.Sprintf("%s/%s", protectionOrgName, protectionRepoName)
	repos := []string{repo}

	devEnv := kube.CreateDefaultDevEnvironment("jx")

	data := make(map[string]string)
	data["domain"] = "dummy.domain.nip.io"
	data["tls"] = "false"

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: kube.IngressConfigConfigmap,
		},
		Data: data,
	}

	_, err = kubeClient.CoreV1().ConfigMaps(ns).Create(cm)
	assert.NoError(t, err)

	err = prow.AddApplication(kubeClient, repos, ns, "", &devEnv.Spec.TeamSettings)
	defer func() {
		err = prow.DeleteApplication(kubeClient, repos, ns)
		assert.NoError(t, err)
	}()
	assert.NoError(t, err)

	o.Args = []string{protectionContext, repo}
	err = o.Run()
	assert.NoError(t, err)
	prowOptions := prow.Options{
		Kind:       config.Protection,
		KubeClient: kubeClient,
		NS:         ns,
	}
	prowConfig, _, err := prowOptions.GetProwConfig()
	assert.NoError(t, err)
	contexts, err := config.GetBranchProtectionContexts(protectionOrgName, protectionRepoName, prowConfig)
	assert.Equal(t, []string{protectionContext}, contexts)

}

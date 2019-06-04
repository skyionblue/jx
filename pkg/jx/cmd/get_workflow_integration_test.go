// +build integration

package cmd_test

import (
	"github.com/jenkins-x/jx/pkg/jx/cmd/cmd_test_helpers"
	"testing"

	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/helm"
	"github.com/jenkins-x/jx/pkg/jx/cmd"
	"github.com/jenkins-x/jx/pkg/jx/cmd/opts"
	"github.com/jenkins-x/jx/pkg/kube"
	resources_test "github.com/jenkins-x/jx/pkg/kube/resources/mocks"
	"github.com/jenkins-x/jx/pkg/workflow"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGetWorkflow(t *testing.T) {
	originalJxHome, tempJxHome, err := cmd_test_helpers.CreateTestJxHomeDir()
	assert.NoError(t, err)
	defer func() {
		err := cmd_test_helpers.CleanupTestJxHomeDir(originalJxHome, tempJxHome)
		assert.NoError(t, err)
	}()
	originalKubeCfg, tempKubeCfg, err := cmd_test_helpers.CreateTestKubeConfigDir()
	assert.NoError(t, err)
	defer func() {
		err := cmd_test_helpers.CleanupTestKubeConfigDir(originalKubeCfg, tempKubeCfg)
		assert.NoError(t, err)
	}()
	o := &cmd.GetWorkflowOptions{
		GetOptions: cmd.GetOptions{
			CommonOptions: &opts.CommonOptions{},
		},
	}

	staging := kube.NewPermanentEnvironment("staging")
	production := kube.NewPermanentEnvironment("production")
	staging.Spec.Order = 100
	production.Spec.Order = 200

	myFlowName := "myflow"
	cmd_test_helpers.ConfigureTestOptionsWithResources(o.CommonOptions,
		[]runtime.Object{},
		[]runtime.Object{
			staging,
			production,
			kube.NewPreviewEnvironment("jx-jstrachan-demo96-pr-1"),
			kube.NewPreviewEnvironment("jx-jstrachan-another-pr-3"),
			workflow.CreateWorkflow("jx", myFlowName,
				workflow.CreateWorkflowPromoteStep("a"),
				workflow.CreateWorkflowPromoteStep("b"),
				workflow.CreateWorkflowPromoteStep("c"),
				workflow.CreateWorkflowPromoteStep("d"),
			),
		},
		gits.NewGitCLI(),
		nil,
		helm.NewHelmCLI("helm", helm.V2, "", true),
		resources_test.NewMockInstaller(),
	)

	jxClient, ns, err := o.JXClientAndDevNamespace()
	assert.NoError(t, err)
	if err == nil {
		workflow, err := workflow.GetWorkflow("", jxClient, ns)
		assert.NoError(t, err)
		if err == nil {
			assert.Equal(t, "default", workflow.Name, "name")
			spec := workflow.Spec
			assert.Equal(t, 2, len(spec.Steps), "number of steps")
			if len(spec.Steps) > 0 {
				cmd_test_helpers.AssertPromoteStep(t, &spec.Steps[0], "staging")
			}
			if len(spec.Steps) > 1 {
				cmd_test_helpers.AssertPromoteStep(t, &spec.Steps[1], "production")
			}
		}
	}

	o.Name = myFlowName
	err = o.Run()
	assert.NoError(t, err)
}

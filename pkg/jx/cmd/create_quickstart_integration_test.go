// +build integration

package cmd_test

import (
	"github.com/jenkins-x/jx/pkg/jx/cmd/cmd_test_helpers"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/helm"
	"github.com/jenkins-x/jx/pkg/jx/cmd"
	"github.com/jenkins-x/jx/pkg/jx/cmd/opts"
	"github.com/jenkins-x/jx/pkg/quickstarts"
	"github.com/jenkins-x/jx/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestCreateQuickstartProjects(t *testing.T) {
	// TODO lets skip this test for now as it often fails with rate limits
	t.SkipNow()

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

	testDir, err := ioutil.TempDir("", "test-create-quickstart")
	assert.NoError(t, err)

	appName := "myvets"

	o := &cmd.CreateQuickstartOptions{
		CreateProjectOptions: cmd.CreateProjectOptions{
			ImportOptions: cmd.ImportOptions{
				CommonOptions: &opts.CommonOptions{},
			},
		},
		GitHubOrganisations: []string{"petclinic-gcp"},
		Filter: quickstarts.QuickstartFilter{
			Text:        "petclinic-gcp/spring-petclinic-vets-service",
			ProjectName: appName,
		},
	}
	cmd_test_helpers.ConfigureTestOptions(o.CommonOptions, gits.NewGitCLI(), helm.NewHelmCLI("helm", helm.V2, testDir, true))
	o.Dir = testDir
	o.OutDir = testDir
	o.DryRun = true
	o.DisableMaven = true
	o.Verbose = true
	o.IgnoreTeam = true
	o.Repository = appName

	err = o.Run()
	assert.NoError(t, err)
	if err == nil {
		appDir := filepath.Join(testDir, appName)
		jenkinsfile := filepath.Join(appDir, "Jenkinsfile")
		tests.AssertFileExists(t, jenkinsfile)
		tests.AssertFileExists(t, filepath.Join(appDir, "Dockerfile"))
		tests.AssertFileExists(t, filepath.Join(appDir, "charts", appName, "Chart.yaml"))
		tests.AssertFileExists(t, filepath.Join(appDir, "charts", appName, "Makefile"))
		tests.AssertFileDoesNotExist(t, filepath.Join(appDir, "charts", "spring-petclinic-vets-service", "Chart.yaml"))
	}
}

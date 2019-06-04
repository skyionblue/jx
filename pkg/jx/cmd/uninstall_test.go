package cmd_test

import (
	"github.com/jenkins-x/jx/pkg/jx/cmd/cmd_test_helpers"
	"testing"

	expect "github.com/Netflix/go-expect"
	gits_test "github.com/jenkins-x/jx/pkg/gits/mocks"
	helm_test "github.com/jenkins-x/jx/pkg/helm/mocks"
	"github.com/jenkins-x/jx/pkg/jx/cmd"
	clients_mocks "github.com/jenkins-x/jx/pkg/jx/cmd/clients/mocks"
	"github.com/jenkins-x/jx/pkg/jx/cmd/opts"
	kuber_mocks "github.com/jenkins-x/jx/pkg/kube/mocks"
	"github.com/jenkins-x/jx/pkg/tests"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd/api"

	. "github.com/petergtz/pegomock"
)

func setupUninstall(contextName string) *kuber_mocks.MockKuber {
	kubeMock := kuber_mocks.NewMockKuber()
	fakeKubeConfig := api.NewConfig()
	fakeKubeConfig.CurrentContext = contextName
	When(kubeMock.LoadConfig()).ThenReturn(fakeKubeConfig, nil, nil)
	return kubeMock
}

func TestUninstallOptions_Run_ContextSpecifiedAsOption_FailsWhenContextNamesDoNotMatch(t *testing.T) {
	originalJxHome, tempJxHome, err := cmd_test_helpers.CreateTestJxHomeDir()
	assert.NoError(t, err)
	defer func() {
		err := cmd_test_helpers.CleanupTestJxHomeDir(originalJxHome, tempJxHome)
		assert.NoError(t, err)
	}()
	kubeMock := setupUninstall("current-context")

	o := &cmd.UninstallOptions{
		CommonOptions: &opts.CommonOptions{},
		Namespace:     "ns",
		Context:       "target-context",
	}
	o.SetKube(kubeMock)
	cmd_test_helpers.ConfigureTestOptions(o.CommonOptions, gits_test.NewMockGitter(), helm_test.NewMockHelmer())

	err = o.Run()
	assert.EqualError(t, err, "The context 'target-context' must match the current context to uninstall")
}

func TestUninstallOptions_Run_ContextSpecifiedAsOption_PassWhenContextNamesMatch(t *testing.T) {
	originalJxHome, tempJxHome, err := cmd_test_helpers.CreateTestJxHomeDir()
	assert.NoError(t, err)
	defer func() {
		err := cmd_test_helpers.CleanupTestJxHomeDir(originalJxHome, tempJxHome)
		assert.NoError(t, err)
	}()
	kubeMock := setupUninstall("correct-context-to-delete")

	o := &cmd.UninstallOptions{
		CommonOptions: &opts.CommonOptions{},
		Namespace:     "ns",
		Context:       "correct-context-to-delete",
	}
	o.SetKube(kubeMock)
	cmd_test_helpers.ConfigureTestOptions(o.CommonOptions, gits_test.NewMockGitter(), helm_test.NewMockHelmer())

	// Create fake namespace (that we will uninstall from)
	err = createNamespace(o, "ns")

	// Run the uninstall
	err = o.Run()
	assert.NoError(t, err)

	// Assert that the namespace has been deleted
	client, err := o.KubeClient()
	assert.NoError(t, err)
	_, err = client.CoreV1().Namespaces().Get("ns", metav1.GetOptions{})
	assert.Error(t, err)
}

func TestUninstallOptions_Run_ContextSpecifiedAsOption_PassWhenForced(t *testing.T) {
	originalJxHome, tempJxHome, err := cmd_test_helpers.CreateTestJxHomeDir()
	assert.NoError(t, err)
	defer func() {
		err := cmd_test_helpers.CleanupTestJxHomeDir(originalJxHome, tempJxHome)
		assert.NoError(t, err)
	}()
	kubeMock := setupUninstall("correct-context-to-delete")

	o := &cmd.UninstallOptions{
		CommonOptions: &opts.CommonOptions{},
		Namespace:     "ns",
		Force:         true,
	}
	o.SetKube(kubeMock)
	cmd_test_helpers.ConfigureTestOptions(o.CommonOptions, gits_test.NewMockGitter(), helm_test.NewMockHelmer())

	// Create fake namespace (that we will uninstall from)
	err = createNamespace(o, "ns")

	// Run the uninstall
	err = o.Run()
	assert.NoError(t, err)

	// Assert that the namespace has been deleted
	client, err := o.KubeClient()
	assert.NoError(t, err)
	_, err = client.CoreV1().Namespaces().Get("ns", metav1.GetOptions{})
	assert.Error(t, err)
}

func TestUninstallOptions_Run_ContextSpecifiedViaCli_FailsWhenContextNamesDoNotMatch(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	originalJxHome, tempJxHome, err := cmd_test_helpers.CreateTestJxHomeDir()
	assert.NoError(t, err)
	defer func() {
		err := cmd_test_helpers.CleanupTestJxHomeDir(originalJxHome, tempJxHome)
		assert.NoError(t, err)
	}()
	kubeMock := setupUninstall("current-context")

	// mock terminal
	console := tests.NewTerminal(t)
	defer console.Cleanup()

	// Test interactive IO
	donec := make(chan struct{})
	go func() {
		defer close(donec)
		console.ExpectString("Enter the current context name to confirm uninstallation of the Jenkins X platform from the ns namespace:")
		console.SendLine("target-context")
		console.ExpectEOF()
	}()

	commonOpts := opts.NewCommonOptionsWithFactory(clients_mocks.NewMockFactory())
	commonOpts.In = console.In
	commonOpts.Out = console.Out
	commonOpts.Err = console.Err
	o := &cmd.UninstallOptions{
		CommonOptions: &commonOpts,
		Namespace:     "ns",
	}

	o.SetKube(kubeMock)

	err = o.Run()
	assert.EqualError(t, err, "The context 'target-context' must match the current context to uninstall")

	console.Close()
	<-donec

	// Dump the terminal's screen.
	t.Logf(expect.StripTrailingEmptyLines(console.CurrentState()))
}

func TestUninstallOptions_Run_ContextSpecifiedViaCli_PassWhenContextNamesMatch(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	originalJxHome, tempJxHome, err := cmd_test_helpers.CreateTestJxHomeDir()
	assert.NoError(t, err)
	defer func() {
		err := cmd_test_helpers.CleanupTestJxHomeDir(originalJxHome, tempJxHome)
		assert.NoError(t, err)
	}()
	kubeMock := setupUninstall("correct-context-to-delete")

	// mock terminal
	console := tests.NewTerminal(t)
	defer console.Cleanup()

	// Test interactive IO
	donec := make(chan struct{})
	//noinspection GoUnhandledErrorResult
	go func() {
		defer close(donec)
		console.ExpectString("Enter the current context name to confirm uninstallation of the Jenkins X platform from the ns namespace:")
		console.SendLine("correct-context-to-delete")
		console.ExpectEOF()
	}()

	commonOpts := opts.NewCommonOptionsWithFactory(clients_mocks.NewMockFactory())
	commonOpts.In = console.In
	commonOpts.Out = console.Out
	commonOpts.Err = console.Err
	o := &cmd.UninstallOptions{
		CommonOptions: &commonOpts,
		Namespace:     "ns",
	}

	o.SetKube(kubeMock)
	cmd_test_helpers.ConfigureTestOptions(o.CommonOptions, gits_test.NewMockGitter(), helm_test.NewMockHelmer())
	o.BatchMode = false // The above line sets batch mode to true. Set it back here :-(

	// Create fake namespace (that we will uninstall from)
	err = createNamespace(o, "ns")

	// Run the uninstall
	err = o.Run()
	assert.NoError(t, err)

	// Assert that the namespace has been deleted
	client, err := o.KubeClient()
	assert.NoError(t, err)
	_, err = client.CoreV1().Namespaces().Get("ns", metav1.GetOptions{})
	assert.Error(t, err)

	console.Close()
	<-donec

	// Dump the terminal's screen.
	t.Logf(expect.StripTrailingEmptyLines(console.CurrentState()))
}

func createNamespace(o *cmd.UninstallOptions, ns string) error {
	client, err := o.KubeClient()
	if err != nil {
		return err
	}
	_, err = client.CoreV1().Namespaces().Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	})
	return err
}

package cmd

import (
	"github.com/hashicorp/go-version"
	"github.com/jenkins-x/jx/pkg/cmd/create"
	"github.com/jenkins-x/jx/pkg/cmd/deletecmd"
	"github.com/jenkins-x/jx/pkg/cmd/helper"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/pkg/errors"

	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/jenkins-x/jx/pkg/cmd/templates"
	"github.com/jenkins-x/jx/pkg/kube"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	upgradeAddonProwLong = templates.LongDesc(`
		Upgrades the Jenkins X platform if there is a newer release
`)

	upgradeAddonProwExample = templates.Examples(`
		# Upgrades the Jenkins X platform 
		jx upgrade addon prow
	`)
)

// UpgradeAddonProwOptions the options for the create spring command
type UpgradeAddonProwOptions struct {
	UpgradeAddonsOptions

	newKnativeBuildVersion string
	Tekton                 bool
	ExternalDNS            bool
}

// NewCmdUpgradeAddonProw defines the command
func NewCmdUpgradeAddonProw(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &UpgradeAddonProwOptions{
		UpgradeAddonsOptions: UpgradeAddonsOptions{
			CreateOptions: create.CreateOptions{
				CommonOptions: commonOpts,
			},
		},
	}

	cmd := &cobra.Command{
		Use:     "prow",
		Short:   "Upgrades any AddonProw added to Jenkins X if there are any new releases available",
		Aliases: []string{"addon"},
		Long:    upgradeAddonProwLong,
		Example: upgradeAddonProwExample,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}

	options.UpgradeAddonsOptions.addFlags(cmd)
	options.InstallFlags.AddCloudEnvOptions(cmd)
	cmd.Flags().StringVarP(&options.newKnativeBuildVersion, "new-knative-build-version", "", "0.1.1", "The new kanative build verion that prow needs to work with")
	cmd.Flags().BoolVarP(&options.Tekton, "tekton", "t", true, "Enables Knative Build Pipeline. Otherwise we default to use Knative Build")
	cmd.Flags().BoolVarP(&options.ExternalDNS, "external-dns", "", true, "Installs external-dns into the cluster. ExternalDNS manages service DNS records for your cluster, providing you've setup your domain record")
	return cmd
}

// Run implements the command
func (o *UpgradeAddonProwOptions) Run() error {
	err := o.Helm().UpdateRepo()
	if err != nil {
		return err
	}
	ns := o.Namespace
	if ns == "" {
		_, ns, err = o.JXClientAndDevNamespace()
		if err != nil {
			return err
		}
	}

	kubeClient, _, err := o.KubeClientAndDevNamespace()
	if err != nil {
		return err
	}

	releases, _, err := o.Helm().ListReleases(ns)
	if err != nil {
		return err
	}

	newKnativeBuildVersion, err := version.NewVersion(o.newKnativeBuildVersion)
	if err != nil {
		return err
	}

	// first lets get the existing hmac and oauth tokens so we can use them when reinstalling
	oauthSecret, err := kubeClient.CoreV1().Secrets(ns).Get("oauth-token", metav1.GetOptions{})
	if err != nil {
		return err
	}
	oauthToken := string(oauthSecret.Data["oauth"])

	hmacSecret, err := kubeClient.CoreV1().Secrets(ns).Get("hmac-token", metav1.GetOptions{})
	if err != nil {
		return err
	}
	hmacToken := string(hmacSecret.Data["hmac"])

	for _, release := range releases {
		if release.ReleaseName == "knative-build" && (release.Status == "DEPLOYED" || release.Status == "FAILED") {
			currentVersion, err := version.NewVersion(release.ChartVersion)
			if err != nil {
				return err
			}
			// if we currently have less than 0.2.x of Knative build chart we need to reinstall as there's issues with
			// the CRDs when upgrading
			if currentVersion.LessThan(newKnativeBuildVersion) {
				message := "The version of Knative Build you are running is too old to support the latest Prow, would " +
					"you like to install the latest Knative Build?\nWARNING: this will remove the previous version and " +
					"install the latest, any existing builds or custom changes to BuildTemplate resources will be lost"

				if !util.Confirm(message, false, "", o.In, o.Out, o.Err) {
					return nil
				}

				// delete knative build
				deleteKnativeBuildOpts := &deletecmd.DeleteKnativeBuildOptions{
					DeleteAddonOptions: deletecmd.DeleteAddonOptions{
						CommonOptions: o.CommonOptions,
					},
				}
				deleteKnativeBuildOpts.ReleaseName = kube.DefaultKnativeBuildReleaseName

				err = deleteKnativeBuildOpts.Run()
				if err != nil {
					return err
				}

			}
		}
	}
	// now let's reinstall prow
	err = o.DeleteChart("jx-prow", true)
	if err != nil {
		return err
	}

	o.OAUTHToken = oauthToken
	o.HMACToken = hmacToken
	isGitOps, _ := o.GetDevEnv()

	_, pipelineUser, err := o.GetPipelineGitAuth()
	if err != nil {
		return errors.Wrap(err, "retrieving the pipeline Git Auth")
	}
	pipelineUserName := ""
	if pipelineUser != nil {
		pipelineUserName = pipelineUser.Username
	}

	return o.InstallProw(o.Tekton, o.ExternalDNS, isGitOps, "", "", pipelineUserName, nil)
}

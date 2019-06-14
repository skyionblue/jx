package deletecmd

import (
	"github.com/jenkins-x/jx/pkg/cmd/create"
	"github.com/jenkins-x/jx/pkg/cmd/helper"
	"github.com/jenkins-x/jx/pkg/log"

	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/jenkins-x/jx/pkg/cmd/templates"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/spf13/cobra"
)

var (
	deleteAddonCloudBeesLong = templates.LongDesc(`
		Deletes the CloudBees addon
`)

	deleteAddonCloudBeesExample = templates.Examples(`
		# Deletes the CloudBees addon
		jx delete addon cloudbees
	`)
)

// DeleteAddonGiteaOptions the options for the create spring command
type DeleteAddoncoreOptions struct {
	DeleteAddonOptions

	ReleaseName string
}

// NewCmdDeleteAddonCloudBees defines the command
func NewCmdDeleteAddonCloudBees(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &DeleteAddoncoreOptions{
		DeleteAddonOptions: DeleteAddonOptions{
			CommonOptions: commonOpts,
		},
	}

	cmd := &cobra.Command{
		Use:     "cloudbees",
		Short:   "Deletes the CloudBees app for Kubernetes addon",
		Aliases: []string{"cloudbee", "cb", "core"},
		Long:    deleteAddonCloudBeesLong,
		Example: deleteAddonCloudBeesExample,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}

	cmd.Flags().StringVarP(&options.ReleaseName, opts.OptionRelease, "r", create.DefaultCloudBeesReleaseName, "The chart release name")
	options.addFlags(cmd)
	return cmd
}

// Run implements the command
func (o *DeleteAddoncoreOptions) Run() error {
	if o.ReleaseName == "" {
		return util.MissingOption(opts.OptionRelease)
	}
	err := o.DeleteChart(o.ReleaseName, o.Purge)
	if err != nil {
		return err
	}
	log.Logger().Infof("Addon %s deleted successfully", util.ColorInfo(o.ReleaseName))

	return nil
}

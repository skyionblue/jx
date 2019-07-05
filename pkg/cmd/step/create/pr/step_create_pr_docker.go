package pr

import (
	"github.com/jenkins-x/jx/pkg/docker"
	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/log"
	"github.com/pkg/errors"

	"github.com/jenkins-x/jx/pkg/cmd/helper"

	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/jenkins-x/jx/pkg/cmd/templates"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/spf13/cobra"
)

var (
	createPullRequestDockerLong = templates.LongDesc(`
		Creates a Pull Request on a git repository updating any lines in the Dockerfile that start with FROM, ENV or ARG=
`)

	createPullRequestDockerExample = templates.Examples(`
					`)
)

// StepCreatePullRequestDockersOptions contains the command line flags
type StepCreatePullRequestDockersOptions struct {
	StepCreatePrOptions

	Name string
}

// NewCmdStepCreatePullRequestDocker Creates a new Command object
func NewCmdStepCreatePullRequestDocker(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &StepCreatePullRequestDockersOptions{
		StepCreatePrOptions: StepCreatePrOptions{
			StepCreateOptions: opts.StepCreateOptions{
				StepOptions: opts.StepOptions{
					CommonOptions: commonOpts,
				},
			},
		},
	}

	cmd := &cobra.Command{
		Use:     "docker",
		Short:   "Creates a Pull Request on a git repository updating the docker file",
		Long:    createPullRequestDockerLong,
		Example: createPullRequestDockerExample,
		Aliases: []string{"version pullrequest"},
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}
	AddStepCreatePrFlags(cmd, &options.StepCreatePrOptions)
	cmd.Flags().StringVarP(&options.Name, "name", "n", "", "The name of the property to update")
	return cmd
}

// Run implements this command
func (o *StepCreatePullRequestDockersOptions) Run() error {
	if err := o.ValidateOptions(); err != nil {
		return errors.WithStack(err)
	}
	if o.Name == "" {
		return util.MissingOption("name")
	}
	if o.Version == "" {
		return util.MissingOption("version")
	}
	if o.SrcGitURL == "" {
		log.Logger().Warnf("srcRepo is not provided so generated PR will not be correctly linked in release notesPR")
	}
	err := o.CreatePullRequest("docker",
		func(dir string, gitInfo *gits.GitRepository) ([]string, error) {
			oldVersions, err := docker.UpdateVersions(dir, o.Version, o.Name)
			if err != nil {
				return nil, errors.Wrapf(err, "updating %s to %s", o.Name, o.Version)
			}
			return oldVersions, nil
		})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

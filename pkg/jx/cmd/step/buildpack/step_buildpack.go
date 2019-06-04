package buildpack

import (
	"github.com/jenkins-x/jx/pkg/jx/cmd/helper"
	"github.com/jenkins-x/jx/pkg/jx/cmd/opts"
	"github.com/spf13/cobra"
)

// StepBuildPackOptions contains the command line flags
type StepBuildPackOptions struct {
	opts.StepOptions
}

// NewCmdStepBuildPack Steps a command object for the "step" command
func NewCmdStepBuildPack(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &StepBuildPackOptions{
		StepOptions: opts.StepOptions{
			CommonOptions: commonOpts,
		},
	}

	cmd := &cobra.Command{
		Use:   "buildpack",
		Short: "buildpack [command]",
		Aliases: []string{
			"build pack", "pack", "bp",
		},
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}
	cmd.AddCommand(NewCmdStepBuildPackApply(commonOpts))
	return cmd
}

// Run implements this command
func (o *StepBuildPackOptions) Run() error {
	return o.Cmd.Help()
}

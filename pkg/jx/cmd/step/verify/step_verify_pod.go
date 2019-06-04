package verify

import (
	"github.com/jenkins-x/jx/pkg/jx/cmd/helper"
	"github.com/jenkins-x/jx/pkg/jx/cmd/opts"
	"github.com/spf13/cobra"
)

// StepVerifyPodOptions contains the command line flags
type StepVerifyPodOptions struct {
	opts.StepOptions
}

// NewCmdStepVerifyPod creates the `jx step verify pod` command
func NewCmdStepVerifyPod(commonOpts *opts.CommonOptions) *cobra.Command {

	options := &StepVerifyPodOptions{
		StepOptions: opts.StepOptions{
			CommonOptions: commonOpts,
		},
	}

	cmd := &cobra.Command{
		Use:   "pod",
		Short: "pod [command]",
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}
	cmd.AddCommand(NewCmdStepVerifyPodCount(commonOpts))
	cmd.AddCommand(NewCmdStepVerifyPodReady(commonOpts))
	return cmd
}

// Run implements this command
func (o *StepVerifyPodOptions) Run() error {
	return o.Cmd.Help()
}

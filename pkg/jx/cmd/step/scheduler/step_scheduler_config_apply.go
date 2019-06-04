package scheduler

import (
	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/jx/cmd/helper"
	"github.com/jenkins-x/jx/pkg/jx/cmd/opts"
	"github.com/jenkins-x/jx/pkg/jx/cmd/templates"
	"github.com/jenkins-x/jx/pkg/pipelinescheduler"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// StepSchedulerConfigApplyOptions contains the command line flags
type StepSchedulerConfigApplyOptions struct {
	opts.StepOptions
	Agent string
	// allow git to be configured externally before a PR is created
	ConfigureGitCallback gits.ConfigureGitFn
}

var (
	stepSchedulerConfigApplyLong = templates.LongDesc(`
        This command will transform your pipeline schedulers in to prow config. 
        If you are using gitops the prow config will be added to your environment repository. 
        For non-gitops environments the prow config maps will applied to your dev environment.
`)
	stepSchedulerConfigApplyExample = templates.Examples(`
	
	jx step scheduler config apply
`)
)

// NewCmdStepSchedulerConfigApply Steps a command object for the "step" command
func NewCmdStepSchedulerConfigApply(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &StepSchedulerConfigApplyOptions{
		StepOptions: opts.StepOptions{
			CommonOptions: commonOpts,
		},
	}

	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "scheduler config apply",
		Long:    stepSchedulerConfigApplyLong,
		Example: stepSchedulerConfigApplyExample,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}
	options.AddCommonFlags(cmd)
	cmd.Flags().StringVarP(&options.Agent, "agent", "", "prow", "The scheduler agent to use e.g. Prow")
	return cmd
}

// Run implements this command
func (o *StepSchedulerConfigApplyOptions) Run() error {
	gitOps, devEnv := o.GetDevEnv()
	switch o.Agent {
	case "prow":
		jxClient, ns, err := o.JXClient()
		if err != nil {
			return errors.WithStack(err)
		}
		teamSettings, err := o.TeamSettings()
		if err != nil {
			return err
		}
		cfg, plugs, err := pipelinescheduler.GenerateProw(gitOps, jxClient, ns, teamSettings.DefaultScheduler.Name, devEnv)
		if err != nil {
			return errors.Wrapf(err, "generating Prow config")
		}
		if gitOps {
			opts := pipelinescheduler.GitOpsOptions{
				Verbose: o.Verbose,
				DevEnv:  devEnv,
			}
			environmentsDir, err := o.EnvironmentsDir()
			if err != nil {
				return errors.Wrapf(err, "getting environments dir")
			}
			opts.EnvironmentsDir = environmentsDir

			gitProvider, _, err := o.CreateGitProviderForURLWithoutKind(devEnv.Spec.Source.URL)
			if err != nil {
				return errors.Wrapf(err, "creating git provider for %s", devEnv.Spec.Source.URL)
			}
			opts.GitProvider = gitProvider
			opts.ConfigureGitFn = o.ConfigureGitCallback
			opts.Gitter = o.Git()
			opts.Helmer = o.Helm()
			err = opts.AddToEnvironmentRepo(cfg, plugs)
			if err != nil {
				return errors.Wrapf(err, "adding Prow config to environment repo")
			}
		} else {
			kubeClient, ns, err := o.KubeClientAndNamespace()
			if err != nil {
				return errors.WithStack(err)
			}
			err = pipelinescheduler.ApplyDirectly(kubeClient, ns, cfg, plugs)
			if err != nil {
				return errors.Wrapf(err, "applying Prow config")
			}
		}
	default:
		return errors.Errorf("%s is an unsupported agent. Available agents are: prow", o.Agent)
	}
	return nil
}

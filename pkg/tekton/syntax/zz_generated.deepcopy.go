// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package syntax

import (
	v1 "k8s.io/api/core/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Agent) DeepCopyInto(out *Agent) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Agent.
func (in *Agent) DeepCopy() *Agent {
	if in == nil {
		return nil
	}
	out := new(Agent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Loop) DeepCopyInto(out *Loop) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]Step, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Loop.
func (in *Loop) DeepCopy() *Loop {
	if in == nil {
		return nil
	}
	out := new(Loop)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParsedPipeline) DeepCopyInto(out *ParsedPipeline) {
	*out = *in
	if in.Agent != nil {
		in, out := &in.Agent, &out.Agent
		if *in == nil {
			*out = nil
		} else {
			*out = new(Agent)
			**out = **in
		}
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		if *in == nil {
			*out = nil
		} else {
			*out = new(RootOptions)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Stages != nil {
		in, out := &in.Stages, &out.Stages
		*out = make([]Stage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Post != nil {
		in, out := &in.Post, &out.Post
		*out = make([]Post, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.WorkingDir != nil {
		in, out := &in.WorkingDir, &out.WorkingDir
		if *in == nil {
			*out = nil
		} else {
			*out = new(string)
			**out = **in
		}
	}
	if in.Environment != nil {
		in, out := &in.Environment, &out.Environment
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParsedPipeline.
func (in *ParsedPipeline) DeepCopy() *ParsedPipeline {
	if in == nil {
		return nil
	}
	out := new(ParsedPipeline)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineOverride) DeepCopyInto(out *PipelineOverride) {
	*out = *in
	if in.Step != nil {
		in, out := &in.Step, &out.Step
		if *in == nil {
			*out = nil
		} else {
			*out = new(Step)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]*Step, len(*in))
		for i := range *in {
			if (*in)[i] == nil {
				(*out)[i] = nil
			} else {
				(*out)[i] = new(Step)
				(*in)[i].DeepCopyInto((*out)[i])
			}
		}
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		if *in == nil {
			*out = nil
		} else {
			*out = new(StepOverrideType)
			**out = **in
		}
	}
	if in.Agent != nil {
		in, out := &in.Agent, &out.Agent
		if *in == nil {
			*out = nil
		} else {
			*out = new(Agent)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineOverride.
func (in *PipelineOverride) DeepCopy() *PipelineOverride {
	if in == nil {
		return nil
	}
	out := new(PipelineOverride)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Post) DeepCopyInto(out *Post) {
	*out = *in
	if in.Actions != nil {
		in, out := &in.Actions, &out.Actions
		*out = make([]PostAction, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Post.
func (in *Post) DeepCopy() *Post {
	if in == nil {
		return nil
	}
	out := new(Post)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostAction) DeepCopyInto(out *PostAction) {
	*out = *in
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostAction.
func (in *PostAction) DeepCopy() *PostAction {
	if in == nil {
		return nil
	}
	out := new(PostAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RootOptions) DeepCopyInto(out *RootOptions) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		if *in == nil {
			*out = nil
		} else {
			*out = new(Timeout)
			**out = **in
		}
	}
	if in.ContainerOptions != nil {
		in, out := &in.ContainerOptions, &out.ContainerOptions
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.Container)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RootOptions.
func (in *RootOptions) DeepCopy() *RootOptions {
	if in == nil {
		return nil
	}
	out := new(RootOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Stage) DeepCopyInto(out *Stage) {
	*out = *in
	if in.Agent != nil {
		in, out := &in.Agent, &out.Agent
		if *in == nil {
			*out = nil
		} else {
			*out = new(Agent)
			**out = **in
		}
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		if *in == nil {
			*out = nil
		} else {
			*out = new(StageOptions)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]Step, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Stages != nil {
		in, out := &in.Stages, &out.Stages
		*out = make([]Stage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Parallel != nil {
		in, out := &in.Parallel, &out.Parallel
		*out = make([]Stage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Post != nil {
		in, out := &in.Post, &out.Post
		*out = make([]Post, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.WorkingDir != nil {
		in, out := &in.WorkingDir, &out.WorkingDir
		if *in == nil {
			*out = nil
		} else {
			*out = new(string)
			**out = **in
		}
	}
	if in.Environment != nil {
		in, out := &in.Environment, &out.Environment
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Stage.
func (in *Stage) DeepCopy() *Stage {
	if in == nil {
		return nil
	}
	out := new(Stage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageOptions) DeepCopyInto(out *StageOptions) {
	*out = *in
	if in.RootOptions != nil {
		in, out := &in.RootOptions, &out.RootOptions
		if *in == nil {
			*out = nil
		} else {
			*out = new(RootOptions)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Stash != nil {
		in, out := &in.Stash, &out.Stash
		if *in == nil {
			*out = nil
		} else {
			*out = new(Stash)
			**out = **in
		}
	}
	if in.Unstash != nil {
		in, out := &in.Unstash, &out.Unstash
		if *in == nil {
			*out = nil
		} else {
			*out = new(Unstash)
			**out = **in
		}
	}
	if in.Workspace != nil {
		in, out := &in.Workspace, &out.Workspace
		if *in == nil {
			*out = nil
		} else {
			*out = new(string)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageOptions.
func (in *StageOptions) DeepCopy() *StageOptions {
	if in == nil {
		return nil
	}
	out := new(StageOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Stash) DeepCopyInto(out *Stash) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Stash.
func (in *Stash) DeepCopy() *Stash {
	if in == nil {
		return nil
	}
	out := new(Stash)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Step) DeepCopyInto(out *Step) {
	*out = *in
	if in.Arguments != nil {
		in, out := &in.Arguments, &out.Arguments
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Loop != nil {
		in, out := &in.Loop, &out.Loop
		if *in == nil {
			*out = nil
		} else {
			*out = new(Loop)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Agent != nil {
		in, out := &in.Agent, &out.Agent
		if *in == nil {
			*out = nil
		} else {
			*out = new(Agent)
			**out = **in
		}
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]*Step, len(*in))
		for i := range *in {
			if (*in)[i] == nil {
				(*out)[i] = nil
			} else {
				(*out)[i] = new(Step)
				(*in)[i].DeepCopyInto((*out)[i])
			}
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Step.
func (in *Step) DeepCopy() *Step {
	if in == nil {
		return nil
	}
	out := new(Step)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Timeout) DeepCopyInto(out *Timeout) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Timeout.
func (in *Timeout) DeepCopy() *Timeout {
	if in == nil {
		return nil
	}
	out := new(Timeout)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Unstash) DeepCopyInto(out *Unstash) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Unstash.
func (in *Unstash) DeepCopy() *Unstash {
	if in == nil {
		return nil
	}
	out := new(Unstash)
	in.DeepCopyInto(out)
	return out
}
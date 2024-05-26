/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package installcmd

import (
	"context"

	//docs "github.com/pkgserver-dev/pkgserver/internal/docs/generated/initdocs"

	"github.com/kubenet-dev/kubenetctl/pkg/run"
	"github.com/spf13/cobra"
)

func NewCommand(ctx context.Context, version string) *cobra.Command {
	return NewRunner(ctx, version).Command
}

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, version string) *Runner {
	r := &Runner{}
	cmd := &cobra.Command{
		Use:  "install [flags]",
		Args: cobra.ExactArgs(0),
		//Short:   docs.InitShort,
		//Long:    docs.InitShort + "\n" + docs.InitLong,
		//Example: docs.InitExamples,
		PreRunE: r.preRunE,
		RunE:    r.runE,
	}

	r.Command = cmd

	return r
}

type Runner struct {
	Command *cobra.Command
}

func (r *Runner) preRunE(_ *cobra.Command, _ []string) error {
	return nil
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	ctx := c.Context()
	//log := log.FromContext(ctx)
	//log.Info("create packagerevision", "src", args[0], "dst", args[1])

	x := run.NewRun("Install kubenet Components")

	x.Step(
		run.S("install package server: (tool to interact with git from k8s using packages (KRM manifests))"),
		run.S("kubectl apply -f https://raw.githubusercontent.com/kubenet-dev/kubenet/v0.0.1/artifacts/out/pkgserver.yaml"),	
	)

	x.Step(
		run.S("install sdc: (tool to interact with yang devices from k8s)"),
		run.S("kubectl apply -f https://raw.githubusercontent.com/kubenet-dev/kubenet/v0.0.1/artifacts/out/sdc.yaml"),
	)

	x.Step(
		run.S("install kuid-server: (tool for inventory and identity (IPAM/VLAN/AS/etc) using k8s api"),
		run.S("kubectl apply -f https://raw.githubusercontent.com/kubenet-dev/kubenet/v0.0.1/artifacts/out/kuid-server.yaml"),
	)

	x.Step(
		run.S("install kuid-apps: (apps leveraging kuid-server focussed on networking"),
		run.S("kubectl apply -f https://raw.githubusercontent.com/kubenet-dev/kubenet/v0.0.1/artifacts/out/kuidapps.yaml"),

	)

	x.Step(
		run.S("install kuid-nokia-srl: (vendor specific app for specific nokia srl artifacts "),
		run.S("kubectl apply -f https://raw.githubusercontent.com/kubenet-dev/kubenet/v0.0.1/artifacts/out/kuid-nokia-srl.yaml"),

	)

	return x.Run(ctx)
}

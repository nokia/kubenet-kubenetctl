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

package sdccmd

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
		Use:  "sdc [flags]",
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

	x := run.NewRun("Configue sdc")

	x.Step(
		run.S("apply the schema for srlinux 24.3.2"),
		run.S("kubectl apply -f https://raw.githubusercontent.com/kubenet-dev/kubenet/main/sdc/schemas/srl24-3-2.yaml"),	
	)

	x.Step(
		run.S("apply the gnmi profile to connect to the target (clab node)"),
		run.S("kubectl apply -f https://raw.githubusercontent.com/kubenet-dev/kubenet/main/sdc/profiles/conn-gnmi-skipverify.yaml"),
	)

	x.Step(
		run.S("apply the gnmi sync profile to sync config from the target (clab node)"),
		run.S("kubectl apply -f https://raw.githubusercontent.com/kubenet-dev/kubenet/main/sdc/profiles/sync-gnmi-get.yaml"),
	)

	x.Step(
		run.S("apply the srl secret with credentials to authenticate to the target (clab node)"),
		run.S("kubectl apply -f https://raw.githubusercontent.com/kubenet-dev/kubenet/main/sdc/profiles/secret.yaml"),
	)

	x.Step(
		run.S("apply the discovery rule to discover the srl devices deployed by containerlab"),
		run.S("kubectl apply -f https://raw.githubusercontent.com/kubenet-dev/kubenet/main/sdc/drrules/dr-dynamic.yaml"),
	)

	return x.Run(ctx)
}

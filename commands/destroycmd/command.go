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

package destroycmd

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
		Use:  "destroy [flags]",
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

	x := run.NewRun("Destroy kubenet Environment")

	x.Step(
		run.S("Drop the iptables rule"),
		run.S("sudo iptables -D DOCKER-USER -o br-$(docker network inspect -f '{{ printf \"%.12s\" .ID }}' kind) -j ACCEPT"),
	)

	x.Step(
		run.S("Delete the kind cluster"),
		run.S("kind delete cluster --name kubenet"),
	)

	x.Step(
		run.S("Destroy Containerlab topology"),
		run.S("sudo containerlab destroy -t https://github.com/kubenet-dev/kubenet/blob/main/lab/3node.yaml --reconfigure"),
	)
	
	return x.Run(ctx)
}

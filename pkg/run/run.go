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

package run

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gookit/color"
)

type Run struct {
	title       string
	description []string
	steps       []step
	out         io.Writer
	setup       func() error
	cleanup     func() error
	options     *Options
}

// Options specify the run options.
type Options struct {
	AutoTimeout      time.Duration
	Auto             bool
	BreakPoint       bool
	ContinueOnError  bool
	HideDescriptions bool
	DryRun           bool
	NoColor          bool
	Immediate        bool
	SkipSteps        int
	Shell            string
}

func emptyFn() error { return nil }

// NewRun creates a new run for the provided description string.
func NewRun(title string, description ...string) *Run {
	return &Run{
		title:       title,
		description: description,
		steps:       nil,
		out:         os.Stdout,
		setup:       emptyFn,
		cleanup:     emptyFn,
		options:     &Options{},
	}
}

func (r *Run) Step(text, command []string) {
	r.steps = append(r.steps, step{r, text, command, false, false})
}

func (r *Run) Run(ctx context.Context) error {
	if r.options.Shell == "" {
		r.options.Shell = "bash"
	}

	//r.options.Auto = getContextValue[bool](ctx, CtxKeyAutomatic)
	r.options.Auto = true // always run in automatic mode

	if err := r.setup(); err != nil {
		return err
	}

	if err := r.printTitleAndDescription(); err != nil {
		return err
	}

	for i, step := range r.steps {

		if err := step.run(i+1, len(r.steps)); err != nil {
			return err
		}

	}

	return r.cleanup()
}

func (r *Run) printTitleAndDescription() error {
	p := color.Cyan.Sprintf
	if err := write(r.out, p("%s\n", r.title)); err != nil {
		return err
	}
	for range r.title {
		if err := write(r.out, p("=")); err != nil {
			return err
		}
	}
	if err := write(r.out, "\n"); err != nil {
		return err
	}

	return nil
}

func write(w io.Writer, str string) error {
	_, err := w.Write([]byte(str))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

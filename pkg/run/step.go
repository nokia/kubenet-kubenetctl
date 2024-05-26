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
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gookit/color"
)

type step struct {
	r                     *Run
	text, command         []string
	canFail, isBreakPoint bool
}

func (s *step) run(current, max int) error {
	if err := s.waitOrSleep(); err != nil {
		return fmt.Errorf("unable to run step: %v: %w", s, err)
	}
	if len(s.text) > 0 && !s.r.options.HideDescriptions {
		s.echo(current, max)
	}
	if s.isBreakPoint {
		return s.wait()
	}
	if len(s.command) > 0 {
		return s.execute()
	}

	return nil
}

func (s *step) waitOrSleep() error {
	if s.r.options.Auto {
		time.Sleep(s.r.options.AutoTimeout)
	} else {
		if err := write(s.r.out, "…"); err != nil {
			return err
		}
		_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			return fmt.Errorf("unable to read newline: %w", err)
		}
		// Move cursor up again
		if err := write(s.r.out, "\x1b[1A"); err != nil {
			return err
		}
	}

	return nil
}

func (s *step) wait() error {
	if !s.r.options.BreakPoint {
		return nil
	}

	if err := write(s.r.out, "bp"); err != nil {
		return err
	}
	_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("unable to read newline: %w", err)
	}
	// Move cursor up again
	if err := write(s.r.out, "\x1b[1A"); err != nil {
		return err
	}

	return nil
}

func (s *step) echo(current, max int) {
	p := color.White.Darken().Sprintf
	if s.r.options.NoColor {
		p = fmt.Sprintf
	}

	prepared := []string{}
	for i, x := range s.text {
		if i == len(s.text)-1 {
			colon := ":"
			if s.command == nil {
				// Do not set the expectation that there is more if no command
				// provided.
				colon = ""
			}
			prepared = append(
				prepared,
				p(
					"# %s [%d/%d]%s\n",
					x, current, max, colon,
				),
			)
		} else {
			m := p("# %s", x)
			prepared = append(prepared, m)
		}
	}
	s.print(prepared...)
}

func (s *step) print(msg ...string) error {
	for _, m := range msg {
		for _, c := range m {
			if !s.r.options.Immediate {
				//nolint:gosec,gomnd // the sleep has no security implications and is randomly chosen
				time.Sleep(time.Duration(rand.Intn(40)) * time.Millisecond)
			}
			if err := write(s.r.out, fmt.Sprintf("%c", c)); err != nil {
				return err
			}
		}
		if err := write(s.r.out, "\n"); err != nil {
			return err
		}
	}

	return nil
}

func (s *step) execute() error {
	joinedCommand := strings.Join(s.command, " ")
	cmd := exec.Command(s.r.options.Shell, "-c", joinedCommand) //nolint:gosec // we purposefully run user-provided code

	cmd.Stderr = s.r.out
	cmd.Stdout = s.r.out

	p := color.Green.Sprintf
	if s.r.options.NoColor {
		p = fmt.Sprintf
	}

	cmdString := p("> %s", strings.Join(s.command, " \\\n    "))
	s.print(cmdString)
	if err := s.waitOrSleep(); err != nil {
		return fmt.Errorf("unable to execute step: %v: %w", s, err)
	}
	if s.r.options.DryRun {
		return nil
	}
	err := cmd.Run()
	if s.canFail {
		return nil
	}
	s.print("")

	if err != nil {
		return fmt.Errorf("step command failed: %w", err)
	}

	return nil
}

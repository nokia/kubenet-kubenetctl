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

package commands

import (
	"context"
	"os"
	"os/signal"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/kubenet-dev/kubenetctl/commands/destroycmd"
	"github.com/kubenet-dev/kubenetctl/commands/installcmd"
	"github.com/kubenet-dev/kubenetctl/commands/invcmd"
	"github.com/kubenet-dev/kubenetctl/commands/networkbridgedcmd"
	"github.com/kubenet-dev/kubenetctl/commands/networkconfigcmd"
	"github.com/kubenet-dev/kubenetctl/commands/networkdefaultcmd"
	"github.com/kubenet-dev/kubenetctl/commands/networkirbcmd"
	"github.com/kubenet-dev/kubenetctl/commands/networkroutedcmd"
	"github.com/kubenet-dev/kubenetctl/commands/sdccmd"
	"github.com/kubenet-dev/kubenetctl/commands/setupcmd"
	"github.com/kubenet-dev/kubenetctl/pkg/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigFileSubDir  = "kubenet"
	defaultConfigFileName    = "kubenet"
	defaultConfigFileNameExt = "yaml"
	defaultConfigEnvPrefix   = "KUBENETCTL"
	//defaultDBPath            = "package_db"
)

var (
	configFile string
)

func GetMain(ctx context.Context) *cobra.Command {
	//var auto bool
	var shell string
	//showVersion := false
	cmd := &cobra.Command{
		Use:          "kubenet",
		Short:        "kubenet is a cli tool for kubenet exercises",
		Long:         "kubenet is a cli tool for kubenet exercises",
		SilenceUsage: true,
		// We handle all errors in main after return from cobra so we can
		// adjust the error message coming from libraries
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// initialize viper settings
			ctx := cmd.Context()
			//ctx = context.WithValue(ctx, run.CtxKeyAutomatic, auto)
			ctx = context.WithValue(ctx, run.CtxKeyShell, shell)
			cmd.SetContext(ctx)
			initConfig()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := cmd.Flags().GetBool("help")
			if err != nil {
				return err
			}
			if h {
				return cmd.Help()
			}

			return cmd.Usage()
		},
	}

	//pf := cmd.PersistentFlags()
	// Catch interrupts and cleanup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			os.Exit(0)
		}
	}()

	// ensure the viper config directory exists
	cobra.CheckErr(os.MkdirAll(path.Join(xdg.ConfigHome, defaultConfigFileSubDir), 0700))
	// initialize viper settings
	initConfig()

	cmd.AddCommand(setupcmd.NewCommand(ctx, version))
	cmd.AddCommand(destroycmd.NewCommand(ctx, version))
	cmd.AddCommand(installcmd.NewCommand(ctx, version))
	cmd.AddCommand(sdccmd.NewCommand(ctx, version))
	cmd.AddCommand(invcmd.NewCommand(ctx, version))
	cmd.AddCommand(networkconfigcmd.NewCommand(ctx, version))
	cmd.AddCommand(networkdefaultcmd.NewCommand(ctx, version))
	cmd.AddCommand(networkbridgedcmd.NewCommand(ctx, version))
	cmd.AddCommand(networkroutedcmd.NewCommand(ctx, version))
	cmd.AddCommand(networkirbcmd.NewCommand(ctx, version))
	cmd.AddCommand(GetVersionCommand(ctx))
	//cmd.PersistentFlags().BoolVarP(&auto, "interactive", "i", true, "run in interacti mode")
	cmd.PersistentFlags().StringVar(&shell, "shell", "bash", "shell to be used to execute the commands")

	return cmd
}

type Runner struct {
	Command *cobra.Command
	//Ctx     context.Context
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {

		viper.AddConfigPath(filepath.Join(xdg.ConfigHome, defaultConfigFileSubDir))
		viper.SetConfigType(defaultConfigFileNameExt)
		viper.SetConfigName(defaultConfigFileName)

		_ = viper.SafeWriteConfig()
	}

	//viper.Set("kubecontext", kubecontext)
	//viper.Set("kubeconfig", kubeconfig)

	viper.SetEnvPrefix(defaultConfigEnvPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_ = 1
	}
}

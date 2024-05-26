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
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/kubenet-dev/kubenetctl/commands/setupcmd"
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

	// ensure the viper config directory exists
	cobra.CheckErr(os.MkdirAll(path.Join(xdg.ConfigHome, defaultConfigFileSubDir), 0700))
	// initialize viper settings
	initConfig()

	cmd.AddCommand(setupcmd.NewCommand(ctx, version))
	cmd.AddCommand(GetVersionCommand(ctx))
	//cmd.PersistentFlags().StringVar(&configFile, "config", "c", fmt.Sprintf("Default config file (%s/%s/%s.%s)", xdg.ConfigHome, defaultConfigFileSubDir, defaultConfigFileName, defaultConfigFileNameExt))

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

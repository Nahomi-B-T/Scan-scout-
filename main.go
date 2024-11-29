/*
   Copyright 2020 Docker Inc.

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

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func main() {
	_, closeFunc := newSigContext()
	defer closeFunc()
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		cmd := newScoutCmd(dockerCli)
		return cmd
	}, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "Docker Inc.",
		Version:       "1.0.0", // Cambiar según sea necesario
	})
}

type options struct {
	showVersion bool
}

func newScoutCmd(dockerCli command.Cli) *cobra.Command {
	var flags options
	cmd := &cobra.Command{
		Short:       "Docker Scout",
		Long:        "Analyze Docker images for vulnerabilities using Docker Scout", // Corregir la cadena de documentación
		Use:         "scout [OPTIONS] IMAGE",
		Annotations: map[string]string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.showVersion {
				return runVersion()
			}
			if len(args) < 1 {
				return fmt.Errorf("error: no image specified")
			}
			image := args[0]
			return runScoutAnalysis(dockerCli, image)
		},
	}
	cmd.Flags().BoolVar(&flags.showVersion, "version", false, "Display version of the scout plugin")

	return cmd
}

func runVersion() error {
	version := "Docker Scout Plugin Version: 1.0.0" // Cambiar según sea necesario
	fmt.Println(version)
	return nil
}

func runScoutAnalysis(dockerCli command.Cli, image string) error {
	// Ejecutar el comando de Docker Scout
	cmd := exec.Command("docker", "scout", "analyze", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Analyzing image '%s' with Docker Scout...\n", image)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running Docker Scout analysis: %w", err)
	}

	fmt.Println("Analysis complete.")
	return nil
}

func newSigContext() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-s
		cancel()
	}()
	return ctx, cancel
}

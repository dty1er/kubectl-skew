// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package ver

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func NewVerCmd() *cobra.Command {
	ver := &cobra.Command{
		Use:   "ver [options] [flags]",
		Short: "kubectl version which can update itself",
		// always pass flags to kubectl without check by cobra
		// to show the same error msg with kubectl one
		DisableFlagParsing: true,
		SilenceUsage:       true,
		Run:                RunKubectlVersion(),
	}
	return ver
}

func RunKubectlVersion() func(*cobra.Command, []string) {
	return func(c *cobra.Command, args []string) {
		args = append([]string{"version"}, args...)
		cmd := exec.Command("kubectl", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package main

import (
	"os"

	kubectlver "github.com/dty1er/kubectl-ver"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-skew", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := kubectlver.New(genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

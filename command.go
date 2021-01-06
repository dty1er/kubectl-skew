// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package ver

import (
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
}

func New(streams genericclioptions.IOStreams) *cobra.Command {
	o := &Options{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}

	skew := NewSkewCmd()

	o.configFlags.AddFlags(skew.Flags())

	return skew
}

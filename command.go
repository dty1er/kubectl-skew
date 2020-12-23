// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package ver

import (
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Options struct {
	configFlags *genericclioptions.ConfigFlags

	resultingContext     *api.Context
	resultingContextName string

	userSpecifiedCluster   string
	userSpecifiedContext   string
	userSpecifiedAuthInfo  string
	userSpecifiedNamespace string

	rawConfig      api.Config
	listNamespaces bool
	args           []string

	genericclioptions.IOStreams
}

func New(streams genericclioptions.IOStreams) *cobra.Command {
	o := &Options{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}

	ver := NewVerCmd()
	check := NewCheckCmd()
	skew := NewSkewCmd()
	install := NewInstallCmd()

	o.configFlags.AddFlags(ver.Flags())
	ver.AddCommand(check)
	ver.AddCommand(skew)
	ver.AddCommand(install)

	return ver
}

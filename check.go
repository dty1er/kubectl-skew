// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package ver

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	debugCheckClient string
	debugCheckServer string
)

func NewCheckCmd() *cobra.Command {
	skew := &cobra.Command{
		Use:   "check [flags]",
		Short: "checks kubectl update",
		RunE:  RunCheck(),
	}

	// flags for debug
	skew.Flags().StringVarP(&debugClient, "debug-client", "c", "", "param for debug: inject client version")
	skew.Flags().MarkHidden("debug-client")

	return skew
}

func RunCheck() func(c *cobra.Command, args []string) error {
	return func(c *cobra.Command, args []string) error {
		versions, err := InspectCurrentVersion()
		if err != nil {
			return err
		}

		latest, err := InspectLatestVersion()
		if err != nil {
			return err
		}

		template := "kubectl: %s\nlatest:  %s\n"

		fmt.Fprintf(os.Stdout, template, versions.Client, latest)

		if latest.Compare(versions.Client) != 0 {
			template = "kubectl update is available.\nYou can update it by \"kubectl ver install\" command. \nFor more details, see \"kubectl ver install --help\"\n"
			fmt.Fprintf(os.Stdout, yellow(template))
		} else {
			template = "kubectl is already up-to-date.\n"
			fmt.Fprintf(os.Stdout, green(template))
		}

		return nil
	}
}

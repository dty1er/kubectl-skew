// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package ver

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

var (
	targetVer = ""
)

func NewInstallCmd() *cobra.Command {
	install := &cobra.Command{
		Use:   "install [flags]",
		Short: "install specified kubectl version",
		RunE:  RunInstall(),
	}

	install.Flags().StringVarP(&targetVer, "version", "v", "", "a version string to install. \"latest\" is also acceptable")

	return install
}

func RunInstall() func(c *cobra.Command, args []string) error {
	return func(c *cobra.Command, args []string) error {
		versions, err := InspectCurrentVersion()
		if err != nil {
			return err
		}

		latest, err := InspectLatestVersion()
		if err != nil {
			return err
		}

		var target *semver.Version
		if targetVer == "" || targetVer == "latest" {
			target = latest
		} else {
			v, err := semver.NewVersion(targetVer)
			if err != nil {
				return err
			}

			target = v
		}

		// if err := checkVersion(latest, target, versions); err != nil {
		// 	return err
		// }
		_ = target
		_ = versions

		binURL := fmt.Sprintf(
			"https://dl.k8s.io/release/v%s/bin/%s/%s/kubectl",
			target.String(), runtime.GOOS, runtime.GOARCH,
		)
		resp, err := http.Get(binURL)
		if err != nil {
			return fmt.Errorf("failed to get binary: %w", err)
		}
		defer resp.Body.Close()

		// duplicate body to buffer to use twice
		var buff bytes.Buffer
		r := io.TeeReader(resp.Body, &buff)

		h := sha256.New()
		if _, err := io.Copy(h, r); err != nil {
			return err
		}

		checksumURL := fmt.Sprintf(
			"https://dl.k8s.io/v%s/bin/%s/%s/kubectl.sha256",
			target.String(), runtime.GOOS, runtime.GOARCH,
		)
		resp, err = http.Get(checksumURL)
		if err != nil {
			return fmt.Errorf("failed to get checksum: %w", err)
		}
		defer resp.Body.Close()

		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Compare checksum. If they are different,
		// installation must not go on
		given := strings.Trim(hex.EncodeToString(h.Sum(nil)), "\n")
		expected := strings.Trim(string(out), "\n")
		if given != expected {
			return fmt.Errorf("binary checksum did not match. binary installation URL: %s, checksum URL: %s", binURL, checksumURL)
		}

		// write binary body to the file then move into $PATH
		dir := os.TempDir()
		saveLocation := path.Join(dir, "kubectl")
		f, err := os.OpenFile(saveLocation, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}

		io.Copy(f, &buff)

		if err := os.Rename(saveLocation, "/usr/local/bin/kubectl"); err != nil {
			return err
		}

		return nil
	}
}

func checkVersion(latest, target *semver.Version, versions *Versions) error {
	skew := CalcKubeVerSkew(latest, versions.Server, target)
	if skew.ClientNeedsUpdate || skew.ClientNeedsDowngradeOrServerCanBeUpdated {
		fmt.Fprintf(os.Stdout, yellow("WARN: given version %s is unsupported according to your kubernetes cluster version and\nkubernetes version skew policy.\nFor more details, run \"kubectl ver skew\".\n"), target)
	}

	if target.Compare(versions.Client) == 0 {
		return fmt.Errorf("your kubectl is already %s. Installation is over.\n", target)
	}

	return nil
}

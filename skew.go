// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package skew

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

var (
	verTemplate = `cluster: %s
kubectl: %s
latest:  %s`

	resultTemplate = `Check result
  Server version: %s
  Client version: %s`

	serverTooOldTemplate = `Your kubernetes cluster version is unsupported.
There are %d minor version skew with the latest which must be within 2.`

	clientTooOldTemplate = `Your kubectl version is unsupported.
There are %d minor version skew with the server which must be between -1 and 1.`

	clientTooNewOrServerTooOldTemplate = `Your kubernetes cluster version is supported, but your kubectl version is too new. 
kubectl and kubernetes cluster version skew must be between -1 and 1, but it's %d.
You can update kubernetes cluster or downgrade kubectl to follow the version skew policy.`
)

type Versions struct {
	Client *semver.Version
	Server *semver.Version
}

type VersionSkew struct {
	ServerAndLatestDelta int
	ServerNeedsUpdate    bool

	ServerAndClientDelta int
	ClientNeedsUpdate    bool

	ClientNeedsDowngradeOrServerCanBeUpdated bool
}

var (
	debugClient string
	debugServer string
)

func NewSkewCmd() *cobra.Command {
	skew := &cobra.Command{
		Use:   "skew [flags]",
		Short: "checks kubernetes cluster and kubectl version skew",
		RunE: func(c *cobra.Command, args []string) error {
			return RunSkew(os.Stdout)
		},
	}

	// flags for debug
	skew.Flags().StringVarP(&debugClient, "debug-client", "", "", "param for debug: inject client version")
	skew.Flags().MarkHidden("debug-client")
	skew.Flags().StringVarP(&debugServer, "debug-server", "", "", "param for debug: inject server version")
	skew.Flags().MarkHidden("debug-server")

	return skew
}

func RunSkew(w io.Writer) error {
	versions, err := InspectCurrentVersion()
	if err != nil {
		return err
	}

	latest, err := InspectLatestVersion()
	if err != nil {
		return err
	}

	fmt.Fprintf(
		w, verTemplate+"\n",
		versions.Server, versions.Client, latest,
	)

	skew := CalcKubeVerSkew(latest, versions.Server, versions.Client)

	ok := green("OK")
	ng := red("NG")
	serverCheckResult := ok
	if skew.ServerNeedsUpdate {
		serverCheckResult = ng
	}

	clientCheckResult := ok
	if skew.ClientNeedsUpdate || skew.ClientNeedsDowngradeOrServerCanBeUpdated {
		clientCheckResult = ng
	}

	fmt.Fprintf(
		w, resultTemplate+"\n",
		serverCheckResult, clientCheckResult,
	)

	if skew.ServerNeedsUpdate {
		fmt.Fprintf(
			w, yellow(serverTooOldTemplate)+"\n",
			skew.ServerAndLatestDelta,
		)
	}
	if skew.ClientNeedsUpdate {
		fmt.Fprintf(
			w, yellow(clientTooOldTemplate)+"\n",
			skew.ServerAndClientDelta,
		)
		fmt.Println("")
	}
	if skew.ClientNeedsDowngradeOrServerCanBeUpdated {
		fmt.Fprintf(
			w, yellow(clientTooNewOrServerTooOldTemplate)+"\n",
			skew.ServerAndClientDelta,
		)
		fmt.Println("")
	}

	return nil
}

// InspectCurrentVersion runs "kubectl version --short" and parses the result
// to inspect the kubectl and kubernetes cluster version.
// This function highly depends on the kubectl implementation, so
// it might be broken when kubectl makes breaking changes.
var InspectCurrentVersion = func() (*Versions, error) {
	versions := &Versions{}

	cmd := exec.Command("kubectl", []string{"version", "--short"}...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	idx := 0
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		if idx >= 2 {
			return nil, fmt.Errorf("something wrong: kubectl version --short result is more than 3 lines")
		}

		// kubectl version --short format:
		// ------
		// Client Version: v1.20.2
		// Server Version: v1.20.0
		line := scanner.Text()
		splitted := strings.Split(line, ": ")
		if len(splitted) != 2 {
			return nil, fmt.Errorf("something wrong: kubectl version --short result has unexpected format")
		}

		switch idx {
		case 0:
			if splitted[0] != "Client Version" {
				return nil, fmt.Errorf("something wrong: kubectl version --short result first line is unexpected")
			}
			v, err := semver.NewVersion(splitted[1])
			if err != nil {
				return nil, err
			}
			versions.Client = v
		case 1:
			if splitted[0] != "Server Version" {
				return nil, fmt.Errorf("something wrong: kubectl version --short result second line is unexpected")
			}
			v, err := semver.NewVersion(splitted[1])
			if err != nil {
				return nil, err
			}
			versions.Server = v
		}

		idx++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if debugClient != "" {
		dcv, err := semver.NewVersion(debugClient)
		if err != nil {
			return nil, err
		}
		versions.Client = dcv
	}

	if debugServer != "" {
		dsv, err := semver.NewVersion(debugServer)
		if err != nil {
			return nil, err
		}
		versions.Server = dsv
	}

	return versions, nil
}

var InspectLatestVersion = func() (*semver.Version, error) {
	u := "https://dl.k8s.io/release/stable.txt"

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	v, err := semver.NewVersion(string(b))
	if err != nil {
		return nil, err
	}

	return v, nil
}

// CalcKubeVerSkewClient compares given 2 versions and
// check if there's too big skew.
// First argument must be server version and second must be client's.
// First return is minor ver diff, second one is if it's too much.
// It compares their minor versions and checks if there are 2 or more difference.
// This is following the kubernetes official version skew policy. For more details,
// see the official documentation.
// https://kubernetes.io/docs/setup/release/version-skew-policy/
func CalcKubeVerSkew(latest, server, client *semver.Version) *VersionSkew {
	vs := &VersionSkew{}

	vs.ServerAndLatestDelta = int(latest.Minor() - server.Minor())
	// server versions must be within 3 minor version compared to
	// the lastet one
	vs.ServerNeedsUpdate = 2 < vs.ServerAndLatestDelta

	vs.ServerAndClientDelta = int(server.Minor() - client.Minor())
	vs.ClientNeedsUpdate = 1 < vs.ServerAndClientDelta

	// If client version is too new && server version is ok, it comes here.
	// e.g. latest: 1.22.0, server: v1.20.0, client: 1.22.0
	// In this case, server version is ok but client version is not
	// because server and client minor ver must be within 1.
	// There are 2 recommendations: update server or downgrade client.
	if vs.ServerAndClientDelta < -1 && !vs.ServerNeedsUpdate {
		vs.ClientNeedsDowngradeOrServerCanBeUpdated = true
	}

	return vs
}

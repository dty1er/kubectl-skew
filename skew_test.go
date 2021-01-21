// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package skew

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Masterminds/semver"
)

func TestRunSkew(t *testing.T) {
	tests := []struct {
		name                      string
		inspectCurrentVersionMock func() (*Versions, error)
		inspectLatestVersionMock  func() (*semver.Version, error)
		want                      []string
		wantErr                   bool
	}{
		{
			name: "both server and client meets the policy",
			inspectCurrentVersionMock: func() (*Versions, error) {
				return &Versions{
					Server: semver.MustParse("v1.18.0"),
					Client: semver.MustParse("v1.17.0"),
				}, nil
			},
			inspectLatestVersionMock: func() (*semver.Version, error) {
				return semver.MustParse("v1.20.0"), nil
			},
			want: []string{
				fmt.Sprintf(verTemplate, "1.18.0", "1.17.0", "1.20.0"),
				fmt.Sprintf(resultTemplate, green("OK"), green("OK")),
			},
		},
		{
			name: "server is too old",
			inspectCurrentVersionMock: func() (*Versions, error) {
				return &Versions{
					Server: semver.MustParse("v1.17.0"),
					Client: semver.MustParse("v1.17.0"),
				}, nil
			},
			inspectLatestVersionMock: func() (*semver.Version, error) {
				return semver.MustParse("v1.20.0"), nil
			},
			want: []string{
				fmt.Sprintf(verTemplate, "1.17.0", "1.17.0", "1.20.0"),
				fmt.Sprintf(resultTemplate, red("NG"), green("OK")),
				fmt.Sprintf(serverTooOldTemplate, 3),
			},
		},
		{
			name: "server and client is too old",
			inspectCurrentVersionMock: func() (*Versions, error) {
				return &Versions{
					Server: semver.MustParse("v1.17.0"),
					Client: semver.MustParse("v1.15.0"),
				}, nil
			},
			inspectLatestVersionMock: func() (*semver.Version, error) {
				return semver.MustParse("v1.20.0"), nil
			},
			want: []string{
				fmt.Sprintf(verTemplate, "1.17.0", "1.15.0", "1.20.0"),
				fmt.Sprintf(resultTemplate, red("NG"), red("NG")),
				fmt.Sprintf(serverTooOldTemplate, 3),
				fmt.Sprintf(clientTooOldTemplate, 2),
			},
		},
		{
			name: "client is too old",
			inspectCurrentVersionMock: func() (*Versions, error) {
				return &Versions{
					Server: semver.MustParse("v1.18.0"),
					Client: semver.MustParse("v1.16.0"),
				}, nil
			},
			inspectLatestVersionMock: func() (*semver.Version, error) {
				return semver.MustParse("v1.20.0"), nil
			},
			want: []string{
				fmt.Sprintf(verTemplate, "1.18.0", "1.16.0", "1.20.0"),
				fmt.Sprintf(resultTemplate, green("OK"), red("NG")),
				fmt.Sprintf(clientTooOldTemplate, 2),
			},
		},
		{
			name: "server version is OK, but client is too new",
			inspectCurrentVersionMock: func() (*Versions, error) {
				return &Versions{
					Server: semver.MustParse("v1.18.0"),
					Client: semver.MustParse("v1.20.0"),
				}, nil
			},
			inspectLatestVersionMock: func() (*semver.Version, error) {
				return semver.MustParse("v1.20.0"), nil
			},
			want: []string{
				fmt.Sprintf(verTemplate, "1.18.0", "1.20.0", "1.20.0"),
				fmt.Sprintf(resultTemplate, green("OK"), red("NG")),
				fmt.Sprintf(clientTooNewOrServerTooOldTemplate, -2),
			},
		},
	}

	for _, c := range tests {
		c := c
		t.Run(c.name, func(t *testing.T) {
			InspectCurrentVersion = c.inspectCurrentVersionMock
			InspectLatestVersion = c.inspectLatestVersionMock

			var buff bytes.Buffer
			err := RunSkew(&buff)
			if (err != nil) != c.wantErr {
				t.Errorf("wantErr does not much")
			}

			got := buff.String()
			for _, w := range c.want {
				if !strings.Contains(got, w) {
					t.Errorf("got does not contain %s. got is %s", w, got)
				}
			}
		})
	}
}

func TestInspectLatestVersion(t *testing.T) {
	// Cannot verify what value is returned because
	// the latest version changes
	_, err := InspectLatestVersion()
	if err != nil {
		t.Errorf("cannot inspect latest version: %s", err)
	}
}

func TestCalcKubeVerSkew(t *testing.T) {
	tests := []struct {
		name                   string
		latest, server, client string
		want                   *VersionSkew
	}{
		{
			name:   "both server and client are ok",
			latest: "v1.20.2",
			server: "v1.18.2",
			client: "v1.17.2",
			want: &VersionSkew{
				ServerAndLatestDelta:                     2,
				ServerNeedsUpdate:                        false,
				ServerAndClientDelta:                     1,
				ClientNeedsUpdate:                        false,
				ClientNeedsDowngradeOrServerCanBeUpdated: false,
			},
		},
		{
			name:   "server is too old",
			latest: "v1.20.2",
			server: "v1.17.2",
			client: "v1.18.2",
			want: &VersionSkew{
				ServerAndLatestDelta:                     3,
				ServerNeedsUpdate:                        true,
				ServerAndClientDelta:                     -1,
				ClientNeedsUpdate:                        false,
				ClientNeedsDowngradeOrServerCanBeUpdated: false,
			},
		},
		{
			name:   "client is too old",
			latest: "v1.20.2",
			server: "v1.18.2",
			client: "v1.16.2",
			want: &VersionSkew{
				ServerAndLatestDelta:                     2,
				ServerNeedsUpdate:                        false,
				ServerAndClientDelta:                     2,
				ClientNeedsUpdate:                        true,
				ClientNeedsDowngradeOrServerCanBeUpdated: false,
			},
		},
		{
			name:   "client is too new",
			latest: "v1.20.2",
			server: "v1.18.2",
			client: "v1.20.2",
			want: &VersionSkew{
				ServerAndLatestDelta:                     2,
				ServerNeedsUpdate:                        false,
				ServerAndClientDelta:                     -2,
				ClientNeedsUpdate:                        false,
				ClientNeedsDowngradeOrServerCanBeUpdated: true,
			},
		},
	}

	for _, c := range tests {
		c := c
		t.Run(c.name, func(t *testing.T) {
			latest := semver.MustParse(c.latest)
			server := semver.MustParse(c.server)
			client := semver.MustParse(c.client)
			v := CalcKubeVerSkew(latest, server, client)

			if !reflect.DeepEqual(v, c.want) {
				t.Errorf("failed: got: %+v, want: %+v", v, c.want)
			}
		})
	}
}

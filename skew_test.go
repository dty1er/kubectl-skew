// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package skew

import (
	"reflect"
	"testing"

	"github.com/Masterminds/semver"
)

func TestInspectLatestVersion(t *testing.T) {
	// Cannot verify what value is returned because
	// the latest version changes
	_, err := InspectLatestVersion()
	if err != nil {
		t.Errorf("cannot inspect latest version: %s", err)
	}
}

func TestCalcKubeVerSkew(t *testing.T) {
	parseVer := func(t *testing.T, s string) *semver.Version {
		v, err := semver.NewVersion(s)
		if err != nil {
			t.Errorf("invalid version is given: %s", s)
		}

		return v
	}
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
			latest := parseVer(t, c.latest)
			server := parseVer(t, c.server)
			client := parseVer(t, c.client)
			v := CalcKubeVerSkew(latest, server, client)

			if !reflect.DeepEqual(v, c.want) {
				t.Errorf("failed: got: %+v, want: %+v", v, c.want)
			}
		})
	}
}

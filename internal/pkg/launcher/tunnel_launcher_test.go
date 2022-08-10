package launcher

import (
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTunnelLauncherSetObviousDefaults(t *testing.T) {
	testCases := map[string]string{
		config.LaunchRDP: utils.RDP,
		config.LaunchSSH: utils.SSH,
	}
	for launch, scheme := range testCases {
		t.Run(launch, func(t *testing.T) {
			params := config.FromValues(map[string]string{
				launch: "1",
			})
			tl, err := NewTunnelLauncher(params)
			require.NoError(t, err)
			assert.Equal(t, scheme, tl.Scheme)
		})
	}
}

func TestValidation(t *testing.T) {
	testCases := []struct {
		desc          string
		params        map[string]string
		expectedError string
	}{
		{
			desc: "bad rdp scheme",
			params: map[string]string{
				config.LaunchRDP: "1",
				config.Scheme:    "ssh",
			},
			expectedError: "launching the remote desktop client on scheme 'ssh' is not supported",
		},
		{
			desc: "bad ssh scheme",
			params: map[string]string{
				config.LaunchSSH: "-l root",
				config.Scheme:    "http",
			},
			expectedError: "launching the ssh client on scheme 'http' is not supported",
		},
		{
			desc: "unsupported schemes",
			params: map[string]string{
				config.LaunchURIHandler: "1",
				config.Scheme:           "svn",
			},
			expectedError: "scheme 'svn' is not supported to be handeled by a default app. supported schemes: vnc, http, https, realvnc",
		},
		{
			desc: "ssh port provided",
			params: map[string]string{
				config.LaunchSSH: "-p 2222",
			},
			expectedError: "do not pass a port with '-p': port will be set dynamically",
		},
		{
			desc: "too many launchers 1",
			params: map[string]string{
				config.LaunchRDP: "1",
				config.LaunchSSH: "-l root",
			},
			expectedError: "conflict: only one launch parameter of '--launch-uri, --launch-rdp, --launch-ssh' allowed",
		},
		{
			desc: "too many launchers 1",
			params: map[string]string{
				config.LaunchRDP:        "1",
				config.LaunchURIHandler: "1",
			},
			expectedError: "conflict: only one launch parameter of '--launch-uri, --launch-rdp, --launch-ssh' allowed",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			params := config.FromValues(tc.params)
			_, err := NewTunnelLauncher(params)
			t.Logf("Got error: %s", err)
			assert.EqualError(t, err, tc.expectedError)
		})
	}
}

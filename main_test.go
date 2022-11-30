package packages_metadata

import (
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

func Test_GetPackagesVersion(t *testing.T) {
	packagesNameVersion := map[string]string{
		"bitbucket.org/releaseband/buffalo-dale":       "v2.1.0",
		"bitbucket.org/releaseband/fortune-five-math":  "v0.1.0",
		"bitbucket.org/releaseband/fortune-three-math": "v0.4.7",
	}
	packagesAliases := map[string]string{
		"buffalo_dale_grandways": "bitbucket.org/releaseband/buffalo-dale",
		"fortune_five":           "bitbucket.org/releaseband/fortune-five-math",
		"fortune_three":          "bitbucket.org/releaseband/fortune-three-math",
		"fortune_three_xmas":     "bitbucket.org/releaseband/fortune-three-math",
	}
	res := GetPackagesMetadata(packagesNameVersion, packagesAliases)
	mp := res.GetMap()

	tests := []struct {
		name       string
		expVersion string
	}{
		{
			name:       "not required",
			expVersion: "",
		},
		{
			name:       "buffalo_dale_grandways",
			expVersion: "v2.1.0",
		},
		{
			name:       "fortune_five",
			expVersion: "v0.1.0",
		},
		{
			name:       "fortune_three",
			expVersion: "v0.4.7",
		},
		{
			name:       "fortune_three_xmas",
			expVersion: "v0.4.7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := res.GetVersion(tt.name)
			require.Equal(t, tt.expVersion, version)
			require.Equal(t, mp[tt.name], version)

		})
	}
}

func Test_GetPackagesMetadata(t *testing.T) {
	var reader io.Reader
	reader = strings.NewReader(`
bitbucket.org/releaseband/buffalo-dale v2.1.0
bitbucket.org/releaseband/fortune-five-math v0.1.0
`)
	packagesVersion, err := GetPackagesVersion(reader)
	require.Nil(t, err)

	tests := []struct {
		name       string
		isExist    bool
		expVersion string
	}{
		{
			name:       "not required",
			isExist:    false,
			expVersion: "",
		},
		{
			name:       "bitbucket.org/releaseband/buffalo-dale",
			isExist:    true,
			expVersion: "v2.1.0",
		},
		{
			name:       "bitbucket.org/releaseband/fortune-five-math",
			isExist:    true,
			expVersion: "v0.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, ok := packagesVersion[tt.name]
			require.Equal(t, tt.isExist, ok)
			require.Equal(t, tt.expVersion, version)

		})
	}
}

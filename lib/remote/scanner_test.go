package remote

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidatePlugin_Errors(t *testing.T) {
	invalidPluginAbsPath, err := filepath.Abs("testdata/invalid.so")
	if err != nil {
		assert.FailNow(t, "failed to get the absolute path of test file: %v", err)
	}

	testcases := []struct {
		name    string
		path    string
		wantErr string
	}{
		{
			name:    "test with invalid path",
			path:    "testdata/doesnotexist",
			wantErr: "given path testdata/doesnotexist does not exist",
		},
		{
			name:    "test with invalid plugin",
			path:    "testdata/invalid.so",
			wantErr: fmt.Sprintf("given plugin testdata/invalid.so is not valid: plugin.Open(\"testdata/invalid.so\"): %s: file too short", invalidPluginAbsPath),
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenPlugin(tt.path)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

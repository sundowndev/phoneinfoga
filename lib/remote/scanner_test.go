package remote

import (
	"fmt"
	"os"
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

func TestScannerOptions(t *testing.T) {
	testcases := []struct {
		name  string
		opts  ScannerOptions
		check func(*testing.T, ScannerOptions)
	}{
		{
			name: "test GetStringEnv with simple options",
			opts: map[string]interface{}{
				"foo": "bar",
			},
			check: func(t *testing.T, opts ScannerOptions) {
				assert.Equal(t, opts.GetStringEnv("foo"), "bar")
				assert.Equal(t, opts.GetStringEnv("bar"), "")
			},
		},
		{
			name: "test GetStringEnv with env vars",
			opts: map[string]interface{}{
				"foo_bar": "bar",
			},
			check: func(t *testing.T, opts ScannerOptions) {
				_ = os.Setenv("FOO_BAR", "secret")
				defer os.Unsetenv("FOO_BAR")

				assert.Equal(t, opts.GetStringEnv("FOO_BAR"), "secret")
				assert.Equal(t, opts.GetStringEnv("foo_bar"), "bar")
				assert.Equal(t, opts.GetStringEnv("foo"), "")
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, tt.opts)
		})
	}
}

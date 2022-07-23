package remote

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidatePlugin_Errors(t *testing.T) {
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
			wantErr: "given plugin testdata/invalid.so is not valid",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenPlugin(tt.path)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

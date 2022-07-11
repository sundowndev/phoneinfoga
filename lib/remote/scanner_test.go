package remote

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/mocks"
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
			_, err := OpenPlugin(tt.path)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func Test_parseEntryFunc(t *testing.T) {
	testcases := []struct {
		name     string
		mocks    func(*mocks.Plugin)
		wantErr  string
		expected Scanner
	}{
		{
			name: "test with invalid plugin",
			mocks: func(p *mocks.Plugin) {
				p.On("Lookup", "NewScanner").Return(nil, errors.New("dummy error"))
			},
			wantErr: "exported function NewScanner not found",
		},
		{
			name: "test with invalid exported function",
			mocks: func(p *mocks.Plugin) {
				fn := func() interface{} {
					return &mocks.Scanner{}
				}

				p.On("Lookup", "NewScanner").Return(fn, nil)
			},
			wantErr: "exported function NewScanner does not follow the remote.Scanner interface",
		},
		{
			name: "test with valid plugin",
			mocks: func(p *mocks.Plugin) {
				fn := func() Scanner {
					return &mocks.Scanner{}
				}

				p.On("Lookup", "NewScanner").Return(fn, nil)
			},
			wantErr:  "",
			expected: &mocks.Scanner{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			fakePlugin := &mocks.Plugin{}

			tt.mocks(fakePlugin)

			s, err := parseEntryFunc(fakePlugin)
			if err != nil {
				if tt.wantErr == "" {
					t.Fatal(err)
				}
				return
			}
			assert.Equal(t, tt.expected, s)

			fakePlugin.AssertExpectations(t)
		})
	}
}

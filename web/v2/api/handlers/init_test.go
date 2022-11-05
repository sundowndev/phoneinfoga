package handlers_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/web/v2/api/handlers"
	"testing"
)

func TestInit(t *testing.T) {
	handlers.Init(filter.NewEngine())
	assert.NotNil(t, handlers.RemoteLibrary)
	assert.Greater(t, len(handlers.RemoteLibrary.GetAllScanners()), 0)
}

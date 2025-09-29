package pkg

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEmptyInt(t *testing.T) {
	var empty Opt[int]
	require.Equal(t, empty, NewEmpty[int]())
}

type filter struct {
	id         Opt[uuid.UUID]
	locationID Opt[uuid.UUID]
	limit      Opt[uint64]
	offset     Opt[uint64]
}

func TestEmptyInFilter(t *testing.T) {
	f := filter{}

	require.Equal(t, f.id, NewEmpty[uuid.UUID]())
	require.Equal(t, f.locationID, NewEmpty[uuid.UUID]())
	require.Equal(t, f.limit, NewEmpty[uint64]())
	require.Equal(t, f.offset, NewEmpty[uint64]())
}

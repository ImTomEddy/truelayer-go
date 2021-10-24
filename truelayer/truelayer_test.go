package truelayer_test

import (
	"testing"

	"github.com/ImTomEddy/truelayer-go/truelayer"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tl := truelayer.New("abc", "abc")
	require.NotNil(t, tl)
}

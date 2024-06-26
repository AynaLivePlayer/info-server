package updaters

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var bilibili = new(Bilibili)

func TestBilibili_GetStatus(t *testing.T) {
	status, err := bilibili.GetStatus("3819533")
	require.NoError(t, err)
	require.Equal(t, bilibili.Platform(), status.Platform)
	require.Equal(t, "Aynakeya", status.Username)
	require.Equal(t, "10003632", status.UserID)
}

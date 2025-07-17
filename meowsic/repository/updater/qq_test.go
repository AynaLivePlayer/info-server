package updater

import (
	_ "github.com/AynaLivePlayer/miaosic/providers/netease"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewQQ(t *testing.T) {
	updater := NewQQ()
	_, err := updater.Run()
	require.NoError(t, err)
}

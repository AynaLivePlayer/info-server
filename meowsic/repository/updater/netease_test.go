package updater

import (
	_ "github.com/AynaLivePlayer/miaosic/providers/netease"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewNetease(t *testing.T) {
	updater := NewNetease()
	_, err := updater.Run()
	require.NoError(t, err)
}

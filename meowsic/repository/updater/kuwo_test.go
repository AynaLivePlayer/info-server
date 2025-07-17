package updater

import (
	_ "github.com/AynaLivePlayer/miaosic/providers/kuwo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewKuwo(t *testing.T) {
	updater := NewKuwo()
	_, err := updater.Run()
	require.NoError(t, err)
}

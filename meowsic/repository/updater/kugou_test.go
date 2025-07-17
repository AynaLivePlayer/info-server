package updater

import (
	kugou "github.com/AynaLivePlayer/miaosic/providers/kugou"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	kugou.UseInstrumental()
}

func TestNewKugou(t *testing.T) {
	updater := NewKugou()
	_, err := updater.Run()
	require.NoError(t, err)
}

func TestNewKugouInstrumental(t *testing.T) {
	updater := NewKugouInstrumental()
	_, err := updater.Run()
	require.NoError(t, err)
}

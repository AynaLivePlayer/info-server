package updater

import (
	_ "github.com/AynaLivePlayer/miaosic/providers/bilivideo"
	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBiliVideo(t *testing.T) {
	updater := NewBiliVideo()
	status, err := updater.Run()
	require.NoError(t, err)
	pp.Println(status)
}

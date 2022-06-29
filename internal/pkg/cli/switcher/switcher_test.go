package switcher_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/micronull/i3rotonda/internal/pkg/cli/switcher"
	"github.com/micronull/i3rotonda/internal/pkg/cli/switcher/mocks"
	"github.com/micronull/i3rotonda/internal/pkg/types"
)

func TestCommand_Run(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	connMock := mocks.NewMockWriteCloser(mockCtrl)

	cmd := switcher.NewCommand(connMock)
	err := cmd.Init([]string{"-a=prev"})
	require.NoError(t, err)

	gomock.InOrder(
		connMock.EXPECT().Write([]byte{types.ActionPrev}).Return(0, nil),
		connMock.EXPECT().Close().Return(nil),
	)

	err = cmd.Run()
	require.NoError(t, err)
}

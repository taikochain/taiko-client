package utils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/taikoxyz/taiko-client/internal/testutils"
	"github.com/taikoxyz/taiko-client/internal/utils"
)

func TestEncodeDecodeBytes(t *testing.T) {
	b := testutils.RandomBytes(1024)

	compressed, err := utils.Compress(b)
	require.Nil(t, err)
	require.NotEmpty(t, compressed)

	decompressed, err := utils.Decompress(compressed)
	require.Nil(t, err)
	fmt.Println(1, decompressed)

	require.Equal(t, b, decompressed)
}

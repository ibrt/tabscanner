package tabscanner_test

import (
	"testing"

	"github.com/ibrt/tabscanner"
	"github.com/stretchr/testify/require"
)

func TestParseMoney(t *testing.T) {
	_, err := tabscanner.ParseMoney("", 2)
	require.Error(t, err)

	_, err = tabscanner.ParseMoney("-1", 2)
	require.Error(t, err)

	_, err = tabscanner.ParseMoney("a", 2)
	require.Error(t, err)

	m, err := tabscanner.ParseMoney("10", 2)
	require.NoError(t, err)
	require.EqualValues(t, 1000, m)

	m, err = tabscanner.ParseMoney("10.000", 2)
	require.NoError(t, err)
	require.EqualValues(t, 1000, m)

	m, err = tabscanner.ParseMoney("10.009", 2)
	require.NoError(t, err)
	require.EqualValues(t, 1000, m)

	m, err = tabscanner.ParseMoney("10.009", 3)
	require.NoError(t, err)
	require.EqualValues(t, 10009, m)

	m, err = tabscanner.ParseMoney("10.009", 0)
	require.NoError(t, err)
	require.EqualValues(t, 10, m)

	m, err = tabscanner.ParseMoney("10.99", 3)
	require.NoError(t, err)
	require.EqualValues(t, 10990, m)

	m, err = tabscanner.ParseMoney("10.99", 2)
	require.NoError(t, err)
	require.EqualValues(t, 1099, m)

	m, err = tabscanner.ParseMoney("0.99", 2)
	require.NoError(t, err)
	require.EqualValues(t, 99, m)

	m, err = tabscanner.ParseMoney("0.99", 3)
	require.NoError(t, err)
	require.EqualValues(t, 990, m)

	m, err = tabscanner.ParseMoney("0.99", 1)
	require.NoError(t, err)
	require.EqualValues(t, 9, m)
}

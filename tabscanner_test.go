package tabscanner_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ibrt/tabscanner"
)

func TestTabScanner_001(t *testing.T) {
	t.Parallel()
	p := getProcessor(t)

	buf, err := ioutil.ReadFile("testdata/001_ok.jpg")
	require.NoError(t, err)

	result, err := p.Process(context.Background(), &tabscanner.ProcessRequest{
		ReceiptImage:  buf,
		DecimalPlaces: tabscanner.IntPtr(2),
		Language:      tabscanner.LanguageEnglish,
		LineExtract:   tabscanner.BoolPtr(true),
		DocumentType:  tabscanner.DocumentTypeReceipt,
	})
	require.NoError(t, err)

	require.True(t, result.ValidatedTotal)
}

func TestTabScanner_002(t *testing.T) {
	t.Parallel()
	p := getProcessor(t)

	buf, err := ioutil.ReadFile("testdata/002_small.jpg")
	require.NoError(t, err)

	result, err := p.Process(context.Background(), &tabscanner.ProcessRequest{
		ReceiptImage:  buf,
		DecimalPlaces: tabscanner.IntPtr(2),
		Language:      tabscanner.LanguageEnglish,
		LineExtract:   tabscanner.BoolPtr(true),
		DocumentType:  tabscanner.DocumentTypeReceipt,
	})
	require.NoError(t, err)

	require.False(t, result.ValidatedEstablishment)
	require.False(t, result.ValidatedTotal)
	require.False(t, result.ValidatedSubTotal)
}

func TestTabScanner_003(t *testing.T) {
	t.Parallel()
	p := getProcessor(t)

	buf, err := ioutil.ReadFile("testdata/003_bad.jpg")
	require.NoError(t, err)

	result, err := p.Process(context.Background(), &tabscanner.ProcessRequest{
		ReceiptImage:  buf,
		DecimalPlaces: tabscanner.IntPtr(2),
		Language:      tabscanner.LanguageEnglish,
		LineExtract:   tabscanner.BoolPtr(true),
		DocumentType:  tabscanner.DocumentTypeReceipt,
	})
	require.NoError(t, err)

	require.False(t, result.ValidatedEstablishment)
	require.False(t, result.ValidatedTotal)
	require.False(t, result.ValidatedSubTotal)
}

func TestTabScanner_004(t *testing.T) {
	t.Parallel()
	p := getProcessor(t)

	buf, err := ioutil.ReadFile("testdata/004_ok.png")
	require.NoError(t, err)

	result, err := p.Process(context.Background(), &tabscanner.ProcessRequest{
		ReceiptImage:  buf,
		DecimalPlaces: tabscanner.IntPtr(2),
		Language:      tabscanner.LanguageEnglish,
		LineExtract:   tabscanner.BoolPtr(true),
		DocumentType:  tabscanner.DocumentTypeReceipt,
	})
	require.NoError(t, err)

	require.True(t, result.ValidatedTotal)
}

func TestTabScanner_005(t *testing.T) {
	t.Parallel()
	p := getProcessor(t)

	buf, err := ioutil.ReadFile("testdata/005_bad.tif")
	require.NoError(t, err)

	_, err = p.Process(context.Background(), &tabscanner.ProcessRequest{
		ReceiptImage:  buf,
		DecimalPlaces: tabscanner.IntPtr(2),
		Language:      tabscanner.LanguageEnglish,
		LineExtract:   tabscanner.BoolPtr(true),
		DocumentType:  tabscanner.DocumentTypeReceipt,
	})
	require.Error(t, err)
}

func getProcessor(t *testing.T) *tabscanner.Processor {
	apiKey := os.Getenv("TEST_API_KEY")
	require.NotEmpty(t, apiKey)

	c := tabscanner.NewClient(apiKey)
	return tabscanner.NewProcessor(c)
}

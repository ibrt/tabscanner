# tabscanner [![Build Status](https://travis-ci.org/ibrt/tabscanner.svg?branch=master)](https://travis-ci.org/ibrt/tabscanner) [![Go Report Card](https://goreportcard.com/badge/github.com/ibrt/tabscanner)](https://goreportcard.com/report/github.com/ibrt/tabscanner) [![Test Coverage](https://codecov.io/gh/ibrt/tabscanner/branch/master/graph/badge.svg)](https://codecov.io/gh/ibrt/tabscanner) [![Go Docs](https://godoc.org/github.com/ibrt/tabscanner?status.svg)](http://godoc.org/github.com/ibrt/tabscanner)

This package is a Go client for the [TabScanner](https://tabscanner.com). Use the straight `Client` for low level access to the API, or a `Processor` for simplified access (it automatically parses errors and polls for results).

```go
c := tabscanner.NewClient("my_api_key")
p := tabscanner.NewProcessor(c)

buf, err := ioutil.ReadFile("my_receipt.jpg")
require.NoError(t, err)

result, err := p.Process(context.Background(), &tabscanner.ProcessRequest{
  ReceiptImage:  buf,
  DecimalPlaces: tabscanner.IntPtr(2),
  Language:      tabscanner.LanguageEnglish,
  LineExtract:   tabscanner.BoolPtr(true),
  DocumentType:  tabscanner.DocumentTypeReceipt,
})
if err != nil {
  return err
}

fmt.Println(result)

```

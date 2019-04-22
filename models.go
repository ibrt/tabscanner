package tabscanner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

// StatusCode represents the status of an operation.
type StatusCode int

// Known status codes.
const (
	StatusCodePending StatusCode = 1
	StatusCodeSuccess StatusCode = 2
	StatusCodeDone    StatusCode = 3
	StatusCodeFailed  StatusCode = 4
)

// Code represents an operation result code.
type Code int

// Known codes.
const (
	CodeImageUploadedSuccessfully                         Code = 200
	CodeAPIKeyAuthenticated                               Code = 201
	CodeResultAvailable                                   Code = 202
	CodeImageUploadedButDidNotMeetTheRecommendedDimension Code = 300
	CodeResultNotYetAvailable                             Code = 301
	CodeAPIKeyNotFound                                    Code = 400
	CodeNotEnoughCredit                                   Code = 401
	CodeTokenNotFound                                     Code = 402
	CodeNoFileDetected                                    Code = 403
	CodeMultipleFilesDetected                             Code = 404
	CodeUnsupportedMimeType                               Code = 405
	CodeFormParserError                                   Code = 406
	CodeUnsupportedFileExtension                          Code = 407
	CodeFileSystemError                                   Code = 408
	CodeOCRFailure                                        Code = 500
	CodeServerError                                       Code = 510
	CodeDatabaseConnectionError                           Code = 520
	CodeDatabaseQueryError                                Code = 521
)

// Language describes the language of a receipt.
type Language string

// Known languages.
const (
	LanguageEnglish Language = "english"
	LanguageSpanish Language = "spanish"
)

// DocumentType describes a document type.
type DocumentType string

// Known document types.
const (
	DocumentTypeReceipt DocumentType = "receipt"
	DocumentTypeInvoice DocumentType = "invoice"
	DocumentTypeAuto    DocumentType = "auto"
)

// IntPtr returns a pointer to the given int.
func IntPtr(i int) *int {
	return &i
}

// BoolPtr returns a pointer to the given bool.
func BoolPtr(b bool) *bool {
	return &b
}

// ResponseHeader describes the common fields in responses.
type ResponseHeader struct {
	Message    string     `json:"message"`
	Status     string     `json:"status"`
	StatusCode StatusCode `json:"status_code"`
	Token      string     `json:"token"`
	Success    bool       `json:"success"`
	Code       Code       `json:"code"`
}

// ProcessRequest describes a process request. Specify only one of ReceiptImageHeader or ReceiptImage data.
type ProcessRequest struct {
	ReceiptImage  []byte       // mime-type is automatically detected
	DecimalPlaces *int         // valid values: nil, 0, 2, 3
	Language      Language     // valid values: "", LanguageEnglish, LanguageSpanish
	Cents         *bool        // valid values: nil, false, true
	LineExtract   *bool        // valid values: nil, false, true
	DocumentType  DocumentType // valid values: "", DocumentTypeReceipt, DocumentTypeInvoice, DocumentTypeAuto
	TestMode      *bool        // valid values: nil, false, true
}

func (r *ProcessRequest) toFormBody() (io.Reader, string, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	if err := r.writeReceipt(w); err != nil {
		return nil, "", err
	}

	if r.DecimalPlaces != nil {
		if err := w.WriteField("decimalPlaces", fmt.Sprintf("%v", *r.DecimalPlaces)); err != nil {
			return nil, "", err
		}
	}

	if r.Language != "" {
		if err := w.WriteField("language", string(r.Language)); err != nil {
			return nil, "", err
		}
	}

	if r.Cents != nil {
		if err := w.WriteField("cents", fmt.Sprintf("%v", *r.Cents)); err != nil {
			return nil, "", err
		}
	}

	if r.LineExtract != nil {
		if err := w.WriteField("lineExtract", fmt.Sprintf("%v", *r.LineExtract)); err != nil {
			return nil, "", err
		}
	}

	if r.DocumentType != "" {
		if err := w.WriteField("documentType", string(r.DocumentType)); err != nil {
			return nil, "", err
		}
	}

	if r.TestMode != nil {
		if err := w.WriteField("testMode", fmt.Sprintf("%v", *r.TestMode)); err != nil {
			return nil, "", err
		}
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return body, w.FormDataContentType(), nil
}

func (r *ProcessRequest) writeReceipt(w *multipart.Writer) error {
	receiptMimeType, receiptExt, err := r.getType()
	if err != nil {
		return err
	}

	receiptHeader := textproto.MIMEHeader{}
	receiptHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="receiptImage"; filename="receipt.%v"`, receiptExt))
	receiptHeader.Set("Content-Type", receiptMimeType)

	receiptPart, err := w.CreatePart(receiptHeader)
	if err != nil {
		return err
	}

	_, err = io.Copy(receiptPart, bytes.NewReader(r.ReceiptImage))
	return err
}

func (r *ProcessRequest) getType() (string, string, error) {
	mimeType := http.DetectContentType(r.ReceiptImage)

	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", "", err
	}

	if len(extensions) == 0 {
		return mimeType, "bin", nil
	}

	return mimeType, extensions[0], nil
}

// ProcessResponse describes a process response.
type ProcessResponse struct {
	*ResponseHeader
	Duplicate      bool   `json:"duplicate"`
	DuplicateToken string `json:"duplicateToken"`
}

// ResultResponse describes a result response.
type ResultResponse struct {
	*ResponseHeader
	Result *ResultResponseResult `json:"result"`
}

// ResultResponseResult describes the result field in a result response.
type ResultResponseResult struct {
	Establishment          string                          `json:"establishment"`
	ValidatedEstablishment bool                            `json:"validatedEstablishment"`
	Date                   string                          `json:"date"`
	Total                  string                          `json:"total"`
	URL                    string                          `json:"url"`
	PhoneNumber            string                          `json:"phoneNumber"`
	PaymentMethod          string                          `json:"paymentMethod"`
	Address                string                          `json:"address"`
	ValidatedTotal         bool                            `json:"validatedTotal"`
	SubTotal               string                          `json:"subTotal"`
	ValidatedSubTotal      bool                            `json:"validatedSubTotal"`
	Cash                   string                          `json:"cash"`
	Change                 string                          `json:"change"`
	Tax                    string                          `json:"tax"`
	Taxes                  []json.Number                   `json:"taxes"`
	Discount               string                          `json:"discount"`
	Discounts              []json.Number                   `json:"discounts"`
	LineItems              []*ResultResponseResultLineItem `json:"lineItems"`
}

// ResultResponseResultLineItem describes a line item in a result.
type ResultResponseResultLineItem struct {
	Quantity         json.Number `json:"qty"`
	Description      string      `json:"desc"`
	Unit             string      `json:"unit"`
	CleanDescription string      `json:"descClean"`
	LineTotal        string      `json:"lineTotal"`
	ProductCode      string      `json:"productCode"`
}

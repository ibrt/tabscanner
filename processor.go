package tabscanner

import (
	"context"
	"fmt"
	"time"
)

// ErrTimedOut is returned when the Processor cannot obtain a result after having exhausted all polling retries.
type ErrTimedOut struct {
	// intentionally empty
}

// Error implements the error interface.
func (e *ErrTimedOut) Error() string {
	return "timed out"
}

// ErrAPI wraps the ResponseHeader for a failure response as an error.
type ErrAPI struct {
	responseHeader *ResponseHeader
}

// ResponseHeader returns the ResponseHeader for this error.
func (e *ErrAPI) ResponseHeader() *ResponseHeader {
	return e.responseHeader
}

// Error implements the error interface.
func (e *ErrAPI) Error() string {
	return fmt.Sprintf("tabscanner error %v: %v", e.responseHeader.Code, e.responseHeader.Message)
}

// PollingStrategy describes a polling strategy for results. Delay is called before each call to the result endpoint,
// with monotonically increasing values of retry, starting from 0. If the returned value is >= 0, the processor waits
// the returned duration and retries, otherwise it stops and returns an error.
type PollingStrategy func(retry int) time.Duration

// DefaultPollingStrategy tries every five seconds, for up to three minutes.
var DefaultPollingStrategy = func(retry int) time.Duration {
	switch {
	case retry == 0:
		return 0
	case retry > 180:
		return -1
	default:
		return 5 * time.Second
	}
}

// ProcessorOption describes a configuration option for the Processor.
type ProcessorOption func(*Processor)

// PollingStrategyProcessorOption allows to inject a PollingStrateg.
func PollingStrategyProcessorOption(pollingStrategy PollingStrategy) ProcessorOption {
	return func(p *Processor) {
		p.pollingStrategy = pollingStrategy
	}
}

// Processor provides a simplified way to scan receipts, implementing automatic error handling an polling for results.
type Processor struct {
	client          *Client
	pollingStrategy PollingStrategy
}

// NewProcessor initializes a new Processor.
func NewProcessor(client *Client, options ...ProcessorOption) *Processor {
	p := &Processor{
		client:          client,
		pollingStrategy: DefaultPollingStrategy,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

// Process uploads a receipt for processing, blocks until a result is available or an error occurs. Polling for results
// is performed according to the Processor's PollingStrategy.
func (p *Processor) Process(ctx context.Context, req *ProcessRequest) (*ResultResponseResult, error) {
	processResp, err := p.client.Process(ctx, req)
	if err != nil {
		return nil, err
	}

	if processResp.StatusCode == StatusCodeFailed {
		return nil, fmt.Errorf(processResp.Message)
	}

	for i, delay := 0, p.pollingStrategy(0); delay >= 0; i, delay = i+1, p.pollingStrategy(i+1) {
		time.Sleep(delay)

		resultResp, err := p.client.Result(ctx, processResp.Token)
		if err != nil {
			return nil, err
		}

		switch resultResp.StatusCode {
		case StatusCodeDone:
			return resultResp.Result, nil
		case StatusCodePending:
			continue
		case StatusCodeFailed:
			return nil, &ErrAPI{responseHeader: resultResp.ResponseHeader}
		default:
			return nil, fmt.Errorf("unexpected status code '%v'", resultResp.StatusCode)
		}
	}

	return nil, &ErrTimedOut{}
}

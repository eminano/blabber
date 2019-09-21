package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockSpeechProcessor struct {
	text string
	err  error
}

func (m *mockSpeechProcessor) ProcessText(input io.Reader) error {
	if m.err != nil {
		return m.err
	}

	return nil
}

func (m *mockSpeechProcessor) GenerateRandomText(_ uint) string {
	return m.text
}

func Test_Learn(t *testing.T) {
	t.Parallel()

	var validProcessor = &mockSpeechProcessor{}

	var tests = []struct {
		name      string
		req       *http.Request
		processor *mockSpeechProcessor

		wantCode     int
		wantResponse string
	}{
		{
			name:         "ok",
			req:          getValidRequest(t, "POST", "/learn", "this is my test body"),
			processor:    validProcessor,
			wantCode:     http.StatusOK,
			wantResponse: "",
		},
		{
			name:         "invalid method",
			req:          getValidRequest(t, "GET", "/generate"),
			processor:    validProcessor,
			wantCode:     http.StatusMethodNotAllowed,
			wantResponse: "Invalid method\n",
		},
		{
			name: "invalid content type",
			req: func() *http.Request {
				var r = getValidRequest(t, "POST", "/learn")
				r.Header.Set("Content-Type", "text/html")
				return r
			}(),
			processor:    validProcessor,
			wantCode:     http.StatusUnprocessableEntity,
			wantResponse: "Invalid content type\n",
		},
		{
			name:         "error - processing text",
			req:          getValidRequest(t, "POST", "/learn", "this is my test body"),
			processor:    &mockSpeechProcessor{err: errors.New("can't compute")},
			wantCode:     http.StatusUnprocessableEntity,
			wantResponse: "Error processing text: can't compute\n",
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var recorder = httptest.NewRecorder()

			var h = Handlers{
				processor: tt.processor,
			}

			h.Learn(recorder, tt.req)

			if recorder.Code != tt.wantCode {
				t.Errorf("got: %v, want: %v", recorder.Code, tt.wantCode)
			}

			if !reflect.DeepEqual(recorder.Body.String(), tt.wantResponse) {
				t.Errorf("got %v, want %v", recorder.Body.String(), tt.wantResponse)
			}
		})
	}
}

func Test_Generate(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name string
		req  *http.Request

		wantCode     int
		wantResponse string
	}{
		{
			name:         "ok",
			req:          getValidRequest(t, "GET", "/generate"),
			wantCode:     http.StatusOK,
			wantResponse: "random text is random",
		},
		{
			name:         "invalid method",
			req:          getValidRequest(t, "POST", "/generate"),
			wantCode:     http.StatusMethodNotAllowed,
			wantResponse: "Invalid method\n",
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var recorder = httptest.NewRecorder()

			var h = Handlers{
				processor: &mockSpeechProcessor{
					text: "random text is random",
				},
			}

			h.Generate(recorder, tt.req)

			if recorder.Code != tt.wantCode {
				t.Errorf("got: %v, want: %v", recorder.Code, tt.wantCode)
			}

			if !reflect.DeepEqual(recorder.Body.String(), tt.wantResponse) {
				t.Errorf("got %v, want %v", recorder.Body.String(), tt.wantResponse)
			}
		})
	}
}

func getValidRequest(t *testing.T, method string, path string, input ...string) *http.Request {
	t.Helper()

	var body io.Reader

	if len(input) > 0 && input[0] != "" {
		var data = &bytes.Buffer{}
		if _, err := data.WriteString(input[0]); err != nil {
			t.Fatalf("error writing request body data: %v", err)
		}

		body = data
	}

	var validReq, err = http.NewRequest(method, path, body)
	if err != nil {
		t.Fatalf("error creaing request: %v", err)
	}

	if body != nil {
		validReq.Header.Set("Content-Type", textPlainContentType)
	}

	return validReq
}

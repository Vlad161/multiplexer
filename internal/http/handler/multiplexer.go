package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/vlad161/affise_test_task/internal/logger"
)

const (
	maxUrlsSize = 20
)

type request = []string
type response = map[string]map[string]interface{}

type handler struct {
	multiplexer Multiplexer
}

func NewMultiplexer(multiplexer Multiplexer) *handler {
	return &handler{
		multiplexer: multiplexer,
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if !h.validate(w, r) {
		return
	}

	ctx := r.Context()
	log := logger.FromContext(ctx)
	resp, errHandle := h.handle(ctx, r.Body)
	defer func() {
		if errCloseBody := r.Body.Close(); errCloseBody != nil {
			log.Error("can't close request body: %w", errCloseBody)
		}
	}()
	if errHandle == nil {
		if errEncode := json.NewEncoder(w).Encode(resp); errEncode == nil {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			return
		}
	}

	h.handlerError(log, errHandle, w)
}

func (h handler) validate(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return false
	}
	return true
}

func (h handler) handle(ctx context.Context, body io.ReadCloser) (response, error) {
	var urls request
	if errDecode := json.NewDecoder(body).Decode(&urls); errDecode != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecodeBody, errDecode)
	}
	if len(urls) > maxUrlsSize {
		return nil, ErrTooMuchUrls
	}

	return h.multiplexer.Urls(ctx, urls)
}

func (h handler) handlerError(log logger.Logger, err error, w http.ResponseWriter) {
	var (
		code    int
		message string
	)

	switch {
	case errors.Is(err, ErrTooMuchUrls):
		code = http.StatusRequestEntityTooLarge
		message = fmt.Sprintf("Max array size is %d", maxUrlsSize)
	case errors.Is(err, ErrDecodeBody):
		code = http.StatusBadRequest
		message = http.StatusText(http.StatusBadRequest)
	default:
		code = http.StatusInternalServerError
		message = http.StatusText(http.StatusInternalServerError)
	}
	log.Error("handle error: %w", err)

	w.WriteHeader(code)
	if _, errWrite := w.Write([]byte(message)); errWrite != nil {
		log.Error("can't write handler error body %w", errWrite)
	}
}

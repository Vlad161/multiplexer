package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/vlad161/multiplexer/internal/logger"
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
	if !h.validate(w, r) {
		return
	}

	ctx := r.Context()
	log := logger.FromContext(ctx)
	resp, errHandle := h.handle(ctx, r.Body)
	defer func() {
		if errCloseBody := r.Body.Close(); errCloseBody != nil {
			log.Error("can't close request body: %v", errCloseBody)
		}
	}()
	if errHandle == nil {
		if errEncode := json.NewEncoder(w).Encode(resp); errEncode == nil {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			return
		}
	}

	h.handlerError(errHandle, w)
	return
}

func (h handler) validate(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return false
	}
	return true
}

func (h handler) handle(ctx context.Context, body io.ReadCloser) (response, error) {
	var urls request
	if errDecode := json.NewDecoder(body).Decode(&urls); errDecode != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecodeBody, errDecode)
	}
	if len(urls) == 0 {
		return nil, ErrEmptyUrls
	}
	if len(urls) > maxUrlsSize {
		return nil, ErrTooMuchUrls
	}

	return h.multiplexer.Urls(ctx, urls)
}

func (h handler) handlerError(err error, w http.ResponseWriter) {
	var (
		code    int
		message string
	)

	switch {
	case errors.Is(err, ErrTooMuchUrls):
		code = http.StatusRequestEntityTooLarge
		message = fmt.Sprintf("Max array size is %d", maxUrlsSize)
	case errors.Is(err, ErrEmptyUrls), errors.Is(err, ErrDecodeBody):
		code = http.StatusBadRequest
		message = http.StatusText(http.StatusBadRequest)
	default:
		code = http.StatusInternalServerError
		message = http.StatusText(http.StatusInternalServerError)
	}

	http.Error(w, message, code)
}

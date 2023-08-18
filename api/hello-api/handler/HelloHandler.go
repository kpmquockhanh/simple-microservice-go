package handler

import (
	"context"
	"net/http"
	logger2 "simple-micro/core/logger"
	"simple-micro/core/transhttp"
	sample_services "simple-micro/exmsg/services"
)

type HelloHandler struct {
	SampleClient sample_services.SampleClient
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := logger2.NewLogger()
	resp, err := h.SampleClient.GetNumber(context.Background(), &sample_services.SampleRequest{})
	if err != nil {
		logger.Infow("Error when GetNumber", "error", err)
		transhttp.Json(w, map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	transhttp.Json(w, resp)
}

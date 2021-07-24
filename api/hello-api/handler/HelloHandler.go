package handler

import (
	"context"
	"log"
	"net/http"
	"simple-micro/core/transhttp"
	sample_services "simple-micro/exmsg/services"
)

type HelloHandler struct {
	SampleClient sample_services.SampleClient
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	resp, err := h.SampleClient.GetNumber(context.Background(), &sample_services.SampleRequest{})
	if err != nil {
		log.Println("Error when GetNumber", err)
		transhttp.Json(w, map[string]interface{}{
			"error": true,
			"message": err.Error(),
		})
		return
	}
	transhttp.Json(w, resp)
}

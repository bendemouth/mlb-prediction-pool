package handlers

import "net/http"

func (h *Handler) HandleHealthCheck(writer http.ResponseWriter, request *http.Request) {
	status := h.healthcheckService.HealthCheck(request.Context())
	h.respondJson(writer, http.StatusOK, status)
}

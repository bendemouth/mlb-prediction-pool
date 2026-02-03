package handlers

import "net/http"

func (h *Handler) HandleHealthCheck(writer http.ResponseWriter, request *http.Request) {
	status := h.healthcheckService.CheckHealth()
	h.respondJson(writer, http.StatusOK, map[string]string{"status": status})
}

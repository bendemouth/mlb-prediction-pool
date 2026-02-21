package handlers

import "net/http"

func (h *Handler) AdminCompleteSeeded(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		h.respondError(writer, http.StatusBadRequest, "Invalid request method")
		return
	}

}

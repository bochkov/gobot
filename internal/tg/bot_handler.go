package tg

import (
	"encoding/json"
	"net/http"
	"strings"

	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) BotHandler(w http.ResponseWriter, req *http.Request) {
	var upd Update
	if err := json.NewDecoder(req.Body).Decode(&upd); err != nil {
		return
	}
	slog.Debug(fmt.Sprintf("%+v", upd.Message))
	token := chi.URLParam(req, "token")
	if h.shouldAnswer(upd.Message) {
		go h.sendAnswer(token, upd.Message)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) shouldAnswer(msg *Message) bool {
	return msg != nil &&
		(strings.HasPrefix(msg.Text, "@resnyx") || strings.HasPrefix(msg.Text, "/"))
}

func (h *Handler) sendAnswer(token string, msg *Message) {
	methods := h.Service.GetAnswers(msg)
	for _, method := range methods {
		res, err := h.Service.Execute(method, token)
		if err != nil {
			slog.Warn(err.Error())
			return
		}
		if res.Ok {
			slog.Info(fmt.Sprintf("response: %+v", res.Result))
		} else {
			slog.Info(fmt.Sprintf("response: %d %s", res.ErrorCode, res.Description))
		}
	}
}

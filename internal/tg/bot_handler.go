package tg

import (
	"encoding/json"
	"net/http"
	"strings"

	"fmt"
	"log/slog"

	"github.com/go-chi/chi/v5"
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
	slog.Debug(fmt.Sprintf("%+v", upd.InlineQuery))
	slog.Debug(fmt.Sprintf("%+v", upd.Message))
	token := chi.URLParam(req, "token")

	iq := upd.InlineQuery
	if iq != nil {
		go h.sendAnswer(token, iq.User.Id, iq.Query)
	}
	msg := upd.Message
	if msg != nil && h.shouldAnswer(msg) {
		go h.sendAnswer(token, msg.Chat.Id, msg.Text)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) shouldAnswer(msg *Message) bool {
	return strings.HasPrefix(msg.Text, "@resnyx") || strings.HasPrefix(msg.Text, "/")
}

func (h *Handler) sendAnswer(token string, chatId int64, txt string) {
	methods := h.Service.GetAnswers(chatId, txt)
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

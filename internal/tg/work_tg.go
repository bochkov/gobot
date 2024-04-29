package tg

import (
	"strings"
)

type TgWorker struct {
}

func NewTgWorker() *TgWorker {
	return &TgWorker{}
}

func (a *TgWorker) Description() string {
	return "Служебная служба"
}

func (a *TgWorker) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "start")
}

func (a *TgWorker) Answer(chatId int64, txt string) []Method {
	sm := SendMessage[int64]{ChatId: chatId}
	sm.Text = "С чего начнем"
	sm.InlineMarkup = &InlineKeyboardMarkup{}

	kb1 := InlineKbButton{Text: "Цитату", Callback: "цитат"}
	kb2 := InlineKbButton{Text: "Анекдот", Callback: "анек"}
	kb3 := InlineKbButton{Text: "Курс валют", Callback: "курс"}
	row := []InlineKbButton{kb1, kb2, kb3}
	sm.InlineMarkup.Keyboard = append(sm.InlineMarkup.Keyboard, row)

	res := make([]Method, 0)
	res = append(res, &sm)
	return res
}

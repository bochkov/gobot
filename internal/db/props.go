package db

import (
	"context"
)

const (
	TgBotTokenKey    string = "tg.bot.token"
	ChatAutoSend     string = "chat.auto.send"
	ChatIdKey        string = "chat.id"
	AnekdotScheduler string = "schedule.anekdot"
)

func GetProp(key string, def string) string {
	query := `select p.value from props p where p.key = $1`
	var res string
	if err := GetPool().QueryRow(context.Background(), query, key).Scan(&res); err != nil {
		return def
	}
	return res
}

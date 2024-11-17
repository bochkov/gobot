package adapters

import (
	"regexp"
	"strings"

	"github.com/bochkov/gobot/internal/services/rutor"
	"github.com/bochkov/gobot/internal/tg"
)

type RutorAdapter struct {
	service rutor.Service
}

func NewRutorAdapter(s rutor.Service) RutorAdapter {
	return RutorAdapter{service: s}
}

func (r RutorAdapter) Description() string {
	return r.service.Description()
}

func (r RutorAdapter) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "rutor")
}

func (r RutorAdapter) Answer(chatId int64, txt string) []tg.Method {
	re := regexp.MustCompile(`https?://\S+`)
	matches := re.FindStringSubmatch(txt)
	if len(matches) == 0 {
		return []tg.Method{
			&tg.SendMessage[int64]{
				ChatId: chatId,
				Text:   "не смог выделить URL",
			},
		}
	}
	var url = matches[0]
	torrent, err := r.service.FetchTorrent(url)
	if err != nil {
		return []tg.Method{
			&tg.SendMessage[int64]{
				ChatId: chatId,
				Text:   err.Error(),
			},
		}
	}
	return []tg.Method{
		&tg.SendMessage[int64]{
			ChatId: chatId,
			Text:   torrent.MagnetUrl,
		},
		&tg.SendDocument[int64]{
			ChatId:                      chatId,
			InputFile:                   tg.InputFile{Bytes: torrent.Bytes, Filename: torrent.Name},
			Caption:                     torrent.Name,
			DisableContentTypeDetection: true,
		},
	}
}

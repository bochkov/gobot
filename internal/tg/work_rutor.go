package tg

import (
	"regexp"
	"strings"

	"github.com/bochkov/gobot/internal/rutor"
)

type RutorWorker struct {
	rutor.Service
}

func NewRutorWorker(s rutor.Service) *RutorWorker {
	return &RutorWorker{Service: s}
}

func (r *RutorWorker) Description() string {
	return "rutor"
}

func (r *RutorWorker) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "rutor")
}

func (r *RutorWorker) Answer(msg *Message) []Method {
	re := regexp.MustCompile(`https?://\S+`)
	matches := re.FindStringSubmatch(msg.Text)
	if len(matches) == 0 {
		return []Method{
			&SendMessage[int64]{
				ChatId: msg.Chat.Id,
				Text:   "не смог выделить URL",
			},
		}
	}
	var url = matches[0]
	torrent, err := r.Service.FetchTorrent(url)
	if err != nil {
		return []Method{
			&SendMessage[int64]{
				ChatId: msg.Chat.Id,
				Text:   err.Error(),
			},
		}
	}
	return []Method{
		&SendMessage[int64]{
			ChatId: msg.Chat.Id,
			Text:   torrent.MagnetUrl,
		},
		&SendDocument[int64]{
			ChatId:                      msg.Chat.Id,
			InputFile:                   InputFile{bytes: torrent.Bytes, filename: torrent.Name},
			Caption:                     torrent.Name,
			DisableContentTypeDetection: true,
		},
	}
}

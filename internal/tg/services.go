package tg

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bochkov/gobot/internal/resnyx/anekdot"
	"github.com/bochkov/gobot/internal/resnyx/auto"
	"github.com/bochkov/gobot/internal/resnyx/cbr"
	"github.com/bochkov/gobot/internal/resnyx/forismatic"
	"github.com/bochkov/gobot/internal/resnyx/rutor"
)

type Service interface {
	Description() string
	IsMatch(text string) bool
	Answer(msg *Message) []Method
}

/// Anek

type Anekdot struct {
}

func (a *Anekdot) Description() string {
	return "Случайный анекдот от https://baneks.ru/"
}

func (a *Anekdot) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "анек")
}

func (a *Anekdot) Answer(msg *Message) []Method {
	sm := SendMessage[int64]{ChatId: msg.Chat.Id}
	anek, err := anekdot.NewBaneks().GetRandom()
	if err != nil {
		log.Print(err)
		sm.Text = "не смог найти анекдот"
	} else {
		sm.Text = anek.Text
	}
	res := make([]Method, 0)
	res = append(res, &sm)
	return res
}

/// Auto

type Auto struct {
}

func (a *Auto) Description() string {
	return "Автомобильные коды регионов РФ"
}

func (a *Auto) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "avto") ||
		strings.Contains(strings.ToLower(text), "авто")
}

func (a *Auto) Answer(msg *Message) []Method {
	re := regexp.MustCompile(`(?P<code>\d+)`)
	matches := re.FindAllStringSubmatch(msg.Text, -1)
	if len(matches) == 0 {
		name := msg.Text[strings.Index(msg.Text, " ")+1:]
		regions, err := auto.FindRegionByName(name)
		if err != nil {
			log.Print(err)
		} else {
			return a.createMessages(msg.Chat.Id, regions)
		}
	} else {
		regions := make([]auto.Region, 0)
		for _, digits := range matches {
			region, err := auto.FindRegionByCode(digits[1])
			if err != nil {
				log.Print(err)
			} else {
				regions = append(regions, *region)
			}
		}
		return a.createMessages(msg.Chat.Id, regions)
	}
	return a.createMessages(msg.Chat.Id, []auto.Region{})
}

func (a *Auto) createMessages(chatId int64, regions []auto.Region) []Method {
	messages := make([]Method, 0)
	if len(regions) == 0 {
		sm := SendMessage[int64]{
			ChatId: chatId,
			Text:   "ничего не нашел (",
		}
		messages = append(messages, &sm)
	} else {
		for _, r := range regions {
			sm := SendMessage[int64]{
				ChatId: chatId,
				Text:   fmt.Sprintf("%s = %s", r.Name, strings.Join(r.Codes, ", ")),
			}
			messages = append(messages, &sm)
		}
	}
	return messages
}

/// CBR

type Cbr struct {
}

func (c Cbr) Description() string {
	return "Официальный курс валют"
}

func (c Cbr) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "курс")
}

func (c Cbr) Answer(msg *Message) []Method {
	text := strings.ToUpper(msg.Text)
	currencies := regexp.MustCompile(`\s+`).Split(text, -1)[1:]
	if len(currencies) == 0 {
		currencies = append(currencies, "USD")
		currencies = append(currencies, "EUR")
	}
	var methods = make([]Method, 0)
	for _, r := range cbr.NewCbr().LatestRange(currencies) {
		sm := SendMessage[int64]{
			ChatId: msg.Chat.Id,
			Text:   r.String(),
		}
		methods = append(methods, &sm)
	}
	return methods
}

/// FORISMATIC

type Forismatic struct {
}

func (f Forismatic) Description() string {
	return "Случайная цитата от forismatic.com"
}

func (f Forismatic) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "цитат")
}

func (f Forismatic) Answer(msg *Message) []Method {
	cite, err := forismatic.NewCiteService().GetQuote()
	sm := SendMessage[int64]{ChatId: msg.Chat.Id}
	if err != nil {
		sm.Text = "не получилось ("
	}
	sm.Text = cite.String()
	return []Method{&sm}
}

/// RUTOR

type Rutor struct {
}

func (r Rutor) Description() string {
	return "Rutor"
}

func (r Rutor) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "rutor")
}

func (r Rutor) Answer(msg *Message) []Method {
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
	torrent, err := rutor.NewService().GetTorrent(url)
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

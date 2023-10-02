package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/bochkov/gobot/internal/util"
	"github.com/carlmjohnson/requests"
	"github.com/gorilla/mux"
)

var services = []Service{
	&Anekdot{},
	&Auto{},
	&Cbr{},
	&Forismatic{},
	&Rutor{},
}

type Bot struct {
	token string
}

func NewBot(token string) *Bot {
	return &Bot{token: token}
}

func (bot *Bot) Execute(method Method) (*TypedResult[any], error) {
	methodName, methodResponse := method.Describe()
	log.Printf("request : method='%s', body=%+v", methodName, util.ToJson(method))
	res := &TypedResult[any]{
		Result: methodResponse,
	}
	reqBuilder := requests.URL(fmt.Sprintf("https://api.telegram.org/bot%s/%s", bot.token, methodName))
	if reflect.ValueOf(method).Elem().FieldByName("InputFile").IsValid() {
		buf := &bytes.Buffer{}
		contentType, err := ToMultipart(method.(MultipartMethod), buf)
		if err != nil {
			return nil, err
		}
		reqBuilder.BodyBytes(buf.Bytes()).ContentType(contentType)
	} else {
		reqBuilder.BodyJSON(&method)
	}
	if err := reqBuilder.
		Method(http.MethodPost).
		ToJSON(res).
		Fetch(context.Background()); err != nil {
		return nil, err
	}
	return res, nil
}

func shouldAnswer(msg *Message) bool {
	return msg != nil &&
		(strings.HasPrefix(msg.Text, "@resnyx") || strings.HasPrefix(msg.Text, "/"))
}

func getAnswers(msg *Message) []Method {
	for _, serv := range services {
		if serv.IsMatch(msg.Text) {
			log.Printf("choosed = %s", serv.Description())
			return serv.Answer(msg)
		}
	}
	return nil
}

func sendAnswer(token string, msg *Message) {
	bot := NewBot(token)
	methods := getAnswers(msg)
	for _, method := range methods {
		res, err := bot.Execute(method)
		if err != nil {
			log.Print(err)
			return
		}
		if res.Ok {
			log.Printf("response: %+v", res.Result)
		} else {
			log.Printf("response: %d %s", res.ErrorCode, res.Description)
		}
	}
}

func BotHandler(_ http.ResponseWriter, req *http.Request) {
	var upd Update
	if err := json.NewDecoder(req.Body).Decode(&upd); err != nil {
		return
	}
	log.Printf("%+v", upd.Message)
	var token = mux.Vars(req)["token"]
	if shouldAnswer(upd.Message) {
		go sendAnswer(token, upd.Message)
	}
}

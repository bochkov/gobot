package tg

type Update struct {
	Message *Message `json:"message"`
}

type User struct {
	Id        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

type ChatType string

const (
	PRIVATE ChatType = "private"
)

type Document struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

type Chat struct {
	Id        int64    `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	UserName  string   `json:"username"`
	Type      ChatType `json:"type"`
}

type Message struct {
	Id       int64    `json:"message_id"`
	User     User     `json:"from"`
	Chat     Chat     `json:"chat"`
	Document Document `json:"document"`
	Date     int64    `json:"date"`
	Text     string   `json:"text"`
}

type TypedResult[result any] struct {
	Ok          bool   `json:"ok"`
	Result      result `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

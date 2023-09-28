package tg

type TypedResult[result any] struct {
	Ok          bool   `json:"ok"`
	Result      result `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

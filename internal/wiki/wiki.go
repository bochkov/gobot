package wiki

import "github.com/bochkov/gobot/internal/push"

type Service interface {
	push.Push
	Today() (string, error)
	Description() string
}

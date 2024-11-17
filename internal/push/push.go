package push

// Сервис, который выполняет пуш. Сейчас это только телеграм-бот.
type Service interface {
	Push()
}

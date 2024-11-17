package services

// базовый интерфейс, который должны реализовывать все сервисы
type Service interface {
	Description() string
}
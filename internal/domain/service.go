package domain

// ServiceCenter решает, можно ли принять транспорт в парк.
type ServiceCenter interface {
	// Inspect возвращает решение о приёмке транспорта.
	Inspect(Transport) bool
}

// BasicServiceCenter принимает только исправный транспорт.
type BasicServiceCenter struct{}

// Inspect возвращает решение по техническому состоянию транспорта.
func (BasicServiceCenter) Inspect(transport Transport) bool {
	return transport.IsOperational()
}

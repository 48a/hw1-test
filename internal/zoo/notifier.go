package zoo

import "fmt"

// Notifier — канал уведомлений сотрудников зоопарка (email, мессенджер, лог).
// Интерфейс введён ради подмены в тестах: мок генерируется mockery (см. .mockery.yml).
type Notifier interface {
	Notify(message string) error
}

// FeedAll «кормит» животных, сообщая о каждом через Notifier.
// При первой же ошибке уведомления возвращает её, обернув через %w.
func FeedAll(n Notifier, animals []Animal) error {
	for _, a := range animals {
		if err := n.Notify("кормлю: " + a.Name); err != nil {
			return fmt.Errorf("уведомление о кормлении: %w", err)
		}
	}
	return nil
}

package domain

import (
	"fmt"
	"slices"
	"strings"
)

// Park хранит принятый транспорт и вещи в порядке добавления.
// Объекты передаются ненулевыми значениями и не изменяются после приёмки.
type Park struct {
	service ServiceCenter
	items   []Inventory
}

// NewPark создаёт пустой парк с заданным сервисом техосмотра.
func NewPark(service ServiceCenter) (*Park, error) {
	if service == nil {
		return nil, ErrServiceRequired
	}
	return &Park{service: service}, nil
}

// AddTransport принимает транспорт только после положительного техосмотра.
func (p *Park) AddTransport(transport Transport) error {
	if err := p.validate(transport); err != nil {
		return err
	}
	if transport.BeginnerEase() < minimumEase || transport.BeginnerEase() > maximumEase {
		return ErrInvalidEase
	}
	if !p.service.Inspect(transport) {
		return fmt.Errorf("%w: %s", ErrInspectionRejected, transport.InventoryNumber())
	}
	p.items = append(p.items, transport)
	return nil
}

// AddItem добавляет вещь без техосмотра и запрещает передачу транспорта.
func (p *Park) AddItem(item Inventory) error {
	if _, ok := item.(Transport); ok {
		return ErrTransportAsItem
	}
	if err := p.validate(item); err != nil {
		return err
	}
	p.items = append(p.items, item)
	return nil
}

func (p *Park) validate(item Inventory) error {
	if item == nil {
		return ErrInvalidAsset
	}
	a, err := newAsset(item.Name(), item.InventoryNumber())
	if err != nil {
		return err
	}
	if consumer, ok := item.(EnergyConsumer); ok {
		if _, err := newEnergy(consumer.DailyEnergy()); err != nil {
			return err
		}
	}
	if slices.ContainsFunc(p.items, func(existing Inventory) bool {
		return strings.TrimLeft(strings.TrimSpace(existing.InventoryNumber()), "0") == a.number
	}) {
		return fmt.Errorf("%w: %s", ErrDuplicateNumber, a.number)
	}
	return nil
}

// Inventory возвращает копию списка всего баланса.
func (p *Park) Inventory() []Inventory { return slices.Clone(p.items) }

// Beginners возвращает принятый транспорт с простотой не ниже 6.
func (p *Park) Beginners() []Transport {
	var result []Transport
	for _, item := range p.items {
		if transport, ok := item.(Transport); ok && transport.BeginnerEase() >= beginnerEase {
			result = append(result, transport)
		}
	}
	return result
}

// Report содержит количество объектов и суточную энергию всего баланса.
type Report struct {
	// Transport содержит количество принятого транспорта.
	Transport int
	// Items содержит количество вещей.
	Items int
	// Total содержит количество объектов всего баланса.
	Total int
	// Energy содержит суточный расход всего баланса в кВт·ч.
	Energy float64
}

// Report рассчитывает показатели по всем принятым объектам.
func (p *Park) Report() Report {
	report := Report{Total: len(p.items)}
	for _, item := range p.items {
		if _, ok := item.(Transport); ok {
			report.Transport++
		} else {
			report.Items++
		}
		if consumer, ok := item.(EnergyConsumer); ok {
			report.Energy += consumer.DailyEnergy()
		}
	}
	return report
}

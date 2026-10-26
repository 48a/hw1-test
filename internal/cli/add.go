package cli

import (
	"fmt"

	"example.com/zoo/internal/domain"
)

func (c *console) addTransport() error {
	kind, err := c.readKind("Вид транспорта (1 — электросамокат, 2 — велосипед, 3 — электровелосипед): ")
	if err != nil {
		return err
	}
	data, err := c.readTransport(kind)
	if err != nil {
		return err
	}
	transport, err := createTransport(kind, data)
	if err != nil {
		return err
	}
	if err := c.park.AddTransport(transport); err != nil {
		return err
	}
	c.printAccepted(transport)
	return nil
}

func createTransport(kind string, data transportInput) (domain.Transport, error) {
	switch kind {
	case "1":
		return domain.NewScooter(data.name, data.number, data.ease, data.daily, data.operational)
	case "2":
		return domain.NewBicycle(data.name, data.number, data.ease, data.operational)
	default:
		return domain.NewElectricBicycle(data.name, data.number, data.ease, data.daily, data.operational)
	}
}

func (c *console) addItem() error {
	kind, err := c.readKind("Вид вещи (1 — шлем, 2 — док-станция без зарядки, 3 — зарядный шкаф): ")
	if err != nil {
		return err
	}
	data, err := c.readAsset()
	if err != nil {
		return err
	}
	item, err := c.createItem(kind, data)
	if err != nil {
		return err
	}
	if err := c.park.AddItem(item); err != nil {
		return err
	}
	c.printAccepted(item)
	return nil
}

func (c *console) createItem(kind string, data assetInput) (domain.Inventory, error) {
	switch kind {
	case "1":
		return domain.NewHelmet(data.name, data.number)
	case "2":
		return domain.NewDockStation(data.name, data.number)
	default:
		daily, err := c.readEnergy()
		if err != nil {
			return nil, err
		}
		return domain.NewChargingCabinet(data.name, data.number, daily)
	}
}

func (c *console) printAccepted(item domain.Inventory) {
	fmt.Fprintf(c.out, "Принято: %s [%s]\n", item.Name(), item.InventoryNumber())
}

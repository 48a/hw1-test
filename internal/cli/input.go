package cli

import (
	"errors"
	"fmt"
	"strconv"

	"example.com/zoo/internal/domain"
)

const (
	minimumEase = 1
	maximumEase = 10
)

var (
	errUnknownKind  = errors.New("неизвестный вид")
	errInvalidState = errors.New("техническое состояние должно быть 1 (исправен) или 2 (неисправен)")
)

type assetInput struct{ name, number string }

type transportInput struct {
	assetInput
	ease        int
	operational bool
	daily       float64
}

func (c *console) readKind(prompt string) (string, error) {
	kind, err := c.read(prompt)
	if err != nil {
		return "", err
	}
	switch kind {
	case "1", "2", "3":
		return kind, nil
	default:
		return "", fmt.Errorf("%w: %s", errUnknownKind, kind)
	}
}

func (c *console) readAsset() (assetInput, error) {
	name, err := c.read("Название: ")
	if err != nil {
		return assetInput{}, err
	}
	if name == "" {
		return assetInput{}, domain.ErrInvalidAsset
	}
	number, err := c.read("Инвентарный номер (положительное целое число): ")
	if err != nil {
		return assetInput{}, err
	}
	if number == "" {
		return assetInput{}, domain.ErrInvalidAsset
	}
	return assetInput{name: name, number: number}, nil
}

func (c *console) readEase() (int, error) {
	text, err := c.read("Простота для новичка (1–10): ")
	if err != nil {
		return 0, err
	}
	ease, err := strconv.Atoi(text)
	if err != nil || ease < minimumEase || ease > maximumEase {
		return 0, domain.ErrInvalidEase
	}
	return ease, nil
}

func (c *console) readState() (bool, error) {
	text, err := c.read("Техническое состояние (1 — исправен, 2 — неисправен): ")
	if err != nil {
		return false, err
	}
	switch text {
	case "1":
		return true, nil
	case "2":
		return false, nil
	default:
		return false, errInvalidState
	}
}

func (c *console) readEnergy() (float64, error) {
	text, err := c.read("Суточный расход энергии (кВт·ч, дробная часть через точку): ")
	if err != nil {
		return 0, err
	}
	daily, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, domain.ErrInvalidEnergy
	}
	return daily, nil
}

func (c *console) readTransport(kind string) (transportInput, error) {
	a, err := c.readAsset()
	if err != nil {
		return transportInput{}, err
	}
	ease, err := c.readEase()
	if err != nil {
		return transportInput{}, err
	}
	operational, err := c.readState()
	if err != nil {
		return transportInput{}, err
	}
	var daily float64
	if kind != "2" {
		daily, err = c.readEnergy()
		if err != nil {
			return transportInput{}, err
		}
	}
	return transportInput{assetInput: a, ease: ease, operational: operational, daily: daily}, nil
}

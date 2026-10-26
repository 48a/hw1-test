// Package domain содержит правила учёта транспорта и вещей без ввода-вывода.
package domain

import "errors"

var (
	// ErrInvalidAsset означает некорректное название или инвентарный номер.
	ErrInvalidAsset = errors.New("название обязательно, инвентарный номер должен быть положительным целым числом")
	// ErrInvalidEase означает простоту вне диапазона 1–10.
	ErrInvalidEase = errors.New("простота должна быть целым числом от 1 до 10")
	// ErrInvalidEnergy означает некорректный суточный расход.
	ErrInvalidEnergy = errors.New("суточный расход энергии должен быть конечным неотрицательным числом")
	// ErrDuplicateNumber означает повторный номер на балансе.
	ErrDuplicateNumber = errors.New("инвентарный номер уже занят")
	// ErrInspectionRejected означает отказ техосмотра.
	ErrInspectionRejected = errors.New("техосмотр отклонил транспорт")
	// ErrTransportAsItem запрещает обход техосмотра через добавление вещи.
	ErrTransportAsItem = errors.New("транспорт нужно добавлять через приёмку транспорта")
	// ErrServiceRequired означает отсутствие сервиса техосмотра.
	ErrServiceRequired = errors.New("нужен сервис техосмотра")
)

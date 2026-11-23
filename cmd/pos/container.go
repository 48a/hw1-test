package main

import (
	"example.com/zoo/internal/app"
	"example.com/zoo/internal/csvexport"
	"example.com/zoo/internal/domain"
	"example.com/zoo/internal/memory"
)

// container хранит зависимости одного сеанса; сборка проверяется компилятором.
type container struct {
	service *app.Service
}

func newContainer() container {
	repository := &memory.Repository{}
	exporter := csvexport.Exporter{}
	current := domain.NewCalculator(domain.MissingAsZero{})
	final := domain.NewCalculator(domain.RequireAll{})
	return container{service: app.New(repository, exporter, current, final)}
}

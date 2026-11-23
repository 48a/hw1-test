# C4-PlantUML — макросы C4 для PlantUML

Лучшая раскладка из текстовых нотаций (Graphviz `dot` маршрутизирует стрелки
вокруг блоков). Цена: нужны Java + Graphviz (или `!pragma layout smetana` —
чисто-Java раскладка чуть хуже). Не рендерится в README.

## Скелет контейнерной диаграммы

```plantuml
@startuml containers
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/v2.12.0/C4_Container.puml

title Контейнеры интернет-банка

Person(customer, "Клиент")

System_Boundary(bank, "Интернет-банк") {
    Container(web, "Web-приложение", "React", "SPA в браузере клиента")
    Container(api, "API", "Go", "REST API")
    ContainerDb(db, "База данных", "PostgreSQL", "Счета и платежи")
}

System_Ext(mainframe, "Мейнфрейм банка")

Rel(customer, web, "Открывает", "HTTPS")
Rel(web, api, "Вызывает", "HTTPS/JSON")
Rel(api, db, "Читает и пишет", "SQL")
Rel(api, mainframe, "Проверяет лимиты", "HTTPS")

@enduml
```

## Файлы макросов по уровням

| Уровень | include |
|---|---|
| L1 Context | `C4_Context.puml` |
| L2 Container | `C4_Container.puml` (включает Context) |
| L3 Component | `C4_Component.puml` |
| Deployment | `C4_Deployment.puml` |
| Dynamic | `C4_Dynamic.puml` |

## Правила

1. **Пиновать версию include** (`.../v2.12.0/...`), не `master`: переименование макроса в апстриме ломает диаграмму.
2. Оффлайн: макросы входят в jar PlantUML — `!include <C4/C4_Container>` без сети.
3. `@startuml <имя>` задаёт имя выходного файла — ставь всегда, иначе коллизии при пакетном рендере.
4. Направление связей для раскладки: `Rel_U/Rel_D/Rel_L/Rel_R`; `Lay_*` для подсказок Graphviz.
5. Проверка в CI: `plantuml -checkonly file.puml` — быстро, без рендера.

## ГЛАВНАЯ ГРАБЛЯ: опечатка молча рисуется

`Rel(api, dbb, ...)` при несуществующем `dbb` не падает: PlantUML рисует
«выдуманный» элемент и выходит с кодом 0. Поэтому:

- всегда `-checkonly` + зрительная проверка на лишние элементы;
- или держи модель в Structurizr DSL / LikeC4 и экспортируй в PlantUML —
  опечатки отловит валидатор модели.

## Рендер

```bash
plantuml -tsvg diagram.puml    # SVG
plantuml -tpng diagram.puml    # PNG
```

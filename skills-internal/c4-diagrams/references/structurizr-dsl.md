# Structurizr DSL — одна модель, много видов

Инструмент от автора C4. Модель описывается один раз в `workspace.dsl`,
виды (диаграммы) декларируются над моделью: переименование контейнера — одна
правка. Требует JVM или Docker. Экспортирует виды в PlantUML, Mermaid, DOT.

## Скелет workspace

```text
workspace "Интернет-банк" "Учебный пример" {

    !identifiers hierarchical

    model {
        customer = person "Клиент банка" "Клиент со счетами"
        mainframe = softwareSystem "Мейнфрейм банка" "Счета и транзакции" {
            tags "External"
        }

        bank = softwareSystem "Интернет-банк" "Просмотр счетов и платежи" {
            webapp = container "Web-приложение" "SPA" "React"
            api = container "API" "REST API" "Go"
            db = container "База данных" "Счета и платежи" "PostgreSQL" {
                tags "Database"
            }
        }

        customer -> webapp "Открывает" "HTTPS"
        webapp -> api "Вызывает" "HTTPS/JSON"
        api -> db "Читает и пишет" "SQL"
        api -> mainframe "Проверяет лимиты" "HTTPS"
    }

    views {
        systemContext bank "context" {
            include *
            autolayout lr
        }

        container bank "containers" {
            include *
            autolayout lr
        }

        styles {
            element "Database" {
                shape Cylinder
            }
            element "External" {
                background "#999999"
            }
        }
    }
}
```

## Ключевые правила DSL

1. `!identifiers hierarchical` — ссылки вида `bank.api.db`; без неё — плоские имена. Включай всегда в новых файлах.
2. Иерархия строгая: `container` только внутри `softwareSystem`, `component` только внутри `container`. DSL это валидирует.
3. Виды: `systemContext`, `container`, `component`, `deployment`, `dynamic`. В каждом виде `include *` + `autolayout lr|tb|bt|rl`.
4. `dynamic` — сценарий: шаги `a -> b "описание"` внутри `dynamic bank "..." { ... }`.
5. Теги через `tags "..."` и оформление в блоке `styles` — не размазывай цвета по элементам.

## CLI

```bash
# Валидация (называет файл и строку при ошибке)
structurizr validate -w workspace.dsl

# Экспорт всех видов в Mermaid / PlantUML (c4plantuml)
structurizr export -w workspace.dsl -f mermaid -o out-mmd
structurizr export -w workspace.dsl -f plantuml/c4plantuml -o out-puml
```

Docker: `docker run -it --rm -v "$PWD":/usr/local/structurizr structurizr/cli validate -w workspace.dsl`

## Грабли

- Сам ничего не рисует: для картинок нужен экспорт + рендерер (PlantUML/Mermaid) или Structurizr Lite.
- Экспорт в Mermaid генерирует `graph TB` с инлайн-стилями, а не синтаксис C4 — это нормально: так раскладка лучше.
- Ссылки на несуществующий элемент = ошибка валидации с номером строки (плюс против PlantUML, который молча рисует).
- Тулинг Structurizr в 2026 в процессе консолидации (CLI/Lite/cloud); сам DSL — де-факто стандарт и читается и другими инструментами. Если модель-как-код нужна без JVM — LikeC4 (см. навык likec4-dsl).

## Когда Structurizr DSL — правильный выбор

- Системе нужно 2+ вида (context + container + deployment), и они должны быть согласованы.
- Диаграммы живут в репозитории годами и переживают смену рендерера.
- Нужна машинная проверка ссылок в CI.

# Mermaid C4 — синтаксис и грабли

Рендерится нативно в GitHub/GitLab Markdown, wiki, статических генераторах.
Статус: experimental — синтаксис может меняться между версиями Mermaid.
Типы: `C4Context`, `C4Container`, `C4Component`, `C4Dynamic`, `C4Deployment`.

## Контекстная диаграмма (L1)

```mermaid
C4Context
    title Контекстная диаграмма интернет-банка

    Person(customer, "Клиент банка", "Клиент, у которого есть счета")
    System(banking, "Интернет-банк", "Просмотр счетов и платежи")
    System_Ext(mainframe, "Мейнфрейм банка", "Хранит счета и транзакции")
    System_Ext(mail, "Почтовая система", "SMTP")

    Rel(customer, banking, "Использует", "HTTPS")
    Rel(banking, mainframe, "Читает счета", "HTTPS")
    Rel(banking, mail, "Отправляет письма", "SMTP")
```

## Контейнерная диаграмма (L2)

```mermaid
C4Container
    title Контейнеры интернет-банка

    Person(customer, "Клиент")

    System_Boundary(bank, "Интернет-банк") {
        Container(web, "Web-приложение", "React", "SPA в браузере клиента")
        Container(api, "API", "Go", "REST API")
        ContainerDb(db, "База данных", "PostgreSQL", "Счета и платежи")
        ContainerQueue(bus, "Шина событий", "Kafka", "События платежей")
    }

    System_Ext(mainframe, "Мейнфрейм банка")

    Rel(customer, web, "Открывает", "HTTPS")
    Rel(web, api, "Вызывает", "HTTPS/JSON")
    Rel(api, db, "Читает и пишет", "SQL")
    Rel(api, bus, "Публикует события", "Kafka")
    Rel(api, mainframe, "Проверяет лимиты", "HTTPS")
```

## Справочник конструкций

- Люди: `Person(alias, "Имя", "Описание")`, внешний: `Person_Ext`
- Системы: `System`, `System_Ext` (внешняя), `SystemDb`, `SystemQueue`
- Контейнеры: `Container(alias, "Имя", "Технология", "Описание")`, `ContainerDb`, `ContainerQueue`; внешний: `Container_Ext`
- Компоненты: `Component`, `ComponentDb`, `Component_Ext`
- Границы: `System_Boundary(alias, "Имя") { ... }`, `Container_Boundary`, `Enterprise_Boundary`
- Связи: `Rel(a, b, "Действие", "Протокол")`; обратная: `BiRel`; варианты: `Rel_U/Rel_D/Rel_L/Rel_R` (направление для раскладки)
- Стили: `UpdateElementStyle(alias, $bgColor="...")`, `UpdateRelStyle(...)`
- Раскладка: `UpdateLayoutConfig($c4ShapeInRow="3", $c4BoundaryInRow="1")`

## Грабли

1. **Нет маршрутизации стрелок.** Mermaid ставит блоки сеткой в порядке объявления, стрелки идут сквозь блоки. Меняй порядок объявления и `Rel_L/Rel_R`/`Rel_U/Rel_D` для управления.
2. **`UpdateLayoutConfig` почти не помогает** на плотных диаграммах — перейди на PlantUML, если раскладка нечитаема.
3. **L3 (Component) с вложенными границами глючит** — fallback: обычный `flowchart` с `subgraph` для границ (теряется словарь C4, но рендер стабилен).
4. Стрелка с опечаткой в alias → ошибка рендера (это хорошо: опечатки видны).
5. Проверка без браузера: `mmdc -i diagram.mmd -o diagram.svg` (Mermaid CLI).

## Динамическая диаграмма (сценарий)

```mermaid
C4Dynamic
    title Оплата заказа

    Container(api, "API", "Go")
    ContainerQueue(bus, "Шина", "Kafka")
    Container(worker, "Воркер платежей", "Go")
    ContainerDb(db, "БД", "PostgreSQL")

    Rel(api, bus, "1. Публикует OrderCreated")
    Rel(worker, bus, "2. Читает событие")
    Rel(worker, db, "3. Списывает средства", "SQL")
```

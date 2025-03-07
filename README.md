# go-rest-template

Проект представляет собой шаблон для создания Go-приложений с расширенным набором функционала.
Go-rest-template реализует повседневные рутины (логирование, graceful shutdown, билд, и т.д.) и дает контроль над работой подпрограмм, оркестрируемых снаружи.

Просто хочется сэкономить немножко времени в создании новых приложений :)

## Быстрый старт

1. Склонируйте репозиторий
2. Запустите скрипт настройки:
   ```bash
   ./setup.sh
   ```
3. Следуйте интерактивным подсказкам для ввода нового имени проекта

Скрипт автоматически:
- Заменит все вхождения `go-rest-template` на ваше новое имя проекта
- Заменит все вхождения `PROJECTNAME` на ваше новое имя проекта
- Обновит все импорты в Go файлах
- Удалит директории `.idea` и `.git`
- Инициализирует новый Git репозиторий
- Удалит сам себя после выполнения

## Функционал

- HTTP сервер на базе chi/v5 с:
    - Graceful shutdown
    - Middleware (CORS, компрессия, логирование)
    - Версионированным API
    - Тестовыми эндпоинтами
- Работа с базами данных:
    - Абстракция датастора
    - MongoDB драйвер
    - Поддержка индексов
- Управление приложением:
    - Конфигурация через переменные окружения
    - Менеджер компонентов
    - Логирование
- Тестирование:
    - Unit тесты
    - Скрипт для проверки покрытия
- Docker поддержка:
    - Multi-stage сборка
    - Compose для разработки
    - Отдельный образ для тестов

## Структура проекта

```
.
├── Dockerfile              # Основной Docker образ
├── Dockerfile-unit         # Docker образ для тестов
├── Makefile               # Основные команды сборки
├── docker-compose.yml     # Настройка окружения разработки
├── app
│   ├── Makefile          # Команды для приложения
│   ├── cmd               # Команды приложения
│   │   ├── config.go     # Конфигурация
│   │   ├── http.go       # HTTP сервер
│   │   ├── log.go        # Настройка логирования
│   │   ├── root.go       # Корневая команда
│   │   └── server.go     # Управление сервером
│   ├── database          # Работа с БД
│   │   ├── drivers       # Драйверы БД
│   │   │   ├── mongo     # MongoDB драйвер
│   │   │   └── ...
│   │   ├── errors.go     # Ошибки БД
│   │   └── factory.go    # Фабрика датасторов
│   ├── director          # Управление компонентами
│   ├── manager           # Менеджер приложения
│   ├── middleware        # HTTP middleware
│   ├── resources         # HTTP хендлеры
│   │   ├── api          # API эндпоинты
│   │   └── version.go   # Информация о версии
│   └── utils            # Утилиты
├── scripts
│   └── coverage.sh      # Проверка test coverage
└── setup.sh             # Скрипт настройки проекта
```

## Ручная настройка

Если вы предпочитаете настроить проект вручную, можно использовать следующие команды:

Для Linux:
```bash
sed -i s/go-rest-template/NEW_NAME/g $(find . -type f -not -path "*/\.git/*" -not -path "*/\.idea/*")
sed -i s/PROJECTNAME/NEW_NAME/g $(find . -type f -not -path "*/\.git/*" -not -path "*/\.idea/*")
```

Для MacOS:
```bash
find . -type f -not -path "*/\.git/*" -not -path "*/\.idea/*" -exec sed -i '' "s/go-rest-template/NEW_NAME/g" {} \;
find . -type f -not -path "*/\.git/*" -not -path "*/\.idea/*" -exec sed -i '' "s/PROJECTNAME/NEW_NAME/g" {} \;
```

## Зависимости

- chi/v5 - HTTP роутер
- cors - CORS middleware
- MongoDB драйвер
- env - Работа с переменными окружения

## Разработка

1. Клонируйте репозиторий
2. Настройте проект через `setup.sh`
3. Запустите зависимости через Docker Compose:
   ```bash
   docker-compose up -d
   ```
4. Запустите приложение:
   ```bash
   make run
   ```

## Тестирование

```bash
# Запуск всех тестов
make test

# Проверка coverage
./scripts/coverage.sh
```
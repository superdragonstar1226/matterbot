## Setup
 - Ensure Docker-Compose are installed for your system
 -  Run docker-compose up -d --build and the mattermost server will be built and will expose the `port 8065` to your system's localhost
 - Start the Bot: 
 ```shell
 go run main.go
 ```

## Команды
 - !помощь - выдает список команд
 - !старт - начинает отчет времени
 - !отчет - следующее сообщение будет распознанно как текст отчета 
 - !конец - останавливает таймер, записывает данные в отчет\

## Config discription:  
"USER_EMAIL"   : email для логина бота\
"USER_PASSWORD": пароль для логина бота\
"USER_NAME"    : Имя бота\
"TEAM_NAME"    : Имя команды (воркспейса)\
"SERVER_ADRES" : адрес сервера\
"SERVER_WS_ADDRESS": адрес веб-сокета\
"TIME_BEFORE_NOTIFICATION": время перед отправкой уведомления о ненаписанном отчете\

# Файл отчета содержит:
- имя пользователя (заголовок)
- время начала работы
- текст отчета
- время завершения работы

# TODOs
- Добавление глобального логгера (врап zap-error, вынести конфиги логгера в конфигурационный файл), применять к сорцам, тестам етц.
- Разбить базовую бизнес-логику по пакетам. Важно: комменты по возможности.

`Пример`: 
```go
// CreateEntity parse and creates new entity instance
// and write it to database.
func CreateEntity() {}
```
- Создание, коннект к базе для хранения репортов бота.
Опционально: хранить репорты и в *.json
Как пример:
таблица репортов
таблица юзеров (привязка юзер-репорт)
- Создание пакета mattermost (слежебные методы, бизнес-логика етц)
- Поднятие `http-сервера` (если нужен - соотнести с `тз`, документацией mattermost API етц)
- По документации `mattermost API`: исследовать вопрос аутентификации (bearer token) и какие методы необходимы.
- `Makefile`: по мере необходимости доработать пайплайн билда, логирования. Включает билд и старт всех докер-контейнеров, создание и дроп баз данных, применение и откат миграций.
билд проекта, служебные комманды
- Для каждого пакета из вышеперечисленных написание `unit-тестов` (Пример в `pkg/mattermost`)

`Пример`:
```go
func someMethod() {}

func Test_someMethod(t *testing.T) {}
```

- Посмотреть работу с `database/sql` для вз-ия с бд.

- Изучить по возможности [migrate-tool](https://github.com/golang-migrate/migrate)

### Команды бота	

- `!План`=старт

```триггерит таймер```

- `!Отчет`=стоп

```command description``` 
```триггерит остановку таймер, запись => redmine```

- `!help` = help list

- `!unknown` = default

```игнор/skip/unknown command```
## Setup
 - Ensure Docker-Compose are installed for your system
 -  Run docker-compose up -d --build and the mattermost server will be built and will expose the `port 8065` to your system's localhost
 - Start the Bot: 
 ```shell
 go run main.go
 ```

## Teams
 - !help - displays a list of commands
 - !start - starts time reporting
 - !report - the next message will be recognized as report text 
 - !end - stops the timer, writes data to the report\

## Config discription:  
"USER_EMAIL"   : 
"USER_PASSWORD": 
"USER_NAME"    : 
"TEAM_NAME"    : 
"SERVER_ADRES" : 
"SERVER_WS_ADDRESS": 
"TIME_BEFORE_NOTIFICATION": 

# The report file contains:
- username (title)
- start time
- report text
- shutdown time

# TODOs
- Adding a global logger (wrap zap-error, put logger configs into a configuration file), apply to sources, tests, etc.
- Break down core business logic into packages. Important: comments whenever possible.

`Example`:
```go
// CreateEntity parse and creates new entity instance
// and write it to database.
func CreateEntity() {}
```
- Creation, connection to the database for storing bot reports.
Optionally: store reports in *.json
As an example:
report table
users table (user-report binding)
- Creation of the mattermost package (tracking methods, business logic, etc.)
- Raising the `http-server` (if needed, correlate with `ts`, documentation of the most important API ETC)
- According to the `mattermost API` documentation: research the issue of authentication (bearer token) and what methods are needed.
- `Makefile`: modify the build pipeline and logging as necessary. Includes building and starting all docker containers, creating and dropping databases, applying and rolling back migrations.
project build, service commands
- For each package from the above, writing `unit tests` (Example in `pkg/mattermost`)

`Example`:
```go
func someMethod() {}

func Test_someMethod(t *testing.T) {}
```

- View work with `database/sql` for accessing from a database.

- Study as much as possible [migrate-tool](https://github.com/golang-migrate/migrate)

### Bot commands	

- `!Plan`=start

```triggers the timer```

- `!Report`=stop

```command description``` 
```trigger stop timer, record => redmine```

- `!help` = help list

- `!unknown` = default

```ignore/skip/unknown command```

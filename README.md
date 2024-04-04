## Running the project

- Make sure you have docker installed.
- Copy `.env.example` to `.env`
- Run `docker-compose up -d`
- Go to `localhost:8080` to verify if the server works.
- [Adminer](https://www.adminer.org/) Database Management runs at `5001` .

If you are running without docker be sure database configuration is provided in `.env` file and run `go run . app:serve`

#### Environment Variables

<details>
    <summary>Variables Defined in the project </summary>

| Key            | Value                    | Desc                                  |
| -------------- | ------------------------ | ------------------------------------- |
| `SERVER_PORT`  | `5000`                   | Port at which app runs                |
| `ENV`          | `development,production` | App running Environment               |
| `LOG_OUTPUT`   | `./server.log`           | Output Directory to save logs         |
| `LOG_LEVEL`    | `info`                   | Level for logging (check lib/logger.go:172) |
| `DB_USER`      | `username`               | Database Username                     |
| `DB_PASS`      | `password`               | Database Password                     |
| `DB_HOST`      | `0.0.0.0`                | Database Host                         |
| `DB_PORT`      | `3306`                   | Database Port                         |
| `DB_NAME`      | `test`                   | Database Name                         |
| `JWT_SECRET`   | `secret`                 | JWT Token Secret key                  |
| `ADMINER_PORT` | `5001`                   | Adminer DB Port                       |
| `DEBUG_PORT`   | `5002`                   | Port that debugger runs in            |

</details>

#### Migration Commands

> ⚓️ &nbsp; Add argument `p=host` if you want to run the migration runner from the host environment instead of docker environment.
> eg; `make p=host migrate-up`

<details>
    <summary>Migration commands available</summary>

| Command             | Desc                                           |
| ------------------- | ---------------------------------------------- |
| `make migrate-up`   | runs migration up command                      |
| `make migrate-down` | runs migration down command                    |
| `make force`        | Set particular version but don't run migration |
| `make goto`         | Migrate to particular version                  |
| `make drop`         | Drop everything inside database                |
| `make create`       | Create new migration file(up & down)           |

</details>

## Implemented Features

- Dependency Injection (go-fx)
- Routing (gin web framework)
- Environment Files
- Logging (file saving on `production`) [zap](https://github.com/uber-go/zap)
- Middlewares (cors)
- Database Setup (mysql)
- Models Setup and Automigrate (gorm)
- Repositories
- Implementing Basic CRUD Operation
- Authentication (JWT)
- Migration Runner Implementation
- Live code refresh
- Dockerize Application with Debugging Support Enabled. Debugger runs at `5002`
- Cobra Commander CLI Support. try: `go run . --help`

## Todos

- [ ] Swagger documentation examples.
- [ ] Unit testing examples. 
- [ ] File upload middelware.
- [ ] Lint and coverage before commit.

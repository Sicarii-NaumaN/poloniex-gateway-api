# poloniex-gateway-api

Как запустить
```shell
$ go mod tidy
$ DB_DSN=postgresql://postgres:postgres@localhost:5432/poloniex IS_SWAGGER_CREATED=false PORT=3000 go run ./cmd
```
**Инициализация**
[main.go](cmd/main.go)

**Синк архивных свечей**

[service.go](internal/syncer/service.go) 

[processor.go](service/processor/processor.go)

**Job по созданию свечей из RecentTrades**

[service.go](internal/syncer/service.go)

[builder.go](service/builder/builder.go)

**Синк по RecentTrades по WS**

[service.go](internal/ws/service.go)

[real_time_saver.go](service/rt_saver/real_time_saver.go)



# poloniex-gateway-api

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



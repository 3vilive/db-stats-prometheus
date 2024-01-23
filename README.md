# db-stats-prometheus


## usage

1. install

```
go get -u github.com/3vilive/db-stats-prometheus
```

2. apply


```go
import "github.com/3vilive/db-stats-prometheus/tracer"

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

tracer.Init(
    ctx,
    tracer.WithCheckInterval(1*time.Second),
    tracer.WithLabels("node", "role", "type"), // custom labels
)

tracer.MustTrace("test", db1)
tracer.MustTraceGormDb("production", db2, map[string]string{
    "node": "01",
    "role": "master",
    "type": "gorm",
})
```

3. add /metrics endpoint

```go
http.Handle("/metrics", promhttp.Handler())
``` 

4. fetch metrics

```
# HELP dbstats_idle The number of idle connections.
# TYPE dbstats_idle gauge
dbstats_idle{name="production",node="01",role="master",type="gorm"} 1
dbstats_idle{name="test",node="",role="",type=""} 1
# HELP dbstats_in_use The number of connections currently in use.
# TYPE dbstats_in_use gauge
dbstats_in_use{name="production",node="01",role="master",type="gorm"} 0
dbstats_in_use{name="test",node="",role="",type=""} 0
# HELP dbstats_max_idle_closed The total number of connections closed due to SetMaxIdleConns.
# TYPE dbstats_max_idle_closed gauge
dbstats_max_idle_closed{name="production",node="01",role="master",type="gorm"} 0
dbstats_max_idle_closed{name="test",node="",role="",type=""} 0
# HELP dbstats_max_lifetime_closed The total number of connections closed due to SetConnMaxLifetime.
# TYPE dbstats_max_lifetime_closed gauge
dbstats_max_lifetime_closed{name="production",node="01",role="master",type="gorm"} 0
dbstats_max_lifetime_closed{name="test",node="",role="",type=""} 0
# HELP dbstats_max_open_connections Maximum number of open connections to the database.
# TYPE dbstats_max_open_connections gauge
dbstats_max_open_connections{name="production",node="01",role="master",type="gorm"} 0
dbstats_max_open_connections{name="test",node="",role="",type=""} 0
# HELP dbstats_open_connections The number of established connections both in use and idle.
# TYPE dbstats_open_connections gauge
dbstats_open_connections{name="production",node="01",role="master",type="gorm"} 1
dbstats_open_connections{name="test",node="",role="",type=""} 1
# HELP dbstats_wait_count The total number of connections waited for.
# TYPE dbstats_wait_count gauge
dbstats_wait_count{name="production",node="01",role="master",type="gorm"} 0
dbstats_wait_count{name="test",node="",role="",type=""} 0
# HELP dbstats_wait_duration The total time blocked waiting for a new connection.
# TYPE dbstats_wait_duration gauge
dbstats_wait_duration{name="production",node="01",role="master",type="gorm"} 0
dbstats_wait_duration{name="test",node="",role="",type=""} 0
```
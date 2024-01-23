# db-stats-prometheus

`db-stats-prometheus` which provides a way to collect and expose database statistics as Prometheus metrics. It includes a tracer that can be used to track and monitor database connections in your application.

## install

To install the `db-stats-prometheus` package, you can use the following command:

```
go get -u github.com/3vilive/db-stats-prometheus
```

## Usage

To use the db-stats-prometheus package, you need to import the tracer package from github.com/3vilive/db-stats-prometheus/tracer. Here's an example of how to initialize and use the tracer:

```go
import (
    "context"
    "time"

    "github.com/3vilive/db-stats-prometheus/tracer"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    tracer.Init(
        ctx,
        tracer.WithCheckInterval(time.Second),
        tracer.WithLabels("node", "role", "type"), // custom labels
    )

    // Init your db instances
    // ...
    
    // Start trace stats of your db instance
    tracer.MustTrace("test", db1)
    tracer.MustTraceGormDb("production", db2, map[string]string{
        "node": "01",
        "role": "master",
        "type": "gorm",
    })

    // Add /metrics endpoint to expose the collected metrics
    http.Handle("/metrics", promhttp.Handler())

    // Start your application server
    // ...
}
```

## Metrics

The db-stats-prometheus package exposes the following metrics:

- `dbstats_idle`: The number of idle connections.
- `dbstats_in_use`: The number of connections currently in use.
- `dbstats_max_idle_closed`: The total number of connections closed due to SetMaxIdleConns.
- `dbstats_max_lifetime_closed`: The total number of connections closed due to SetConnMaxLifetime.
- `dbstats_max_open_connections`: Maximum number of open connections to the database.
- `dbstats_open_connections`: The number of established connections both in use and idle.
- `dbstats_wait_count`: The total number of connections waited for.
- `dbstats_wait_duration`: The total time blocked waiting for a new connection.

These metrics are labeled with additional information such as the database name, node, role, and type.

You can fetch the metrics by making a GET request to the /metrics endpoint of your application.

```
GET /metrics
```

The response will include the metrics in a format similar to the following:

```
# HELP dbstats_idle The number of idle connections.
# TYPE dbstats_idle gauge
dbstats_idle{name="production",node="01",role="master",type="gorm"} 1
dbstats_idle{name="test",node="",role="",type=""} 1
# ...
```

# db-stats-prometheus


## usage

```go
import "github.com/3vilive/db-stats-prometheus/tracer"

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

tracer.Init(
    ctx,
    tracer.WithCheckInterval(1*time.Second),
    tracer.WithLabels("node", "role", "type"), // custom labels
)

tracer.MustTrace(db1, "test")
tracer.MustTraceGormDb(db2, "production", map[string]string{
    "node": "01",
    "role": "master",
    "type": "gorm",
})
```
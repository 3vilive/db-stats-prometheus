package tracer

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/3vilive/db-stats-prometheus/metrics"
)

type Tracer struct {
	ctx              context.Context
	mutex            sync.Mutex
	checkInterval    time.Duration
	dbStatsCollector *metrics.DbStatsCollector
	dbLabelMap       map[*sql.DB]map[string]string
}

func NewTracer(ctx context.Context, configs ...ApplyConfig) *Tracer {
	config := DefaultConfig()

	for _, apply := range configs {
		apply(&config)
	}

	tracer := &Tracer{
		ctx:              ctx,
		checkInterval:    config.CheckInterval,
		dbStatsCollector: metrics.NewDbStatsCollector(config.Labels),
		dbLabelMap:       make(map[*sql.DB]map[string]string),
	}

	tracer.Start()

	return tracer
}

func (t *Tracer) Trace(db *sql.DB, dbName string, labels ...map[string]string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	dbLabels := make(map[string]string)
	dbLabels["name"] = dbName
	for _, labelMap := range labels {
		for k, v := range labelMap {
			dbLabels[k] = v
		}
	}

	t.dbLabelMap[db] = dbLabels
}

func (t *Tracer) Start() {
	go func() {
		ticker := time.NewTicker(t.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-t.ctx.Done():
				return
			case <-ticker.C:
				t.Check()
			}
		}
	}()
}

func (t *Tracer) Check() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for db, labels := range t.dbLabelMap {
		t.dbStatsCollector.Set(db.Stats(), labels)
	}
}

package tracer

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/3vilive/db-stats-prometheus/metrics"
)

type dbWrapper struct {
	db     *sql.DB
	labels map[string]string
}

type Tracer struct {
	ctx              context.Context
	mutex            sync.Mutex
	checkInterval    time.Duration
	dbStatsCollector *metrics.DbStatsCollector
	dbMap            map[string]*dbWrapper
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
		dbMap:            make(map[string]*dbWrapper),
	}

	tracer.Start()

	return tracer
}

func (t *Tracer) Trace(dbName string, db *sql.DB, labels ...map[string]string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	dbLabels := make(map[string]string)
	for _, labelMap := range labels {
		for k, v := range labelMap {
			dbLabels[k] = v
		}
	}
	dbLabels["name"] = dbName

	t.dbMap[dbName] = &dbWrapper{
		db:     db,
		labels: dbLabels,
	}
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

	for _, dbWrapper := range t.dbMap {
		t.dbStatsCollector.Set(dbWrapper.db.Stats(), dbWrapper.labels)
	}
}

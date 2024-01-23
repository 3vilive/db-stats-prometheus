package tracer

import (
	"context"
	"database/sql"
	"sync"
)

var gInitTracerOnce sync.Once
var gTracer *Tracer

type DBGetter interface {
	DB() (*sql.DB, error)
}

func Init(ctx context.Context, configs ...ApplyConfig) {
	gInitTracerOnce.Do(func() {
		gTracer = NewTracer(ctx, configs...)
	})
}

func Trace(dbName string, db *sql.DB, labels ...map[string]string) error {
	if gTracer == nil {
		Init(context.Background())
	}

	gTracer.Trace(dbName, db, labels...)
	return nil
}

func MustTrace(dbName string, db *sql.DB, labels ...map[string]string) {
	if err := Trace(dbName, db, labels...); err != nil {
		panic(err)
	}
}

func TraceGormDb(dbName string, db DBGetter, labels ...map[string]string) error {
	internalDb, err := db.DB()
	if err != nil {
		return err
	}

	return Trace(dbName, internalDb, labels...)
}

func MustTraceGormDb(dbName string, db DBGetter, labels ...map[string]string) {
	if err := TraceGormDb(dbName, db, labels...); err != nil {
		panic(err)
	}
}

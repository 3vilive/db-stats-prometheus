package tracer

import (
	"context"
	"database/sql"
	"errors"
)

var gTracer *Tracer

type DBGetter interface {
	DB() (*sql.DB, error)
}

func Init(ctx context.Context, configs ...ApplyConfig) {
	if gTracer != nil {
		return
	}

	gTracer = NewTracer(ctx, configs...)
}

func Trace(db *sql.DB, dbName string, labels ...map[string]string) error {
	if gTracer == nil {
		return errors.New("tracer not initialized")
	}

	gTracer.Trace(db, dbName, labels...)
	return nil
}

func MustTrace(db *sql.DB, dbName string, labels ...map[string]string) {
	if err := Trace(db, dbName, labels...); err != nil {
		panic(err)
	}
}

func TraceGormDb(db DBGetter, dbName string, labels ...map[string]string) error {
	internalDb, err := db.DB()
	if err != nil {
		return err
	}

	return Trace(internalDb, dbName, labels...)
}

func MustTraceGormDb(db DBGetter, dbName string, labels ...map[string]string) {
	if err := TraceGormDb(db, dbName, labels...); err != nil {
		panic(err)
	}
}

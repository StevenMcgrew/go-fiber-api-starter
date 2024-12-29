package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
)

func ExecuteSqlFile(path string) {
	path, pathErr := filepath.Abs(path)
	if pathErr != nil {
		log.Fatal("Error getting absolute path to sql file:", pathErr)
	}
	bytes, ioErr := os.ReadFile(path)
	if ioErr != nil {
		log.Fatal("Error reading sql file: ", ioErr)
	}
	sql := string(bytes)
	_, execErr := Pool.Exec(Ctx, sql)
	if execErr != nil {
		log.Fatal("Error executing sql from file: ", execErr)
	}
}

func Many[T any](sql string, args pgx.NamedArgs, ptrModel *T) ([]T, error) {
	// Run the query
	rows, err := Pool.Query(Ctx, sql, args)
	if err != nil {
		return nil, err
	}
	// Parse the rows
	parsedRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, err
	}
	return parsedRows, nil
}

func One[T any](sql string, args pgx.NamedArgs, ptrModel *T) (T, error) {
	// Run the query
	rows, err := Pool.Query(Ctx, sql, args)
	if err != nil {
		return *ptrModel, err
	}
	// Parse the rows
	parsedRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return *ptrModel, err
	}
	if len(parsedRows) == 0 {
		return *ptrModel, pgx.ErrNoRows
	}
	return parsedRows[0], nil
}

func None(sql string, args pgx.NamedArgs) error {
	_, err := Pool.Exec(Ctx, sql, args)
	if err != nil {
		return err
	}
	return nil
}

package db

import (
	"fmt"
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

func GetRowCount(tableName string) (uint, error) {
	sql := fmt.Sprintf("SELECT COUNT(*) FROM %s;", tableName)
	// Run the query
	rows, err := Pool.Query(Ctx, sql)
	if err != nil {
		return 0, err
	}
	// Parse the rows
	parsedRows, err := pgx.CollectRows(rows, pgx.RowTo[uint])
	if err != nil {
		return 0, err
	}
	if len(parsedRows) != 1 {
		return 0, fmt.Errorf("error getting count of rows from database")
	}
	count := parsedRows[0]
	return count, nil
}

type QueryBuilder struct {
	page       uint
	perPage    uint
	tableName  string
	fieldNames []string
}

func NewQueryBuilder(page uint, perPage uint, tableName string, fieldNames []string) QueryBuilder {
	return QueryBuilder{
		page:       page,
		perPage:    perPage,
		tableName:  tableName,
		fieldNames: fieldNames,
	}
}

func (qb QueryBuilder) BuildQuery() (string, error) {

	// Start of query string
	q := fmt.Sprintf("SELECT * FROM %s ", qb.tableName)

	// TODO: add other query options

	// append LIMIT and OFFSET
	q += fmt.Sprintf("LIMIT %d OFFSET (%d - 1) * %d;", qb.perPage, qb.page, qb.perPage)

	return q, nil
}

package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

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
	query      string
	tableName  string
	fieldNames []string
	keywords   map[string]string
}

func NewQueryBuilder(page uint, perPage uint, query string, tableName string, fieldNames []string) QueryBuilder {
	return QueryBuilder{
		page:       page,
		perPage:    perPage,
		query:      query,
		tableName:  tableName,
		fieldNames: fieldNames,
		keywords: map[string]string{
			"where":       "WHERE",
			"eq":          "=",
			"not_eq":      "!=",
			"gt":          ">",
			"lt":          "<",
			"gt_eq":       ">=",
			"lt_eq":       "<=",
			"and":         "AND",
			"or":          "OR",
			"not":         "NOT",
			"between":     "BETWEEN",
			"not_between": "NOT BETWEEN",
			"in":          "IN",
			"not_in":      "NOT IN",
			"is_null":     "IS NULL",
			"is_not_null": "IS NOT NULL",
			"like":        "LIKE",
			"ilike":       "ILIKE",
			"order_by":    "ORDER BY",
			"asc":         "ASC",
			"desc":        "DESC",
		},
	}
}

func (qb QueryBuilder) Build() (string, error) {
	// Split by dot
	dotSplitWords := strings.Split(qb.query, ".")

	// Start of query string
	q := fmt.Sprintf("SELECT * FROM %s ", qb.tableName)

	// Transform and append the dotSplitWords to the query string
	for _, word := range dotSplitWords {

		if strings.Contains(word, ",") {
			// Split the word by comma
			commaSplitWords := strings.Split(word, ",")
			// Transform and append the commaSplitWords to the query string
			for _, w := range commaSplitWords {
				if kw, ok := qb.keywords[w]; ok {
					q += kw + ", "
					continue
				}
				if slices.Contains(qb.fieldNames, w) {
					q += w + ", "
					continue
				}
				q += "'" + w + "', "
			}
			// Remove last space and comma and append a single space
			q = q[:len(q)-2] + " "
			continue
		}

		if keyword, ok := qb.keywords[word]; ok {
			q += keyword + " "
			continue
		}
		if slices.Contains(qb.fieldNames, word) {
			q += word + " "
			continue
		}
		q += "'" + word + "' "
	}

	// append LIMIT and OFFSET
	q += fmt.Sprintf("LIMIT %d OFFSET (%d - 1) * %d;", qb.perPage, qb.page, qb.perPage)

	return q, nil
}

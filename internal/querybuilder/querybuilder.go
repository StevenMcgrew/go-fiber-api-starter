package querybuilder

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

type QueryBuilder struct {
	page       int
	perPage    int
	query      string
	tableName  string
	fieldNames []string
	keywords   map[string]string
}

func New(page int, perPage int, query string, tableName string, fieldNames []string) QueryBuilder {
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
			"starts_with": "'%s%'",
			"ends_with":   "'%%s'",
			"contains":    "'%%s%'",
			"order_by":    "ORDER BY",
			"asc":         "ASC",
			"desc":        "DESC",
		},
	}
}

func (qb QueryBuilder) Build() (string, error) {
	// Unescape query (normally the query is from a url)
	unescaped, err := url.QueryUnescape(qb.query)
	if err != nil {
		return "", err
	}

	// Split by dot
	dotSplitWords := strings.Split(unescaped, ".")

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

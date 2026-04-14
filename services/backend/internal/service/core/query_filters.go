package core

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func combineFilters(filters ...bson.M) bson.M {
	parts := make([]bson.M, 0, len(filters))
	for _, filter := range filters {
		if len(filter) == 0 {
			continue
		}
		parts = append(parts, filter)
	}

	switch len(parts) {
	case 0:
		return bson.M{}
	case 1:
		return parts[0]
	default:
		return bson.M{"$and": parts}
	}
}

func buildTextSearchClause(search string, fields ...string) bson.M {
	trimmed := strings.TrimSpace(search)
	if trimmed == "" || len(fields) == 0 {
		return bson.M{}
	}

	clauses := make([]bson.M, 0, len(fields))
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			continue
		}
		clauses = append(clauses, bson.M{field: bson.M{"$regex": trimmed, "$options": "i"}})
	}

	if len(clauses) == 0 {
		return bson.M{}
	}

	return bson.M{"$or": clauses}
}

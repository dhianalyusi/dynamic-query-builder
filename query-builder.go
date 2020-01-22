package dynamic_query_builder

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

type DQB string

type Expression struct {
	Key   string
	Exp   string
	Value interface{}
}

func (d DQB) Table(tableName string) DQB {
	if d == "" {
		return DQB("select * from " + tableName)
	}

	return DQB("select * from " + tableName + d.ToString())
}

func (d DQB) Select(fields string) DQB {
	return DQB(strings.ReplaceAll(d.ToString(), "select *", "select " + fields))
}

func (d DQB) Where(components DQB) DQB {
	if strings.Contains(d.ToString(), "where") {
		return DQB(strings.ReplaceAll(d.ToString(), "where", "where " + string(components) + " and"))
	}

	return DQB(d.ToString() + " where " + string(components))
}

func (d DQB) Limit(limit interface{}) DQB {
	return DQB(d.ToString() + " limit " + fmt.Sprintf("%v", limit))
}

func (d DQB) Offset(offset interface{}) DQB {
	return DQB(d.ToString() + " offset " + fmt.Sprintf("%v", offset))
}

func (d DQB) And(components ...interface{}) DQB {
	return d.Clause("and", components...)
}

func (d DQB) Or(components ...interface{}) DQB {
	return d.Clause("or", components...)
}

func (d DQB) Join(joinType string, table string, foreignKey string, referenceKey string) DQB {
	if strings.Contains(d.ToString(), "where") {
		return DQB(strings.ReplaceAll(d.ToString(), "where",  joinType + " " + table + " on " + foreignKey + " = " + referenceKey + " where"))
	} else if strings.Contains(d.ToString(), "limit") {
		return DQB(strings.ReplaceAll(d.ToString(), "limit",  joinType + " " + table + " on " + foreignKey + " = " + referenceKey + " limit"))
	}

	return DQB(d.ToString() + " " + joinType + " " + table + " on " + foreignKey + " = " + referenceKey)
}

func (d DQB) Order(sort string) DQB {
	return DQB(d.ToString() + " order by " + sort)
}

func (d DQB) Clause(operator string, components ...interface{}) DQB {
	clauses := make([]string, 0)
	for _, i := range components {
		value := componentToString(i)
		if value != "" {
			clauses = append(clauses, string(value))
		}
	}

	if len(clauses) > 0 {
		return DQB("(" + strings.Join(clauses, " " + operator + " ") + ")")
	}

	return ""
}

func (d DQB) NewExpression(key string, assignment string, value interface{}) Expression {
	return Expression{Key: key, Exp: assignment, Value: value}
}

func (d DQB) ToString() string {
	return string(d)
}

func (e Expression) ToString() string {
	switch  e.Value.(type) {
	case int, int16, int32, int64:
		val := strconv.Itoa(e.Value.(int))
		clause := e.Key + e.Exp + e.getReplaceExp()
		return fmt.Sprintf(clause, val)

	default:
		if strings.TrimSpace(e.Value.(string)) == "" {
			return ""
		} else {
			e.Value = template.HTMLEscapeString(e.Value.(string))
			clause := e.Key + " " + e.Exp + " " + e.getReplaceExp()
			val := fmt.Sprintf(clause, e.Value)
			return val
		}
	}
}

func (e Expression) getReplaceExp() string {
	switch e.Value.(type) {
	case int, int64, int32, int16:
		return "%s"
	default:
		return "'%s'"
	}
}

func componentToString(c interface{}) DQB {
	switch v := c.(type) {
	case Expression:
		return DQB(c.(Expression).ToString())
	case string, *string:
		return DQB(c.(string))
	case DQB:
		return v
	default:
		return ""
	}
}

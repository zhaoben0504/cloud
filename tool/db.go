package tool

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

const (
	generateSQLErr = "生成sql语句错误"
	columnErr      = "列数不一致"
)

type SQLBuilder struct {
	Table string
	Model interface{}
	Cond  interface{}
}

// NewSQLBuilder sql builder tag规则 `db:"column,add,set,sort"`
func NewSQLBuilder(table string, model interface{}) *SQLBuilder {
	return &SQLBuilder{Table: table, Model: model}
}

// Condition 生成条件部分
func (builder *SQLBuilder) Condition(c interface{}) *SQLBuilder {
	builder.Cond = c
	return builder
}

// BuildInsertRow 生成单条插入sql
func (builder *SQLBuilder) BuildInsertRow() (string, []interface{}, error) {
	sql, param := builder.generate("add-row")
	if sql == "" {
		fmt.Println(generateSQLErr)
		return "", nil, errors.New(generateSQLErr)
	}

	return "INSERT INTO `" + builder.Table + "`" + sql, param, nil
}

// BuildInsert 生成批量插入sql
func (builder *SQLBuilder) BuildInsert() (string, []interface{}, error) {
	sql, param := builder.generate("add-rows")
	if sql == "" {
		fmt.Println(generateSQLErr)
		return "", nil, errors.New(generateSQLErr)
	}

	return "INSERT INTO " + builder.Table + sql, param, nil
}

// generate 生成sql的列及条件部分，返回 列，条件部分
func (builder *SQLBuilder) generate(action string) (string, []interface{}) {
	switch action {
	case "add-row":
		column := ""
		values := ""
		params := make([]interface{}, 0)
		originType := reflect.TypeOf(builder.Model)
		if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
			fmt.Println("param error")
			return "", params
		}
		originValue := reflect.ValueOf(builder.Model)

		for i := 0; i < originType.Elem().NumField(); i++ {
			tag := originType.Elem().Field(i).Tag.Get("db")
			if tag == "" {
				continue
			}
			if originValue.Elem().Field(i).Kind() == reflect.Ptr {
				if !originValue.Elem().Field(i).IsNil() {
					if column != "" {
						column += ","
						values += ","
					}
					if strings.Index(tag, ",") > 0 {
						column += "`" + tag[:strings.Index(tag, ",")] + "`"
					} else {
						column += "`" + tag + "`"
					}
					values += "?"
					params = append(params, originValue.Elem().Field(i).Interface())
				}
			}
		}
		if len(column) == 0 {
			return "", params
		}
		return fmt.Sprintf("(%s) VALUES (%s)", column, values), params
	case "add-rows":
		column := ""
		columnCount := make(map[int]int, 10)
		values := ""
		params := make([]interface{}, 0)
		originType := reflect.TypeOf(builder.Model)
		if originType.Kind() != reflect.Slice && originType.Kind() != reflect.Array ||
			originType.Elem().Kind() != reflect.Struct {
			fmt.Println("param error")
			return "", params
		}
		originValue := reflect.ValueOf(builder.Model)
		for j := 0; j < originValue.Len(); j++ {
			item := originValue.Index(j)
			itemType := item.Type()
			row := ""
			for i := 0; i < itemType.NumField(); i++ {
				tag := itemType.Field(i).Tag.Get("db")
				if tag == "" {
					continue
				}
				if item.Field(i).Kind() == reflect.Ptr {
					if !item.Field(i).IsNil() {
						if j == 0 {
							if i > 0 {
								column += ","
							}
							if strings.Index(tag, ",") > 0 {
								column += "`" + tag[:strings.Index(tag, ",")] + "`"
							} else {
								column += "`" + tag + "`"
							}
						}

						if num, ok := columnCount[i]; ok {
							columnCount[i] = num + 1
						} else {
							columnCount[i] = 1
						}
						if row != "" {
							row += ","
						}

						row += "?"
						params = append(params, item.Field(i).Interface())
					}
				}
			}
			if len(row) > 0 {
				if j > 0 {
					values += ","
				}
				values += fmt.Sprintf("(%s)", row)
			}
		}

		if len(values) == 0 {
			fmt.Println("no data")
			return "", params
		}
		index := 0
		count := 0
		for _, v := range columnCount {
			if index == 0 {
				count = v
			}
			if count != v {
				log.Println(columnErr)
				return "", nil
			}
			index++
		}
		return fmt.Sprintf("(%s) VALUES %s", column, values), params
	case "set":
		set := ""
		params := make([]interface{}, 0)
		originType := reflect.TypeOf(builder.Model)
		if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
			fmt.Println("param error")
			return set, params
		}
		originValue := reflect.ValueOf(builder.Model)

		for i := 0; i < originType.Elem().NumField(); i++ {
			tag := originType.Elem().Field(i).Tag.Get("db")
			if tag == "" {
				continue
			}
			if strings.Contains(tag, "set") {
				if !originValue.Elem().Field(i).IsNil() {
					if set != "" {
						set += ","
					}
					if strings.Index(tag, ",") > 0 {
						set += "`" + tag[:strings.Index(tag, ",")] + "`=?"
					} else {
						set += "`" + tag + "`=?"
					}
					params = append(params, originValue.Elem().Field(i).Interface())
				}
			}
		}
		return set, params
	default:
		return "", nil
	}
}

// XormUpdateParam generate xorm param
func XormUpdateParam(model interface{}) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	originType := reflect.TypeOf(model)
	if originType.Kind() != reflect.Ptr || originType.Elem().Kind() != reflect.Struct {
		fmt.Println("param error")
		return nil, errors.New("param error")
	}
	originValue := reflect.ValueOf(model)

	for i := 0; i < originType.Elem().NumField(); i++ {
		if originType.Elem().Field(i).Type.Kind() != reflect.Ptr {
			continue
		}
		if !originValue.Elem().Field(i).IsNil() {
			tag := originType.Elem().Field(i).Tag.Get("db")
			if tag == "" {
				continue
			}
			if strings.Index(tag, ",") > 0 {
				params[tag[:strings.Index(tag, ",")]] = originValue.Elem().Field(i).Interface()
			} else {
				params[tag] = originValue.Elem().Field(i).Interface()
			}
		}
	}
	if len(params) == 0 {
		return nil, nil
	}
	return params, nil
}

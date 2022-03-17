package dialect

import (
	"fmt"
	"strings"

	"rxdrag.com/entity-engine/consts"
	"rxdrag.com/entity-engine/meta"
)

type MySQLBuilder struct {
}

func (*MySQLBuilder) BuildFieldExp(fieldName string, fieldArgs map[string]interface{}) (string, []interface{}) {
	var params []interface{}
	queryStr := "true "
	for key, value := range fieldArgs {
		switch key {
		case consts.ARG_EQ:
			queryStr = queryStr + " AND " + fieldName + "=?"
			params = append(params, value)
			break
		case consts.ARG_ISNULL:
			if value == true {
				queryStr = queryStr + " AND ISNULL(" + fieldName + ")"
			}
			break
		default:
			panic("Can not find token:" + key)
		}
	}
	return "(" + queryStr + ")", params
}

func (b *MySQLBuilder) BuildBoolExp(where map[string]interface{}) (string, []interface{}) {
	var params []interface{}
	queryStr := ""
	for key, value := range where {
		switch key {
		case consts.ARG_AND:
			break
		case consts.ARG_NOT:
			break
		case consts.ARG_OR:
			break
		default:
			fiedleStr, fieldParam := b.BuildFieldExp(key, value.(map[string]interface{}))
			queryStr = queryStr + " AND " + fiedleStr
			params = append(params, fieldParam...)
		}
	}
	return queryStr, params
}

func (b *MySQLBuilder) ColumnTypeSQL(column *meta.Column) string {
	typeStr := "text"
	switch column.Type {
	case meta.COLUMN_ID:
		typeStr = "bigint(64)"
		break
	case meta.COLUMN_INT:
		typeStr = "int"
		if column.Length == 1 {
			typeStr = "tinyint"
		} else if column.Length == 2 {
			typeStr = "smallint"
		} else if column.Length == 3 {
			typeStr = "mediumint"
		} else if column.Length == 4 {
			typeStr = "int"
		} else if column.Length > 4 {
			length := column.Length
			if length > 64 {
				length = 64
			}
			typeStr = fmt.Sprintf("bigint(%d)", length)
		}
		if column.Unsigned {
			typeStr = typeStr + " UNSIGNED"
		}
		break
	case meta.COLUMN_FLOAT:
		if column.Length > 4 {
			typeStr = "double"
		} else {
			typeStr = "float"
		}
		if column.FloatM > 0 && column.FloatD > 0 && column.FloatM >= column.FloatD {
			typeStr = fmt.Sprint(typeStr+"(%d,%d)", column.FloatM, column.FloatD)
		}
		if column.Unsigned {
			typeStr = typeStr + " UNSIGNED"
		}
		break
	case meta.COLUMN_BOOLEAN:
		typeStr = "tinyint(1)"
		break
	case meta.COLUMN_STRING:
		typeStr = "text"
		if column.Length > 0 {
			if column.Length <= 255 {
				typeStr = fmt.Sprintf("varchar(%d)", column.Length)
			} else if column.Length <= 65535 {
				typeStr = "text"
			} else if column.Length <= 16777215 {
				typeStr = "mediumtext"
			} else {
				typeStr = "longtext"
			}
		}
		break
	case meta.COLUMN_DATE:
		typeStr = "datetime"
		break
	case meta.COLUMN_SIMPLE_JSON:
		typeStr = "json"
		break
	case meta.COLUMN_SIMPLE_ARRAY:
		typeStr = "json"
		break
	case meta.COLUMN_JSON_ARRAY:
		typeStr = "json"
		break
	case meta.COLUMN_ENUM:
		typeStr = "tinytext"
		break
	}
	return typeStr
}

func (b *MySQLBuilder) BuildColumnSQL(column *meta.Column) string {
	sql := "`" + column.Name + "` " + b.ColumnTypeSQL(column)
	if column.Generated {
		sql = sql + " AUTO_INCREMENT"
	}
	return sql
}

func (b *MySQLBuilder) BuildCreateTableSQL(table *meta.Table) string {
	sql := "CREATE TABLE `%s` (%s)"
	fieldSqls := make([]string, len(table.Columns))
	for i := range table.Columns {
		columnSql := b.BuildColumnSQL(&table.Columns[i])
		if table.Columns[i].Nullable {
			columnSql = columnSql + " NULL"
		} else {
			columnSql = columnSql + " NOT NULL"
		}
		fieldSqls[i] = columnSql
	}
	for _, column := range table.Columns {
		if column.Primary {
			fieldSqls = append(fieldSqls, fmt.Sprintf("PRIMARY KEY (`%s`)", column.Name))
		}
	}

	//建索引
	for _, column := range table.Columns {
		if column.Index {
			indexSql := "INDEX %s ( `%s`)"
			fieldSqls = append(fieldSqls, fmt.Sprintf(indexSql, column.Name+consts.INDEX_SUFFIX, column.Name))
		}
	}

	sql = fmt.Sprintf(sql, table.Name, strings.Join(fieldSqls, ","))

	return sql
}

func (b *MySQLBuilder) BuildDeleteTableSQL(table *meta.Table) string {
	return "DROP TABLE " + table.Name
}

func (b *MySQLBuilder) BuildModifyTableAtoms(diff *meta.TableDiff) []meta.ModifyAtom {
	var atoms []meta.ModifyAtom
	if diff.OldTable.Name != diff.NewTable.Name {
		//修改表名
		atoms = append(atoms, meta.ModifyAtom{
			ExcuteSQL: fmt.Sprintf("ALTER TABLE %s RENAME TO %s ", diff.OldTable.Name, diff.NewTable.Name),
			UndoSQL:   fmt.Sprintf("ALTER TABLE %s RENAME TO %s ", diff.NewTable.Name, diff.OldTable.Name),
		})
	}
	b.appendDeleteColumnAtoms(diff, &atoms)
	b.appendAddColumnAtoms(diff, &atoms)
	b.appendModifyColumnAtoms(diff, &atoms)
	return atoms
}

func (b *MySQLBuilder) appendDeleteColumnAtoms(diff *meta.TableDiff, atoms *[]meta.ModifyAtom) {
	for _, column := range diff.DeleteColumns {
		//删除索引
		if column.Index {
			indexName := column.Name + consts.INDEX_SUFFIX
			*atoms = append(*atoms, meta.ModifyAtom{
				ExcuteSQL: fmt.Sprintf("DROP INDEX %s ON %s ", indexName, diff.NewTable.Name),
				UndoSQL:   fmt.Sprintf("CREATE INDEX %s ON %s (%s)", indexName, diff.NewTable.Name, column.Name),
			})
		}
		//删除列
		*atoms = append(*atoms, meta.ModifyAtom{
			ExcuteSQL: fmt.Sprintf("ALTER TABLE %s DROP  %s ", diff.NewTable.Name, column.Name),
			UndoSQL:   fmt.Sprintf("ALTER TABLE %s ADD COLUMN  %s %s", diff.NewTable.Name, column.Name, b.ColumnTypeSQL(&column)),
		})
	}
}

func (b *MySQLBuilder) appendAddColumnAtoms(diff *meta.TableDiff, atoms *[]meta.ModifyAtom) {
	for _, column := range diff.AddColumns {
		//添加列
		*atoms = append(*atoms, meta.ModifyAtom{
			ExcuteSQL: fmt.Sprintf("ALTER TABLE %s ADD COLUMN  %s %s", diff.NewTable.Name, column.Name, b.ColumnTypeSQL(&column)),
			UndoSQL:   fmt.Sprintf("ALTER TABLE %s DROP  %s ", diff.NewTable.Name, column.Name),
		})
		//添加索引
		if column.Index {
			indexName := column.Name + consts.INDEX_SUFFIX
			*atoms = append(*atoms, meta.ModifyAtom{
				ExcuteSQL: fmt.Sprintf("CREATE INDEX %s ON %s (%s)", indexName, diff.NewTable.Name, column.Name),
				UndoSQL:   fmt.Sprintf("DROP INDEX %s ON %s ", indexName, diff.NewTable.Name),
			})
		}
	}
}

func (b *MySQLBuilder) appendModifyColumnAtoms(diff *meta.TableDiff, atoms *[]meta.ModifyAtom) {
	for _, columnDiff := range diff.ModifyColumns {
		//删除索引
		if columnDiff.OldColumn.Index {

		}
		//更改列
		if columnDiff.OldColumn.Name != columnDiff.OldColumn.Name ||
			columnDiff.OldColumn.Type != columnDiff.OldColumn.Type ||
			columnDiff.OldColumn.Length != columnDiff.OldColumn.Length ||
			columnDiff.OldColumn.FloatD != columnDiff.OldColumn.FloatD ||
			columnDiff.OldColumn.FloatM != columnDiff.OldColumn.FloatM ||
			columnDiff.OldColumn.Unsigned != columnDiff.OldColumn.Unsigned {

		}
		//添加索引
		if columnDiff.NewColumn.Index {

		}
	}
}

package model

type Table struct {
	MetaUuid string
	Name     string
	Columns  []*Column
}

func FindTable(metaUuid string, tables []*Table) *Table {
	for i := range tables {
		if tables[i].MetaUuid == metaUuid {
			return tables[i]
		}
	}
	return nil
}
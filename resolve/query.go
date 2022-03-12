package resolve

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"rxdrag.com/entity-engine/meta"
	"rxdrag.com/entity-engine/repository"
)

func QueryOneResolveFn(entity *meta.Entity) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		return repository.QueryOne(entity, p.Args)
	}
}

func QueryResolveFn(entity *meta.Entity) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		names := entity.ColumnNames()
		queryStr := "select %s from %s "
		for _, iSelection := range p.Info.Operation.GetSelectionSet().Selections {
			switch selection := iSelection.(type) {
			case *ast.Field:
				fmt.Println(selection.Directives[len(selection.Directives)-1].Name.Value)
			case *ast.InlineFragment:
			case *ast.FragmentSpread:
			}
		}

		queryStr = fmt.Sprintf(queryStr, strings.Join(names, ","), entity.GetTableName())
		//err = db.Select(&instances, queryStr)
		return repository.Query(entity, queryStr)
	}
}
package resolve

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/authorization"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/utils"
)

func QueryOneInterfaceResolveFn(intf *graph.Interface) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		instance := repository.QueryOneInterface(intf, p.Args)
		return instance, nil
	}
}

func QueryInterfaceResolveFn(intf *graph.Interface) graphql.FieldResolveFn {
	var uuids []string
	for i := range intf.Children {
		uuids = append(uuids, intf.Children[i].Uuid())
	}
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		makeQueryVerifier(p, uuids)
		return repository.QueryInterface(intf, p.Args), nil
	}
}

func QueryOneEntityResolveFn(entity *graph.Entity) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		instance := repository.QueryOneEntity(entity, p.Args)
		return instance, nil
	}
}

func QueryEntityResolveFn(entity *graph.Entity) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		// for _, iSelection := range p.Info.Operation.GetSelectionSet().Selections {
		// 	switch selection := iSelection.(type) {
		// 	case *ast.Field:
		// 		//fmt.Println(selection.Directives[len(selection.Directives)-1].Name.Value)
		// 	case *ast.InlineFragment:
		// 	case *ast.FragmentSpread:
		// 	}
		// }
		makeQueryVerifier(p, []string{entity.Uuid()})
		return repository.QueryEntity(entity, p.Args), nil
	}
}

func QueryAssociationFn(asso *graph.Association) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		var (
			source      = p.Source.(map[string]interface{})
			v           = p.Context.Value
			loaders     = v(consts.LOADERS).(*Loaders)
			handleError = func(err error) error {
				return fmt.Errorf(err.Error())
			}
		)
		defer utils.PrintErrorStack()

		if loaders == nil {
			panic("Data loaders is nil")
		}
		loader := loaders.GetLoader(asso)
		thunk := loader.Load(p.Context, NewKey(source[consts.ID].(uint64)))
		return func() (interface{}, error) {
			data, err := thunk()
			if err != nil {
				return nil, handleError(err)
			}

			var retValue interface{}
			if data == nil {
				if asso.IsArray() {
					retValue = []map[string]interface{}{}
				} else {
					retValue = nil
				}
			} else {
				retValue = data
			}
			return retValue, nil
		}, nil
	}
}

func makeQueryVerifier(p graphql.ResolveParams, entityUuids []string) *authorization.AbilityVerifier {
	verifier := authorization.NewVerifier()

	if verifier == nil {
		panic("Can not find Ability Verifier")
	}

	verifier.Init(p, entityUuids, meta.META_ABILITY_TYPE_READ)

	// if !verifier.CanReadEntity() && !node.IsInterface() {
	// 	panic("No permission to read: " + node.Name())
	// }

	// args := verifier.WeaveAuthInArgs(inputArgs)
	return verifier
}

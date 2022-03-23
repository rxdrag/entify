package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entity-engine/meta"
	"rxdrag.com/entity-engine/model"
)

var Cache TypeCache

type TypeCache struct {
	ObjectTypeMap        map[string]*graphql.Object
	EnumTypeMap          map[string]*graphql.Enum
	InterfaceTypeMap     map[string]*graphql.Interface
	UpdateInputMap       map[string]*graphql.InputObject
	SaveInputMap         map[string]*graphql.InputObject
	WhereExpMap          map[string]*graphql.InputObject
	DistinctOnEnumMap    map[string]*graphql.Enum
	OrderByMap           map[string]*graphql.InputObject
	EnumComparisonExpMap map[string]*graphql.InputObjectFieldConfig
	MutationResponseMap  map[string]*graphql.Output
	AggregateMap         map[string]*graphql.Output
}

//where表达式缓存，query跟mutation都用

func (c *TypeCache) MakeCache() {
	c.clearCache()
	enums, interfaces, normals := meta.Metas.SplitEntities()
	c.makeEnums(enums)
	c.makeOutputInterfaces(interfaces)
	c.makeOutputObjects(normals)
	c.makeRelations()
	c.makeArgs()
	c.makeInputs()
}

func (c *TypeCache) OutputType(entity *model.Entity) graphql.Type {
	if entity.EntityType == meta.ENTITY_ENUM {
		return c.EnumTypeMap[entity.Name]
	} else if entity.EntityType == meta.ENTITY_INTERFACE {
		return c.InterfaceTypeMap[entity.Name]
	} else {
		return c.ObjectTypeMap[entity.Name]
	}
}

func (c *TypeCache) WhereExp(entity *model.Entity) *graphql.InputObject {
	return c.WhereExpMap[entity.Name]
}

func (c *TypeCache) OrderByExp(entity *model.Entity) *graphql.InputObject {
	return c.OrderByMap[entity.Name]
}

func (c *TypeCache) DistinctOnEnum(entity *model.Entity) *graphql.Enum {
	return c.DistinctOnEnumMap[entity.Name]
}

func (c *TypeCache) SaveInput(entity *model.Entity) *graphql.InputObject {
	return c.SaveInputMap[entity.Name]
}

func (c *TypeCache) UpdateInput(entity *model.Entity) *graphql.InputObject {
	return c.UpdateInputMap[entity.Name]
}

func (c *TypeCache) MutationResponse(entity *model.Entity) *graphql.Output {
	return c.MutationResponseMap[entity.Name]
}

func (c *TypeCache) mapInterfaces(entities []*model.Entity) []*graphql.Interface {
	interfaces := []*graphql.Interface{}
	for i := range entities {
		interfaces = append(interfaces, c.InterfaceTypeMap[entities[i].Name])
	}

	return interfaces
}

func (c *TypeCache) clearCache() {
	c.ObjectTypeMap = make(map[string]*graphql.Object)
	c.EnumTypeMap = make(map[string]*graphql.Enum)
	c.InterfaceTypeMap = make(map[string]*graphql.Interface)
	c.UpdateInputMap = make(map[string]*graphql.InputObject)
	c.SaveInputMap = make(map[string]*graphql.InputObject)
	c.WhereExpMap = make(map[string]*graphql.InputObject)
	c.DistinctOnEnumMap = make(map[string]*graphql.Enum)
	c.OrderByMap = make(map[string]*graphql.InputObject)
	c.EnumComparisonExpMap = make(map[string]*graphql.InputObjectFieldConfig)
	c.MutationResponseMap = make(map[string]*graphql.Output)
	c.AggregateMap = make(map[string]*graphql.Output)
}

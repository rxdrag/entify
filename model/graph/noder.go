package graph

type Noder interface {
	Uuid() string
	InnerId() uint64
	Name() string
	TableName() string
	Description() string
	IsInterface() bool
	Interface() *Interface
	Entity() *Entity
	AddAssociation(a *Association)
	GetAssociationByName(name string) *Association
	//Attributes() []*Attribute
	AllAssociations() []*Association
	AllAttributes() []*Attribute
	AllMethods() []*Method
	AllAttributeNames() []string
	GetAttributeByName(name string) *Attribute
	IsEmperty() bool
}

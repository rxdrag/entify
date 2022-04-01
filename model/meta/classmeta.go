package meta

const (
	CLASSS_ENTITY      string = "Normal"
	CLASSS_ENUM        string = "Enum"
	CLASSS_ABSTRACT    string = "Abstract"
	CLASS_VALUE_OBJECT string = "ValueObject"
)

type ClassMeta struct {
	Uuid        string          `json:"uuid"`
	InnerId     uint64          `json:"innerId"`
	Name        string          `json:"name"`
	StereoType  string          `json:"stereoType"`
	Attributes  []AttributeMeta `json:"attributes"`
	Methods     []MethodMeta    `json:"methods"`
	Root        bool            `json:"root"`
	Description string          `json:"description"`
}
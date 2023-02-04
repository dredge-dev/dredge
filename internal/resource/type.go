package resource

import (
	"strings"

	"github.com/dredge-dev/dredge/internal/exec"
)

type Type struct {
	Name    string
	IsArray bool
	Fields  []Field
}

func GetType(e *exec.DredgeExec, typeName string) (*Type, error) {
	isArray := false
	if strings.HasPrefix(typeName, "[]") {
		isArray = true
		typeName = strings.TrimPrefix(typeName, "[]")
	}

	if typeName == "string" || typeName == "date" || typeName == "object" {
		return &Type{
			typeName,
			isArray,
			nil,
		}, nil
	}

	resourceType, err := GetResourceDefinition(e, typeName)
	if err != nil {
		return nil, err
	}

	return &Type{
		typeName,
		isArray,
		resourceType.Fields,
	}, nil
}

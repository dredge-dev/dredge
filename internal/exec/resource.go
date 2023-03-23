package exec

import (
	"fmt"
	"strings"

	"github.com/dredge-dev/dredge/internal/api"
	"github.com/dredge-dev/dredge/internal/resource"
)

func (e *DredgeExec) GetResources() ([]string, error) {
	var resources []string
	for name, _ := range e.DredgeFile.Resources {
		_, err := e.GetResourceDefinition(name)
		if err != nil {
			return nil, err
		}
		resources = append(resources, name)
	}
	return resources, nil
}

func (e *DredgeExec) GetResourceDefinition(resourceName string) (*api.ResourceDefinition, error) {
	for _, rd := range e.ResourceDefinitions {
		if rd.Name == resourceName {
			return &rd, nil
		}
	}
	return nil, fmt.Errorf("could not find resource definition for %s", resourceName)
}

func (e *DredgeExec) GetType(typeName string) (*api.Type, error) {
	isArray := false
	if strings.HasPrefix(typeName, "[]") {
		isArray = true
		typeName = strings.TrimPrefix(typeName, "[]")
	}

	if typeName == "string" || typeName == "date" || typeName == "object" {
		return &api.Type{
			Name:    typeName,
			IsArray: isArray,
			Fields:  nil,
		}, nil
	}

	resourceType, err := e.GetResourceDefinition(typeName)
	if err != nil {
		return nil, err
	}

	return &api.Type{
		Name:    typeName,
		IsArray: isArray,
		Fields:  resourceType.Fields,
	}, nil
}

func (e *DredgeExec) GetResource(resourceName string) (*resource.Resource, error) {
	rd, err := e.GetResourceDefinition(resourceName)
	if err != nil {
		return nil, err
	}

	r, ok := e.DredgeFile.Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("could not find resource %s", resourceName)
	}
	if len(r) == 0 {
		return nil, fmt.Errorf("no provider specified for this resource")
	}

	return resource.NewResource(rd, r)
}

func (e *DredgeExec) GetProviders() ([]resource.Provider, error) {
	return resource.GetProviders()
}

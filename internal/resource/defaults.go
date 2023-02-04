package resource

func GetDefaultResourceDefinitions() []ResourceDefinition {
	return []ResourceDefinition{
		{
			Name: "release",
			Fields: []Field{
				{
					Name:        "name",
					Description: "Release name",
					Type:        "string",
				},
				{
					Name:        "date",
					Description: "Release date",
					Type:        "date",
				},
				{
					Name:        "title",
					Description: "Release title",
					Type:        "string",
				},
			},
			Commands: []Command{
				{
					Name:       "get",
					Inputs:     []string{},
					OutputType: "[]release",
				},
				{
					Name:       "search",
					Inputs:     []string{"text"},
					OutputType: "[]release",
				},
				{
					Name:       "describe",
					Inputs:     []string{"name"},
					OutputType: "object",
				},
			},
		},
		{
			Name: "issue",
			Fields: []Field{
				{
					Name:        "name",
					Description: "Issue name",
					Type:        "string",
				},
				{
					Name:        "title",
					Description: "Issue title",
					Type:        "string",
				},
				{
					Name:        "type",
					Description: "Issue type",
					Type:        "string",
				},
				{
					Name:        "state",
					Description: "Issue state",
					Type:        "string",
				},
				{
					Name:        "date",
					Description: "Issue creation date",
					Type:        "date",
				},
			},
			Commands: []Command{
				{
					Name:       "get",
					Inputs:     []string{},
					OutputType: "[]issue",
				},
				{
					Name:       "create",
					Inputs:     []string{},
					OutputType: "issue",
				},
			},
		},
		{
			Name: "doc",
			Fields: []Field{
				{
					Name:        "name",
					Description: "Name",
					Type:        "string",
				},
				{
					Name:        "author",
					Description: "Author",
					Type:        "string",
				},
				{
					Name:        "location",
					Description: "Location",
					Type:        "string",
				},
				{
					Name:        "date",
					Description: "Last updated date",
					Type:        "date",
				},
			},
			Commands: []Command{
				{
					Name:       "get",
					Inputs:     []string{},
					OutputType: "[]doc",
				},
				{
					Name:       "search",
					Inputs:     []string{},
					OutputType: "[]doc",
				},
			},
		},
		{
			Name: "deploy",
			Fields: []Field{
				{
					Name:        "name",
					Description: "Name",
					Type:        "string",
				},
				{
					Name:        "version",
					Description: "Version",
					Type:        "string",
				},
				{
					Name:        "instances",
					Description: "Number of instances",
					Type:        "string",
				},
				{
					Name:        "type",
					Description: "Instance type",
					Type:        "string",
				},
			},
			Commands: []Command{
				{
					Name:       "get",
					Inputs:     []string{},
					OutputType: "[]deploy",
				},
				{
					Name:       "describe",
					Inputs:     []string{},
					OutputType: "object",
				},
				{
					Name:       "update",
					Inputs:     []string{},
					OutputType: "deploy",
				},
			},
		},
	}
}

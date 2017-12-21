package template

import (
	"html/template"
)

func NewFuncMap() []template.FuncMap {
	return []template.FuncMap{
		map[string]interface{}{
			"AppSubUrl": func() string {
				return ""
			},
		},
	}
}

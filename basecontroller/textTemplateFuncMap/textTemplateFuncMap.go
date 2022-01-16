package textTemplateFuncMap

import (
	"text/template"
)

func New() textTemplateFuncMap {
	return textTemplateFuncMap{}
}

type textTemplateFuncMap = template.FuncMap

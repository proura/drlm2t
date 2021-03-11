package model

type TemplateTest struct {
	Name    string `mpastructure:"name"`
	Content string `mpastructure:"content"`
}

type Template struct {
	Name          string         `mpastructure:"name"`
	TemplateTests []TemplateTest `mpastructure:"templatetests"`
}

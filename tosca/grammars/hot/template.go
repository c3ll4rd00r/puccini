package hot

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Template
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#template-structure]
//

type Template struct {
	*Entity `name:"template"`

	HeatTemplateVersion  *string                `read:"heat_template_version" require:"heat_template_version"`
	Description          *string                `read:"description"`
	ParameterGroups      []*ParameterGroup      `read:"parameter_groups,[]ParameterGroup"`
	Parameters           Parameters             `read:"parameters,Parameter"`
	Resources            []*Resource            `read:"resources,Resource"`
	Outputs              Outputs                `read:"outputs,Output"`
	ConditionDefinitions []*ConditionDefinition `read:"conditions,ConditionDefinition"`
}

func NewTemplate(context *tosca.Context) *Template {
	return &Template{
		Entity:     NewEntity(context),
		Parameters: make(Parameters),
		Outputs:    make(Outputs),
	}
}

// tosca.Reader signature
func ReadTemplate(context *tosca.Context) interface{} {
	self := NewTemplate(context)

	context.ImportScript("tosca.resolve", "internal:/tosca/simple/1.1/js/resolve.js")
	context.ImportScript("tosca.coerce", "internal:/tosca/simple/1.1/js/coerce.js")
	context.ImportScript("tosca.visualize", "internal:/tosca/simple/1.1/js/visualize.js")
	context.ImportScript("tosca.utils", "internal:/tosca/simple/1.1/js/utils.js")
	context.ImportScript("tosca.helpers", "internal:/tosca/simple/1.1/js/helpers.js")

	context.ScriptNamespace.Merge(DefaultScriptNamespace)

	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))

	return self
}

// parser.Importer interface
func (self *Template) GetImportSpecs() []*tosca.ImportSpec {
	var importSpecs []*tosca.ImportSpec
	return importSpecs
}

// parser.HasInputs interface
func (self *Template) SetInputs(inputs map[string]interface{}) {
	context := self.Context.FieldChild("parameters", nil)
	for name, data := range inputs {
		childContext := context.MapChild(name, data)
		parameter, ok := self.Parameters[name]
		if !ok {
			childContext.ReportUndefined("parameter")
			continue
		}

		parameter.Value = ReadValue(childContext).(*Value)
		if parameter.Type != nil {
			parameter.Value.Fix(*parameter.Type)
		}
	}
}

// tosca.Normalizable interface
func (self *Template) Normalize() *normal.ServiceTemplate {
	log.Info("{normalize} template")

	s := normal.NewServiceTemplate()

	if self.Description != nil {
		s.Description = *self.Description
	}

	s.ScriptNamespace = self.Context.ScriptNamespace

	self.Parameters.Normalize(s.Inputs, self.Context.FieldChild("parameters", nil))
	self.Outputs.Normalize(s.Outputs, self.Context.FieldChild("outputs", nil))

	for _, resource := range self.Resources {
		s.NodeTemplates[resource.Name] = resource.Normalize(s)
	}

	for _, resource := range self.Resources {
		resource.NormalizeDependencies(s)
	}

	return s
}
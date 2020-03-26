package cloudify_v1_3

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/yamlkeys"
)

//
// Value
//

type Value struct {
	*Entity `name:"value"`
	Name    string

	Description *string
	TypeName    string

	rendered bool
}

func NewValue(context *tosca.Context) *Value {
	return &Value{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadValue(context *tosca.Context) interface{} {
	ToFunctionCalls(context)
	return NewValue(context)
}

// tosca.Mappable interface
func (self *Value) GetKey() string {
	return self.Name
}

func (self *Value) RenderProperty(dataType *DataType, definition *PropertyDefinition) {
	if definition == nil {
		self.RenderParameter(dataType, nil, true, false)
	} else {
		self.RenderParameter(dataType, definition.ParameterDefinition, true, false)
	}
}

func (self *Value) RenderParameter(dataType *DataType, definition *ParameterDefinition, validateRequire bool, allowNil bool) {
	if self.rendered {
		// Avoid rendering more than once (can happen if we use the "default" value)
		return
	}
	self.rendered = true

	if (definition != nil) && (definition.Description != nil) {
		self.Description = definition.Description
	} else {
		self.Description = dataType.Description
	}

	self.TypeName = tosca.GetCanonicalName(dataType)

	if _, ok := self.Context.Data.(*tosca.FunctionCall); ok {
		return
	}

	if allowNil && (self.Context.Data == nil) {
		return
	}

	// TODO: dataType.Complete(self.Context.Data, self.Context)

	// Internal types
	if internalTypeName, typeValidator, reader, ok := dataType.GetInternal(); ok {
		if typeValidator != nil {
			if self.Context.Data == nil {
				// Nil data only happens when a parameter is added despite not having a
				// "default" value; we will give it a valid zero value instead
				self.Context.Data = ard.TypeZeroes[internalTypeName]
			}

			// Primitive types
			if !typeValidator(self.Context.Data) {
				self.Context.ReportValueWrongType(internalTypeName)
			}
		} else {
			// Special types
			self.Context.Data = reader(self.Context)
		}
	} else if self.Context.ValidateType("!!map") {
		// Complex data types

		map_ := self.Context.Data.(ard.Map)

		// All properties must be defined in type
		for key := range map_ {
			name := yamlkeys.KeyString(key)
			if _, ok := dataType.PropertyDefinitions[name]; !ok {
				self.Context.MapChild(name, nil).ReportUndeclared("property")
				delete(map_, key)
			}
		}

		// Render properties
		for key, definition := range dataType.PropertyDefinitions {
			if data, ok := map_[key]; ok {
				var value *Value
				if value, ok = data.(*Value); !ok {
					// Convert to value
					value = NewValue(self.Context.MapChild(key, data))
					map_[key] = value
				}
				if definition.DataType != nil {
					value.RenderProperty(definition.DataType, definition)
				}
			} else if validateRequire {
				if definition.IsRequired() {
					self.Context.MapChild(key, data).ReportPropertyRequired("property")
				}
			}
		}
	}
}

func (self *Value) Normalize() normal.Constrainable {
	var normalConstrainable normal.Constrainable

	switch data := self.Context.Data.(type) {
	case ard.List:
		normalList := normal.NewList(len(data))
		for index, value := range data {
			normalList.Set(index, NewValue(self.Context.ListChild(index, value)).Normalize())
		}
		normalConstrainable = normalList

	case ard.Map:
		normalMap := normal.NewMap()
		normalMap.Type = self.TypeName
		for key, value := range data {
			if _, ok := key.(string); !ok {
				// Cloudify DSL does not support complex keys
				self.Context.MapChild(key, yamlkeys.KeyData(key)).ReportValueWrongType("!!str")
			}
			name := yamlkeys.KeyString(key)
			normalMap.Put(name, NewValue(self.Context.MapChild(name, value)).Normalize())
		}
		normalConstrainable = normalMap

	case *tosca.FunctionCall:
		NormalizeFunctionCallArguments(data, self.Context)
		normalConstrainable = normal.NewFunctionCall(data)

	default:
		value := normal.NewValue(data)
		value.Type = self.TypeName
		normalConstrainable = value
	}

	if self.Description != nil {
		normalConstrainable.SetDescription(*self.Description)
	}

	return normalConstrainable
}

//
// Values
//

type Values map[string]*Value

func (self Values) SetIfNil(context *tosca.Context, key string, data interface{}) {
	if _, ok := self[key]; !ok {
		self[key] = NewValue(context.MapChild(key, data))
	}
}

func (self Values) RenderMissingValue(definition *ParameterDefinition, kind string, required bool, context *tosca.Context) {
	if definition.Default != nil {
		self[definition.Name] = definition.Default
	} else if required {
		context.MapChild(definition.Name, nil).ReportPropertyRequired(kind)
	} else if kind == "attribute" {
		// Parameters should always appear, even if they have no default value
		self[definition.Name] = NewValue(context.MapChild(definition.Name, nil))
	}
}

func (self Values) RenderProperties(definitions PropertyDefinitions, kind string, context *tosca.Context) {
	for key, definition := range definitions {
		if value, ok := self[key]; !ok {
			self.RenderMissingValue(definition.ParameterDefinition, kind, definition.IsRequired(), context)
			// (If the above assigns the "default" value -- it has already been rendered elsewhere)
		} else if definition.DataType != nil {
			value.RenderProperty(definition.DataType, definition)
		}
	}
}

func (self Values) RenderParameters(definitions ParameterDefinitions, kind string, context *tosca.Context) {
	for key, definition := range definitions {
		if _, ok := self[key]; !ok {
			self.RenderMissingValue(definition, kind, false, context)
		}
	}

	for key, value := range self {
		if definition, ok := definitions[key]; !ok {
			value.Context.ReportUndeclared(kind)
			delete(self, key)
		} else if definition.DataType != nil {
			value.RenderParameter(definition.DataType, definition, true, true)
		}
	}
}

func (self Values) Normalize(normalConstrainables normal.Constrainables, prefix string) {
	for key, value := range self {
		normalConstrainables[prefix+key] = value.Normalize()
	}
}

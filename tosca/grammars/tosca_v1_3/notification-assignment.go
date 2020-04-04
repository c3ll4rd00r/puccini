package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// NotificationAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.19
//

type NotificationAssignment struct {
	*Entity `name:"notification"`
	Name    string

	Description    *string                  `read:"description"`
	Implementation *InterfaceImplementation `read:"implementation,InterfaceImplementation"`
	Outputs        NotificationOutputs      `read:"outputs,NotificationOutput"`
}

func NewNotificationAssignment(context *tosca.Context) *NotificationAssignment {
	return &NotificationAssignment{
		Entity:  NewEntity(context),
		Name:    context.Name,
		Outputs: make(NotificationOutputs),
	}
}

// tosca.Reader signature
func ReadNotificationAssignment(context *tosca.Context) tosca.EntityPtr {
	self := NewNotificationAssignment(context)

	if context.Is("!!map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("!!map", "!!str") {
		// Short notation
		self.Implementation = ReadInterfaceImplementation(context.FieldChild("implementation", context.Data)).(*InterfaceImplementation)
	}

	return self
}

// tosca.Mappable interface
func (self *NotificationAssignment) GetKey() string {
	return self.Name
}

func (self *NotificationAssignment) Normalize(normalInterface *normal.Interface) *normal.Notification {
	log.Debugf("{normalize} notification: %s", self.Name)

	normalNotification := normalInterface.NewNotification(self.Name)

	if self.Description != nil {
		normalNotification.Description = *self.Description
	}

	if self.Implementation != nil {
		self.Implementation.NormalizeNotification(normalNotification)
	}

	self.Outputs.Normalize(normalInterface.NodeTemplate, normalNotification.Outputs)

	return normalNotification
}

//
// NotificationAssignments
//

type NotificationAssignments map[string]*NotificationAssignment

func (self NotificationAssignments) CopyUnassigned(assignments NotificationAssignments) {
	for key, assignment := range assignments {
		if selfAssignment, ok := self[key]; ok {
			selfAssignment.Outputs.CopyUnassigned(assignment.Outputs)
		} else {
			self[key] = assignment
		}
	}
}

func (self NotificationAssignments) Render(definitions NotificationDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		assignment, ok := self[key]

		if !ok {
			assignment = NewNotificationAssignment(context.FieldChild(key, nil))
			self[key] = assignment
		}

		if assignment.Description == nil {
			assignment.Description = definition.Description
		}

		if (assignment.Implementation == nil) && (definition.Implementation != nil) {
			// If the definition has an implementation then we must have one, too
			assignment.Implementation = NewInterfaceImplementation(assignment.Context.FieldChild("implementation", nil))
		}

		if assignment.Implementation != nil {
			assignment.Implementation.Render(definition.Implementation)
		}

		assignment.Outputs.Inherit(definition.Outputs)
	}

	for key, assignment := range self {
		if _, ok := definitions[key]; !ok {
			assignment.Context.ReportUndeclared("notification")
			delete(self, key)
		}
	}
}

func (self NotificationAssignments) Normalize(normalInterface *normal.Interface) {
	for key, notification := range self {
		normalInterface.Notifications[key] = notification.Normalize(normalInterface)
	}
}

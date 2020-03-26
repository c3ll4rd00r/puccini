// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/simple/1.0/data.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_0

metadata:

  puccini.scriptlet.import|tosca.comparer.version: internal:/tosca/simple/1.3/js/comparers/version.js
  puccini.scriptlet.import|tosca.constraint._primitive: internal:/tosca/simple/1.3/js/constraints/_primitive.js

data_types:

  #
  # Primitive
  #

  boolean:
    metadata:
      puccini.type: '!!bool'
    constraints:
    - _primitive: '!!bool'

  integer:
    metadata:
      puccini.type: '!!int'
    constraints:
    - _primitive: '!!int'

  float:
    metadata:
      puccini.type: '!!float'
    constraints:
    - _primitive: '!!float'

  string:
    metadata:
      puccini.type: '!!str'
    constraints:
    - _primitive: '!!str'

  timestamp:
    metadata:
      puccini.type: timestamp

  #
  # Special
  #

  version:
    metadata:
      puccini.type: version
      puccini.comparer: tosca.comparer.version
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 3.2.2

  range:
    metadata:
      puccini.type: range
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 3.2.3

  #
  # With entry schema
  #

  list:
    metadata:
      puccini.type: '!!seq'
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 3.2.4

  map:
    metadata:
      puccini.type: '!!map'
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 3.2.5

  #
  # Scalar
  #

  scalar-unit.size:
    metadata:
      puccini.type: scalar-unit.size
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 3.2.6.4

  scalar-unit.time:
    metadata:
      puccini.type: scalar-unit.time
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 3.2.6.5

  scalar-unit.frequency:
    metadata:
      puccini.type: scalar-unit.frequency
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 3.2.6.6

  #
  # Complex
  #

  tosca.datatypes.Root:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 5.2.1
    description: >-
      This is the default (root) TOSCA Root Type definition that all complex TOSCA Data Types derive
      from.

  tosca.datatypes.Credential:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 5.2.2
    description: >-
      The Credential type is a complex TOSCA data Type used when describing authorization
      credentials used to access network accessible resources.
    derived_from: tosca.datatypes.Root
    properties:
      protocol:
        description: >-
          The optional protocol name.
        type: string
        required: false
      token_type:
        description: >-
          The required token type.
        type: string
        default: password
      token:
        description: >-
          The required token used as a credential for authorization or access to a networked
          resource.
        type: string
      keys:
        description: >-
          The optional list of protocol-specific keys or assertions.
        type: map
        entry_schema:
          type: string
        required: false
      user:
        description: >-
          The optional user (name or ID) used for non-token based credentials.
        type: string
        required: false

  tosca.datatypes.network.NetworkInfo:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 5.2.3
    description: >-
      The Network type is a complex TOSCA data type used to describe logical network information.
    derived_from: tosca.datatypes.Root
    properties:
      network_name:
        description: >-
          The name of the logical network. e.g., "public", "private", "admin". etc.
        type: string
        required: false
      network_id:
        description: >-
          The unique ID of for the network generated by the network provider.
        type: string
        required: false
      addresses:
        description: >-
          The list of IP addresses assigned from the underlying network.
        type: list
        entry_schema:
          type: string
        required: false

  tosca.datatypes.network.PortInfo:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 5.2.4
    description: >-
      The PortInfo type is a complex TOSCA data type used to describe network port information.
    derived_from: tosca.datatypes.Root
    properties:
      port_name:
        description: >-
          The logical network port name.
        type: string
        required: false
      port_id:
        description: >-
          The unique ID for the network port generated by the network provider.
        type: string
        required: false
      network_id:
        description: >-
          The unique ID for the network.
        type: string
        required: false
      mac_address:
        description: >-
          The unique media access control address (MAC address) assigned to the port.
        type: string
        required: false
      addresses:
        description: >-
          The list of IP address(es) assigned to the port.
        type: list
        entry_schema:
          type: string
        required: false

  tosca.datatypes.network.PortDef:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 5.2.5
    description: >-
      The PortDef type is a TOSCA data Type used to define a network port.
    derived_from: integer
    constraints:
    - in_range: [ 1, 65535 ]

  tosca.datatypes.network.PortSpec:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 5.2.6
    description: >-
      The PortSpec type is a complex TOSCA data Type used when describing port specifications for a
      network connection.
    derived_from: tosca.datatypes.Root
    properties:
      protocol:
        description: >-
          The required protocol used on the port.
        type: string
        constraints:
        - valid_values: [ udp, tcp, igmp ]
        default: tcp
      source:
        description: >-
          The optional source port.
        type: tosca.datatypes.network.PortDef
        required: false
      source_range:
        description: >-
          The optional range for source port.
        type: range
        constraints:
        - in_range: [ 1, 65535 ]
        required: false
      target:
        description: >-
          The optional target port.
        type: tosca.datatypes.network.PortDef
        required: false
      target_range:
        description: >-
          The optional range for target port.
        type: range
        constraints:
        - in_range: [ 1, 65535 ]
        required: false
`
}

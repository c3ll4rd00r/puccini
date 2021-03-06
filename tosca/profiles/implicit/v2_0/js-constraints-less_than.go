// This file was auto-generated from a YAML file

package v2_0

func init() {
	Profile["/tosca/implicit/2.0/js/constraints/less_than.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

clout.exec('tosca.lib.utils');

function validate(v1, v2) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	return tosca.compare(v1, v2) < 0;
}
`
}

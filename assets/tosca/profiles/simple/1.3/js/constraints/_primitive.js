
function validate(v, type) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	if (v === null)
		return true;
	return puccini.validateType(v, type);
}
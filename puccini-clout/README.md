puccini-clout
=============

`scriptlet exec`
----------------

Execute a JavaScript scriptlet embedded in a Clout. The optional input is a filesystem path or URL
to a Clout file. If no input is provided will attempt to read YAML from stdin. By default the output
(assuming the scriptlet generates output) will be output to stdout, but you can use the `--output/-o`
switch to specify a file (or direct to a file in the shell via `>`).

`exec` creates a specialized JavaScript environment in which to run the code, providing  access to
the parsed Clout structure as well as a few helper functions.

The default format for output is YAML, but you can switch to JSON or XML using `--format/-f`. Note
that Clout in JSON may lose some type information (e.g. JSON doesn't distinguish between an integer
and a float).

`scriptlet list`
----------------

Lists all available JavaScript scriptlets in the Clout.

`scriptlet get`
---------------

Prints out JavaScript scriptlet source code extracted from the Clout.

`scriptlet put`
---------------

Embeds/replaces JavaScript scriptlets in the Clout and outputs a new Clout. This can be used to add
scriptlets "on the fly" via piping (e.g. to add a plugin).

#!/bin/bash

if [ -f "/.bzkenv" ]; then
	source /.bzkenv
fi

{{if not .ContinueOnCmdError}}
set -e
{{end}}

echo "<PHASE:{{.Name}}>"

{{range .Commands}}
echo "<CMD:{{.}}>"
{{.}}
{{end}}

echo "cd \"$(pwd)\"" > /.bzkenv

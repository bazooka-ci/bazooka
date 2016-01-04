package main

import (
	"testing"
	"text/template"
)

func TestParseTemplate(t *testing.T) {
	template.Must(template.ParseFiles("template/bazooka_phase.sh"))
	template.Must(template.ParseFiles("template/bazooka_run.sh"))
	template.Must(template.ParseFiles("template/Dockerfile"))
}

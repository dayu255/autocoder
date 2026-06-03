package template

import _ "embed"

//go:embed assets/template.cpp
var DefaultCPP string

//go:embed assets/template.py
var DefaultPY string

package vconfig

import (
	"errors"
	"fmt"
)

type Section struct {
	data map[string][]string
	name string
}

// Create empty stand-alone section
func NewSection(name string) *Section {
	data := make(map[string][]string)
	var sec = Section{data: data, name: name}
	return &sec
}

func (sec *Section) Name() string {
	return sec.name
}

// assign values to variable; rewrite old values
func (sec *Section) SetValues(name string, values []string) {
	sec.data[name] = values
}

// add values to variable. extend list of old values with new ones
func (sec *Section) AddValues(name string, values []string) {
	if _, ok := sec.data[name]; ok { // if variable with this name already exists
		sec.data[name] = append(sec.data[name], values...)
	} else {
		sec.data[name] = values
	}
}

// create a copy of existing section
func (sec *Section) Duplicate() *Section {
	data := make(map[string][]string)
	newsec := Section{name: sec.name, data: data}

	for name, values := range sec.data {
		newvalues := make([]string, len(values))
		copy(newvalues, values)
		newsec.data[name] = newvalues
	}
	return &newsec
}

// return string array containing values of particular variable
func (sec *Section) GetValues(name string) ([]string, error) {
	if _, ok := sec.data[name]; ok {
		return sec.data[name], nil
	} else {
		return  nil, errors.New(fmt.Sprintf("Variable '%v' does not exist in this section", name))
	}
}

//returns single value
func (sec *Section) GetSingleValue(name string, default_value string) (string, error) {
	values, err := sec.GetValues(name)
	if err != nil {
		return default_value, err
	}
	if len(values) > 1 {
		return "", errors.New(fmt.Sprintf("Variable '%v' has multiple values", name))
	}
	return values[0], nil
}

// returns list of variables from section
func (sec *Section) GetVariables() []string {
	variables := make([]string, len(sec.data))
	idx := 0
	for variable := range sec.data {
		variables[idx] = variable
		idx++
	}
	return variables
}

// return section's content as string
func (sec *Section) ToString() string {
	var str string
	if sec.name != "__globalvars__" {
		str = fmt.Sprintf("[%v]", sec.name)
	}
	vars := sec.GetVariables()
	for _, variable := range vars {
		vals, _ := sec.GetValues(variable)
		//fmt.Println(vals)
		for _, val := range vals {
			str = fmt.Sprintf("%v\n%v=%v", str, variable, val)
		}
	}
	return str
}

//Merges content of secfrom to secto section
func (secto *Section) Merge(secfrom *Section) {
	for _, variable := range secfrom.GetVariables() {
		vals, _ := secfrom.GetValues(variable)
		secto.AddValues(variable, vals)
	}
}

// clears content of section keeping name unchanged
func (sec *Section) ClearContent() {
	sec.data = make(map[string][]string)
}

package vconfig
//+
import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type VConfig []*Section

// Creates new VConfig object
func New(filename string) (VConfig, error) {
	var vc VConfig

	if filename != "" {
		err := vc.fromFile(filename)
		if err != nil {
			return nil, err
		}
	}
	return vc, nil
}

// Loads data to VCOnfig from file
func FromFile(filename string) (VConfig, error) {
	vc, err := New(filename)
	if err != nil {
		return nil, err
	}
	return vc, nil
}

// Loads data to VConfig from string
func (vc *VConfig) FromString(str string) error {
	var current_section *Section
	gsection := NewSection("__globalvars__") // make __globalvars__ const
	vc.AddSection(gsection)
	current_section = gsection

	for idx, line := range strings.Split(string(str), "\n") {
        //TODO add strip in case of case extra spaces
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") { // if [] then beginnin of section
			if len(line) > 2 {
				key := line[1 : len(line)-1]
				sec := NewSection(key)
				vc.AddSection(sec)
				current_section = sec
			} else {
				return errors.New(fmt.Sprintf("Empty section at line %d", idx+1))
			}
			continue
		}
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
            splitline := strings.SplitN(line, "=", 2)

			if len(splitline) != 2 {
				return errors.New(fmt.Sprintf("Line %d: %s\nMissing equal sign?", idx+1, line))
			}

			variable := strings.Trim(splitline[0], " ")

			if len(variable) == 0 {
				return errors.New(fmt.Sprintf("Line %d: %s\nZero-length variable names are not allowed", idx+1, line))
			}

			value := strings.Trim(splitline[1], " ")

			if len(value) == 0 {
				return errors.New(fmt.Sprintf("Line %d: %s\nZero-length values are not allowed", idx+1, line))
			}

			current_section.AddValues(variable, []string{value})
		}
	}
	return nil
}

// Adds stand-alone section to existing VConfig
func (vc *VConfig) AddSection(sec *Section) error {
	if vc.sectionInVConfig(sec) {
		return errors.New("Section already exists in VConfig. Duplicate section before adding an existing one")
	}
	if sec.name == "__globalvars__" {
		gsec, err := vc.GetSections("__globalvars__")
		if err == nil { // if vc contains global variables
			variables := sec.GetVariables()
			for _, variable := range variables {
				val, _ := sec.GetValues(variable)
				gsec[0].AddValues(variable, val)

			}
            return nil

		}
    }
    //TODO check logic
    *vc = append(*vc, sec)

	return nil
}

// Creates new section and adds it to VConfig
func (vc *VConfig) NewSection(name string) *Section {
	nsec := NewSection(name)
	vc.AddSection(nsec)
	return nsec

}

// return section(s) based on its name
func (vc *VConfig) GetSections(name string) ([]*Section, error) {
	if name == "" {
		return vc.GetSections("__globalvars__")
	}
	sections := make([]*Section, 0)
	for _, sec := range *vc {
		if sec.name == name {
			sections = append(sections, sec)
		}
	}
	if len(sections) == 0 {
		return nil, errors.New(fmt.Sprintf("Section '%v' does not exist in this VConfig", name))
	}
	return sections, nil
}

// return content of VConfig as string
func (vc *VConfig) ToString() string {
	var str string
	gsections, err := vc.GetSections("__globalvars__")
	if err == nil {
		str = fmt.Sprintf("%v\n", gsections[0].ToString())
	}
	for _, sec := range *vc {
		if sec.name != "__globalvars__" {
			str = fmt.Sprintf("%v\n\n%v", str, sec.ToString())
		}
	}
	return str
}

// Write content of VCOnfig to file
func (vc *VConfig) ToFile(fname string) error {
	if err := ioutil.WriteFile(fname, []byte(vc.ToString()), 0666); err != nil { // check encoding:
		// 1. convert string to ascii file using excape characters: %q in sprinf
		//2.
        //TODO
		return err
	}
	return nil
}

// Return single value
func (vc *VConfig) GetSingleValue(sec_name string, var_name string, default_val string) (string, error) {
	if sec_name == "" {
		sec_name = "__globalvars__"
	}
	sec, err := vc.GetSections(sec_name)
	if err != nil {
		return "", err
	}
	if len(sec) > 1 {
		return "", errors.New(fmt.Sprintf("Multiple sections with name '%v'", sec_name))
	}
	val, err := sec[0].GetSingleValue(var_name, default_val)
	return val, err
}

func (vcto *VConfig) Merge(vcfrom VConfig) {
	for _, sec := range vcfrom {
		vcto.AddSection(sec.Duplicate())
	}
}

func (vc *VConfig) GetSectionsByVar(secName, varName, varValue string) ([]*Section, error){
    var sections []*Section
    err := errors.New("")
    if secName == ""{
        sections = *vc
    } else {
        sections, err = vc.GetSections(secName)
        if err != nil {return nil, err}
    }
    secs := make([]*Section, 0)
    for _, section := range sections {
        v, err := section.GetSingleValue(varName, "")
        if err != nil {return nil, err}
        if v == varValue {secs = append(secs, section)}
    }
    if len(secs) == 0 {return nil, errors.New(fmt.Sprintf("Cannot find section %v where %v=%v", secName, varName, varValue))}
    return secs, nil
}

//**************************   internals *****************************
// read config file and fill VConfig with data
func (vc *VConfig) fromFile(filename string) error {
	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = vc.FromString(string(str))
	if err != nil {
		return err
	}
	return nil
}

// Checks if given section belongs to VCOnfig
func (vc *VConfig) sectionInVConfig(sec *Section) bool {
	isSectionInVConfig := false
	for _, s := range *vc {
		if s == sec {
			isSectionInVConfig = true
		}
	}
	return isSectionInVConfig
}

package vconfig

import ("testing"
	"fmt"
    	"io/ioutil"
    	"os"
	"reflect"
	)

func TestReadingCorrectConfigFile(t *testing.T) {
	cfg := `
[testSection]
testVar = var1
testVar = var2
testVarSingle = var3

[testSection]
testVar = var1
testVar = var3

[testSectionSingle]
testVar4 = var4`
	currentDir,  err := os.Getwd()
	_ = os.RemoveAll(fmt.Sprintf("%v/test/", currentDir))

	if err := os.MkdirAll(fmt.Sprintf("%v/test", currentDir), 0700); err != nil {t.Fatal(err)}
    	defer func(){
            if err := os.RemoveAll(fmt.Sprintf("%v/test/", currentDir)); err != nil {
                t.Fatal(err)
            	}
        }()

	if err := ioutil.WriteFile(fmt.Sprintf("%v/test/config.conf", currentDir), []byte(cfg), 0700); err != nil {t.Fatal(err)}

	//read config with lib
	_,err = FromFile(fmt.Sprintf("%v/test/config.conf", currentDir))

	if err !=nil {
		t.Fatalf("Config was not sucessfully read, but should be...")
		}

}


func TestReadingInCorrectConfigFile(t *testing.T) {
	cfg := `
testSection]
testVar = var1
testVar = var2
testVarSingle = var3

[testSection
testVar = var1
testVar = var3

[testSectionSingle]
testVar4 = var4`
	currentDir,  err := os.Getwd()
	_ = os.RemoveAll(fmt.Sprintf("%v/test/", currentDir))

	if err := os.MkdirAll(fmt.Sprintf("%v/test", currentDir), 0700); err != nil {t.Fatal(err)}
    	defer func(){
            if err := os.RemoveAll(fmt.Sprintf("%v/test/", currentDir)); err != nil {
                t.Fatal(err)
            	}
        }()

	if err := ioutil.WriteFile(fmt.Sprintf("%v/test/config.conf", currentDir), []byte(cfg), 0700); err != nil {t.Fatal(err)}

	//read config with lib
	_,err = FromFile(fmt.Sprintf("%v/test/config.conf", currentDir))

	if err ==nil {
		t.Fatalf("Config should not be successfully read, but no error was reported...")
		}






}

func TestReadingValues(t *testing.T) {
	cfg := `
[testSection]
testVar = var1
testVarSingle = var3

[testSection]
testVar = var1


[testSection1]
testVar = var1
testVar = var3

[testSectionSingle]
testVar4 = var4`

	currentDir,  err := os.Getwd()
	_ = os.RemoveAll(fmt.Sprintf("%v/test/", currentDir))

	if err := os.MkdirAll(fmt.Sprintf("%v/test", currentDir), 0700); err != nil {t.Fatal(err)}
    	defer func(){
            if err := os.RemoveAll(fmt.Sprintf("%v/test/", currentDir)); err != nil {
                t.Fatal(err)
            	}
        }()

	if err := ioutil.WriteFile(fmt.Sprintf("%v/test/config.conf", currentDir), []byte(cfg), 0700); err != nil {t.Fatal(err)}

	//read config with lib
	vc,err := FromFile(fmt.Sprintf("%v/test/config.conf", currentDir))

	if err !=nil {
		t.Fatalf("Config was not sucessfully read, but should be...")
		}

	// check that lib will return errors when geting single value, but we have 2 or more sections or 2 or more variables...
	_, err =vc.GetSingleValue("testSection","testVar","test")
	if err == nil {
		t.Fatalf("When we read single value - lib should raise an error if we have several sections")
		}


	_, err =vc.GetSingleValue("testSection1","testVar","test")
	if err == nil {
		t.Fatalf("When we read single value - lib should raise an error if we have several variables in section")
		}


	// now we should successfully read single variable
	val, err :=vc.GetSingleValue("testSectionSingle","testVar4","test")
	if err != nil {
		t.Fatalf("We should have been succesfully read single variable but we did not")
		}

	if val !="var4" {
		t.Fatalf("Variable should be eq 'var4', but it is %v",val)
		}

	// test reading section and then values from it
	sections, err := vc.GetSections("testSection1")
	if err != nil {
		t.Fatalf("We should have been succesfully read section but we did not")
		}

	resStr:=""
	for _, section := range sections {
		resStr += section.ToString()
	}
	if resStr !=`[testSection1]
testVar = var1
testVar = var3` {
		t.Fatalf("Section was read, but has incorrect data - %v",resStr)
	}

	vals,err := sections[0].GetValues("testVar")
	if err != nil {
		t.Fatalf("We should have been succesfully read variables from section but we did not")
		}

	if vals[0]!="var1" &&  vals[1]!="var3" {
		t.Fatalf("We did not read variables from section correctly")
	}

	// check that getSingleValue from section will return error if more then 1 value
	_,err = sections[0].GetSingleValue("testVar","test")
	if err == nil {
		t.Fatalf("We should have raised an error when geting single value from section, but have more than 1 value")
		}

}

func TestEditingConfig(t *testing.T) {
	section := NewSection("testSection")

	section.AddValues("testValue",[] string {"1","2","3"})
	section.AddValues("testValue",[] string {"4"})

	vals,err := section.GetValues("testValue")
	if err != nil {
		t.Fatalf("We should have been succesfully read variables from section but we did not")
		}

	if vals[0]!="1" &&  vals[1]!="2" &&  vals[2]!="3" &&  vals[3]!="4" {
		t.Fatalf("We did not write variables to section correctly")
	}

	section.SetValues("testValue",[] string {"4"})
	vals,err = section.GetValues("testValue")
	if err != nil {
		t.Fatalf("We should have been succesfully read variables from section but we did not")
		}

	if vals[0]!="4" &&  len(vals)>1 {
		t.Fatalf("We did not overwrite variables in section correctly, values are %v", vals)
	}

	// test merge

	section_to_merge := NewSection("testSection_merge")
	section_to_merge.AddValues("testValue",[] string {"2"})
	section.Merge(section_to_merge)


	vals,err = section.GetValues("testValue")
	if err != nil {
		t.Fatalf("We should have been succesfully read variables from section but we did not")
		}

	if vals[0]!="4" && vals[1]!="2" &&  len(vals)>2 {
		t.Fatalf("We did not merge 2 sections correctly, values are %v", vals)
	}

	//GetVariables
	section.AddValues("testValue_1",[] string {"4"})
	variables := section.GetVariables()

	if variables[0]!="testValue" && variables[1]!="testValue_1" &&  len(variables)>2 {
		t.Fatalf("We did not get variables from section correctly, variables are - %v", variables)
	}

	//Duplicate
	section_dup := section.Duplicate()
	section_dup.SetValues("testValue_1",[] string {"3"})

	vals1,_ := section_dup.GetValues("testValue")
	vals2,_ := section.GetValues("testValue")
	// now we should check that testValue is the same for sections but testValue_1 is different
	if ! reflect.DeepEqual(vals1,vals2)  {
		t.Fatalf("Section duplicate did not work correctly, variable is not the same")
	}

	vals1,_ = section_dup.GetValues("testValue_1")
	vals2,_ = section.GetValues("testValue_1")
	// now we should check that testValue is the same for sections but testValue_1 is different
	if reflect.DeepEqual(vals1,vals2)  {
		t.Fatalf("Section duplicate did not work correctly, changed variable is the same")
	}

}
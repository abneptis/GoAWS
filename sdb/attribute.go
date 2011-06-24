package sdb

import (
	"http"
	"strconv"
)

type Item struct {
	Name      string
	Attribute []Attribute
}

type Attribute struct {
	Name    string
	Value   string
	Exists  *bool
	Replace *bool
}

type AttributeList []Attribute

type AttrListType string

const (
	ATTRIBUTE_LIST AttrListType = "Attribute."
	EXPECTED_LIST  AttrListType = "Expected."
)

func (self AttributeList) Values(afix AttrListType) (v http.Values) {
	v = http.Values{}
	for i := range self {
		prefix := string(afix) + strconv.Itoa(i+1) + "."
		v.Set(prefix+"Name", self[i].Name)
		if self[i].Value != "" {
			v.Set(prefix+"Value", self[i].Value)
		}
		if self[i].Replace != nil {
			if *self[i].Replace {
				v.Set(prefix+"Replace", "true")
			} else {
				v.Set(prefix+"Replace", "false")
			}
		}
		if self[i].Exists != nil {
			if *self[i].Exists {
				v.Set(prefix+"Exists", "true")
			} else {
				v.Set(prefix+"Exists", "false")
			}
		}
	}
	return
}


// miscelanious helper functions.
func AttrMissing(name string) Attribute {
	f := false
	return Attribute{Name: name, Exists: &f}
}

func AttrExists(name string) Attribute {
	t := true
	return Attribute{Name: name, Exists: &t}
}

func AttrEquals(name string, value string) Attribute {
	t := true
	return Attribute{Name: name, Value: value, Exists: &t}
}

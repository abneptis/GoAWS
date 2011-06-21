package simpledb

type Attribute struct {
  Name string
  Value string
  Exists *bool
}

type AttributeList []Attribute

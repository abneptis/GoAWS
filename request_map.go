package aws
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "com.abneptis.oss/maptools"

import "os"

/* 
  RequestMap is a generalized map structure with only
  one feature over a standard map:
    - An "Allowed" field, which is a map of [string]bool
  indicating whether a field is allowed at all (if present),
  and then if the field is required to pass validation (bool=true)
  
    This will not allow bad fields to be caught at compile time,
  however it will allow an error to be thrown up at runtime if
  something is amiss (e.g., missspellings, CapitalIzation, ...)
*/

type RequestMap struct {
  Allowed map[string]bool
  Values map[string]string
}

/* 
  Returns an error iff a non-allowed value is within the instance,
 or a required value is not specified.
  We assume that required values may NOT be the empty string in
 validation.
*/

func (self RequestMap)Validate()(err os.Error){
  for name, required := range(self.Allowed) {
    if required && self.Values[name] == "" {
      return os.NewError("Required value not found: " + name)
    }
  }
  for name, _ := range(self.Values) {
    err = self.ValidKey(name)
    if err != nil { return }
  }
  return
}

/* 
  Returns an error if the key is not allowed according
  to the internal Allowed map, otherwise nil.
*/
func (self RequestMap)ValidKey(key string)(err os.Error){
  _, valid := self.Allowed[key]
  if ! valid { err = os.NewError("Invalid key name " + key) }
  return
}

/*
    Sets a key / value pair on the object, unless the key is
  not permitted, in which case it will return an os.Error.
*/

func (self RequestMap)Set(key string, value string)(err os.Error){
  err = self.ValidKey(key)
  if err == nil {
    self.Values[key] = value
  }
  return
}

/*
  Returns the keystring.

  An error is only returned if the key name is invalid, NOT
  if the string is empty.
*/
func (self RequestMap)Get(key string)(value string, err os.Error){
  err = self.ValidKey(key)
  if err == nil {
    value = self.Values[key]
  }
  return
}

/*
  Returns a bool indicating whether the key is currently
  set.
*/
func (self RequestMap)IsSet(key string)(out bool){
  _, out = self.Values[key]
  return
}

/*
  Converts a requestmap to a string using the specified
  seperators, escaping, and optional sorting by keys.  

  Useful for canonicalization of the request map to various AWS services.
*/
func (self RequestMap)ToString(kvsep, pairsep string, kesc, vesc func(string)(string), sorted bool)(string){
  eMap := maptools.StringStringEscape(self.Values, kesc, vesc)
  return maptools.StringStringJoin(eMap, kvsep, pairsep, sorted)
}

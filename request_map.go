package goaws
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "com.abneptis.oss/maptools"

import "os"

type RequestMap struct {
  Allowed map[string]bool
  Values map[string]string
}

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

func (self RequestMap)ValidKey(key string)(err os.Error){
  _, valid := self.Allowed[key]
  if ! valid { err = os.NewError("Invalid key name " + key) }
  return
}

func (self RequestMap)Set(key string, value string)(err os.Error){
  err = self.ValidKey(key)
  if err == nil {
    self.Values[key] = value
  }
  return
}

func (self RequestMap)Get(key string)(value string, err os.Error){
  err = self.ValidKey(key)
  if err == nil {
    value = self.Values[key]
  }
  return
}

func (self RequestMap)IsSet(key string)(out bool){
  _, out = self.Values[key]
  return
}

func (self RequestMap)ToString(kvsep, pairsep string, kesc, vesc func(string)(string), sorted bool)(string){
  eMap := maptools.StringStringEscape(self.Values, kesc, vesc)
  return maptools.StringStringJoin(eMap, kvsep, pairsep, sorted)
}

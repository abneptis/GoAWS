package aws

import (
  "http"
  "sort"
)


func SortedEscape(v http.Values)(out string){
  keys := []string{}
  for k, _ := range(v){
    keys = append(keys, k)
  }
  sort.SortStrings(keys)
  for k := range(keys) {
    if k > 0 {
      out += "&"
    }
    out += http.URLEscape(keys[k]) + "=" + http.URLEscape(v.Get(keys[k]))
  }
  return
}

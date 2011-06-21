package aws

import (
  "com.abneptis.oss/urltools"
)

import (
  "http"
  "sort"
)

// (2011-06-21) - The standard go http.Values.Escape
// works properly for SQS  and S3, but it should be
// noted that at least SDB requiers more to be escaped
// than is officially standard.
//
// Sorted Escape also sorts the keys before joining them (needed
// for canonicalization).
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
    // out += http.URLEscape(keys[k]) + "=" + http.URLEscape(v.Get(keys[k]))
    out += escape(keys[k]) + "=" + escape(v.Get(keys[k]))
  }
  return
}

func escapeTest(b byte)(out bool){
  switch b {
    case 'a','b','c','d','e','f','g','h','i','j','k','l','m',
         'A','B','C','D','E','F','G','H','I','J','K','L','M',
         'n','o','p','q','r','s','t','u','v','w','x','y','z',
         'N','O','P','Q','R','S','T','U','V','W','X','Y','Z',
         '0','1','2','3','4','5','6','7','8','9','-','.','_':
      out = false
    default:
      out = true
  }
  return
}

func escape(in string)(out string){
  return urltools.Escape(in, escapeTest, urltools.PercentUpper)
}


package goaws
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

type mapFunc func(int)(int)

var NilMapping = func(in int)(int){return in}
var NilTest    = func(in int)(bool){return false}

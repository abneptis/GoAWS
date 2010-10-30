package goaws

type mapFunc func(int)(int)

var NilMapping = func(in int)(int){return in}

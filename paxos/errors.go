package paxos

import "errors"

var ErrorEmptyKey = errors.New("empty key error")
var ErrorEmptyValue = errors.New("empty value error")
var ErrorNoPeers = errors.New("nil peers present")
var ErrorNotFound = errors.New("not found error")
var ErrorUnkonwCommand = errors.New("unkonw command")

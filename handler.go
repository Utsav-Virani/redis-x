// -----------------------------------------------------
// REDIS-X
// © Utsav Virani
// Written by: (Utsav Virani)
// -----------------------------------------------------

package main

import "sync"

var SETs = map[string]string{}
var SETsmu = sync.RWMutex{} // because server is concurrent

var HSETs = map[string]map[string]string{}
var HSETsmu = sync.RWMutex{}

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].bulk
	value := args[1].bulk

	SETsmu.Lock()
	SETs[key] = value
	SETsmu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}
	key := args[0].bulk

	SETsmu.RLock()
	value, ok := SETs[key]
	SETsmu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}
	key := args[0].bulk
	field := args[1].bulk
	value := args[2].bulk

	HSETsmu.Lock()
	if _, ok := HSETs[key]; !ok {
		HSETs[key] = map[string]string{}
	}
	HSETs[key][field] = value
	HSETsmu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}
	key := args[0].bulk
	field := args[1].bulk

	SETsmu.RLock()
	value, ok := HSETs[key][field]
	SETsmu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}
	key := args[0].bulk

	SETsmu.RLock()
	fields, ok := HSETs[key]
	SETsmu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	var array []Value
	for field, value := range fields {
		array = append(array, Value{typ: "bulk", bulk: field})
		array = append(array, Value{typ: "bulk", bulk: value})
	}

	return Value{typ: "array", array: array}
}

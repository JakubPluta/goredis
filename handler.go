package main

import "sync"

var Handlers = map[string]func([]RedisMessage) RedisMessage{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"INFO":    info,
	"HSET":    hset,
	"HGETALL": hgetall,
	"HGET":    hget,
}

var SetsMap = map[string]string{}
var SetsMutex = sync.RWMutex{} // The sync.RWMutex type provides a reader/writer mutual exclusion lock.
var HSetsMap = map[string]map[string]string{}
var HSetsMutex = sync.RWMutex{}

func info(args []RedisMessage) RedisMessage {

	return RedisMessage{typ: "string", str: "This is a Redis server"}
}

func ping(args []RedisMessage) RedisMessage {

	if len(args) == 0 {
		return RedisMessage{typ: "string", str: "PONG"}
	}

	return RedisMessage{typ: "string", str: args[0].bulk}
}

// set sets the value of a key in a Redis map.
// args []RedisMessage
// RedisMessage
func set(args []RedisMessage) RedisMessage {
	// defines a set function that sets a key-value pair in a Redis map
	// and returns an "OK" message if successful, or an error message if the input arguments are invalid.

	if len(args) != 2 {
		return RedisMessage{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}
	k, v := args[0].bulk, args[1].bulk
	SetsMutex.Lock()
	SetsMap[k] = v
	SetsMutex.Unlock()
	return RedisMessage{typ: "string", str: "OK"}
}

// get returns the value associated with the given key.
//
// args []RedisMessage
// RedisMessage
func get(args []RedisMessage) RedisMessage {
	// defines a get function that retrieves the value associated with a given key from a map called SetsMap.
	// If the key is not found, it returns null; if the number of arguments is not 1, it returns an error message.

	if len(args) != 1 {
		return RedisMessage{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	k := args[0].bulk
	SetsMutex.Lock()
	v, ok := SetsMap[k]
	SetsMutex.Unlock()
	if !ok {
		return RedisMessage{typ: "null"}
	}
	return RedisMessage{typ: "bulk", bulk: v}

}

func hset(args []RedisMessage) RedisMessage {
	// defines a hset function that sets a key-value pair in a Redis hash map.
	// and returns an "OK" message if successful, or an error message if the input arguments are invalid.
	if len(args) != 3 {
		return RedisMessage{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}
	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSetsMutex.Lock()
	if _, ok := HSetsMap[hash]; !ok {
		HSetsMap[hash] = map[string]string{}
	}
	HSetsMap[hash][key] = value
	HSetsMutex.Unlock()

	return RedisMessage{typ: "string", str: "OK"}

}

func hget(args []RedisMessage) RedisMessage {
	if len(args) != 2 {
		return RedisMessage{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSetsMutex.Lock()
	value, ok := HSetsMap[hash][key]
	HSetsMutex.Unlock()
	if !ok {
		return RedisMessage{typ: "null"}
	}
	return RedisMessage{typ: "bulk", bulk: value}
}

func hgetall(args []RedisMessage) RedisMessage {
	if len(args) != 1 {
		return RedisMessage{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	HSetsMutex.Lock()
	value, ok := HSetsMap[hash]
	HSetsMutex.Unlock()

	if !ok {
		return RedisMessage{typ: "null"}
	}

	values := []RedisMessage{}

	for k, v := range value {
		values = append(values, RedisMessage{typ: "bulk", bulk: k})
		values = append(values, RedisMessage{typ: "bulk", bulk: v})

	}

	return RedisMessage{typ: "array", array: values}
}

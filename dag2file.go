package main

import (
	"encoding/json"
	"strings"
)

// sec/LEVI_104/gm.rs
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	pathSegments := strings.Split(path, "/")
	if len(pathSegments) == 0 { // 没东西
		return nil
	}

	objBytes, _ := store.Get(hash)

	var obj Object
	json.Unmarshal(objBytes, &obj)

	return recursiveSearch(store, obj, pathSegments, hp)
}

func recursiveSearch(store KVStore, obj Object, pathSegments []string, hp HashPool) []byte {
	if len(pathSegments) == 0 { // 没东西
		return nil
	}

	for _, value := range obj.Links {
		switch pathSegments[0] {
			case "blob":
				blob_value, _ := store.Get(CalHash(value, hp))
				return store.Get(blob_value)
			case "link":
				// 所有入stack
				return recursiveSearch(store, getObject(store, value.Hash), pathSegments, hp)
			case "tree":
				// 取[1:]的原因是tree入stack，剩下递归
				return recursiveSearch(store, getObject(store, value.Hash), pathSegments[1:], hp)
		}
	}

	return nil
}

////////////////////////////////////// [Helper functions] ////////////////////////////////

func CalHash(data []byte, h hash.Hash) []byte {
	h.Reset()
	hash := h.Sum(data)
	h.Reset()
	return hash
}

func getObject(store KVStore, hash []byte) Object {
	objBytes, _ := store.Get(hash)

	var obj Object
	json.Unmarshal(objBytes, &obj)
	return obj
}
package merkledag

import (
	"encoding/json"
	"hash"
)

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

func Add(store KVStore, node Node, h hash.Hash) []byte {
	var obj *Object

	if node.Type() == FILE {
		file := node.(File)
		obj = StoreFile(store, file, h)
	} else if node.Type() == DIR {
		dir := node.(Dir)
		obj = StoreDir(store, dir, h)
	}

	jsonMarshal, _ := json.Marshal(obj)
	hash := calculateHash(jsonMarshal)
	store.Put(hash, jsonMarshal)

	return hash
}

func calculateHash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func StoreFile(store KVStore, file File, h hash.Hash) *Object {
	if len(file.Bytes()) <= 256*1024 {
		data := file.Bytes()
		blob := Object{Data: data, Links: nil}
		return &blob
	}

	// Implement chunking logic for large files here
}

func StoreDir(store KVStore, dir Dir, h hash.Hash) *Object {
	it := dir.It()
	treeObject := &Object{}
	for it.Next() {
		n := it.Node()
		switch n.Type() {
		case FILE:
			file := n.(File)
			tmp := StoreFile(store, file, h)
			treeObject.Links = append(treeObject.Links, Link{
				Hash: calculateHash(json.Marshal(tmp)),
				Size: int(file.Size()),
				Name: file.Name(),
			})
		case DIR:
			dir := n.(Dir)
			tmp := StoreDir(store, dir, h)
			treeObject.Links = append(treeObject.Links, Link{
				Hash: calculateHash(json.Marshal(tmp)),
				Size: int(dir.Size()),
				Name: dir.Name(),
			})
		}
	}
	return treeObject
}

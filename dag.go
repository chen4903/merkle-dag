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
	// TODO 将分片写入到KVStore中，并返回Merkle Root

	var object Object

	if node.Type() == FILE { // 处理文件节点
		fileNode := node.(File)
		object.Data = fileNode.Bytes()
	} else if node.Type() == DIR { // 处理目录节点
		dirNode := node.(Dir)
		it := dirNode.It()

		for it.Next() {
			childNode := it.Node()
			// 递归
			childHash := Add(store, childNode, h)
			object.Links = append(object.Links, Link{
				Name: childNode.Name(),
				Hash: childHash,
				Size: int(childNode.Size()),
			})
		}
	}

	// 序列化对象
	objectBytes, err := json.Marshal(object)
	if err != nil {
		panic(err) // 处理错误
	}

	// 计算哈希
	h.Reset()
	h.Write(objectBytes)
	hashBytes := h.Sum(nil)

	// 存储对象: Key是hash，value是我们存储的节点对象
	if err := store.Put(hashBytes, objectBytes); err != nil {
		panic(err) // 处理错误
	}

	return hashBytes
}

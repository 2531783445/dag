package merkledag

import (
	"encoding/json"
	"strings"
)

func RetrieveList(store KVStore, hash []byte) []byte {
	var data []byte
	var nextHash []byte = hash

	for nextHash != nil {
		value, err := store.Get(nextHash)
		if err != nil {
			return nil
		}

		var node ListNode
		err = json.Unmarshal(value, &node)
		if err != nil {
			return nil
		}

		fileData, err := store.Get(node.Hash)
		if err != nil {
			return nil
		}

		data = append(data, fileData...)

		nextHash = node.Next
	}

	return data
}

func HashFile(store KVStore, hash []byte, path string, hp HashPool) []byte {
	value, err := store.Get(hash)
	if err != nil {
		return nil
	}

	var obj Object
	err = json.Unmarshal(value, &obj)
	if err != nil {
		return nil
	}

	parts := strings.Split(path, "/")

	if len(obj.Links) > 0 {
		for i, link := range obj.Links {
			if link.Name == parts[0] {
				if len(parts) == 1 {
					if string(obj.Data[i][0]) == "list" {
						return RetrieveList(store, link.Hash)
					} else if string(obj.Data[i][0]) == "blob" {
						value, err := store.Get(link.Hash)
						if err != nil {
							return nil
						}
						return value
					}
					return link.Hash
				}
				return HashFile(store, link.Hash, strings.Join(parts[1:], "/"), hp)
			}
		}
		return nil
	}
	return nil
}

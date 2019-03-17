package main

// FNVHash 计算key的hash值, prime为hash种子
func FNVHash(key string, prime uint64) uint64 {
	var hash uint64 = 2166136261
	// var prime uint64 = 16777619
	for i := 0; i < len(key); i++ {
		hash = (hash * prime) ^ uint64(key[i])
	}

	return hash

}

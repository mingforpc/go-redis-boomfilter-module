package main

// calculateFilterSize 给定的布隆过滤器是bit位的大小，但是申请内存时需要的是多少byte
func calculateFilterSize(filtersize int64) int64 {

	byteCount := filtersize / 8

	if byteCount == 0 || filtersize%byteCount > 0 {
		byteCount++
	}

	return byteCount
}

package main

/*
#include "../redismodule.h"

// 返回long long格式的响应
int replyWithLongLong(RedisModuleCtx *ctx, long long ll);
// 返回错误信息
int replyWithError(RedisModuleCtx *ctx, const char *err);
// 返回字符串信息
int replyWithSimpleString(RedisModuleCtx *ctx, const char *msg);

// 命令的初始化
int initCommonds(RedisModuleCtx *ctx);

// 从RedisModuleString **argv获取参数
RedisModuleString *getArgvString(RedisModuleString **argv, int index);
// 从RedisModuleString中获取字符串
const char *getModuleStringPtr(RedisModuleString *argv, size_t *len);
// 从RedisModuleString中获取整型
int stringToLongLong(const RedisModuleString *str, long long *ll);
// 将*char转为RedisModuleString
RedisModuleString *createString(RedisModuleCtx *ctx, const char *ptr, size_t len);
// long long 转为RedisModuleString
RedisModuleString *createStringFromLongLong(RedisModuleCtx *ctx, long long ll);
// 释放 RedisModuleString
void freeModuleString(RedisModuleCtx *ctx, RedisModuleString *str);

// 打开一个RedisModuleKey结构
RedisModuleKey *openKey(RedisModuleCtx *ctx, RedisModuleString *keyname, int mode);
// 关闭RedisModuleKey结构
void closeKey(RedisModuleKey *kp);
// 检查RedisModuleKey中的类型
int keyType(RedisModuleKey *kp);
// Truncate字符串
int stringTruncate(RedisModuleKey *key, size_t newlen);
// Redis Set 方法
int stringSet(RedisModuleKey *key, RedisModuleString *str);
// Redis Get 方法(只读)
char *stringGet(RedisModuleKey *key, size_t *len);
// Redis Del 方法
int delKey(RedisModuleKey *key);
// 获取hash字典的key数量
size_t getHashLength(RedisModuleKey *key);

// hash set函数
int hashSet(RedisModuleKey *key, int flags, RedisModuleString *vk , RedisModuleString *vv);
// hash get函数
int hashGet(RedisModuleKey *key, int flags, RedisModuleString *vk , RedisModuleString **vv);
// setbit 命令, 使用 RedisModule_Call 实现
RedisModuleCallReply *setBit(RedisModuleCtx *ctx, char *key, char *offset, char *val);
// getbit 命令, 使用 RedisModule_Call 实现
RedisModuleCallReply *getBit(RedisModuleCtx *ctx, char *key, char *offset);
// getReplyType 获取 RedisModuleCallReply 的类型
int getReplyType(RedisModuleCallReply *reply);
// 从 Reply 中获取整型结果
long long getIntegerFromReply(RedisModuleCallReply *reply);
// 释放 RedisModuleCallReply
void freeReply(RedisModuleCallReply *reply);

*/
import (
	"C"
)
import (
	"fmt"
	"math/rand"
	"strconv"
)

const moduleName = "boomfilter"
const version = 1

//export RedisModule_OnLoad
func RedisModule_OnLoad(ctx *C.RedisModuleCtx, argv **C.RedisModuleString, argc C.int) C.int {

	if C.RedisModule_Init(ctx, C.CString(moduleName), C.int(version), C.REDISMODULE_APIVER_1) == C.REDISMODULE_ERR {
		return C.REDISMODULE_ERR
	}

	if C.initCommonds(ctx) == C.REDISMODULE_ERR {
		return C.REDISMODULE_ERR
	}

	return C.REDISMODULE_OK
}

//export boomfilterCreate
// 创建布隆过滤器,及其相关key:
// 布隆过滤器key: {用户输入的key name}
// 保存hash种子的key: boomfilter.{key name}.hashseek.set
// 保存布隆过滤器长度的key: boomfilter.{key name}.total.size
//                       因为布隆过滤器长度是x个bit, 但内存申请是y个byte
// usage: boomfilter.createboomfilter {key name} {hash func count} {filter size}
func boomfilterCreate(ctx *C.RedisModuleCtx, argv **C.RedisModuleString, argc C.int) C.int {

	gargc := int(argc)

	// boomfilter.createboomfilter {key name} {hash func count} {filter size}
	if gargc != 4 {

		C.replyWithError(ctx, C.CString("ERR invalid arguments"))

		return C.REDISMODULE_ERR
	}

	keystr := C.getArgvString(argv, C.int(1))
	hashCountStr := C.getArgvString(argv, C.int(2))
	filterSizeStr := C.getArgvString(argv, C.int(3))

	var hashCount C.longlong
	var filterSize C.longlong

	// golang string key
	gkeystr := C.GoString(C.getModuleStringPtr(keystr, nil))

	// hashCount
	if C.stringToLongLong(hashCountStr, &hashCount) != C.REDISMODULE_OK {
		C.replyWithError(ctx, C.CString("ERR invalid arguments"))
		return C.REDISMODULE_ERR
	}
	// filterSize
	if C.stringToLongLong(filterSizeStr, &filterSize) != C.REDISMODULE_OK {
		C.replyWithError(ctx, C.CString("ERR invalid arguments"))
		return C.REDISMODULE_ERR
	}

	// 检查key是否已经被占用
	key := getRedisKey(ctx, gkeystr, C.REDISMODULE_WRITE)
	if checkKeyExist(ctx, key) {
		C.replyWithError(ctx, C.CString("boomfilter key existed"))
		return C.REDISMODULE_ERR
	}
	defer C.closeKey(key)

	// 检查 boomfilter.{key name}.hashseek.set 是否存在
	// boomfilter.{key name}.hashseek.set 该键用来保存 {key name} 的布隆过滤器用了多少个hash和每个hash的seek
	hashSeekListKey := getRedisKey(ctx, fmt.Sprintf("boomfilter.%s.hashseek.set", gkeystr), C.REDISMODULE_WRITE)
	if checkKeyExist(ctx, hashSeekListKey) {
		C.replyWithError(ctx, C.CString("boomfilter hash seek key existed"))
		return C.REDISMODULE_ERR
	}
	defer C.closeKey(hashSeekListKey)

	// 检查 boomfilter.{key name}.total.size 是否存在
	// boomfilter.{key name}.total.size 是用来保存布隆过滤器的大小
	filterSizeKey := getRedisKey(ctx, fmt.Sprintf("boomfilter.%s.total.size", gkeystr), C.REDISMODULE_WRITE)
	if checkKeyExist(ctx, filterSizeKey) {
		C.replyWithError(ctx, C.CString("boomfilter total size key existed"))
		return C.REDISMODULE_ERR
	}
	defer C.closeKey(filterSizeKey)

	// 创建指定长度的内存区域
	byteCount := calculateFilterSize(int64(filterSize))
	C.stringTruncate(key, C.size_t(byteCount))

	// 创建指定数量的hash seek,并保存到 boomfilter.{key name}.hashseek.set 中
	ghashCount := int64(hashCount)
	hashSeekSet := make(map[int64]int64)
	for ghashCount > 0 {

		seek := rand.Int63n(33555238)

		if seek > 0 && hashSeekSet[seek] == 0 {
			hashSeekSet[seek] = seek
			seekStr := C.createStringFromLongLong(ctx, C.longlong(seek))

			ghashCount--
			// 存入hash字典中，key: index, value: hash seek
			hsahCountStr := C.createStringFromLongLong(ctx, C.longlong(ghashCount))
			C.hashSet(hashSeekListKey, C.REDISMODULE_HASH_NONE, hsahCountStr, seekStr)

			C.freeModuleString(ctx, seekStr)
			C.freeModuleString(ctx, hsahCountStr)
		}
	}

	// 创建 boomfilter.{key name}.total.size， 用来保存布隆过滤器的总长度
	C.stringSet(filterSizeKey, filterSizeStr)

	C.replyWithSimpleString(ctx, C.CString("OK"))
	return C.REDISMODULE_OK
}

//export boomfilterClean
// 删除指定布隆过滤器对应的key
// usage: boomfilter.cleanboomfilter {key}
func boomfilterClean(ctx *C.RedisModuleCtx, argv **C.RedisModuleString, argc C.int) C.int {

	res := 0

	keystr := C.getArgvString(argv, C.int(1))
	key := C.openKey(ctx, keystr, C.REDISMODULE_WRITE)
	defer C.closeKey(key)

	// 删除布隆过滤器的key
	if C.delKey(key) == C.REDISMODULE_OK {
		res++
	}

	// golang string key
	gkeystr := C.GoString(C.getModuleStringPtr(keystr, nil))

	// 删除保存hash seek的key
	hashSeekListKey := getRedisKey(ctx, fmt.Sprintf("boomfilter.%s.hashseek.set", gkeystr), C.REDISMODULE_WRITE)
	defer C.closeKey(hashSeekListKey)

	if C.delKey(hashSeekListKey) == C.REDISMODULE_OK {
		res++
	}

	// 删除保存总长度的key
	filterSizeKey := getRedisKey(ctx, fmt.Sprintf("boomfilter.%s.total.size", gkeystr), C.REDISMODULE_WRITE)
	defer C.closeKey(filterSizeKey)

	if C.delKey(filterSizeKey) == C.REDISMODULE_OK {
		res++
	}

	// 返回影响的key数量
	C.replyWithLongLong(ctx, C.longlong(res))

	return C.REDISMODULE_OK
}

//export boomfilterAdd
// 向指定布隆过滤器添加元素
// usage: boomfilter.add {key} {val}
func boomfilterAdd(ctx *C.RedisModuleCtx, argv **C.RedisModuleString, argc C.int) C.int {

	gargc := int(argc)

	if gargc < 3 {
		C.replyWithError(ctx, C.CString("ERR invalid arguments"))
		return C.REDISMODULE_ERR
	}

	keystr := C.getArgvString(argv, C.int(1))
	gkeystr := C.GoString(C.getModuleStringPtr(keystr, nil))
	key := getRedisKey(ctx, gkeystr, C.REDISMODULE_WRITE)
	defer C.closeKey(key)

	if C.keyType(key) != C.REDISMODULE_KEYTYPE_STRING {
		// key指定的布隆过滤器不存在
		C.replyWithLongLong(ctx, C.longlong(0))
	}

	// 检查 boomfilter.{key name}.hashseek.set 是否存在
	// boomfilter.{key name}.hashseek.set 该键用来保存 {key name} 的布隆过滤器用了多少个hash和每个hash的seek
	hashSeekListKey := getRedisKey(ctx, fmt.Sprintf("boomfilter.%s.hashseek.set", gkeystr), C.REDISMODULE_WRITE)
	if C.keyType(hashSeekListKey) != C.REDISMODULE_KEYTYPE_HASH {
		C.replyWithError(ctx, C.CString("boomfilter hash seek key not exist"))
		return C.REDISMODULE_ERR
	}
	defer C.closeKey(hashSeekListKey)

	// 检查 boomfilter.{key name}.total.size 是否存在
	// boomfilter.{key name}.total.size 是用来保存布隆过滤器的大小
	filterSizeKey := getRedisKey(ctx, fmt.Sprintf("boomfilter.%s.total.size", gkeystr), C.REDISMODULE_WRITE)
	if C.keyType(filterSizeKey) != C.REDISMODULE_KEYTYPE_STRING {
		C.replyWithError(ctx, C.CString("boomfilter total size not exist"))
		return C.REDISMODULE_ERR
	}
	defer C.closeKey(filterSizeKey)

	// 获取boomfilter的总长度
	len := C.size_t(0)
	sizeStr := C.stringGet(filterSizeKey, &len)
	if sizeStr == nil {
		C.replyWithError(ctx, C.CString("Err get boomfilter total size"))
		return C.REDISMODULE_ERR
	}
	gsize, err := strconv.ParseUint(C.GoString(sizeStr), 10, 0)
	if err != nil {
		C.replyWithError(ctx, C.CString("Err parse boomfilter total size"))
		return C.REDISMODULE_ERR
	}

	// 获取hash seek, 并保存入seekList中
	seekList := make([]int64, 0)
	hashSeekCount := C.getHashLength(hashSeekListKey)
	ghashSeekCount := int(hashSeekCount)
	for i := 0; i < ghashSeekCount; i++ {
		// 获取hash seek, 并计算hash，设置对应位
		seekHkey := C.createStringFromLongLong(ctx, C.longlong(i))
		seekHval := &C.RedisModuleString{}
		if C.hashGet(hashSeekListKey, C.REDISMODULE_HASH_NONE, seekHkey, &seekHval) == C.REDISMODULE_ERR {
			// 键不存在
			C.replyWithError(ctx, C.CString("ERR invalid data, please check"))
			return C.REDISMODULE_ERR
		}

		var cSeek C.longlong
		if C.stringToLongLong(seekHval, &cSeek) == C.REDISMODULE_ERR {
			// seek的值不是整型
			C.replyWithError(ctx, C.CString("ERR invalid data, please check"))
			return C.REDISMODULE_ERR
		}
		seek := int64(cSeek)
		seekList = append(seekList, seek)

		C.freeModuleString(ctx, seekHkey)
	}

	// 处理多个添加的val
	for i := 2; i < gargc; i++ {

		val := C.getArgvString(argv, C.int(i))
		gval := C.GoString(C.getModuleStringPtr(val, nil))

		for _, seek := range seekList {
			// 计算hash值，并通过setbit标志位
			hashVal := FNVHash(gval, uint64(seek))
			offset := hashVal % gsize
			offsetStr := strconv.FormatUint(offset, 10)

			callReply, err := C.setBit(ctx, C.CString(gkeystr), C.CString(offsetStr), C.CString("1"))
			if callReply == nil || err != nil {
				C.replyWithError(ctx, C.CString("ERR set bit"))
				return C.REDISMODULE_ERR
			}

			C.freeReply(callReply)
		}

	}

	// 返回处理的val个数
	C.replyWithLongLong(ctx, C.longlong(gargc-2))

	return C.REDISMODULE_OK
}

//export boomfilterExists
// 查询传入的值是否已经存在于指定的boomfilter中
// usage: boomfilter.exists {key} {val}
func boomfilterExists(ctx *C.RedisModuleCtx, argv **C.RedisModuleString, argc C.int) C.int {

	gargc := int(argc)
	if gargc < 3 {
		C.replyWithError(ctx, C.CString("ERR invalid arguments"))
		return C.REDISMODULE_ERR
	}

	keystr := C.getArgvString(argv, C.int(1))
	gkeystr := C.GoString(C.getModuleStringPtr(keystr, nil))
	key := getRedisKey(ctx, gkeystr, C.REDISMODULE_WRITE)
	defer C.closeKey(key)

	if C.keyType(key) != C.REDISMODULE_KEYTYPE_STRING {
		// key指定的布隆过滤器不存在
		C.replyWithLongLong(ctx, C.longlong(0))
	}

	// 检查 boomfilter.{key name}.hashseek.set 是否存在
	// boomfilter.{key name}.hashseek.set 该键用来保存 {key name} 的布隆过滤器用了多少个hash和每个hash的seek
	hashSeekListKey := getRedisKey(ctx, fmt.Sprintf("boomfilter.%s.hashseek.set", gkeystr), C.REDISMODULE_WRITE)
	if C.keyType(hashSeekListKey) != C.REDISMODULE_KEYTYPE_HASH {
		C.replyWithError(ctx, C.CString("boomfilter hash seek key not exist"))
		return C.REDISMODULE_ERR
	}
	defer C.closeKey(hashSeekListKey)

	// 检查 boomfilter.{key name}.total.size 是否存在
	// boomfilter.{key name}.total.size 是用来保存布隆过滤器的大小
	filterSizeKey := getRedisKey(ctx, fmt.Sprintf("boomfilter.%s.total.size", gkeystr), C.REDISMODULE_WRITE)
	if C.keyType(filterSizeKey) != C.REDISMODULE_KEYTYPE_STRING {
		C.replyWithError(ctx, C.CString("boomfilter total size not exist"))
		return C.REDISMODULE_ERR
	}
	defer C.closeKey(filterSizeKey)

	// 获取boomfilter的总长度
	len := C.size_t(0)
	sizeStr := C.stringGet(filterSizeKey, &len)
	if sizeStr == nil {
		C.replyWithError(ctx, C.CString("Err get boomfilter total size"))
		return C.REDISMODULE_ERR
	}
	gsize, err := strconv.ParseUint(C.GoString(sizeStr), 10, 0)
	if err != nil {
		C.replyWithError(ctx, C.CString("Err parse boomfilter total size"))
		return C.REDISMODULE_ERR
	}

	// 获取hash seek, 并保存入seekList中
	seekList := make([]int64, 0)
	hashSeekCount := C.getHashLength(hashSeekListKey)
	ghashSeekCount := int(hashSeekCount)
	for i := 0; i < ghashSeekCount; i++ {
		// 获取hash seek, 并计算hash，设置对应位
		seekHkey := C.createStringFromLongLong(ctx, C.longlong(i))
		seekHval := &C.RedisModuleString{}
		if C.hashGet(hashSeekListKey, C.REDISMODULE_HASH_NONE, seekHkey, &seekHval) == C.REDISMODULE_ERR {
			// 键不存在
			C.replyWithError(ctx, C.CString("ERR invalid data, please check"))
			return C.REDISMODULE_ERR
		}

		var cSeek C.longlong
		if C.stringToLongLong(seekHval, &cSeek) == C.REDISMODULE_ERR {
			// seek的值不是整型
			C.replyWithError(ctx, C.CString("ERR invalid data, please check"))
			return C.REDISMODULE_ERR
		}
		seek := int64(cSeek)
		seekList = append(seekList, seek)

		C.freeModuleString(ctx, seekHkey)
	}

	val := C.getArgvString(argv, C.int(2))
	gval := C.GoString(C.getModuleStringPtr(val, nil))

	for _, seek := range seekList {
		// 计算hash值，并通过setbit标志位
		hashVal := FNVHash(gval, uint64(seek))
		offset := hashVal % gsize
		offsetStr := strconv.FormatUint(offset, 10)

		callReply, err := C.getBit(ctx, C.CString(gkeystr), C.CString(offsetStr))
		if callReply == nil || err != nil {
			C.replyWithError(ctx, C.CString("ERR get bit"))
			return C.REDISMODULE_ERR
		}

		cVal := C.getIntegerFromReply(callReply)

		if int64(cVal) == 0 {
			// 返回0，则表明不存在

			C.replyWithLongLong(ctx, C.longlong(0))
			return C.REDISMODULE_OK
		}

		C.freeReply(callReply)

	}

	// val 存在
	C.replyWithLongLong(ctx, C.longlong(1))

	return C.REDISMODULE_OK
}

func getRedisKey(ctx *C.RedisModuleCtx, keystr string, openmode C.int) *C.RedisModuleKey {
	redisKeyStr := C.createString(ctx, C.CString(keystr), C.size_t(len(keystr)))
	defer C.freeModuleString(ctx, redisKeyStr)
	key := C.openKey(ctx, redisKeyStr, openmode)

	return key
}

// checkKeyExist 检查键是否存在
func checkKeyExist(ctx *C.RedisModuleCtx, key *C.RedisModuleKey) bool {

	if C.keyType(key) != C.REDISMODULE_KEYTYPE_EMPTY {
		return true
	}

	return false
}

func main() {
}

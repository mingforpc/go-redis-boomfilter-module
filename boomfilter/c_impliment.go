package main

/*
#include "../redismodule.h"


int replyWithLongLong(RedisModuleCtx *ctx, long long ll) {
	return RedisModule_ReplyWithLongLong(ctx, ll);
}

int replyWithError(RedisModuleCtx *ctx, const char *err) {
	return RedisModule_ReplyWithError(ctx, err);

}

int replyWithSimpleString(RedisModuleCtx *ctx, const char *msg) {
	return RedisModule_ReplyWithSimpleString(ctx, msg);
}

RedisModuleString *getArgvString(RedisModuleString **argv, int index) {
	return argv[index];
}

const char *getModuleStringPtr(RedisModuleString *argv, size_t *len) {
	return RedisModule_StringPtrLen(argv, len);
}

int stringToLongLong(const RedisModuleString *str, long long *ll) {
	return RedisModule_StringToLongLong(str, ll);
}

RedisModuleKey *openKey(RedisModuleCtx *ctx, RedisModuleString *keyname, int mode) {
	RedisModuleKey *key;
	key = RedisModule_OpenKey(ctx, keyname, mode);

	return key;
}

RedisModuleString *createString(RedisModuleCtx *ctx, const char *ptr, size_t len) {
	return RedisModule_CreateString(ctx, ptr, len);
}

RedisModuleString *createStringFromLongLong(RedisModuleCtx *ctx, long long ll) {
	return RedisModule_CreateStringFromLongLong(ctx, ll);
}

void freeModuleString(RedisModuleCtx *ctx, RedisModuleString *str) {
	return RedisModule_FreeString(ctx, str);
}

void closeKey(RedisModuleKey *kp) {
	RedisModule_CloseKey(kp);
}

int keyType(RedisModuleKey *kp) {
	return RedisModule_KeyType(kp);
}

int stringTruncate(RedisModuleKey *key, size_t newlen) {
	return RedisModule_StringTruncate(key, newlen);
}

int stringSet(RedisModuleKey *key, RedisModuleString *str) {
	return RedisModule_StringSet(key, str);
}

int hashSet(RedisModuleKey *key, int flags, RedisModuleString *vk , RedisModuleString *vv) {
	return RedisModule_HashSet(key, flags, vk, vv, NULL);

}

int hashGet(RedisModuleKey *key, int flags, RedisModuleString *vk , RedisModuleString **vv) {
	return RedisModule_HashGet(key, flags, vk, vv, NULL);
}

char *stringGet(RedisModuleKey *key, size_t *len) {
	return RedisModule_StringDMA(key, len, REDISMODULE_READ);
}

int delKey(RedisModuleKey *key) {
	return RedisModule_DeleteKey(key);
}

size_t getHashLength(RedisModuleKey *key) {
	return RedisModule_ValueLength(key);
}

RedisModuleCallReply *setBit(RedisModuleCtx *ctx, char *key, char *offset, char *val) {
	return RedisModule_Call(ctx, "SETBIT", "ccc", key, offset, val);
}

RedisModuleCallReply *getBit(RedisModuleCtx *ctx, char *key, char *offset) {
	return RedisModule_Call(ctx, "GETBIT", "cc", key, offset);
}

int getReplyType(RedisModuleCallReply *reply) {
	return RedisModule_CallReplyType(reply);
}

long long getIntegerFromReply(RedisModuleCallReply *reply) {
	return RedisModule_CallReplyInteger(reply);
}

void freeReply(RedisModuleCallReply *reply) {
	return RedisModule_FreeCallReply(reply);
}

extern int boomfilterCreate(RedisModuleCtx *ctx, RedisModuleString  **argv, int argc);
extern int boomfilterClean(RedisModuleCtx *ctx, RedisModuleString  **argv, int argc);
extern int boomfilterAdd(RedisModuleCtx *ctx, RedisModuleString  **argv, int argc);
extern int boomfilterExists(RedisModuleCtx *ctx, RedisModuleString  **argv, int argc);

int initCommonds(RedisModuleCtx *ctx) {


	if (RedisModule_CreateCommand(ctx,"boomfilter.createboomfilter", boomfilterCreate, "write", 1, 0, 0) == REDISMODULE_ERR)
		return REDISMODULE_ERR;

	if (RedisModule_CreateCommand(ctx, "boomfilter.cleanboomfilter", boomfilterClean, "write", 1, 0, 0) == REDISMODULE_ERR)
		return REDISMODULE_ERR;

	if (RedisModule_CreateCommand(ctx, "boomfilter.add", boomfilterAdd, "write", 1, 0, 0) == REDISMODULE_ERR)
		return REDISMODULE_ERR;

	if (RedisModule_CreateCommand(ctx, "boomfilter.exists", boomfilterExists, "readonly", 1, 0, 0) == REDISMODULE_ERR)
		return REDISMODULE_ERR;

	return REDISMODULE_OK;
}

*/
import "C"

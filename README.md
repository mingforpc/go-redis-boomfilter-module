# Redis bloomfilter module

[![License](https://img.shields.io/badge/license-Apache%202-green.svg)](https://www.apache.org/licenses/LICENSE-2.0)

一个使用CGO实现的Redis布隆过滤器module.

PS: `redismodule.h`基于Redis 4.0

## 编译命令

`go build -buildmode=c-shared -o boomfilter.so  ./boomfilter`

##　使用

### 加载

在Redis cli中`module load {PATH}/boomfilter.so`加载Module。

### 命令

#### 创建

`boomfilter.createboomfilter {key name} {hash func count} {filter size}`

* `key name`: 想要创建的过滤器名称
* `hash func count`: hash的次数
* `filter size`: 布隆过滤器的size

将会在Redis保存3个键值:

* `{key name}`: 保存布隆过滤器的数据
* `boomfilter.{key name}.hashseek.set`: 保存指定数量的用来计算hash的随机值
* `boomfilter.{key name}.total.size`: 保存布隆过滤器的大小

#### 删除boomfilter

`boomfilter.cleanboomfilter {key}`

删除整个布隆过滤器, 会把上述创建的3个键值对都删除

#### 添加

`boomfilter.add {key} {val} ...`

往指定的布隆过滤器中添加元素

#### 是否存在

`boomfilter.exists {key} {val}`

判断指定值是否存在于指定的过滤器中
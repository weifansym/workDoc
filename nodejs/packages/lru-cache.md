## lru-cache
lru-cache 用于在内存中管理缓存数据，并且支持LRU算法。可以让程序不依赖任何外部数据库实现缓存管理。
* LRU算法：尽量保留最近使用过的项
* 可指定缓存大小
* 可指定缓存项过期时间
```
const LRU = require('lru-cache');

const cache = LRU({
  max: 500,
  maxAge: 1000 * 60 * 60
});

cache.set('key','value');
cache.get('key'); // "value"

cache.reset(); // 清空
```
虽然，lru-cache 使用非常方便，但是lru-cache的缓存数据保存在当前进程内存内，这就决定了依赖lru-cache的项目是有状态的程序，这样就不能够
分布式部署多实例负载均衡，所以如果系统设计需要多实例运行，那么还是需要使用redis。

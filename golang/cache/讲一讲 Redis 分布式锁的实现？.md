### Redis 分布式锁最简单的实现

要实现分布式锁，确实需要使用具备互斥性的 Redis 操作。其中一种常用的方式是使用 `SETNX` 命令，该命令表示”SET if Not Exists”，即只有在 key 不存在时才设置其值，否则不进行任何操作。通过这种方式，两个客户端进程可以执行 `SETNX` 命令来实现互斥，从而达到分布式锁的目的。

下面是一个示例：

客户端 1 申请加锁，加锁成功：
```
SETNX lock_key 1
```

客户端 2 申请加锁，由于它处于较晚的时间，加锁失败：
```
SETNX lock_key 1
```
通过这种方式，您可以使用 Redis 的互斥性来实现简单的分布式锁机制。
![image](https://github.com/user-attachments/assets/5524e945-a7c6-4ad8-913c-2096663a41e7)

对于加锁成功的客户端，可以执行对共享资源的操作，比如修改 MySQL 的某一行数据或调用 API 请求。

操作完成后，需要及时释放锁，以便后续的请求能够访问共享资源。释放锁非常简单，只需使用 `DEL` 命令来删除相应的锁键（key）即可。

下面是释放锁的示例逻辑：
```
DEL lock_key
```

通过执行以上 `DEL` 命令，成功释放锁，以让后续的请求能够获得锁并执行操作共享资源的逻辑。

这样，通过使用 `SETNX` 命令进行加锁，然后使用 `DEL` 命令释放锁，您就可以实现基本的分布式锁机制。
![image](https://github.com/user-attachments/assets/a358dabe-f59e-4ce4-ae87-8e4abff09672)


但是，它存在一个很大的问题，当客户端 1 拿到锁后，如果发生下面的场景，就会造成「死锁」：

1、程序处理业务逻辑异常，没有及时释放锁。

2、进程崩溃或意外停止，无法释放锁。

在这种情况下，客户端将永远占用该锁，其他客户端将无法获取该锁。如何解决这个问题呢？

### 如何避免死锁？

当考虑在申请锁时为其设置一个「租期」时，可以在 Redis 中通过设置「过期时间」来实现。假设我们假设操作共享资源的时间不会超过 10 秒，在加锁时，可以给该 key 设置一个 10 秒的过期时间即可。这样做可以确保在申请锁后的一段时间内，如果锁的持有者在该时间内没有更新锁的过期时间，锁将会自动过期，从而防止锁被永久占用

```
SETNX lock 1 // 加锁
EXPIRE lock 10 // 10s后自动过期
```
![image](https://github.com/user-attachments/assets/79776646-331e-4b0f-82a9-0c79510fd6fa)

这样一来，无论客户端是否异常，这个锁都可以在 10s 后被「自动释放」，其它客户端依旧可以拿到锁。

但现在还是有问题：

当前的操作是将加锁和设置过期时间作为两个独立的命令执行，存在一个问题，即可能只执行了第一条命令而第二条命令却未能及时执行，从而导致问题。例如：

* SETNX 命令执行成功后，由于网络问题导致 EXPIRE 命令执行失败。
* SETNX 命令执行成功后，Redis 异常宕机，导致 EXPIRE 命令没有机会执行。
* SETNX 命令执行成功后，客户端异常崩溃，同样导致 EXPIRE 命令没有机会执行。
    
总之，这两条命令不能保证是原子操作（一起成功），就有潜在的风险导致过期时间设置失败，依旧发生「死锁」问题。

幸运的是，在 Redis 2.6.12 版本之后，Redis 扩展了 SET 命令的参数。用这一条命令就可以了：
```
SET lock 1 EX 10 NX
```
![image](https://github.com/user-attachments/assets/34a4484b-c281-474b-96ac-bb4a17a3ef11)

### 锁被别人释放怎么办？

上面的命令执行时，每个客户端在释放锁时，并没有进行严格的验证，存在释放别人锁的潜在风险。为了解决这个问题，可以在加锁时为每个客户端设置一个唯一的标识符（unique identifier），并在解锁时对比标识符来验证是否有权释放锁。

例如，可以是自己的线程 ID，也可以是一个 UUID（随机且唯一），这里我们以 UUID 举例：
```
SET lock $uuid EX 20 NX
```

之后，在释放锁时，要先判断这把锁是否还归自己持有，伪代码可以这么写：

```php
if redis.get("lock") == $uuid:
   redis.del("lock")
```

这里释放锁使用的是 GET + DEL 两条命令，这时，又会遇到我们前面讲的原子性问题了。这里可以使用 lua 脚本来解决。

安全释放锁的 Lua 脚本如下：
```
if redis.call("GET",KEYS[1]) == ARGV[1]
then
    return redis.call("DEL",KEYS[1])
else
    return 0
end
```

好了，这样一路优化，整个的加锁、解锁的流程就更「严谨」了。

这里我们先小结一下，基于 Redis 实现的分布式锁，一个严谨的的流程如下：

1、加锁

```
SET lock_key $unique_id EX $expire_time NX
```

2、操作共享资源

3、释放锁：Lua 脚本，先 GET 判断锁是否归属自己，再 DEL 释放锁

### go 代码实现分布式锁

```
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/google/uuid"
)

const (
    LockTime         = 5 * time.Second
    RS_DISTLOCK_NS   = "tdln:"
    RELEASE_LOCK_LUA = `
        if redis.call('get',KEYS[1])==ARGV[1] then
            return redis.call('del', KEYS[1])
        else
            return 0
        end
    `
)

type RedisDistLock struct {
    id          string
    lockName    string
    redisClient *redis.Client
    m           sync.Mutex
}

func NewRedisDistLock(redisClient *redis.Client, lockName string) *RedisDistLock {
    return &RedisDistLock{
        lockName:    lockName,
        redisClient: redisClient,
    }
}

func (this *RedisDistLock) Lock() {
    for !this.TryLock() {
        time.Sleep(100 * time.Millisecond)
    }
}

func (this *RedisDistLock) TryLock() bool {
    if this.id != "" {
        // 处于加锁中
        return false
    }
    this.m.Lock()
    defer this.m.Unlock()
    if this.id != "" {
        // 处于加锁中
        return false
    }
    ctx := context.Background()
    id := uuid.New().String()
    reply := this.redisClient.SetNX(ctx, RS_DISTLOCK_NS+this.lockName, id, LockTime)
    if reply.Err() == nil && reply.Val() {
        this.id = id
        return true
    }

    return false
}

func (this *RedisDistLock) Unlock() {
    if this.id == "" {
        // 未加锁
        panic("解锁失败，因为未加锁")
    }
    this.m.Lock()
    defer this.m.Unlock()
    if this.id == "" {
        // 未加锁
        panic("解锁失败，因为未加锁")
    }
    ctx := context.Background()
    reply := this.redisClient.Eval(ctx, RELEASE_LOCK_LUA, []string{RS_DISTLOCK_NS + this.lockName}, this.id)
    if reply.Err() != nil {
        panic("释放锁失败！")
    } else {
        this.id = ""
    }
}

func main() {

    client := redis.NewClient(&redis.Options{
        Addr: "172.16.11.111:64495",
    })
    const LOCKNAME = "百家号：福大大架构师每日一题"

    lock := NewRedisDistLock(client, LOCKNAME)

    lock.Lock()
    fmt.Println("加锁main")
    ch := make(chan struct{})
    go func() {
        lock := NewRedisDistLock(client, LOCKNAME)
        lock.Lock()
        fmt.Println("加锁go程")
        lock.Unlock()
        fmt.Println("解锁go程")
        ch <- struct{}{}
    }()
    time.Sleep(time.Second * 2)
    lock.Unlock()
    fmt.Println("解锁main")
    <-ch
}

```

![image](https://github.com/user-attachments/assets/6e810569-e17a-46b7-9345-6b6c02b296c5)

### 锁过期时间不好评估怎么办？

![image](https://github.com/user-attachments/assets/a7335bfd-a265-4a9c-b9d0-c11ff570b4fa)

看上面这张图，加入 key 的失效时间是 10s，但是客户端 C 在拿到分布式锁之后，然后业务逻辑执行超过 10s，那么问题来了，在客户端 C 释放锁之前，其实这把锁已经失效了，那么客户端 A 和客户端 B 都可以去拿锁，这样就已经失去了分布式锁的功能了！！！

比较简单的妥协方案是，尽量「冗余」过期时间，降低锁提前过期的概率，但是这个并不能完美解决问题，那怎么办呢？

### 分布式锁加入看门狗

在加锁过程中，可以设置一个过期时间，并启动一个守护线程（也称为「看门狗」线程），定时检测锁的剩余有效时间。如果锁即将过期，但共享资源操作尚未完成，守护线程可以自动对锁进行续期，重新设置过期时间。

为什么要使用守护线程：

![image](https://github.com/user-attachments/assets/3964bbf4-ab1f-4fa3-8a80-723e637fbfb8)


### go 中的红锁

```go
package main

import (
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/go-redsync/redsync/v4"
    "github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func main() {
    client := redis.NewClient(&redis.Options{
        Addr:     "172.16.11.111:64495",
        Password: "", // 如果有密码，请提供密码
        DB:       0,  // 如果使用不同的数据库，请修改为准确的数据库编号
    })

    pool := goredis.NewPool(client)

    const LOCKNAME = "百家号：福大大架构师每日一题"

    redsync := redsync.New(pool)

    mutex := redsync.NewMutex(LOCKNAME)

    if err := mutex.Lock(); err != nil {
        fmt.Println("加锁失败:", err)
        return
    }

    fmt.Println("加锁main")

    ch := make(chan struct{})

    go func() {
        mutex := redsync.NewMutex(LOCKNAME)

        if err := mutex.Lock(); err != nil {
            fmt.Println("加锁失败:", err)
            return
        }

        fmt.Println("加锁go程")
        mutex.Unlock()
        fmt.Println("解锁go程")

        ch <- struct{}{}
    }()

    time.Sleep(time.Second * 2)
    mutex.Unlock()
    fmt.Println("解锁main")

    <-ch
}
```

参考：
1: [讲一讲 Redis 分布式锁的实现](https://learnku.com/articles/78132)

2: [Go中使用Redis分布式锁的最佳实践](https://fantasticbin.com/archives/32/)

3: [用 Go + Redis 实现分布式锁](https://www.cnblogs.com/kevinwan/p/15688489.html)

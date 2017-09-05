## jsonwebtoken

 这个模块是对 [JSON Web Tokens](https://tools.ietf.org/html/rfc7519)的扩展

###   安装

```
$ npm install jsonwebtoken
```

###  使用

#### jwt.sign(payload, secretOrPrivateKey, [options, callback])

(异步)：如果提供一个回调函数，回调函数将会调用带有**err**或**JWT**.

(同步)：返回字符类型的JsonWebToken

payload：可以是一个对象，buffer或者是字符串，注意：只有当payload是对象的时候才可以设置exp属性。

secretOrPrivateKey： 可以是一个字符，buffer以及对象，包含HMAC算法密码，RSA 和ECDSA编码的私钥，如果使用一个带有密码的私钥对象，例如：{ key, passphrase }，在这个例子中确保你传递了algorithm选项。

options：

​    algorithm：(默认：HS256)

​    expiresIn：以秒级表示，或者是描述一个时间范围的字符串：例如：60，"2 days"，"10h"，"7d"

​    notBefore：以秒级表示，或者是描述一个时间范围的字符串：例如：60，"2 days"，"10h"，"7d"

​    audience：

​    jwtid：

​    subject：

​    noTimestamp：

​    header：

​    keyid：

如果payload不是buffer或者字符串，会强制使用`JSON.stringify`来转换。

`expiresIn`, `notBefore`, `audience`, `subject`, `issuer`都是没有默认值的，这些声明可以在payload上通过`exp`, `nbf`, `aud`, `sub` 和`iss`来代表，但是你不能在两个地方同时设置。

注意：exp, nbf, iat都是**NumericDate**类型。

可以通过option.header来设置header。

创建的jwts默认将会包含一个**iat** 声明，除非你指定了noTimestamp，如果iat插入到了payload中，他将会代替真的timestamp，用在估算像exp这个在options.expiresIn中指定的时间段。

举例：

```
// sign with default (HMAC SHA256)
var jwt = require('jsonwebtoken');
var token = jwt.sign({ foo: 'bar' }, 'shhhhh');
//backdate a jwt 30 seconds
var older_token = jwt.sign({ foo: 'bar', iat: Math.floor(Date.now() / 1000) - 30 }, 'shhhhh');

// sign with RSA SHA256
var cert = fs.readFileSync('private.key');  // get private key
var token = jwt.sign({ foo: 'bar' }, cert, { algorithm: 'RS256'});

// sign asynchronously
jwt.sign({ foo: 'bar' }, cert, { algorithm: 'RS256' }, function(err, token) {
  console.log(token);
});
```

#### Token Expiration (exp claim)

标准的JWT在处理过期时间的时候使用exp声明。过期时间使用**NumericDate**代表。

> A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds. This is equivalent to the IEEE Std 1003.1, 2013 Edition [POSIX.1] definition "Seconds Since the Epoch", in which each day is accounted for by exactly 86400 seconds, other than that non-integer values can be represented. See RFC 3339 [RFC3339] for details regarding date/times in general and UTC in particular.

表明exp字段应该包含秒级的数值。

使用过期时间为一小时为一个token签名：

```
jwt.sign({
  exp: Math.floor(Date.now() / 1000) + (60 * 60),
  data: 'foobar'
}, 'secret');
```

另一种生成token的方式：

```
jwt.sign({
  data: 'foobar'
}, 'secret', { expiresIn: 60 * 60 });

//or even better:

jwt.sign({
  data: 'foobar'
}, 'secret', { expiresIn: '1h' });
```


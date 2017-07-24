## joi
[joi](https://github.com/hapijs/joi)用于对象模型描述以及对象验证
## Introduction(介绍)
想象一下你是Facebook的开发人员，你想要用户使用真实的名字登录网站，而不想让用户在姓中输入**l337_p@nda**这种字符。你就要思考怎么定义来限制用户的输入，
并根据限制条件来验证用户的输入。

接下来我们将介绍一种解决方式，那就是来使用[joi](https://github.com/hapijs/joi),[joi](https://github.com/hapijs/joi)允许你创建javascript
对象模式，创建的对象模式中包含了验证规则信息，我们将会使用这些规则信息来效验数据。

## Example(实例)

     const Joi = require('joi');

     const schema = Joi.object().keys({
        username: Joi.string().alphanum().min(3).max(30).required(),
        password: Joi.string().regex(/^[a-zA-Z0-9]{3,30}$/),
        access_token: [Joi.string(), Joi.number()],
        birthyear: Joi.number().integer().min(1900).max(2013),
        email: Joi.string().email()
     }).with('username', 'birthyear').without('password', 'access_token');

    // Return result.
    const result = Joi.validate({ username: 'abc', birthyear: 1994 }, schema);
    // result.error === null -> valid

    // You can also pass a callback which will be called synchronously with the validation result.
    Joi.validate({ username: 'abc', birthyear: 1994 }, schema, function (err, value) { });  // err === null -> valid
  
 上面的模式定义了下面的限制条件：
 
* username:
  * 必填的字符
  * 只能包含字母数字组成的字符
  * 大于3个字符小于30个字符
  * 必须和**birthyear**组合使用

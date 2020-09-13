## 根据用户id生成唯一邀请码
### 错误思路
随机生成一个字符串，再将用户id拼接到字符串后面，但是这样id就太明显了，容易暴露，而且如果id很长的话，会导致邀请码很长，不利于用户使用。

所以可以将用户id插入到生成的字符串中，隔一个字符插入一个id的数字，这样id混合在字符串中，不容易暴露，但是长度问题并没有得到优化，于是把隔一个字符插入一个id的数字改为隔一个字符插入两个id的数字。然而长度好像并没有受到太大的影响。
正解
### 思考：一个10进制的数字短还是一个16进制的数字短？
肯定是16进制相对短一些，所以我们可以直接把用户id转成10+26=36进制的不就可以了吗？具体代码如下：
```
function createCode($user_id)
{
    static $source_string = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    $num = $user_id;
    $code = '';
    while($num)
    {
        $mod = $num % 36;  
        $num = ($num - $mod) / 36;
        $code = $source_string[$mod].$code;
    }
    return $code;
}
```
邀请码保证了唯一性，并且长度不会太长，用户id也能够根据邀请码反推出来，但是有一点不好的是，别人也可以根据邀请码去反推出user_id，因此，我们需要做一些优化。
### 优化
把0剔除，当做补位符号，比如小于六位的邀请码在高位补0，不要0、1、O、I这些容易混淆的字符这样36进制就变成了32进制，然后把字符串顺序打乱，这样，在不知道$source_string的情况下，是没办法解出正确的user_id的。

代码如下：
```
function createReferralCode($user_id) 
{
	static $source_string = 'E5FCDG3HQA4BNPJ2RSTUV67MWX89KLYZ';
	$num = $user_id;
	$code = '';
	while ( $num > 0) {
		$mod = $num % 32;
		$num = ($num - $mod) / 32;
		$code = $source_string[$mod].$code;
	}
	if(empty($code[5]))
		$code = str_pad($code,6,'0',STR_PAD_LEFT);
	return $code;

}
```
对应user_id的唯一邀请码就生成了，再附一个解码函数：
```
function deReferralCode($code) 
{
	static $source_string = 'E5FCDG3HQA4BNPJ2RSTUV67MWX89KLYZ';
	if (strrpos($code, '0') !== false)
		$code = substr($code, strrpos($code, '0')+1);
	$len = strlen($code);
	$code = strrev($code);
	$num = 0;
	for ($i=0; $i < $len; $i++) {
		$num += strpos($source_string, $code[$i]) * pow(32, $i);
	}
	return $num;

}
```


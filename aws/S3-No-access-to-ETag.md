在浏览器中使用AWS的SDK直接上传文件到S3时，需要在S3 Bucket上配置CORS才能成功上传，否则ajax请求会被浏览器拦截。

### 普通CORS访问配置
官方文档Cross-Origin Resource Sharing (CORS)中提供了开启CORS的范例，摘录如下:
```
<CORSConfiguration>
 <CORSRule>
   <AllowedOrigin>http://www.example1.com</AllowedOrigin>
   <AllowedMethod>PUT</AllowedMethod>
   <AllowedMethod>POST</AllowedMethod>
   <AllowedMethod>DELETE</AllowedMethod>
   <AllowedHeader>*</AllowedHeader>
 </CORSRule>
 <CORSRule>
   <AllowedOrigin>http://www.example2.com</AllowedOrigin>
   <AllowedMethod>PUT</AllowedMethod>
   <AllowedMethod>POST</AllowedMethod>
   <AllowedMethod>DELETE</AllowedMethod>
   <AllowedHeader>*</AllowedHeader>
 </CORSRule>
 <CORSRule>
   <AllowedOrigin>*</AllowedOrigin>
   <AllowedMethod>GET</AllowedMethod>
 </CORSRule>
</CORSConfiguration>
```
### 支持Multipart Upload的配置
当上传的文件比较大的时候，AWS的javascript的SDK会使用Multipart upload的方式来上传, 而Multipart upload的机制中是需要用到Header中的Etag的，因此需要在S3的CORS的rule中配置允许暴露ETag,
也即需要添加**<ExposeHeader>ETag</ExposeHeader>**

示例Rule:
```
<?xml version="1.0" encoding="UTF-8"?>
<CORSConfiguration xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
<CORSRule>
    <AllowedOrigin>http://localhost:3000</AllowedOrigin>
    <AllowedMethod>PUT</AllowedMethod>
    <AllowedMethod>POST</AllowedMethod>
    <ExposeHeader>ETag</ExposeHeader>
    <AllowedHeader>*</AllowedHeader>
</CORSRule>
</CORSConfiguration>
```
S3 Multipart Upload的原理在官方博客Amazon S3: Multipart Upload中有相关的说明。

如下是一个Parts Uploaded的PUT请求的response的示例。PUT成功上传分片到S3后，S3返回的本次请求的Etag为0a2f92d61cdc4682ba52adb9e077991f

参考：https://www.cnblogs.com/duhuo/p/14828021.html


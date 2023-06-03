## 文件上传
在NOS中用户的基本操作单元是对象，亦可以理解为文件，S3 GO SDK提供了丰富的上传接口，可以通过以下的方式上传文件:
* 流式上传
* 本地文件上传
* 大对象上传

字符串上传、本地文件上传最大为100M，大对象上传对文件大小没有限制
### 流式上传
通过PutObject方法上传对象，该方法只支持小于100M的对象。示例代码如下：
```
func putObjectByContent(s3Cleint *s3.S3,bucketName ,objectName string,content io.ReadSeeker){
    putObjectInput := &s3.PutObjectInput{
        Bucket:aws.String(bucketName),
        Key:aws.String(objectName),
        Body:content,
    }
    _,err := s3Cleint.PutObject(putObjectInput)

    if err != nil {
        fmt.Println("putObject : ",err.Error())
    } else {
        fmt.Println("upload file ok")
    }
}

//使用示例
content := "hello world"
putObjectByContent(getS3Client(),SrcBucket,"main",bytes.NewReader([]byte(content))
```
> 上传的字符串内容不超过100M
### 本地文件上传
通过PutObjectByFile方法上传本地文件，该方法只支持小于100M的文件。示例代码如下：
```
//content传入文件content的ReadSeeker即可
func putObjectByFile(s3Cleint *s3.S3,bucketName ,objectName string,content io.ReadSeeker){
    putObjectInput := &s3.PutObjectInput{
        Bucket:aws.String(bucketName),
        Key:aws.String(objectName),
        Body:content,
    }
    _,err := s3Cleint.PutObject(putObjectInput)

    if err != nil {
        fmt.Println("putObject : ",err.Error())
    } else {
        fmt.Println("upload file ok")
    }
}
//使用示例
bts,err := ioutil.ReadFile("src/objecttest/main")
if err != nil {
    fmt.Println("read file error : ",err.Error())
    return;
}
putObjectByFile(getS3Client(),bucketName,objectName,bytes.NewReader(bts))
```
> 上传的文件内容不超过100M

### 分片上传
除了通过putObject接口上传文件到NOS之外，NOS还提供了另外一种上传模式-分片上传,用户可以在如下应用场景内（但不限于此），使用分片上传模式，如：
* 需支持断点上传
* 上传超过100M的文件
* 网络条件较差，经常和NOS服务器断开连接
* 上传前无法确定文件大小

####  初始化分块
通过InitMultiUpload方法实现分块上传的初始化。示例代码如下：
```
func InitMultipartUpload(s3Client *s3.S3,bucketName,objectName string)(uploadId string,err error){
    resp,err := s3Client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{Bucket:aws.String(bucketName),
    Key:aws.String(objectName)})
    if err != nil {
        fmt.Println(err.Error())
        return
    } else {
        uploadId = *resp.UploadId
        fmt.Println("uploadId : " , uploadId)
        return
    }
}
//使用示例
uploadId,err := InitMultipartUpload(getS3Client(),bucketName,objectName)
```
### 分块上传
通过UploadPart方法实现分块上传。示例代码如下：
```
partSize := 2 * 1024 * 1024//  16KB =< size >= 100MB
var partNum int64 = 0
buffer := make([]byte,partSize)
var parts []*s3.CompletedPart
for ;; {
    partNum++
    readLen,err := file.Read(buffer)
    if err != nil || readLen == 0 {
        break
    }
    uploadPartResp,err := UploadPart(s3Client,bucketName,objectName,uploadId,partNum,bytes.NewReader(buffer[0:readLen]))
    if err != nil {
        fmt.Println("UploadPart : " + err.Error())
        break;
    }
    part := s3.CompletedPart{ETag:uploadPartResp.ETag,PartNumber:aws.Int64(partNum)}
    fmt.Println("partETage : " , part.ETag , " , partNum : " , part.PartNumber)
    parts = append(parts,&part)
```
#### 分块终止上传
通过AbortMultiUpload方法终止分块上传。示例代码如下：
```
s3Client.AbortMultipartUpload(&s3.AbortMultipartUploadInput{Bucket:aws.String(bucketName),
            Key:aws.String(objectName),UploadId:aws.String(uploadId)})
```
#### 完成分块上传
通过CompleteMultiUpload方法完成分块上传。示例代码如下：
```
comp := s3.CompleteMultipartUploadInput{Bucket:aws.String(bucketName),Key:aws.String(objectName),UploadId:aws.String(uploadId),
    MultipartUpload:&s3.CompletedMultipartUpload{Parts:parts}}
    _,err  = s3Client.CompleteMultipartUpload(&comp)
if err != nil {
    fmt.Println(err.Error())
}
```
#### 列出所有上传的分块
通过ListMultiUploads方法罗列出所有执行中的Multipart Upload事件，即已经被初始化的Multipart Upload但是未被Complete或者Abort的Multipart Upload事件。可以设置的参数为：
|参数|作用|
|--|--|
|KeyMarker|指定某一uploads key，只有大于该key-marker的才会被列出|
|MaxUploads|最多返回max-uploads条记录，取值范围[0-1000]，默认1000|

示例代码如下：
```
resp,err := s3Client.ListMultipartUploads(&s3.ListMultipartUploadsInput{Bucket:aws.String(SrcBucket),
    MaxUploads:aws.Int64(10)})
if err != nil {
    fmt.Println(err.Error())
    return;
}
```
#### 完整的分块上传
```
func MultipartUpload(s3Client *s3.S3,bucketName,objectName,filePath string){
    uploadId,err := InitMultipartUpload(s3Client,bucketName,objectName)
    if err != nil {
        fmt.Println("InitMultipartUpload : " + err.Error())
        return
    }
    file,err := os.Open(filePath)
    if err != nil {
        fmt.Println("open file : " + err.Error())
        return
    }
    partSize := 2 * 1024 * 1024
    var partNum int64 = 0
    buffer := make([]byte,partSize)
    var parts []*s3.CompletedPart
    for ;; {
        partNum++
        readLen,err := file.Read(buffer)
        if err != nil || readLen == 0 {
            break
        }
        uploadPartResp,err := UploadPart(s3Client,bucketName,objectName,uploadId,partNum,bytes.NewReader(buffer[0:readLen]))
        if err != nil {
            fmt.Println("UploadPart : " + err.Error())
            break;
        }
        part := s3.CompletedPart{ETag:uploadPartResp.ETag,PartNumber:aws.Int64(partNum)}
        fmt.Println("partETage : " , part.ETag , " , partNum : " , part.PartNumber)
        parts = append(parts,&part)
    }


    comp := s3.CompleteMultipartUploadInput{Bucket:aws.String(bucketName),Key:aws.String(objectName),UploadId:aws.String(uploadId),
    MultipartUpload:&s3.CompletedMultipartUpload{Parts:parts}}
    _,err  = s3Client.CompleteMultipartUpload(&comp)
    if err != nil {
        fmt.Println("CompleteMultipartUpload : " + err.Error())
        _,err := s3Client.AbortMultipartUpload(&s3.AbortMultipartUploadInput{Bucket:aws.String(bucketName),
        Key:aws.String(objectName),UploadId:aws.String(uploadId)})
        if err != nil {
            fmt.Println("AbortMultipartUpload failed ")
            return;
        }
        return
    } else {
        fmt.Println("CompleteMultipartUpload Ok")
    }
}
```
转自：http://public-cloud-doc.nos-eastchina1.126.net/s3golangsdk/uploadobject.html


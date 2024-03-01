今天考虑一个mcu混合的实现，也就是接收多路过来的rtp流，然后转发出去一路的rtmp流，使用ffmpeg测试做的记录，刚开始一直通过ffmpeg推送的文件流不能满足要求，还是对参数配置不熟悉；

### 0、ffmpeg 命令格式：
```
$ ffmpeg \

-y \ # 全局参数

-c:a libfdk_aac -c:v libx264 \ # 输入文件参数

-i input.mp4 \ # 输入文件

-c:v libvpx-vp9 -c:a libvorbis \ # 输出文件参数

output.webm # 输出文件
```


#### 下列为较常使用的参数：
-i——设置输入文件名。

-f——设置输出格式。

-y——若输出文件已存在时则覆盖文件。

-fs——超过指定的文件大小时则结束转换。

-t——指定输出文件的持续时间，以秒为单位。

-ss——从指定时间开始转换，以秒为单位。

-t从-ss时间开始转换（如-ss 00:00:01.00 -t 00:00:10.00即从00:00:01.00开始到00:00:11.00）。

-title——设置标题。

-timestamp——设置时间戳。

-vsync——增减Frame使影音同步。

-c——指定输出文件的编码。

-metadata——更改输出文件的元数据。

-help——查看帮助信息

##### 影像参数：

-b:v——设置影像流量，默认为200Kbit/秒。（单位请引用下方注意事项）

-r——设置帧率值，默认为25。

-s——设置画面的宽与高。

-aspect——设置画面的比例。

-vn——不处理影像，于仅针对声音做处理时使用。

-vcodec( -c:v )——设置影像影像编解码器，未设置时则使用与输入文件相同之编解码器。

##### 声音参数：

-b:a——设置每Channel（最近的SVN版为所有Channel的总合）的流量。（单位请引用下方注意事项）

-ar——设置采样率。

-ac——设置声音的Channel数。


-acodec ( -c:a ) ——设置声音编解码器，未设置时与影像相同，使用与输入文件相同之编解码器。

-an——不处理声音，于仅针对影像做处理时使用。

-vol——设置音量大小，256为标准音量。（要设置成两倍音量时则输入512，依此类推。）

-preset：指定输出的视频质量，会影响文件的生成速度，有以下几个可用的值 ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow。 


转自：
1. https://it3q.com/article/59
2. https://cloud.tencent.com/developer/article/1409507
3. https://it3q.com/tag/39
4. https://blog.csdn.net/zjun1001/article/details/107496158
5. https://cloud.tencent.com/developer/article/1409537?from_column=20421&from=20421

6. https://developer.aliyun.com/article/1365562

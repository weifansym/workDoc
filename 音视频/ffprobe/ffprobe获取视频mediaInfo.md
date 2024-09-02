ffmpeg官网：[ffmpeg.org](https://ffmpeg.org/)/ ffmpeg是用来处理音视频图像的最流行的工具，它的功能非常强大

## 首先介绍一下视频基础服务知识
视频中的视频和音频信息称之为stream，分为看下视频和音频有哪些属性
### 视频信息
* 常见的视频格式format：mp4、mov、flv、avi、hls等
* 常见的视频编码格式codec：codec_name 如H264、H265、vp8、vp8、av1
* 像素格式：pix_fmt 比如 yuv420p
* 宽、高、时长：width、height、duration
* 码率：bitrate= size*8/duration
* 帧率：fps

### 音频信息
* 音频编码codec：codec_name 如aac、pcm、mp3
* 音频声道数：channels
* 音频采样：sample_rate

### 获取视频format
```
ffprobe -loglevel error -print_format json -show_format input.mp4
{
    "format": {
        "filename": "ounewse111.mp4",
        "nb_streams": 2,
        "nb_programs": 0,
        "format_name": "mov,mp4,m4a,3gp,3g2,mj2",
        "format_long_name": "QuickTime / MOV",
        "start_time": "0.000000",
        "duration": "9202.370000",
        "size": "528909162",
        "bit_rate": "459802",
        "probe_score": 100,
        "tags": {
            "major_brand": "isom",
            "minor_version": "512",
            "compatible_brands": "isomiso2avc1mp41",
            "encoder": "Lavf58.29.100"
        }
    }
}
```
输出结果格式为json，读取到文本后可以转换成json进行解析

### 获取视频宽高，读取stream
```
ffprobe -loglevel error -print_format json -show_format -show_streams input.mp4
{
    "streams": [
        {
            "index": 0,
            "codec_name": "h264",
            "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
            "profile": "High",
            "codec_type": "video",
            "codec_time_base": "1/60",
            "codec_tag_string": "avc1",
            "codec_tag": "0x31637661",
            "width": 854,
            "height": 480,
            "coded_width": 864,
            "coded_height": 480,
            "has_b_frames": 2,
            "sample_aspect_ratio": "1280:1281",
            "display_aspect_ratio": "16:9",
            "pix_fmt": "yuv420p",
            "level": 31,
            "chroma_location": "left",
            "refs": 1,
            "is_avc": "true",
            "nal_length_size": "4",
            "r_frame_rate": "30/1",
            "avg_frame_rate": "30/1",
            "time_base": "1/90000",
            "start_pts": 5940,
            "start_time": "0.066000",
            "duration_ts": 828201000,
            "duration": "9202.233333",
            "bit_rate": "385468",
            "bits_per_raw_sample": "8",
            "nb_frames": "276067",
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "language": "und",
                "handler_name": "VideoHandler"
            }
        },
        {
            "index": 1,
            "codec_name": "aac",
            "codec_long_name": "AAC (Advanced Audio Coding)",
            "profile": "LC",
            "codec_type": "audio",
            "codec_time_base": "1/44100",
            "codec_tag_string": "mp4a",
            "codec_tag": "0x6134706d",
            "sample_fmt": "fltp",
            "sample_rate": "44100",
            "channels": 2,
            "channel_layout": "stereo",
            "bits_per_sample": 0,
            "r_frame_rate": "0/0",
            "avg_frame_rate": "0/0",
            "time_base": "1/44100",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 405824512,
            "duration": "9202.369887",
            "bit_rate": "65325",
            "max_bit_rate": "65325",
            "nb_frames": "396313",
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "language": "und",
                "handler_name": "SoundHandler"
            }
        }
    ],
    "format": {
        "filename": "ounewse111.mp4",
        "nb_streams": 2,
        "nb_programs": 0,
        "format_name": "mov,mp4,m4a,3gp,3g2,mj2",
        "format_long_name": "QuickTime / MOV",
        "start_time": "0.000000",
        "duration": "9202.370000",
        "size": "528909162",
        "bit_rate": "459802",
        "probe_score": 100,
        "tags": {
            "major_brand": "isom",
            "minor_version": "512",
            "compatible_brands": "isomiso2avc1mp41",
            "encoder": "Lavf58.29.100"
        }
    }
}
```
同时读取视频的format和stream，整个视频meta信息已经获取到了，index位流的序号，codec_type=video是视频信息，codec_type=audio是音频信息
### 说一下视频Stream的SAR/DAR/PAR
sample_aspect_ratio  sar像素的采样宽高比，这个值不是固定的1:1

display_aspect_ratio dar视频播放窗口宽高比

par 视频宽和高，像素数量之比，就是看到width和height

视频信息中的width和height，并不一定是视频真实播放的实际大小，需要根据sample_aspect_ratio进行调整

可以得出：DAR=PAR*SAR，dar才是视频播放的真实效果

### 说一下视频r_frame_rate和avg_frame_rate
avg_frame_rate是实际平均帧率，也就是总帧数/总时长

r_frame_rate是最大的帧率，也就是说这个视频中帧率最高的时间点就是这个值

nb_frames 则为总帧数

这两个值不一样的原因是动态帧率，帧率不是恒定的，可以根据视频画面的运动程度调整，运动程度大或镜头场景切换快的时候帧率就高，镜头场景切换慢的时候帧率就低，这样优点就是既能优化视频的大小，也能保留视频的画质

### 只读取一个stream
```
ffprobe -loglevel error -print_format json -select_streams 0 -show_streams input.mp4
```
select_streams可以选择其中一个，参数是stream的index，一般视频流为0，音频流为1
### 读取stream中的frame数量
```
ffprobe -loglevel error -print_format json -count_frames -show_streams   input.mp4
{
    "streams": [
        {
            "index": 0,
            "codec_name": "h264",
            "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
            "profile": "High",
            "codec_type": "video",
            "codec_time_base": "1/60",
            "codec_tag_string": "avc1",
            "codec_tag": "0x31637661",
            "width": 854,
            "height": 480,
            "coded_width": 864,
            "coded_height": 480,
            "has_b_frames": 2,
            "sample_aspect_ratio": "1280:1281",
            "display_aspect_ratio": "16:9",
            "pix_fmt": "yuv420p",
            "level": 31,
            "chroma_location": "left",
            "refs": 4,
            "is_avc": "true",
            "nal_length_size": "4",
            "r_frame_rate": "30/1",
            "avg_frame_rate": "30/1",
            "time_base": "1/15360",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 460800,
            "duration": "30.000000",
            "bit_rate": "842240",
            "bits_per_raw_sample": "8",
            "nb_frames": "900",
            "nb_read_frames": "900",
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "language": "und",
                "handler_name": "VideoHandler"
            }
        },
        {
            "index": 1,
            "codec_name": "aac",
            "codec_long_name": "AAC (Advanced Audio Coding)",
            "profile": "LC",
            "codec_type": "audio",
            "codec_time_base": "1/44100",
            "codec_tag_string": "mp4a",
            "codec_tag": "0x6134706d",
            "sample_fmt": "fltp",
            "sample_rate": "44100",
            "channels": 2,
            "channel_layout": "stereo",
            "bits_per_sample": 0,
            "r_frame_rate": "0/0",
            "avg_frame_rate": "0/0",
            "time_base": "1/44100",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 1323000,
            "duration": "30.000000",
            "bit_rate": "128538",
            "max_bit_rate": "128538",
            "nb_frames": "1293",
            "nb_read_frames": "1292",
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "language": "und",
                "handler_name": "SoundHandler"
            }
        }
    ]
}
```
nb_frames是meta信息中的frame数量，可能存在不可解码的

nb_read_frames是实际读取可解码的frame数量

下一章介绍stream下的packet和frame信息

转自：https://juejin.cn/post/7201363918363000890

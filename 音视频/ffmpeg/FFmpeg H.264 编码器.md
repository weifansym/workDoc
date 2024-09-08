## FFmpeg H.264 编码器
H.264 / MPEG-4 AVC 是目前最被广泛被应用的视讯编码格式，它的压缩效率比MPEG-2、MPEG-4、RV40 …等旧视讯编码格式还要高许多。

如果要输出H.264 / AVC 视讯编码，则需要libx264编码器，FFmpeg 的组态设定之中必须有--enable-libx264则才可以使用。

H.264 有多个版本，版本越高压缩比就越高（对应到profile）

### 码率控制（Rate control）
Rate control 是指控制每個畫格用了多少個位元的方法，这将影响档案大小和品质分布。这里就先不探讨Rate control 的种类有哪些。
通常会用下列两种模式： Constant Rate Factor (CRF)、Two-Pass ABR，若没有控制输出大小的需求则使用CRF 即可
#### Constant Rate Factor（CRF）
固定品质指标，而不在意大小；CRF 会得到最佳的bitrate 分配结果，缺点是你不能直接指定一个目标bitrate 或是档案大小。
设定值范围为0 – 51（可能会依照编译版本而不同），0 为最高品质，预设值为23，建议的范围在18 - 28，17 或18 接近视觉无损，但在技术上来说并不是无损。

#### Two-Pass ABR（Average Bitrate）
以平均值来说ABR 与CBR 相同，但ABR 允许在“适当” 的时候使用更好的bitrate 以取得更好的画质，但ABR 是用动态补偿的方式来计算画面的复杂度，也就是说在影片的起始及开头，或是当画面不如预期变化的时候则可能会影响到产出。
所以在two-pass 的时候，第一次转档时会先分析纪录影片的内容，第二次再依照分析的结果加以转码，但相对的就是非常耗时，要是没有输出大小的需求，则用CRF 即可。

### Preset
为选项集合，用来设定编码速度，相对的也会影响到压缩比；编码越快则压缩比越低。
一速度递减排序为：ultrafast、superfast、veryfast、faster、fast、medium（预设）、slow、slower、veryslow、placebo
### Tune 选项
也是为一个集合选项，可针对特定的影片类型微调参数设定值，以获得更好的品质或压缩率
* Film：film，用于高解析度电影，降低deblocking
* Animation：animation，动画＆卡通，使用deblocking 和较多的reference frames
* Grain：grain，保留旧影片的颗粒感
* Still Image：stillimage，幻灯片效果的影片
* Fast Decode：fastdecode，启用部分禁用的filter 来加速编码
* Zero Latency：zerolatency，快速编码＆低延迟
* PSNR：psnr，优化PSNR
* SSIM：ssim，优化SSIM

### Profile
复合选项，Profile 越好压缩比也越高，但编码复杂度相对提升；越好的Profile 相对启用的功能也越多，所以对播放硬体的需求也较高
<img width="341" alt="截屏2024-09-08 14 20 27" src="https://github.com/user-attachments/assets/38f3936e-4efb-454c-9303-d4409f819f76">

### Level
复合选项，与解码器的效能及容量对应，Level 越高需求越高
<img width="486" alt="截屏2024-09-08 14 21 40" src="https://github.com/user-attachments/assets/97981eb5-fd96-44f5-abd1-8a3383db8fb7">

### Decoded picture buffering
```
capacity = min(floor(MaxDpbMbs / (PicWidthInMbs * FrameHeightInMbs)), 16)
```
<img width="515" alt="截屏2024-09-08 14 22 47" src="https://github.com/user-attachments/assets/9f80b1a2-d3be-4d4d-b611-86ce4b2204db">

### 自定义
以上说的都是复合指令，若是有需求也可利用-x264-params 做个别设定
可参考：
[https://ffmpeg.org/ffmpeg-codecs.html#libx264_002c-libx264rgb](https://ffmpeg.org/ffmpeg-codecs.html#libx264_002c-libx264rgb)

* -I,keyint
设定i-frame 间隔，也就是GOP size

* -i,--min-keyint
最小GOP size

* -r,--ref
控制DPB（Decoded Picture Buffer） 大小，表示P-frame 参照多少个frame。预设值为3，范围：0 - 16。

* --scenecut
配置动态i-frame 的参考值；x264 的frame 中会纪录与参考frame 的差异，当判断差异过大时就是场景变更，此时就会插入一个i-frame 来做纪录，差异的依据可由设定scenecut 来做判别，当设为0 时等同no-scenecut。

* --no-scenecut
停用adaptive I-frame decision

### 参数优先顺序
```
preset -> tune -> "custom" -> profile -> level
```
转自：[FFmpeg H.264 编码器](https://blog.dexiang.me/zh-tw/technologies/ffmpeg-h264-options/)


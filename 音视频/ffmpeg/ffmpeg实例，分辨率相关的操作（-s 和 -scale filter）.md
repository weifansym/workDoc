## ffmpeg实例，分辨率相关的操作（-s 和 -scale filter）
### 调整视频分辨率-s
```
1、用-s参数设置视频分辨率，参数值wxh，w宽度单位是像素，h高度单位是像素
ffmpeg -i input_file -s 320x240 output_file

2、预定义的视频尺寸
	下面两条命令有相同效果
	ffmpeg -i input.avi -s 640x480 output.avi
	ffmpeg -i input.avi -s vga output.avi
```
### Scale filter调整分辨率
```
Scale filter的优点是可以使用一些额外的参数
	Scale=width:height[:interl={1|-1}]
	
下面两条命令有相同效果
	ffmpeg -i input.mpg -s 320x240 output.mp4 
	ffmpeg -i input.mpg -vf scale=320:240 output.mp4

对输入视频成比例缩放
改变为源视频一半大小
	ffmpeg -i input.mpg -vf scale=iw/2:ih/2 output.mp4
改变为原视频的90%大小：
	ffmpeg -i input.mpg -vf scale=iw*0.9:ih*0.9 output.mp4
```
### 在未知视频的分辨率时，保证调整的分辨率与源视频有相同的横纵比。
可能会有错误，不推荐使用，最好传入明确的缩放值
另外，scale只能接受偶数，否则height not divisible by 2
```
宽度固定400，高度成比例：
	ffmpeg -i input.avi -vf scale=400:-2

相反地，高度固定300，宽度成比例：
	ffmpeg -i input.avi -vf scale=-2:300
```
转自：[Scale filter调整分辨率](https://blog.csdn.net/yu540135101/article/details/84346505)


## 1\. Pillow

### 1.1 介绍
`Pillow` 是第三方开源的 `Python` 图像处理库，它支持多种图片格式，包括 `BMP、GIF、JPEG、PNG、TIFF` 等。`Pillow` 库包含了大量的图片处理函数和方法，可以进行图片的读取、显示、旋转、缩放、裁剪、转换等操作。在后续的深度学习中也不可或缺，是个比较重要的库。

> Pillow是torchvison使用的图像处理库,**而torchvison是PyTorch中专门用来处理图像的库**

[Github(截止当前 11.1K): https://github.com/python-pillow/Pillow](https://github.com/python-pillow/Pillow)

-   [官方文档: https://pillow.readthedocs.io/en/stable/index.html](https://pillow.readthedocs.io/en/stable/index.html)
-   [中文文档: https://www.osgeo.cn/pillow/](https://www.osgeo.cn/pillow/)

### 1.2 安装
```
$ pip install --upgrade Pillow
``` 

### 1.3 常用子模块

`Pillow`模块中有很多子模块,常用的子模块有：

-   `Image`: 该模块是`Pillow`中最重要的模块之一，用于处理图像文件。它提供了打开、保存、调整大小、旋转、裁剪、滤镜等功能，是图像处理的核心。
-   `ImageDraw`: 该模块提供了在图像上绘制各种形状（如线条、矩形、圆形）和文本的功能。可以使用不同的颜色和宽度绘制，创建自定义的标记或绘制图表等。
-   `ImageFont`: 该模块用于加载和使用`TrueType`字体文件，以便在图像上绘制文本时设置字体样式、大小和颜色。
-   `ImageFilter`: 该模块提供了各种滤镜效果，如模糊、锐化、边缘增强等。这些滤镜可以用于图像增强、特效处理和图像识别等应用。
-   `ImageEnhance`: 该模块用于调整图像的亮度、对比度、颜色饱和度等参数，使得图像更加清晰、明亮或具有特定的调色效果。
-   `ImageChops`: 该模块用于执行图像的逻辑和算术操作，如合并、比较、掩蔽等。可以进行图像合成、混合和提取等操作。
-   `ImageOps`: 该模块提供了各种图像处理操作，如镜像、翻转、自动对比度调整等。可以方便地进行图像变换和增强。
-   `ImageStat`: 该模块用于计算图像的统计信息，如均值、中位数、直方图等。可用于图像质量评估、颜色分析和特征提取等任务。

## 2\. 图像基本操作

### 2.1 读取图片

使用`Image.open()`来打开图像后，可以直接访问其属性信息，属性信息如下:

```
from PIL import Image

if __name__ == '__main__':
    imgPath = "img/test.jpg"
    # 读取图像
    img = Image.open(imgPath)
    # 打印图像属性
    print("读取对象img:", img)
    print("图像文件名:", img.filename)
    print("图像扩展名:", img.format)
    print("图像描述:", img.format_description)
    print("图像尺寸:", img.size)
    print("图像模式:", img.mode)
    print("图像宽度(像素):", img.width)
    print("图像高度(像素):", img.height)
    print("图象有关的数据的字典:", img.info)
    print("---------------------- 计算图片大小 --------------------------")
    # 图片大小(文件大小)
    with open(imgPath, "rb") as f:
        size = len(f.read())
        print("{}图片的大小(按照文件大小方式): {} byte，{} kb，{} Mb".format(img.filename, size, size / 1e3, size / 1e6))


# ---------------------------- 输出 ------------------------------  
读取对象img: <PIL.JpegImagePlugin.JpegImageFile image mode=RGB size=3840x2160 at 0x10F6F2910>
图像文件名: img/test.jpg
图像扩展名: JPEG
图像描述: JPEG (ISO 10918)
图像尺寸: (3840, 2160)
图像模式: RGB
图像宽度(像素): 3840
图像高度(像素): 2160
图象有关的数据的字典: {'jfif': 257, 'jfif_version': (1, 1), 'jfif_unit': 0, 'jfif_density': (72, 72), 'exif': b'Exif\x00\x00II*\x00\x08\x00\x00\x00\x02\x00\x1a\x01\n\x00\x01\x00\x00\x00&\x00\x00\x00\x1b\x01\n\x00\x01\x00\x00\x00.\x00\x00\x00\x00\x00\x00\x00H\x00\x00\x00\x01\x00\x00\x00H\x00\x00\x00\x01\x00\x00\x00', 'dpi': (72, 72), 'photoshop': {1028: b'\x1c\x02A\x00\nTopaz Labs'}}
---------------------- 计算图片大小 --------------------------
img/test.jpg图片的大小(按照文件大小方式): 11422554 byte，11422.554 kb，11.422554 Mb
```
### 2.2 另存图片

```python
from PIL import Image

if __name__ == '__main__':
    imgPath = "img/test.jpg"
    # 读取图像
    # ImageUse.openImg(imgPath)
    img = Image.open(imgPath)
    # 文件名 —> quality值
    imgNameMap = {
        "test_default": 75,
        "test_compress": 1,
        "test_quality": 100,
    }
    # 不修改文件格式
    for fileName, qualityVal in imgNameMap.items():
        img.save("./img/{}.jpg".format(fileName), quality=qualityVal)

    # 修改文件格式
    for fileName, qualityVal in imgNameMap.items():
        img.save("./img/{}.png".format(fileName), quality=qualityVal)
```

**1\. 查看另存的图像大小信息:**

```shell
# 使用ls -lh
$ ls -lh img/*
-rw-r--r--@ 1 liuqh  staff    11M Aug 29 10:04 img/test.jpg
-rw-r--r--  1 liuqh  staff   134K Sep  6 14:54 img/test_compress.jpg
-rw-r--r--  1 liuqh  staff    13M Sep  6 14:54 img/test_compress.png
-rw-r--r--  1 liuqh  staff   570K Sep  6 14:54 img/test_default.jpg
-rw-r--r--  1 liuqh  staff    13M Sep  6 14:54 img/test_default.png
-rw-r--r--  1 liuqh  staff   5.0M Sep  6 14:54 img/test_quality.jpg
-rw-r--r--  1 liuqh  staff    13M Sep  6 14:54 img/test_quality.png
``` 
**2.使用说明**

-   **另存为同类型格式时:** `save`方法会默认进行压缩(`quality=75`),可以通过调整`quality`(1-100)值的大小来调整压缩后的图片质量。
-   **另存为不同类型格式时:** 上面示例发现图片从`jpg`转到`png`后，图片大小比原来的大了`2M`,具体原因还未搜索到….

### 2.3 调整图片

`Pillow`模块还提供对图片进行大小调整、逆时针方向旋转、上下翻转、左右翻转等方法

```python
from PIL import Image

if __name__ == '__main__':
    img = Image.open("./img/a.jpg")
    # 更改图片大小为500 * 500
    img.resize((500, 500), Image.LANCZOS).save("./img/small.jpg")
    # ------------------------ 图片逆时针旋转 ----------------------
    # 图片逆时针旋转90
    img.rotate(90).save("./img/90.jpg")
    # 图片逆时针旋转120
    img.rotate(120).save("./img/120.jpg")
    # 图片逆时针旋转180
    img.rotate(180).save("./img/180.jpg")
    # ------------------------ 图片翻转 ----------------------
    # 图片左右翻转
    img.transpose(Image.FLIP_LEFT_RIGHT).save("./img/flip_left_right.jpg")
    # 图片上下翻转
    img.transpose(Image.FLIP_TOP_BOTTOM).save("./img/flip_top_bottom.jpg")
``` 
**`resize`函数中的第二个参数说明:**

**参数:resample** 一个可选的重采样过滤器,常用的常量含义如下:

-   `Image.NEAREST`: 从输入图像中选取一个最近的像素。忽略所有其他输入像素
-   `Image.BILINEAR`： 双线采样法
-   `Image.LANCZOS`: 力求输出最高质量像素的过滤器，只可用于 [`resize()`](https://www.osgeo.cn/pillow/reference/Image.html#PIL.Image.Image.resize) 和 [`thumbnail()`](https://www.osgeo.cn/pillow/reference/Image.html#PIL.Image.Image.thumbnail) 方法。

<img width="1968" height="842" alt="image" src="https://github.com/user-attachments/assets/0ad5e647-6b8b-4330-8a07-928340435883" />


### 2.4 编辑图片

```python
from PIL import Image

if __name__ == '__main__':
    # 打开图片
    img = Image.open("./img/spider_man.jpg")
    # 复制原对象
    imgCopy = img.copy()
    # 创建缩略图
    size = (500, 500)
    img.thumbnail(size)
    img.save("./img/spider_man_thumb.jpg")
    # 裁剪左上角英文
    cropImg = imgCopy.crop((0, 0, 500, 200))
    cropImg.save("./img/crop.jpg")
    # 裁剪蜘蛛侠头
    cropHeadImg = imgCopy.crop((2100, 200, 2900, 1000))
    cropHeadImg.save("./img/crop_head.jpg")

    # 图像合成-一次
    # cropHeadImg.show()
    onePasteImg = imgCopy.copy()
    onePasteImg.paste(cropHeadImg, (1300, 0))
    onePasteImg.save("./img/paste_once.jpg")
    # onePasteImg.show()

    # 图像合成-铺满
    imgPaste = imgCopy.copy()
    for x in range(0, 5000, 800):
        for y in range(1000, 3400, 800):
            imgPaste.paste(cropHeadImg, (x, y))
    imgPaste.save("./img/paste_spread_out.jpg")
```
> 图像裁剪函数`crop`,接受的元组四个数字代表含义如下:`(左上角的x轴坐标,左上角的y轴坐标,左上角的y轴坐标,左上角的y轴坐标)`

<img width="1956" height="1292" alt="image" src="https://github.com/user-attachments/assets/3ceeb7a1-5781-44ff-8b02-4f84e709b9c6" />


## 3\. 图像绘制

### 3.1 绘制图形

> 图像绘制的功能基本都在 `ImageDraw`包内

```python
from PIL import ImageDraw
from PIL import Image

if __name__ == '__main__':
    # 创建一个图像
    img = Image.new("RGBA", (1000, 1000), "Cyan")
    # 获取绘制对象
    draw = ImageDraw.Draw(img)
    # 画点 黑色
    for x in range(0, 200, 5):
        for y in range(0, 200, 5):
            draw.point([(x, y)], fill="black")
    # 划线(十字架)
    draw.line([(500, 0), (500, 1000)], fill="red")
    draw.line([(0, 500), (1000, 500)], fill="blue")

    # 画圆
    draw.ellipse((180, 180, 480, 480), fill="green")
    # 画椭圆
    draw.ellipse((710, 100, 900, 410), fill="red")
    # 画矩形(蓝色底层、红色边框线)
    draw.rectangle((100, 520, 400, 800), fill="blue", outline="black")
    # 画多变形
    draw.polygon([(650, 650), (700, 510), (930, 880), (508, 711), (666, 999)], fill="Purple")
    # img.show()
    img.save("./img/draw.png")
``` 
<img width="1768" height="1814" alt="image" src="https://github.com/user-attachments/assets/f96ba6a3-95f0-4bdf-bdf8-b14e7f3b2285" />


### 3.2 填充文字

```python
from PIL import ImageDraw, ImageFont

from PIL import Image

if __name__ == '__main__':
    img = Image.open("./img/spider_man.jpg")
    # 获取绘制对象
    draw = ImageDraw.Draw(img)
    # 设置字体,mac用户字体位置在:/Library/Fonts
    imgFont = ImageFont.truetype("/Library/Fonts/AdobeKaitiStd-Regular.otf", 150)
    # 写入中文
    draw.text((1500, 1000), text="蜘蛛侠", font=imgFont, fill="red")
    # 写入英文
    imgFont2 = ImageFont.truetype("/Library/Fonts/AdobeKaitiStd-Regular.otf", 80)
    draw.text((1500, 1200), text="Spider Man", font=imgFont2)
    img.thumbnail((1000, 600))
    # img.save("./img/t.jpg")
    img.show()
```
**draw.text() 部分参数说明:**

-   `xy`: 文本的坐标。
-   `text` : 要绘制的字符串。如果它包含任何换行符，文本将被传递给multiline\_text()。
-   `fill`: 用于文本的颜色。
-   `font`: 一个`ImageFont`实例。
-   `stroke_width`: 文本笔画的宽度。
-   `stroke_fill`: 用于文本描边的颜色。如果没有给出，将默认为填充参数。

<img width="1040" height="562" alt="image" src="https://github.com/user-attachments/assets/f87a9a79-6ff8-499b-8ec2-2e489f93e91c" />

## 4\. 图像滤镜

### 4.1 模糊

```python
from PIL import ImageFilter
from PIL import Image

if __name__ == '__main__':
    img = Image.open("./img/girl.jpg")
    # 盒子模糊，数值越大越模糊
    boxObscureImg = img.filter(ImageFilter.BoxBlur(15))
    boxObscureImg.save("./img/boxObscureImg.jpg")
    # 高斯模糊
    gaussianBlur = img.filter(ImageFilter.GaussianBlur(15))
    gaussianBlur.save("./img/gaussianBlur.jpg")
```


<img width="2454" height="656" alt="image" src="https://github.com/user-attachments/assets/ecd20bb8-a555-4e4c-82b6-1ed14e2864e7" />

### 4.2 轮廓和浮雕

```python
from PIL import ImageFilter
from PIL import Image

if __name__ == '__main__':
    img = Image.open("./img/girl.jpg")
    # 寻找图像轮廓信息
    contourImg = img.filter(ImageFilter.CONTOUR)
    contourImg.save("./img/contour.jpg")
    contourImg.show()
    # 寻图像的边界信息
    find_edgesImg = img.filter(ImageFilter.FIND_EDGES)
    find_edgesImg.save("./img/find_edges.jpg")
    find_edgesImg.show()
    # 浮雕滤波
    embossImg = img.filter(ImageFilter.EMBOSS)
    embossImg.save("./img/emboss.jpg")
    embossImg.show()
```
<img width="2444" height="624" alt="image" src="https://github.com/user-attachments/assets/8a5b43d7-ac05-479f-b962-18c01de111b3" />


### 4.3 增强

```python
from PIL import ImageFilter
from PIL import Image

if __name__ == '__main__':
    img = Image.open("./img/girl.jpg")
    img.show()
    # 深度边缘增强滤波
    edgeEnhanceImg = img.filter(ImageFilter.EDGE_ENHANCE_MORE)
    edgeEnhanceImg.save("./img/edge_enhance_more.jpg")
    edgeEnhanceImg.show()
```
增强之后，图片是比之前的更清晰，具体图片这里就不在展示.下面列举几个增强参数

| 参数常量 | 说明 |
|-----------------------------|-----------------|
| `ImageFilter.DETAIL` | 细节滤波，使得图像显示更加精细 |
| `ImageFilter.EDGE_ENHANCE` | 边界增强滤波 |
| `ImageFilter.EDGE_ENHANCE_MORE` | 深度边缘增强滤波 |
| `ImageFilter.SHARPEN` | 锐化滤波 |


文章链接: [http://liuqh.icu/2023/09/08/python/package/5.pil/](http://liuqh.icu/2023/09/08/python/package/5.pil/)

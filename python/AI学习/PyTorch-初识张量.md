# AI学习(3):PyTorch-初识张量 | 猿码记

## 1.介绍

> @注:下面是对`PyTorch`进行了简单的介绍，不喜欢可直接跳过。

### 1.1 什么是PyTorch

`PyTorch`是一个由`Facebook`人工智能研究团队开发的开源机器学习库，用于开发人工智能和深度学习的应用程序。`PyTorch`支持广泛的机器学习和深度学习算法，并基于强大的`GPU`加速计算库`CUDA`，提供了高效的张量计算（如数组计算）和深度神经网络功能。

**`PyTorch`的主要特性:**

-   **易用性**：`PyTorch`提供了一个类似于`NumPy`的编程环境，以及全面的深度学习功能，使得神经网络的构造和训练都变得非常直观。
-   **动态计算图**：`PyTorch`使用动态计算图，这意味着您可以在运行过程中更改图形。这在某些模型（例如循环神经网络或递归神经网络）中非常有用，这些模型的结构可能需要在运行时进行更改。
-   **`Python`支持**：`PyTorch`完全集成在`Python`中，可以与其他`Python`库（如`NumPy`和`Cython`）无缝地协作。
-   **分布式训练**：`PyTorch`支持在多个`GPU上`分布式的训练模型，可以有效地加速大数据集的模型训练过程。
-   **ONNX兼容性**：`PyTorch`支持`Open Neural Network Exchange（ONNX）`模型格式。这意味着您可以在不同的深度学习框架（例如`Caffe2、Microsoft Cognitive Toolkit、MXNet`等）之间轻松迁移模型。

### 1.2 PyTorch发展史

1.  **初始发布（2016年）：** `PyTorch`最初由`Facebook`的人工智能研究实验室（`Facebook AI Research`，简`称FAIR`）开发，并于2016年首次发布。初始版本主要以动态计算图为特点，这使得定义和修改模型变得非常灵活。
2.  **动态计算图（2016-2017年）：** `PyTorch`最初的设计采用动态计算图，这使得用户能够更自由地调试和修改模型。这种灵活性对研究人员和实践者来说是一个吸引点，尤其在处理变化的输入大小时更为方便。
3.  **静态计算图的引入（2017年）：** 随着`TensorFlow`等框架采用静态计算图的方式，`PyTorch`也在2017年引入了静态计算图的支持，这使得`PyTorch`更适用于一些需要性能优化的应用。
4.  **PyTorch 1.0（2018年）：** `PyTorch 1.0`的发布标志着一个重要的里程碑。它引入了`Eager Execution`（即动态计算图）和TorchScript（即静态计算图）的融合，使得用户可以在训练和部署中选择合适的计算图方式。
5.  **TorchServe和TorchElastic（2019年）：** 在2019年，`PyTorch`推出了`TorchServe`和`TorchElastic`，这是用于模型部署和分布式训练的工具，使得将PyTorch模型投入实际应用更为方便。
6.  **PyTorch 1.7和Beyond（2020年以后）：** 后续版本不断改进性能、增加新特性，并推动`PyTorch`在深度学习社区中的广泛采用。`PyTorch`继续保持开源性质，积极响应用户需求和社区贡献。
7.  **PyTorch 2.0(2023年3月)：** 推出新的编译器`torch.compile`。它将`PyTorch`的性能推向了新的高度，并开始将`PyTorch`的部分内容从`C++`中移回到`Python`中。据称，使用`torch.compile`对模型进行编译可以提升模型速度`30%`

## 2.安装环境

### 2.1 安装python3.10

为了保证`pytorch`运行环境的干净，这里单独为其创建一个新环境。

bash

```ruby
# 安装
$ conda create -n pytorch310 python=3.10
# 激活环境
$ conda activate pytorch310
# 查看版本
$ python -V
Python 3.10.13
```

### [](https://liuqh.icu/2024/01/26/ai/basic/3-zhang-liang-shi-yong/#2-2-%E5%AE%89%E8%A3%85%E4%BE%9D%E8%B5%96%E5%8C%85 "2.2 安装依赖包")2.2 安装依赖包

bash

```ruby
$ conda install numpy  matplotlib
```

-   **`numpy`：** 提供了强大的数组和矩阵操作，与 PyTorch 的张量操作兼容，常用于数据处理和转换。
-   **`matplotlib`:** 可视化训练过程中的损失曲线、模型输出、数据分布等;

这两个包的具体使用教程可查看文章:

-   [Python常用库(六):科学计算库Numpy-上篇:创建、访问、赋值 https://mp.weixin.qq.com/s/b0aPs1VMh0l0QM2D\_q1OHw](https://mp.weixin.qq.com/s/b0aPs1VMh0l0QM2D_q1OHw)
-   [Python库学习(七):Numpy-续篇一:结构数组 https://mp.weixin.qq.com/s/ThdIwvSaUFZEWks1D0RYzw](https://mp.weixin.qq.com/s/ThdIwvSaUFZEWks1D0RYzw)
-   [Python库学习(八):Numpy-续篇二:数组操作 https://mp.weixin.qq.com/s/5VXpfL-P8b0Li3wn5BKu4w](https://mp.weixin.qq.com/s/5VXpfL-P8b0Li3wn5BKu4w)
-   [Python库学习(九):Numpy-续篇三:数组运算 https://mp.weixin.qq.com/s/qtGHvB33-KewrUtIDU5JIw](https://mp.weixin.qq.com/s/qtGHvB33-KewrUtIDU5JIw)
-   [Python库学习(十):Matplotlib绘画库 https://mp.weixin.qq.com/s/Pb0kO6R3Q7ejX6x51y4TPw](https://mp.weixin.qq.com/s/Pb0kO6R3Q7ejX6x51y4TPw)

### [](https://liuqh.icu/2024/01/26/ai/basic/3-zhang-liang-shi-yong/#2-3-%E5%AE%89%E8%A3%85PyTorch "2.3 安装PyTorch")2.3 安装PyTorch

安装命令直接访问官网生成: [https://pytorch.org](https://pytorch.org/)

<img width="926" height="415" alt="image" src="https://github.com/user-attachments/assets/91de6dd1-80a2-4f01-8c33-7cc1ee456aa7" />


> @注意: 由于本人使用是Mac，没办法享受CUDA加速，后面在想办法体验~

bash

```ruby
# 运行安装
$ conda install pytorch::pytorch torchvision torchaudio -c pytorch
```

**验证安装结果:**

```python
import torch
if __name__ == '__main__':
    print("torch版本:", torch.__version__)

# torch版本: 2.1.2
```

### 3.1 核心模块

`PyTorch`的核心模块主要包括以下几个部分：

-   **`torch`：** 提供了张量（`tensor`）的基本操作，类似于 `NumPy` 数组。`PyTorch` 中的张量是深度学习模型的基本构建块。
-   **`torch.nn`：** 提供了构建神经网络模型所需的各种类和函数。包括神经网络的层、损失函数、优化器等。
-   **`torch.optim`：** 包含了各种优化算法，例如随机梯度下降 `(SGD)、Adam、RMSprop` 等，用于优化神经网络的参数。
-   **`torch.autograd`：** 实现了自动求导机制，允许用户定义的操作在反向传播过程中自动计算梯度。
-   **`torch.utils.data`：** 提供了用于加载和处理数据的工具，包括 `Dataset` 和 `DataLoader` 类，使得数据在训练时更容易进行批量处理。
-   **`torchvision`：** 提供了处理图像数据集的工具，包括常用的图像变换、数据集加载等。
-   **`torchtext`：** 用于处理文本数据的工具，包括加载文本数据集、文本变换等。
-   **`torch.nn.functional`：** 包含一些不具有内部状态的函数，这些函数通常在神经网络的中间层中使用，例如激活函数、池化操作等。
-   **`torch.distributed`：** 提供了分布式训练的工具，用于在多个 `GPU` 或多台机器上进行模型的训练。
-   **`torchaudio`：** 用于处理音频数据的工具，包括加载音频数据集、音频变换等。

### 3.2 PyTorch2.0新模块

在`PyTorch 2.0`中，引入了一些新的模块和功能：

-   **`TorchDynamo`:** 将`Python`代码`JIT`编译成`FX`图的新特性，可以提高模型训练速度。
-   **`AOTAutograd`:** 预编译自动求导函数的新特性，可以提高模型训练速度。
-   **`PrimTorch`:** 定义更小且更稳定的运算符集的新特性，可以提高模型训练速度。
-   **`TorchInductor`:** 为多个加速器和后端生成快速代码的新特性，可以提高模型训练速度。

## 4.版本介绍

> 上面我们安装的`PyTorch`版本的是`2.1.2`,后面学习也是基于这个版本；

### 4.1 Pytorch2.x Vs Pytorch1.x

以下是`PyTorch 2.0`和`PyTorch 1.x`之间的主要区别：

-   编译器支持： 在`PyTorch 2.0`中，已经支持了编译器模式，可以提高模型训练速度。这是`PyTorch 2.0`与`PyTorch 1.x`之间的一个主要区别。
-   API更新： 在`PyTorch 2.0`中，进行了一些API更新，以便更好地支持深度学习任务。这使得`PyTorch 2.0`与`PyTorch 1.x`之间的API使用有所不同。
-   新功能： 在`PyTorch 2.0`中，添加了一些新功能，如编译器模式、新的数据加载和预处理工具等。这使得`PyTorch 2.0`与`PyTorch 1.x`之间的功能有所不同。
-   性能提升： 在`PyTorch 2.0`中，实现了性能提升，如模型训练速度的提高。这使得`PyTorch 2.0`与`PyTorch 1.x`之间的性能有所不同。

### 4.2 PyTorch 2.x

> `PyTorch 2.0`在 2023.03发布，对之前的`1.x`版本是`100%`兼容。

从`PyTorch`版本发布历史信息中，可以看出`PyTorch`从`1.3`版本之后，后面版本直接就到了`2.0`;为什么会有这么大的跳跃呢？官方解释如下:

properties

```csharp
PyTorch 2.0 is what 1.14 would have been. We were releasing substantial new features that we believe change how you meaningfully use PyTorch, so we are calling it 2.0 instead.

// 译文
PyTorch 2.0是1.14的延续。我们发布了一些重大新功能，我们相信这些功能会改变您对PyTorch的实质性使用方式，因此我们将其称为2.0而不是1.14
```

其中最重要的新功能是：`torch.compile`，据官方描述，其可以大幅提高模型训练速度。而且使用特别简单，仅仅是一行代码:`model = torch.compile(model)`，下面是官方描述(汉字是软件译文):

<img width="1628" height="544" alt="image" src="https://github.com/user-attachments/assets/fb653205-3d24-49f9-b64d-f5e894ce5d7c" />

### [](https://liuqh.icu/2024/01/26/ai/basic/3-zhang-liang-shi-yong/#4-3-PyTorch2-0%E6%80%A7%E8%83%BD%E6%B5%8B%E8%AF%95 "4.3 PyTorch2.0性能测试")4.3 PyTorch2.0性能测试

为了验证`PyTorch2.0`带来的性能提升，官方从机器学习开源社区收集了163个模型，用于验证；

163个模型数据，具体来源如下:

-   46 models from [HuggingFace Transformers](https://github.com/huggingface/transformers)
-   61 models from [TIMM](https://github.com/rwightman/pytorch-image-models): a collection of state-of-the-art PyTorch image models by Ross Wightman
-   56 models from [TorchBench](https://github.com/pytorch/benchmark/): a curated set of popular code-bases from across github

除了使用`torch.compile`对上述模型进行编译，不改其他代码的前提下，测试性能如下:

<img width="1372" height="942" alt="image" src="https://github.com/user-attachments/assets/352e5c49-7479-480b-9aaa-c8da79563f75" />


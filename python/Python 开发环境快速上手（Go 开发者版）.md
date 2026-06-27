# Python 开发环境快速上手（Go/Java 开发者版）

> 本文总结了 Python 开发中最容易混淆的概念：Python
> 版本、pyenv、venv、pip、Conda、PyCharm、Docker 等。

## 1. Python 生态关系

``` text
                 Python生态

          pyenv（Python版本管理）
                    │
        Python3.9  Python3.11 Python3.12
                    │
        ┌───────────┴───────────┐
        │                       │
      venv                  Conda
（官方虚拟环境）         （完整环境管理器）
        │
        ▼
      pip install
        │
        ▼
 requirements.txt
```

### 对应 Node.js

  Node.js        Python
  -------------- --------
  nvm            pyenv
  npm            pip
  node_modules   venv
  nvm + npm      Conda

------------------------------------------------------------------------

## 2. 每个工具负责什么

### Python

运行程序的解释器。

``` bash
python3 --version
```

### pyenv

管理多个 Python 版本。

``` bash
pyenv install 3.9.21
pyenv local 3.9.21
```

### venv

官方虚拟环境，为每个项目隔离依赖。

``` bash
python -m venv .venv
source .venv/bin/activate
```

### pip

安装 Python 包。

``` bash
pip install -r requirements.txt
```

### Conda

同时管理 Python、虚拟环境以及部分系统依赖，AI 项目中常见。

------------------------------------------------------------------------

## 3. 为什么必须使用 venv

多个项目依赖不同版本库时，如果共用全局 site-packages 会互相覆盖。

因此推荐：

> 一个项目，一个 venv。

------------------------------------------------------------------------

## 4. PyCharm 自动创建 .venv

PyCharm 自动执行了：

``` bash
python -m venv .venv
```

`.venv` 就是项目独立运行环境。

不要提交到 Git：

``` gitignore
.venv/
```

------------------------------------------------------------------------

## 5. PyCharm、Terminal、Docker 的关系

三者互相独立：

-   PyCharm Interpreter：IDE 使用的解释器
-   Terminal：Shell PATH 决定使用哪个 Python
-   Docker：容器中的 Python

推荐保持一致：

``` text
Docker 3.9
    │
PyCharm 3.9
    │
Terminal 3.9
```

------------------------------------------------------------------------

## 6. requirements.txt

作用类似 package.json，用于记录依赖：

``` text
fastapi==0.104.1
numpy==1.25.2
requests==2.31.0
```

恢复环境：

``` bash
pip install -r requirements.txt
```

------------------------------------------------------------------------

## 7. 常用命令

``` bash
python --version
python3 --version
which python
which pip
pip list
echo $VIRTUAL_ENV
```

------------------------------------------------------------------------

## 8. 推荐开发流程

1.  查看 Dockerfile 或项目文档要求的 Python 版本。
2.  使用 pyenv 安装对应 Python。
3.  创建虚拟环境：

``` bash
python -m venv .venv
```

4.  激活：

``` bash
source .venv/bin/activate
```

5.  安装依赖：

``` bash
pip install -r requirements.txt
```

6.  PyCharm 使用 `.venv/bin/python` 作为 Project Interpreter。

------------------------------------------------------------------------

## 9. 常见坑

-   不要使用比生产环境更新很多的 Python。
-   切换 Python 大版本后，需要删除并重新创建 `.venv`。
-   `python: command not found` 往往是没有激活虚拟环境或系统只提供
    `python3`。
-   PyCharm Interpreter 与 Terminal Python 不是同一个概念。

------------------------------------------------------------------------

## 10. 最佳实践

``` text
Docker
FROM python:3.9
        │
        ▼
pyenv 安装 Python3.9
        │
        ▼
python -m venv .venv
        │
        ▼
source .venv/bin/activate
        │
        ▼
pip install -r requirements.txt
        │
        ▼
PyCharm 使用 .venv/bin/python
```

遵循以上流程，可以避免绝大多数 Python 环境问题。

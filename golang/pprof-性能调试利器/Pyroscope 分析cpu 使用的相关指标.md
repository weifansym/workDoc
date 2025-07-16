## Pyroscope 分析cpu 使用的相关指标
在 Pyroscope 中分析 Go 应用的 CPU 使用情况时，最关键的是理解采样方式和指标语义。Pyroscope 主要基于 pprof，使用采样分析调用栈，从而可视化 CPU 的使用情况。
### Pyroscope 分析 CPU 的核心指标
Pyroscope 默认在采样模式下收集 CPU Profile，每隔 10ms 进行一次采样（Go 中 runtime/pprof 默认 100Hz 采样频率）。在分析页面中，常见的指标有：
<img width="774" height="336" alt="截屏2025-07-16 18 29 00" src="https://github.com/user-attachments/assets/b12d05e5-a42c-45b9-a51b-ac3d5fec0ccc" />

### 哪个指标最能反映 CPU 使用情况？
#### 选：cpu（即 CPU 火焰图） → 查看函数的 self 值
> 因为它反映的是 函数“本体”自身消耗 CPU 的程度，不含调用链中子函数。

#### 其次：看 total，确定整个调用链的 CPU 占用
> 常用于追踪“谁最终触发了大量 CPU 使用”。

### 实际操作：如何看火焰图判断 CPU 使用情况
1. 打开 Pyroscope Web 控制台；
2. 选择 cpu.profile；
3. 查看最宽的函数栈（图上最宽的块）；
4. 鼠标悬停看 tooltip：
  * self: 420ms (38%) ← 这个函数自身耗费了 38% 的 CPU；
  * total: 700ms (65%) ← 加上它调用的函数，总共耗了 65%。
### 常见热点函数示例
| 函数                                          | 说明                  |
| ------------------------------------------- | ------------------- |
| `runtime.selectgo`                          | goroutine 选择阻塞/调度热点 |
| `syscall.Syscall` / `poll_runtime_pollWait` | 系统调用或 IO 阻塞         |
| `json.Marshal` / `encoding/json`            | 数据序列化开销大            |
| 自定义函数 `handler.ProcessXXX`                  | 应用代码 CPU 重占用        |

### 实用分析思路
| 问题                      | 分析方式                     |
| ----------------------- | ------------------------ |
| CPU 高，但看不出热点            | 看火焰图是否很多小函数堆叠（可能是高并发小请求） |
| 某函数宽，但没有子栈              | `self` 比较高，是热点函数本体，优化重点  |
| 某函数 `total` 高但 `self` 低 | 可能是中间调度函数，不是根源           |
| 火焰图平缓但宽                 | 程序整体忙，但无明显热点，需系统性优化      |

### 补充：结合系统资源看 Pyroscope 效果更好
| 工具                                   | 用法               |
| ------------------------------------ | ---------------- |
| `top` / `htop`                       | 看整体 CPU 使用率      |
| `go tool pprof -http=:8080`          | 本地分析 `.pprof` 文件 |
| `pyroscope server` + `pyroscope CLI` | 实时流式采集分析         |

### 总结
| 目标          | 推荐指标                            |
| ----------- | ------------------------------- |
| 判断谁占了最多 CPU | `cpu.profile` 火焰图，查看 **最宽的函数栈** |
| 判断函数自身消耗    | **`self`**                      |
| 判断调用链整体压力   | **`total`**                     |
| 判断采样强度      | `samples`                       |
| 函数占比排序      | `percent` 或火焰图顶部导航面板            |





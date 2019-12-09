## EventLoop中的事件分类
1. macrotask 与 microtask
一直有面试题以及技术帖提到Node中的事件循环中，有一个被称为macrotask的东西（另一个是microtask）。一直都知道有task和microtask的区分，
但不太清楚什么是macrotask。甚至在spec以及[谷歌V8官方的技术博客](https://v8.dev/blog/fast-async#tasks-vs.-microtasks)中也没有提到macrotask，
最多无非就是task（见下图）。于是就想要稍微查下。

参看：[EventLoop中的事件分类](https://xenojoshua.com/2019/02/event-loop-spec/)

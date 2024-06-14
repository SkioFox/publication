### note
#### 注入 Trace 的实现原理

> 自动注入 Trace 函数的实现，见ast.instrumenter

![img.png](img/img.png)

执行代码

```shell
make Makefile // 生成可执行文件
./instrument -w  examples/demo/demo.go
// 可能需要同步依赖和修改demo中的go mode引用为本地的
```
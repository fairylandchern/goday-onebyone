### version1
#### 功能
```
仅支持半包及粘包的解决方案，以及简单的自定义协议应该如何处理

protocol在解析数据时，长度未被正确返回，导致解析失败，但在解析成功的情况下，
貌似只是读出来一部分从服务端获取到的数据，另一部分数据读取失败了。
```
#### 改进计划
```text
打算在v2版本中，对数据的读写进行分离，正确支持其访问路径等信息,
可以结合bio，nio，aio等数据处理流程进行理解。

目前对原有流程进行了一些小的调整，但都是同步的读写操作，
但是从测试的结果来看，客户端貌似不能够及时的读取出来服务端回传的结果，
或者可以说是服务器传输的数据的语义不能够正确的被理解。
```
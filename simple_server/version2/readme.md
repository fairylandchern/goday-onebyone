### version2
#### 功能及存在的问题
```
支持客户端和服务端的数据读写分离操作，并且需要针对version1中服务端的数据无法被正确读出的现象进行解决，
可能是由于客户端并没有保存相应链接的数据的对应关系，所以导致最后的数据无法根据token找到相应的正确的客
户端，从而导致后续的处理逻辑失败。

目前，所有的代码都放在一个单体之中，暂未考虑处理流程与对象的深度绑定，容错处理等目前暂未完善，需要在后
续的编码中逐步完善该逻辑。

client处理的流程，目前仍然是无状态的，数据被读回之后的后续处理措施仍然不完善，待解决。

client写入数据到服务端的时候，不加锁，貌似数据也是可以被正确的传输到服务端的，这个还需要后续进行一些
测试，对细节支持进行验证
```
#### 改进计划
```text
1.保存请求id和对应的调用方法的对应关系，方便收到响应后进行回调及后续处理，并且处理流程要与单个client相绑
定，后续服务仍然具有大堆的改进空间。

2.应该首先重构client部分，保证流程可以被穿起来
```
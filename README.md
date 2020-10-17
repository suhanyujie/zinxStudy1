# zinx 服务器开发
Go 网络编程实践

## relation
客户端 https://github.com/suhanyujie/zinxDemo1

## zinx 消息 TLV 序列化
TCP 连接建立后，客户端以**流**的方式跟服务端进行通讯，此时，我们无法区分某个数据包的长度，因此，需要以一定的格式定义数据。
客户端按照特定的格式组装数据（封包），发送到服务端，服务端再以约定好的格式拆解数据（拆包）。

### 拆包
* 例如，先读取数据包的前几个字节，这几个字节存储的是数据的长度和类型。然后在根据消息的类型和长度读取实际的数据。

## branch desc
### v0.5.1
根据丹冰前辈的[教程](https://www.bilibili.com/video/BV1wE411d7th?p=21)，这个分支的 feature 就是讲封包、解包的逻辑结合到 zinx 框架中。主要分为 3 点：
* 将 Message 添加到 Request 中
* 修改 Connection 中读取数据的机制，由之前读取字节的方式，改为读 TLV 数据包的方式
* 给链接提供一个发包机制，发送的 TLV 数据包

### v0.6
这个 feature 是为了给 zinx 服务框架新增一个消息管理模块，用于支持多路由。
* 该模块具备一个属性：用于存放 msg id 和 router 的映射关系
* 具备方法，1.增加路由到集合中，2.支持通过 msg id 映射到对应的 router 进行处理

### v0.7
请求连接的读写分离

### v0.8
引入消息队列/协程池，将处理业务的协程的数量固定下来，以免协程数太多导致协程切换资源开销太大

#### step
* 在 MsgHandler 中添加协程池数量 WorkPoolSize，存储请求的队列 TaskQueue 等属性
* 根据属性配置，创建协程
* 将客户端请求接入队列

### v0.9
加入连接管理机制，连接进入后，实例化成一个连接，将其加入连接管理器，连接关闭后，将该连接移除

#### step
- 创建、定义连接管理模块
- 将管理连接模块继承到 zinx 中
- 在创建连接之前，销毁连接之前提供钩子

### v0.9.1
提供 api，供开发者注册连接开始、关闭的 hook

## reference
* 导入本地的包 https://www.cnblogs.com/wind-zhou/p/12824857.html
* zinx 官方教程 https://www.bilibili.com/video/av71067087
* zinx 官方仓库 https://github.com/aceld/zinx

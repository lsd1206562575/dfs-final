# Distributed System - Large File Distributed Storage System

The [link]([https://drive.google.com/drive/folders/1NjjSyvn2dLs7_yWIESPTeVHNM4oTj3CG?usp=sharing](https://drive.google.com/drive/folders/1PZZzg5v7j_IXsCzuKvT4Oj9ADcSfu8JK?usp=drive_link)) for Slides and demo video

### 主要 main 函数都在 cmd/DBFS 下

需要运行三个组件：

1. datanode 组件，主要用来上传文件或者下载文件
2. namenode 组件，将文件的元数据信息存储至 redis 中，以及在 datanode 上传完成后的信息写入到 redis 中
3. client 组件，主要用来生成文件的元数据信息，并对文件进行分发

### 所有公共的部分都在 internal 包下：

每个包都跟是对应的关系，根据作用放在不同的包下

1. 其中 rpc 下定义的都是调用的 args 和 reply
2. client 包含了主要的处理逻辑
3. node 下包含了 namenode 和 datanode 的所提供的服务
4. meta 定义的主要是文件的相关信息
5. db 定义相关的数据库信息
6. conf 其实可以直接搞一个配置中心，每次通过 rpc 获取相关配置

### 如何测试？

所有的测试数据都在 test 目录下

test/chunk 下放的是 文件分片后的数据

1. 先启动 datanode 和 namenode， 去 cmd/DBFS 下的 datanode 和 namenode 下，启动 main.go
2. 启动 client 时，需要指定相关参数 `go run main.go -m=upload -f="test.mp4"`

### TODO List

1. namenode 的相关代码
2. redis 需要上传那些信息
3. 如何将分片后的多个文件 分发到不同的 datanode 中去
4. redis 可以搞一个主备
5. datanode 服务器怎样来进行文件的备份？？？

### 服务器信息：

- 34.86.98.8 （redis 主备服务器）
- 34.150.221.87 （3 个 datanode）
- 35.245.122.13 （3 个 datanode）

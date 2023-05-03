# DFS记录

### 1. 使用docker搭建sftp服务器（也可以使用S3来存，把文件存到S3就好）

参考博客：https://juejin.cn/post/6983896796430860296     https://springboot.io/t/topic/4875

- 拉取镜像

  ```bash
  # 搜索
  docker search sftp
  
  # 拉取镜像
  docker pull atmoz/sftp
  ```

- 启动容器

  ```bash
  # 单个用户启动
  docker run --name sftp -v /data/docker/volumes/sftp/upload:/home/admin/upload --privileged=true -p 2022:22 -d atmoz/sftp admin:pass:1001
  
  # 多用户
  docker run --name sftp -v /data/docker/volumes/sftp/conf/users.conf:/etc/sftp/users.conf:ro -v /data/docker/volumes/sftp/data:/home --privileged=true -p 2022:22 -d atmoz/sftp
  ```

- 自定义用户

  ```bash
  # 添加用户配置
  vi /data/docker/volumes/sftp/conf/users.conf
  
  #user:pass:uid:gid 用户名:密码:用户id:组id
  # 用户信息
  admin:123456:1001:100
  test:123456:1002:100
  
  # 修改配置文件权限
  chmod 755 /data/docker/volumes/sftp/conf/users.conf
  ```

- 权限问题: 通过xftp连接sftp上传文件提示：`Sending the file failed.`

  ```bash
  # 容器启动后，用户会在数据目录下生成对应的目录：/data/docker/volumes/sftp/data/admin
  # 但这个/data/docker/volumes/sftp/data/admin目录并不能直接上传文件，需要将目录权限修改为777
  # 需要在这个目录下单独创建一个目录并修改权限才能上传文件
  mkdir -p /data/docker/volumes/sftp/data/admin/upload
  chmod 777 /data/docker/volumes/sftp/data/admin/upload
  
  # 重启容器
  docker restart sftp
  
  或者直接在容器中修改 upload 目录权限：
  $ docker exec -it sftp bash
  root@35f5c9abeb71:/# cd home/
  root@35f5c9abeb71:/home/admin# ls -lh
  total 0
  drwxr-xr-x. 2 root root 6 May 10 07:05 upload
  root@35f5c9abeb71:/home/admin# chmod 777 upload/ -R
  root@35f5c9abeb71:/home/admin# ls -lh
  total 0
  drwxrwxrwx. 2 root root 21 May 10 07:27 upload
  ```

- sftp常用命令

  ```
  登陆：
  sftp -P <端口> <用户名>@<IP>
  
  ?：查看帮助
  quit：退出
  cd lcd：进入某目录 （注：有l前缀表示是宿主机）
  ls lls：查看目录
  pwd lpwd：查看当前路径
  mdir ：创建目录
  put：上传文件（目录：-r）
  get：下载文件
  ```

所以对于我们搭建只需要先安装docker环境，然后拉镜像，启动容器：

```bash
docker pull atmoz/sftp

docker run --name sftp1 -v /data/docker/volumes/sftp1/upload:/home/admin/upload --privileged=true -p 2021:22 -d atmoz/sftp admin:admin:1001

# 再给每个挂载的目录修改upload的权限
chmod 777 /data/docker/volumes/sftp1/upload
```

***按道理来说，应该再搞一个类似心跳机制的东西，这个sftp服务器应该每隔10s发一个请求，以保证该node是可用的。***

------

### 2. 使用docker搭建redis集群/zookeeper集群来保存元数据信息。

***集群可以实现主备，以此保证元数据不会丢失，当某个节点宕机会自动进行选主。***



------

### 3. 编写相关接口（https://github.com/sivanWu0222/DistributedFileServer）

- 基础版的文件上传服务，基本功能可以使用，例如文件上传下载
- 在文件上传保存之后，在下载那一端部署一个反向代理，然后将文件作为一个静态资源来处理，例如nginx，下载的时候后端服务会提供一个接口，用于构造下载文件的url，客户端获取url之后，就去下载，下载的时候会经过nginx， nginx再做一次静态资源访问将文件download下来，一些限流以及权限访问都可以在nginx做，可以减轻golang实现后端的压力。
- 定义的文件元信息，从而生成文件的meta data
- 
- 
# 通用初始化

> 以maxblog-sgw为例

## 1. 代码仓库

### 1.1 新建repo

1. 选择`go-maxms`模板新建，并clone

### 1.2 代码操作

1. 将所有`go-maxms`替换为`maxblog-sgw`，注意修改jenkins的全局变量`StaticFolderName`
2. 检查所需组件
    + 数据库
    + HTTP/RPC
3. 检查配置文件，注意数据库名称，如果有vault，注意添加kv
4. 检查路由、中间件
5. 检查依赖，运行`make wire`
6. 运行`make lint`
7. 修改makefile中的contract链接，运行`make spec`，获取contract
8. `make run`
9. commit

### 1.3 配置代码仓库

1. 进入`Settings`
2. `General`目录，找到`Pull Requests`，只允许squash merging
3. `Branches`目录，如下配置

![分支配置.png](img/init_common/分支配置.png)

4. 进入用户全局`Settings`，进入`Applications`,找到`jenkins-go-maxms`，点击`Configuration`，在`Repository access`中，将maxblog-sgw加进去

### 1.4 数据库

1. 如果用到数据库模块，需要在数据库中添加对应名称的数据库
2. 在根目录创建.env文件，内容格式如下，填写secret
```shell
# mysql
MYSQL_DB_NAME=
MYSQL_USER_NAME=
MYSQL_PASSWORD=

# wechat
APP_ID=
APP_SECRET=
```

## 2. Jenkins

### 2.1 Job

1. 新建View `MaxBlog`，存在则作罢
2. 新建Multibranch Pipeline，配置如下，`Validate`成功即可

![配置多流水线分支.png](img/init_common/配置多流水线分支.png)

### 2.2 Harbor

新建Project

![Harbor新建项目.png](img/init_common/Harbor新建项目.png)

### 2.3 Prometheus

修改prometheus.yml，增加job

## 3 启动项目

### 3.1 使用docker compose

1. 新建`maxblog-sgw`文件夹
2. 拷贝bin，environment，Dockerfile，docker-compose，注意修改docker-compose文件中的命名
3. cd到`maxblog-sgw`，构建镜像
```shell
docker build -t maxblog-sgw .
```
4. 创建容器
```shell
docker compose up -d
```

# go-deploy-system-server

基于 Go 语言的自动化代码发布（部署）系统后端服务，可将 Git 仓库代码通过 RSYNC 同步发布到远程服务器，支持回滚。

## 技术栈

| 类别 | 技术 |
|---|---|
| 语言 | Go 1.22 |
| Web 框架 | Gin v1.10 |
| ORM | GORM v1.25 + MySQL |
| 认证 | JWT (dgrijalva/jwt-go) |
| 密码加密 | Scrypt + AES |
| 日志 | Logrus + file-rotatelogs（按大小/时间切割） |
| 配置 | INI 文件 (gopkg.in/ini.v1) |
| SSH | golang.org/x/crypto/ssh |
| 部署 | Docker 容器化 |

## 目录结构

```
├── main.go                  # 入口：初始化数据库 → 启动路由
├── config/
│   └── config.ini           # 配置文件
├── routes/
│   └── router.go            # 路由定义
├── api/v1/                  # HTTP Handler 层
│   ├── login.go             # 登录
│   ├── user.go              # 用户管理
│   ├── department.go        # 部门管理
│   ├── engineroom.go        # 机房管理
│   ├── server.go            # 服务器管理
│   ├── deployment.go        # 发布项目配置管理
│   ├── release.go           # 代码发布/拉取/回滚
│   ├── deploymentlog.go     # 发布日志
│   ├── upload.go            # 秘钥上传
│   └── health.go            # 健康检查
├── middleware/              # 中间件
│   ├── cors.go              # 跨域处理
│   ├── jwt.go               # JWT 认证（生成/校验）
│   ├── logger.go            # 请求日志记录
│   └── role.go              # 管理员权限校验
├── model/                   # 数据模型 + 业务逻辑
│   ├── Db.go                # 数据库初始化、自动迁移、初始化默认账户
│   ├── Department.go        # 部门 CRUD
│   ├── User.go              # 用户 CRUD + 修改密码
│   ├── Engineroom.go        # 机房 CRUD
│   ├── Server.go            # 服务器 CRUD + SSH 连接
│   ├── Deployment.go        # 发布项目配置 CRUD
│   ├── Login.go             # 登录校验
│   ├── Release.go           # Git 拉取/发布/回滚核心逻辑
│   ├── DeploymentLog.go     # 发布日志查询
│   └── utils.go             # 锁文件、权限检查、Git 操作、SSH 远程执行
├── utils/                   # 工具包
│   ├── config.go            # 加载 .ini 配置
│   ├── aespwd/              # AES 加解密
│   ├── scryptpwd/           # Scrypt 密码哈希
│   ├── md5/                 # MD5
│   ├── errmsg/              # 统一错误码
│   └── validator/           # 请求参数校验
└── Dockerfile               # 多阶段构建镜像
```

## 快速开始

### 1. 配置

修改 `config/config.ini`：

```ini
[server]
AppMode = debug          # debug | release
HttpPort = :3000
JwtKey = 89js82js72@a=KCAFJWQER012
PwdKey = aoefqCINAETCA
ServerGitKey = a&D*71&FBA12-9P*
KeyFilePath = data/go_deployment_system/upload/key
CodePath = data/go_deployment_system/git

[database]
DbHost = 127.0.0.1
DbPort = 3307
DbUser = root
DbPassWord = root123
DbName = go_deployment_system

[log]
LogPath = data/go_deployment_system/log
LogFileName = ops.log
LogSaveTime = 10          # 最大保存天数
LogSplitSize = 10          # 切割大小 (MB)
```

### 2. 运行

确保 MySQL 已启动并创建好对应数据库，然后：

```bash
go run main.go
```

或者指定配置文件路径：

```bash
go run main.go -config /path/to/config.ini
```

### 3. 默认账户

首次启动会自动创建：

- 部门：管理部
- 用户名：`admin`
- 密码：`123456`
- 角色：管理员

## API 路由

服务监听 `:3000`，所有接口前缀为 `api/v1`。

### 公开接口

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/v1/login/` | 用户登录 |

### 登录用户接口（需要 JWT Token）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/releases` | 获取当前用户拥有的项目列表 |
| PUT | `/api/v1/release/gitpull/:id` | 拉取项目代码到本地 |
| POST | `/api/v1/release/add` | 发布代码到远程服务器 |
| POST | `/api/v1/release/rollback/:id` | 回滚项目 |
| GET | `/api/v1/deploymentlogs` | 发布日志列表/搜索 |
| POST | `/api/v1/user/changepassword` | 修改密码 |

### 管理员接口（需要 JWT Token + 管理员角色）

**部门管理**

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/department` | 部门列表 |
| POST | `/api/v1/department` | 添加部门 |
| PUT | `/api/v1/department/:id` | 修改部门 |
| DELETE | `/api/v1/department/:id` | 删除部门 |

**用户管理**

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/user` | 用户列表 |
| POST | `/api/v1/user` | 添加用户 |
| PUT | `/api/v1/user/:id` | 修改用户 |
| DELETE | `/api/v1/user/:id` | 删除用户 |

**机房管理**

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/engineroom` | 机房列表 |
| POST | `/api/v1/engineroom` | 添加机房 |
| PUT | `/api/v1/engineroom/:id` | 修改机房 |
| DELETE | `/api/v1/engineroom/:id` | 删除机房 |

**服务器管理**

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/server` | 服务器列表 |
| POST | `/api/v1/server` | 添加服务器 |
| PUT | `/api/v1/server/:id` | 修改服务器 |
| DELETE | `/api/v1/server/:id` | 删除服务器 |
| GET | `/api/v1/server/connect/:id` | 测试服务器连接 |

**项目配置管理**

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/deployment` | 项目配置列表 |
| POST | `/api/v1/deployment` | 添加项目配置 |
| PUT | `/api/v1/deployment/:id` | 修改项目配置 |
| DELETE | `/api/v1/deployment/:id` | 删除项目配置 |

**其他**

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/v1/upload` | 秘钥上传 |
| GET | `/api/v1/health` | 健康检查 |

### 请求认证

登录成功后返回 JWT Token，后续请求在 Header 中携带：

```
Authorization: Bearer <token>
```

## 数据库表

| 表名 | 说明 |
|---|---|
| `department` | 部门 |
| `user` | 用户 |
| `engineroom` | 机房 |
| `server` | 服务器 |
| `deployment` | 发布项目配置（Git 仓库地址、分支、目标路径等） |
| `deployment_to_server` | 项目-服务器多对多关联 |
| `deployment_to_user_role` | 项目-发布用户多对多关联 |
| `deployment_log` | 发布日志 |

## 核心业务流程

1. 管理员在后台配置**机房** → **服务器**（SSH 连接信息） → **发布项目**（Git 仓库、分支、目标路径），并指定哪些服务器和哪些用户拥有该项目的发布权限
2. 用户登录后，可查看自己被授权的项目列表
3. 点击"拉取代码"，系统在发布机上执行 `git clone` / `git pull`，展示变更文件列表
4. 用户选择要发布的文件，填写备注，点击发布
5. 系统将变更文件通过 **RSYNC** 同步到目标服务器
6. 支持**回滚**：根据历史发布日志，恢复到之前的 Git HEAD 版本

### 并发控制

通过**锁文件机制**防止同一项目被多人同时发布。

### 密码安全

- 用户登录密码：Scrypt 哈希存储
- Git 密码 / 服务器密码：AES 加密存储，使用时解密

## Docker 部署

```bash
# 构建镜像
docker build -t go-deploy-system-server:v0.1 .

# 运行容器
docker run --name go-deploy-server -d -p 3000:3000 go-deploy-system-server:v0.1
```

## 前端项目

对应的 Vue 2 + Element UI 前端项目位于同级目录 `go-deploy-system-web/`。

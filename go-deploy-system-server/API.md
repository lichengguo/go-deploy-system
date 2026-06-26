# API 接口文档

Base URL: `http://host:3000/api/v1`

## 通用说明

### 认证方式

除登录和健康检查外，所有接口需在请求头携带 JWT Token：

```
Authorization: Bearer <token>
```

管理员接口额外需要用户角色 `role = 1`。

### 响应结构

所有接口统一返回 JSON：

```json
{
  "status": 200,
  "message": "成功",
  "data": {},
  "total": 0
}
```

### 错误码一览

| 码段 | 含义 |
|---|---|
| 200 | 成功 |
| 500 | 失败 |
| 1001~1006 | 部门模块 |
| 2001~2007 | 用户模块 |
| 3001~3003 | Token 相关 |
| 4001~4005 | 机房模块 |
| 5001~5005 | 锁文件相关 |
| 6001~6005 | 服务器模块 |
| 7001~7012 | 发布项目模块 |
| 9001~9004 | Git 相关 |

### 字段校验说明

请求体字段需满足以下通用规则：

- 使用 `validate` tag 标注约束
- 校验失败时 `status` 为 500，`message` 返回具体字段错误信息

---

## 一、公开接口

### 1.1 登录

```
POST /api/v1/login/
```

**请求体**

```json
{
  "user_name": "admin",
  "password": "123456"
}
```

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| user_name | string | 是 | 用户名，4~20 字符 |
| password | string | 是 | 密码，6~20 字符 |

**成功响应**

```json
{
  "status": 200,
  "message": "成功",
  "user_name": "admin",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "role": 1
}
```

| 字段 | 类型 | 说明 |
|---|---|---|
| token | string | JWT Token，后续请求携带 |
| role | int | 1=管理员，2=普通用户 |

---

## 二、登录用户接口（需要 JWT Token）

### 2.1 获取当前用户的项目列表

```
GET /api/v1/releases
```

**响应 data**

```json
[
  {
    "id": 1,
    "deploy_name": "项目A"
  }
]
```

### 2.2 拉取代码到本地

```
PUT /api/v1/release/gitpull/:id
```

| 参数 | 位置 | 类型 | 说明 |
|---|---|---|---|
| id | Path | int | 发布项目 ID |

**响应 data**

```json
{
  "deploy_file_list": ["src/main.go", "config/app.ini"],
  "git_url": "git@github.com:user/repo.git",
  "git_head": "a1b2c3d",
  "commit_email": "dev@example.com",
  "git_info": "fix: bug修复",
  "deploy_ip": ["192.168.1.10", "192.168.1.11"],
  "deploy_path": "/data/www/app"
}
```

### 2.3 发布代码到服务器

```
POST /api/v1/release/add
```

**请求体**

```json
{
  "deployment_id": 1,
  "deployment_file_list": "src/main.go,config/app.ini",
  "deployment_commit": "上线新功能",
  "git_head": "a1b2c3d"
}
```

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| deployment_id | int | 是 | 发布项目 ID |
| deployment_file_list | string | 是 | 要发布的文件列表，逗号分隔 |
| deployment_commit | string | 是 | 发布备注，最长 100 字符 |
| git_head | string | 是 | Git HEAD 指针，最长 100 字符 |

> 发布流程：先检查锁文件 → 校验用户权限 → 找到目标服务器 → 通过 RSYNC 推送代码到各服务器

### 2.4 回滚项目

```
POST /api/v1/release/rollback/:id
```

| 参数 | 位置 | 类型 | 说明 |
|---|---|---|---|
| id | Path | int | 发布日志 ID |

> 根据发布日志记录恢复到之前版本的 Git HEAD

### 2.5 发布日志列表

```
GET /api/v1/deploymentlogs
```

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| page | int | 否 | 页码，默认 1 |
| pagesize | int | 否 | 每页条数，默认 10，最大 100 |
| deployment_name | string | 否 | 按项目名称搜索 |
| deployment_user_name | string | 否 | 按发布人搜索 |

> 管理员可查看所有日志，普通用户只能看自己的

**响应 data**

```json
[
  {
    "id": 1,
    "deployment_id": 1,
    "deployment_name": "项目A",
    "deployment_user_name": "admin",
    "deployment_file_list": "src/main.go,config/app.ini",
    "deployment_commit": "上线新功能",
    "deployment_status": 1,
    "deployment_fail_info": "",
    "git_head": "a1b2c3d",
    "created_at": 1700000000
  }
]
```

`deployment_status`: 1=成功，2=失败

### 2.6 修改密码

```
POST /api/v1/user/changepassword
```

**请求体**

```json
{
  "old_pwd": "123456",
  "new_pwd": "654321"
}
```

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| old_pwd | string | 是 | 旧密码，6~100 字符 |
| new_pwd | string | 是 | 新密码，6~100 字符 |

---

## 三、管理员接口（需要 JWT Token + role=1）

### 3.1 部门管理

#### 3.1.1 部门列表

```
GET /api/v1/department
```

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| page | int | 否 | 页码 |
| pagesize | int | 否 | 每页条数 |
| department_name | string | 否 | 按名称搜索 |

**响应 data**

```json
[
  {
    "id": 1,
    "department_name": "管理部",
    "created_at": 1700000000,
    "updated_at": 1700000000
  }
]
```

#### 3.1.2 添加部门

```
POST /api/v1/department
```

**请求体**

```json
{
  "department_name": "研发部"
}
```

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| department_name | string | 是 | 部门名称，最长 20 字符 |

#### 3.1.3 修改部门

```
PUT /api/v1/department/:id
```

**请求体**

```json
{
  "department_name": "研发二部"
}
```

#### 3.1.4 删除部门

```
DELETE /api/v1/department/:id
```

> 如果部门下存在用户，无法删除

---

### 3.2 用户管理

#### 3.2.1 用户列表

```
GET /api/v1/user
```

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| page | int | 否 | 页码 |
| pagesize | int | 否 | 每页条数 |
| user_name | string | 否 | 按用户名搜索 |

**响应 data**

```json
[
  {
    "id": 1,
    "user_name": "admin",
    "role": 1,
    "status": 1,
    "department_id": 1,
    "department": {
      "id": 1,
      "department_name": "管理部"
    },
    "created_at": 1700000000
  }
]
```

| role | 说明 |
|---|---|
| 1 | 管理员 |
| 2 | 普通用户 |

| status | 说明 |
|---|---|
| 1 | 可用 |
| 2 | 禁用 |

#### 3.2.2 添加用户

```
POST /api/v1/user
```

**请求体**

```json
{
  "user_name": "zhangsan",
  "password": "123456",
  "role": 2,
  "status": 1,
  "department_id": 1
}
```

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| user_name | string | 是 | 用户名，4~20 字符 |
| password | string | 是 | 密码，6~100 字符 |
| role | int | 是 | 1=管理员，2=普通用户 |
| status | int | 是 | 1=可用，2=禁用 |
| department_id | int | 是 | 所属部门 ID |

#### 3.2.3 修改用户

```
PUT /api/v1/user/:id
```

> 请求体同添加用户。密码与数据库不一致时才会更新（已自动 Scrypt 加密）。

#### 3.2.4 删除用户

```
DELETE /api/v1/user/:id
```

> 会同时删除该用户在项目-用户关联表中的记录

---

### 3.3 机房管理

#### 3.3.1 机房列表

```
GET /api/v1/engineroom
```

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| page | int | 否 | 页码 |
| pagesize | int | 否 | 每页条数 |
| engineroom_name | string | 否 | 按名称搜索 |

**响应 data**

```json
[
  {
    "id": 1,
    "engineroom_name": "北京机房",
    "contact": "张三",
    "contact_info": "13800138000",
    "address": "北京市朝阳区xxx",
    "created_at": 1700000000
  }
]
```

#### 3.3.2 添加机房

```
POST /api/v1/engineroom
```

**请求体**

```json
{
  "engineroom_name": "北京机房",
  "contact": "张三",
  "contact_info": "13800138000",
  "address": "北京市朝阳区xxx"
}
```

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| engineroom_name | string | 是 | 机房名称，最长 100 字符 |
| contact | string | 是 | 联系人，最长 100 字符 |
| contact_info | string | 是 | 联系方式，最长 100 字符 |
| address | string | 是 | 机房地址，最长 100 字符 |

#### 3.3.3 修改机房

```
PUT /api/v1/engineroom/:id
```

> 请求体同添加机房

#### 3.3.4 删除机房

```
DELETE /api/v1/engineroom/:id
```

> 如果机房下存在服务器，无法删除

---

### 3.4 服务器管理

#### 3.4.1 服务器列表

```
GET /api/v1/server
```

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| page | int | 否 | 页码 |
| pagesize | int | 否 | 每页条数 |
| server_name | string | 否 | 按名称搜索 |

**响应 data**

```json
[
  {
    "id": 1,
    "engineroom_id": 1,
    "engineroom": {
      "id": 1,
      "engineroom_name": "北京机房"
    },
    "server_name": "web-01",
    "server_ip": "192.168.1.10",
    "server_port": "22",
    "server_user": "root",
    "server_key": "data/keys/abc123",
    "server_status": 1,
    "created_at": 1700000000
  }
]
```

> 响应中不包含密码字段

#### 3.4.2 添加服务器

```
POST /api/v1/server
```

**请求体**

```json
{
  "engineroom_id": 1,
  "server_name": "web-01",
  "server_ip": "192.168.1.10",
  "server_port": "22",
  "server_user": "root",
  "server_pwd": "mypassword",
  "server_key": "data/keys/abc123",
  "server_status": 1
}
```

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| engineroom_id | int | 是 | 所属机房 ID |
| server_name | string | 是 | 服务器名称，最长 100 字符 |
| server_ip | string | 是 | 登录 IP，最长 100 字符 |
| server_port | string | 是 | 登录端口，最长 100 字符 |
| server_user | string | 是 | 登录账号，最长 100 字符 |
| server_pwd | string | 否 | 登录密码（AES 加密存储） |
| server_key | string | 否 | 登录秘钥文件路径 |
| server_status | int | 是 | 1=可用，2=冻结 |

> 同一机房下，服务器名称 + IP + 用户不能同时相同
> `server_pwd` 和 `server_key` 至少填写一个用于 SSH 连接

#### 3.4.3 修改服务器

```
PUT /api/v1/server/:id
```

> 请求体同添加服务器

#### 3.4.4 删除服务器

```
DELETE /api/v1/server/:id
```

#### 3.4.5 测试连接

```
GET /api/v1/server/connect/:id
```

> 使用配置的密码或秘钥尝试 SSH 连接到目标服务器

---

### 3.5 项目配置管理

#### 3.5.1 项目配置列表

```
GET /api/v1/deployment
```

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| page | int | 否 | 页码 |
| pagesize | int | 否 | 每页条数 |
| deployment_name | string | 否 | 按项目名称搜索 |

**响应 data**

```json
[
  {
    "id": 1,
    "deploy_name": "项目A",
    "git_url_http": "https://github.com/user/repo.git",
    "git_url_ssh": "",
    "git_branch": "master",
    "git_user": "gituser",
    "git_key": "data/keys/xyz789",
    "deploy_server_path": "/data/www/app",
    "server_list": [
      { "id": 1, "server_name": "web-01", "server_ip": "192.168.1.10" }
    ],
    "user_list": [
      { "id": 2, "user_name": "zhangsan" }
    ],
    "created_at": 1700000000
  }
]
```

> `git_passwd` 不会返回给前端

#### 3.5.2 添加项目配置

```
POST /api/v1/deployment
```

**请求体**

```json
{
  "deploy_name": "项目A",
  "git_url_http": "https://github.com/user/repo.git",
  "git_url_ssh": "",
  "git_branch": "master",
  "git_user": "gituser",
  "git_passwd": "gitpassword",
  "git_key": "",
  "deploy_server_path": "/data/www/app",
  "server_id": [1, 2],
  "user_id": [2, 3]
}
```

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| deploy_name | string | 是 | 项目名称，最长 100 字符 |
| git_url_http | string | 条件 | HTTP Git 地址，与 git_url_ssh 二选一 |
| git_url_ssh | string | 条件 | SSH Git 地址，与 git_url_http 二选一 |
| git_branch | string | 否 | Git 分支，默认 master，最长 50 字符 |
| git_user | string | 条件 | Git 账号，使用 HTTP 时必填 |
| git_passwd | string | 条件 | Git 密码，使用 HTTP 时必填（AES 加密存储） |
| git_key | string | 条件 | Git 秘钥，使用 SSH 时必填 |
| deploy_server_path | string | 是 | 服务器目标发布目录，最长 100 字符 |
| server_id | []int | 是 | 目标服务器 ID 列表 |
| user_id | []int | 是 | 授权发布用户 ID 列表 |

#### 3.5.3 修改项目配置

```
PUT /api/v1/deployment/:id
```

> 请求体同添加项目配置

#### 3.5.4 删除项目配置

```
DELETE /api/v1/deployment/:id
```

> 会同时删除项目-服务器和项目-用户的关联记录

---

### 3.6 秘钥上传

```
POST /api/v1/upload
```

**请求方式**: `multipart/form-data`

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| keyfile | file | 是 | 上传的秘钥文件 |

**成功响应**

```json
{
  "status": 200,
  "message": "成功",
  "url": "data/go_deployment_system/upload/key/abc123def456"
}
```

> 文件名会被 MD5 化以防重名，文件权限自动设为 600

---

### 3.7 健康检查

```
GET /api/v1/health
```

> 无需认证

**响应**

```json
{
  "status": 200,
  "message": "ok"
}
```

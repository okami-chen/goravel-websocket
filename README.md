# Golang实现的分布式WebSocket微服务

[![Go](https://img.shields.io/badge/Go-1.13-blue.svg)](https://golang.google.cn)
![GitHub release](https://img.shields.io/github/v/release/okami-chen/goravel-websocket)
![Travis (.org)](https://api.travis-ci.com/okami-chen/goravel-websocket.svg?branch=master)
[![star](https://img.shields.io/github/stars/woodylan/go-websocket?style=social)](https://github.com/woodylan/go-websocket/stargazers)

## 简介

本系统基于Golang、ETCD、RPC实现分布式WebSocket微服务，也可以单机部署，单机部署不需要ETCD、RPC。分布式部署可以支持nginx负责均衡、水平扩容部署，程序之间使用RPC通信。

基本流程为：用ws协议连接本服务，得到一个clientId，由客户端上报这个clinetId给服务端，服务端拿到这个clientId之后，可以给这个客户端发送信息，绑定这个客户端都分组，给分组发送消息。

目前实现的功能有，给指定客户端发送消息、绑定客户端到分组、给分组里的客户端批量发送消息、获取在线的客户端、上下线自动通知。适用于长连接的大部分场景，分组可以理解为聊天室，绑定客户端到分组相当于把客户端添加到聊天室，给分组发送信息相当于给聊天室的每个人发送消息。


## 接口

#### 连接接口

**请求地址：**/api/websocket/ws?systemId=xxx

**协议：** websocket

**请求参数**：systemId 系统ID

**响应示例：**

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "clientId": "9fa54bdbbf2778cb"
  }
}
```

#### 注册系统

**请求地址：**/api/websocket/register

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**响应示例：**

```json
{
  "code": 0,
  "msg": "success",
  "data": []
}
```

#### 发送信息给指定客户端

**请求地址：**/api/websocket/send_to_client

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| clientId | string | 是       | 客户端ID |
| sendUserId | string | 是       | 发送者ID |
| code | integer | 是       | 自定义的状态码 |
| msg | string | 是       | 自定义的状态消息 |
| data | sring、array、object | 是       | 消息内容 |

**响应示例：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "messageId": "5b4646dd8328f4b1"
    }
}
```

#### 绑定客户端到分组

**请求地址：**/api/websocket/bind_to_group

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| sendUserId | string | 是       | 发送者ID |
| clientId | string | 是       | 客户端ID |
| groupName | string | 是       | 分组名 |

**响应示例：**

```json
{
  "code": 0,
  "msg": "success",
  "data": []
}
```

#### 发送信息给指定分组

**请求地址：**/api/websocket/send_to_group

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| sendUserId | string | 是       | 发送者ID |
| groupName | string | 是       | 分组名 |
| code | integer | 是       | 自定义的状态码 |
| msg | string | 是       | 自定义的状态消息 |
| data | sring、array、object | 是       | 消息内容 |

**响应示例：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "messageId": "5b4646dd8328f4b1"
    }
}
```


#### 发送信息给指定系统

**请求地址：**/api/websocket/send_to_system

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| sendUserId | string | 是       | 发送者ID |
| code | integer | 是       | 自定义的状态码 |
| msg | string | 是       | 自定义的状态消息 |
| data | sring、array、object | 是       | 消息内容 |

**响应示例：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "messageId": "5b4646dd8328f4b1"
    }
}
```

#### 获取在线的客户端列表

**请求地址：**/api/websocket/get_online_list

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| groupName | string | 是       | 分组名 |
| code | integer | 是       | 自定义的状态码 |
| msg | string | 是       | 自定义的状态消息 |
| data | sring、array、object | 是       | 消息内容 |

**响应示例：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "count": 2,
        "list": [
            "WQReWw6m+wct+eKk/2rDiWcU4maU8JRTRZEX8c7Te6LzCa//VCXr/0KeVyO0sdNt",
            "j6YdsGFH4rfbYN/vS6UavJ5fVclWIB9W+Gqg9R/92cLJqgAp2ZPkvMbQiwQBJmDc"
        ]
    }
}
```

#### 发送指定连接

**请求地址：**/api/websocket/close_client

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| clientId | string | 是       | 客户端ID |

**响应示例：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {}
}
```
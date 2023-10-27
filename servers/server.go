package servers

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/gorilla/websocket"
	"github.com/okami-chen/goravel-websocket/events"
	"github.com/okami-chen/goravel-websocket/tools/util"
	"strconv"
	"time"
)

// channel通道
var ToClientChan chan clientInfo

// channel通道结构体
type clientInfo struct {
	ClientId   string
	SendUserId string
	MessageId  string
	Code       int
	Msg        string
	Data       *string
}

type RetData struct {
	MessageId  string      `json:"messageId"`
	SendUserId string      `json:"sendUserId"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

// 心跳间隔
var heartbeatInterval = 60 * time.Second

func init() {
	ToClientChan = make(chan clientInfo, 1000)
}

var Manager = NewClientManager() // 管理者

// 发送信息到指定客户端
func SendMessage2Client(clientId string, sendUserId string, code int, msg string, data *string) (messageId string) {
	messageId = util.GenUUID()
	SendMessage2LocalClient(messageId, clientId, sendUserId, code, msg, data)
	return
}

// 关闭客户端
func CloseClient(clientId, systemId string) {
	if util.IsCluster() {
		//addr, _, _, isLocal, err := util.GetAddrInfoAndIsLocal(clientId)
		//if err != nil {
		//	log.Errorf("%s", err)
		//	return
		//}
		//
		////如果是本机则发送到本机
		//if isLocal {
		//	CloseLocalClient(clientId, systemId)
		//} else {
		//	//发送到指定机器
		//	CloseRpcClient(addr, clientId, systemId)
		//}
	} else {
		//如果是单机服务，则只发送到本机
		CloseLocalClient(clientId, systemId)
	}

	return
}

// 添加客户端到分组
func AddClient2Group(systemId string, groupName string, clientId string, userId string, extend string) {
	//如果是集群则用redis共享数据
	if util.IsCluster() {
		//判断key是否存在
		//addr, _, _, isLocal, err := util.GetAddrInfoAndIsLocal(clientId)
		//if err != nil {
		//	log.Errorf("%s", err)
		//	return
		//}
		//
		//if isLocal {
		//	if client, err := Manager.GetByClientId(clientId); err == nil {
		//		//添加到本地
		//		Manager.AddClient2LocalGroup(groupName, client, userId, extend)
		//	} else {
		//		log.Error(err)
		//	}
		//} else {
		//	//发送到指定的机器
		//	SendRpcBindGroup(addr, systemId, groupName, clientId, userId, extend)
		//}
	} else {
		if client, err := Manager.GetByClientId(clientId); err == nil {
			//如果是单机，就直接添加到本地group了
			Manager.AddClient2LocalGroup(groupName, client, userId, extend)
		}
	}
}

// 发送信息到指定分组
func SendMessage2Group(systemId, sendUserId, groupName string, code int, msg string, data *string) (messageId string) {
	messageId = util.GenUUID()
	if util.IsCluster() {
		//发送分组消息给指定广播
		//go SendGroupBroadcast(systemId, messageId, sendUserId, groupName, code, msg, data)
	} else {
		//如果是单机服务，则只发送到本机
		Manager.SendMessage2LocalGroup(systemId, messageId, sendUserId, groupName, code, msg, data)
	}
	return
}

// 发送信息到指定系统
func SendMessage2System(systemId, sendUserId string, code int, msg string, data string) {
	messageId := util.GenUUID()
	//如果是单机服务，则只发送到本机
	Manager.SendMessage2LocalSystem(systemId, messageId, sendUserId, code, msg, &data)
}

// 获取分组列表
func GetOnlineList(systemId *string, groupName *string) map[string]interface{} {
	var clientList []string
	retList := Manager.GetGroupClientList(*systemId + ":" + *groupName)
	clientList = append(clientList, retList...)
	return map[string]interface{}{
		"count": len(clientList),
		"list":  clientList,
	}
}

// 通过本服务器发送信息
func SendMessage2LocalClient(messageId, clientId string, sendUserId string, code int, msg string, data *string) {
	//facades.Log().Info("发送到通道：" + clientId)
	ToClientChan <- clientInfo{ClientId: clientId, MessageId: messageId, SendUserId: sendUserId, Code: code, Msg: msg, Data: data}
	return
}

// 发送关闭信号
func CloseLocalClient(clientId, systemId string) {
	if conn, err := Manager.GetByClientId(clientId); err == nil && conn != nil {
		if conn.SystemId != systemId {
			return
		}
		Manager.DisConnect <- conn
		facades.Log().Info("主动踢掉客户端：" + clientId)

		//连接事件
		t := carbon.Now("PRC").ToDateTimeString()
		events.NewClientKillEvent(conn.UserId, conn.UserId, t)
	}
	return
}

// 监听并发送给客户端信息
func WriteMessage() {
	for {
		clientInfo := <-ToClientChan
		if conn, err := Manager.GetByClientId(clientInfo.ClientId); err == nil && conn != nil {
			if err := Render(conn.Socket, clientInfo.MessageId, clientInfo.SendUserId, clientInfo.Code, clientInfo.Msg, clientInfo.Data); err != nil {
				Manager.DisConnect <- conn
				facades.Log().Error("终端设备离线：" + clientInfo.ClientId)
				//设备失败事件
				t := carbon.Now("PRC").ToDateTimeString()
				events.NewClientMessageFailEvent(conn.UserId, conn.UserId, t, clientInfo.MessageId)

				//设备离线
				events.NewClientOffloneEvent(conn.UserId, conn.UserId, t)
			} else {
				facades.Log().Infof("终端设备消息：%s, 消息编号：%s", clientInfo.ClientId, clientInfo.MessageId)
				//设备成功事件
				t := carbon.Now("PRC").ToDateTimeString()
				events.NewClientMessageSuccessEvent(conn.UserId, conn.UserId, t, clientInfo.MessageId)
			}
		}
	}
}

func Render(conn *websocket.Conn, messageId string, sendUserId string, code int, message string, data interface{}) error {
	return conn.WriteJSON(RetData{
		Code:       code,
		MessageId:  messageId,
		SendUserId: sendUserId,
		Msg:        message,
		Data:       data,
	})
}

// 启动定时器进行心跳检测
func PingTimer() {
	go func() {
		interval := facades.Config().Get("websocket.interval", "10")
		num, err := strconv.Atoi(interval.(string))
		if err != nil {
			facades.Log().Errorf("类型转换失败: %s ", err.Error())
			return
		}
		duration := time.Duration(num) * time.Second
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			<-ticker.C
			//发送心跳
			for clientId, conn := range Manager.AllClient() {
				if err := conn.Socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					Manager.DisConnect <- conn
					facades.Log().Errorf("发送心跳失败: %s 总连接数：%d", clientId, Manager.Count())
					t := carbon.Now("PRC").ToDateTimeString()
					events.NewClientOffloneEvent(conn.UserId, conn.UserId, t)
				} else {
					facades.Log().Infof("发送心跳成功: %s 总连接数：%d", clientId, Manager.Count())
					//设备在线事件
					t := carbon.Now("PRC").ToDateTimeString()
					e := events.NewClientKeepLiveEvent(conn.UserId, conn.UserId, t)
					if e != nil {
						facades.Log().Errorf("NewClientKeepLiveEvent: %s ", e.Error())
					}
				}
			}
		}
	}()
}

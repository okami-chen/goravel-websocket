package servers

import (
	"encoding/json"
	"errors"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/okami-chen/goravel-websocket/events"
	"sync"
	"time"
)

// 连接管理
type ClientManager struct {
	ClientIdMap     map[string]*Client // 全部的连接
	ClientIdMapLock sync.RWMutex       // 读写锁

	Connect    chan *Client // 连接处理
	DisConnect chan *Client // 断开连接处理

	GroupLock sync.RWMutex
	Groups    map[string][]string

	SystemClientsLock sync.RWMutex
	SystemClients     map[string][]string
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		ClientIdMap:   make(map[string]*Client),
		Connect:       make(chan *Client, 10000),
		DisConnect:    make(chan *Client, 10000),
		Groups:        make(map[string][]string, 100),
		SystemClients: make(map[string][]string, 100),
	}

	return
}

// 管道处理程序
func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.Connect:
			// 建立连接事件
			manager.EventConnect(client)
		case conn := <-manager.DisConnect:
			// 断开连接事件
			manager.EventDisconnect(conn)
		}
	}
}

// 建立连接事件
func (manager *ClientManager) EventConnect(client *Client) {
	manager.AddClient(client)
	facades.Log().Infof("客户端已连接: %s, 客户端总数: %d", client.ClientId, Manager.Count())

	//连接事件
	t := carbon.Carbon{}.Now("PRC").ToDateTimeString()
	events.NewClientConnectEvent(client.UserId, client.UserId, t)
}

// 断开连接时间
func (manager *ClientManager) EventDisconnect(client *Client) {
	//关闭连接
	_ = client.Socket.Close()
	manager.DelClient(client)

	mJson, _ := json.Marshal(map[string]string{
		"clientId": client.ClientId,
		"userId":   client.UserId,
		"extend":   client.Extend,
	})
	data := string(mJson)
	sendUserId := ""

	//发送下线通知
	if len(client.GroupList) > 0 {
		for _, groupName := range client.GroupList {
			SendMessage2Group(client.SystemId, sendUserId, groupName, 1002, "客户端下线", &data)
		}
	}

	t := uint64(time.Now().Unix()) - client.ConnectTime
	facades.Log().Infof("客户端已断开: %s, 总数: %d, 耗时: %d", client.ClientId, Manager.Count(), t)
	//连接断开事件
	events.NewClientDisConnectEvent(client.UserId, client.UserId, t)

	//标记销毁
	client.IsDeleted = true
	client = nil
}

// 添加客户端
func (manager *ClientManager) AddClient(client *Client) {
	manager.ClientIdMapLock.Lock()
	defer manager.ClientIdMapLock.Unlock()

	manager.ClientIdMap[client.ClientId] = client
}

// 获取所有的客户端
func (manager *ClientManager) AllClient() map[string]*Client {
	manager.ClientIdMapLock.RLock()
	defer manager.ClientIdMapLock.RUnlock()

	return manager.ClientIdMap
}

// 客户端数量
func (manager *ClientManager) Count() int {
	manager.ClientIdMapLock.RLock()
	defer manager.ClientIdMapLock.RUnlock()
	return len(manager.ClientIdMap)
}

// 删除客户端
func (manager *ClientManager) DelClient(client *Client) {
	manager.delClientIdMap(client.ClientId)

	//删除所在的分组
	if len(client.GroupList) > 0 {
		for _, groupName := range client.GroupList {
			manager.delGroupClient(client.SystemId+":"+groupName, client.ClientId)
		}
	}

	// 删除系统里的客户端
	manager.delSystemClient(client)
}

// 删除clientIdMap
func (manager *ClientManager) delClientIdMap(clientId string) {
	manager.ClientIdMapLock.Lock()
	defer manager.ClientIdMapLock.Unlock()

	delete(manager.ClientIdMap, clientId)
}

// 通过clientId获取
func (manager *ClientManager) GetByClientId(clientId string) (*Client, error) {
	manager.ClientIdMapLock.RLock()
	defer manager.ClientIdMapLock.RUnlock()

	if client, ok := manager.ClientIdMap[clientId]; !ok {
		return nil, errors.New("客户端不存在")
	} else {
		return client, nil
	}
}

// 通过userId获取
func (manager *ClientManager) GetByUserId(userId string) ([]*Client, error) {
	var clients []*Client
	for _, client := range manager.AllClient() {
		if client.UserId == userId {
			clients = append(clients, client)
		}
	}
	if len(clients) > 0 {
		return clients, nil
	}
	return nil, errors.New("客户端不存在")
}

// 发送到本机分组
func (manager *ClientManager) SendMessage2LocalGroup(systemId, messageId, sendUserId, groupName string, code int, msg string, data *string) {
	if len(groupName) > 0 {
		clientIds := manager.GetGroupClientList(systemId + ":" + groupName)
		if len(clientIds) > 0 {
			for _, clientId := range clientIds {
				if _, err := Manager.GetByClientId(clientId); err == nil {
					//添加到本地
					SendMessage2LocalClient(messageId, clientId, sendUserId, code, msg, data)
				} else {
					//删除分组
					manager.delGroupClient(systemId+":"+groupName, clientId)
				}
			}
		}
	}
}

// 发送给指定业务系统
func (manager *ClientManager) SendMessage2LocalSystem(systemId, messageId string, sendUserId string, code int, msg string, data *string) {
	if len(systemId) > 0 {
		clientIds := Manager.GetSystemClientList(systemId)
		if len(clientIds) > 0 {
			for _, clientId := range clientIds {
				SendMessage2LocalClient(messageId, clientId, sendUserId, code, msg, data)
			}
		}
	}
}

// 添加到本地分组
func (manager *ClientManager) AddClient2LocalGroup(groupName string, client *Client, userId string, extend string) {
	//标记当前客户端的userId
	client.UserId = userId
	client.Extend = extend

	//判断之前是否有添加过
	for _, groupValue := range client.GroupList {
		if groupValue == groupName {
			return
		}
	}

	// 为属性添加分组信息
	groupKey := client.SystemId + ":" + groupName

	manager.addClient2Group(groupKey, client)

	client.GroupList = append(client.GroupList, groupName)

	mJson, _ := json.Marshal(map[string]string{
		"clientId": client.ClientId,
		"userId":   client.UserId,
		"extend":   client.Extend,
	})
	data := string(mJson)
	sendUserId := ""

	//发送系统通知
	SendMessage2Group(client.SystemId, sendUserId, groupName, 1001, "客户端上线", &data)
}

// 添加到本地分组
func (manager *ClientManager) addClient2Group(groupKey string, client *Client) {
	manager.GroupLock.Lock()
	defer manager.GroupLock.Unlock()
	manager.Groups[groupKey] = append(manager.Groups[groupKey], client.ClientId)
}

// 删除分组里的客户端
func (manager *ClientManager) delGroupClient(groupKey string, clientId string) {
	manager.GroupLock.Lock()
	defer manager.GroupLock.Unlock()

	for index, groupClientId := range manager.Groups[groupKey] {
		if groupClientId == clientId {
			manager.Groups[groupKey] = append(manager.Groups[groupKey][:index], manager.Groups[groupKey][index+1:]...)
		}
	}
}

// 获取本地分组的成员
func (manager *ClientManager) GetGroupClientList(groupKey string) []string {
	manager.GroupLock.RLock()
	defer manager.GroupLock.RUnlock()
	return manager.Groups[groupKey]
}

// 添加到系统客户端列表
func (manager *ClientManager) AddClient2SystemClient(systemId string, client *Client) {
	manager.SystemClientsLock.Lock()
	defer manager.SystemClientsLock.Unlock()
	manager.SystemClients[systemId] = append(manager.SystemClients[systemId], client.ClientId)
}

// 删除系统里的客户端
func (manager *ClientManager) delSystemClient(client *Client) {
	manager.SystemClientsLock.Lock()
	defer manager.SystemClientsLock.Unlock()

	for index, clientId := range manager.SystemClients[client.SystemId] {
		if clientId == client.ClientId {
			manager.SystemClients[client.SystemId] = append(manager.SystemClients[client.SystemId][:index], manager.SystemClients[client.SystemId][index+1:]...)
		}
	}
}

// 获取指定系统的客户端列表
func (manager *ClientManager) GetSystemClientList(systemId string) []string {
	manager.SystemClientsLock.RLock()
	defer manager.SystemClientsLock.RUnlock()
	return manager.SystemClients[systemId]
}

package sse

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	discord "github.com/bwmarrin/discordgo"
	"gopkg.in/antage/eventsource.v1"
)

type DiscordAction string

const (
	Begin  DiscordAction = "Begin"
	Update DiscordAction = "Update"
	End    DiscordAction = "End"
	Error  DiscordAction = "Error"
)

type DiscordActMessage struct {
	Message discord.Message
	Action  DiscordAction
}

var (
	// DiscordMessageCh = make(chan *DiscordActMessage, 1)

	// 并发安全的map，用于存储不同的msgCh
	MsgChManager = MsgChMap{
		Map: &sync.Map{},
	}
)

func SSE(w http.ResponseWriter, r *http.Request, ch chan *DiscordActMessage) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// 设置一个定时器
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	id := 0

	es := eventsource.New(nil, nil)
	defer es.Close()
	es.ServeHTTP(w, r)

	for {
		id++
		select {
		case <-ctx.Done():
			// 打印退出信息
			log.Println("client closed")
			return
		case <-ticker.C:
			// 每秒钟发送一次消息
			es.SendEventMessage(`{"msg":"wait"}`, "data", strconv.Itoa(id))
		case msg := <-ch:
			if msg.Action == End {
				es.SendEventMessage(fmt.Sprintf(`{"url":"%s","id":"%s","msgHash":"%s"}`,
					msg.Message.Attachments[0].URL, msg.Message.ID,
					getLastString(msg.Message.Attachments[0].Filename)), "data", strconv.Itoa(id))
				return
			}
			if msg.Action == Begin {
				es.SendEventMessage(`{"msg":"begin"}`, "data", strconv.Itoa(id))
			}
			if msg.Action == Update {
				es.SendEventMessage(fmt.Sprintf(`{"url":"%s","id":"%s"}`,
					msg.Message.Attachments[0].URL, msg.Message.ID), "data", strconv.Itoa(id))
			}
		}
	}
}

// 将字符串按_以及.分割，然后返回最后一个字符串
func getLastString(s string) string {
	split := strings.Split(s, "_")
	return strings.Split(split[len(split)-1], ".")[0]
}

type MsgChMap struct {
	*sync.Map
}

// 在map中增加一个msgCh，key为一个唯一id
func (m *MsgChMap) AddMsgCh(msgCh chan *DiscordActMessage) string {
	// id 为当前纳秒时间戳+随机数
	id := fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Intn(1000))
	m.Store(id, msgCh)
	return id
}

// 从map中删除一个msgCh
func (m *MsgChMap) DelMsgCh(id string) {
	m.Delete(id)
}

// 从map中获取一个msgCh
func (m *MsgChMap) GetMsgCh(id string) (chan *DiscordActMessage, bool) {
	msgCh, ok := m.Load(id)
	if !ok {
		return nil, false
	}
	return msgCh.(chan *DiscordActMessage), true
}

// 将信息进行包装，加上id(<!id:唯一标识符>)
func WrapMsg(msg string, id string) string {
	return fmt.Sprintf("<!id:%s> %s", id, msg)
}

// 将信息进行解包，返回id以及信息
func UnwrapMsg(msg string) (string, string, error) {
	// 通过正则表达式判断是否符合格式
	rul := `<!id:\d+> .+$`
	reg := regexp.MustCompile(rul)
	if !reg.MatchString(msg) {
		return "", "", fmt.Errorf("msg %s not match rule %s", msg, rul)
	}

	// 截断到第一个空格
	split := strings.SplitN(msg, " ", 2)
	// 去掉前面的<!id:和后面的>
	id := strings.TrimSuffix(strings.TrimPrefix(split[0], "**<!id:"), ">")
	return id, split[1], nil
}
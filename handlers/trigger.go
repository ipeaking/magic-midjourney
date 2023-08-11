package handlers

import (
	"errors"
	"wrap-midjourney/sse"

	"github.com/gin-gonic/gin"
)

type RequestTrigger struct {
	Type         string `json:"type"`
	DiscordMsgId string `json:"discordMsgId,omitempty"`
	MsgHash      string `json:"msgHash,omitempty"`
	Prompt       string `json:"prompt,omitempty"`
	Index        int64  `json:"index,omitempty"`
	SessionID    string `json:"sessionID,omitempty"`
}

func MidjourneyBot(c *gin.Context) {
	var body RequestTrigger
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ch := make(chan *sse.DiscordActMessage, 1)
	id := ""
	if body.SessionID != "" {
		sse.MsgChManager.AddMsgCh1ID(body.SessionID, ch)
		defer sse.MsgChManager.DelMsgCh(body.SessionID)
	} else {
		switch body.Type {
		case "generate":
		case "describe":
		default:
			c.JSON(400, gin.H{"error": "Must have sessionID"})
		}
		id = sse.MsgChManager.AddMsgCh(ch)
		defer sse.MsgChManager.DelMsgCh(id)
	}

	wrapPrompt := sse.WrapMsg(body.Prompt, id)
	var err error
	switch body.Type {
	case "generate":
		err = GenerateImage(wrapPrompt)
	case "upscale":
		err = ImageUpscale(body.Index, body.DiscordMsgId, body.MsgHash)
	case "variation":
		err = ImageVariation(body.Index, body.DiscordMsgId, body.MsgHash)
	case "maxUpscale":
		err = ImageMaxUpscale(body.DiscordMsgId, body.MsgHash)
	case "reset":
		err = ImageReset(body.DiscordMsgId, body.MsgHash)
	case "describe":
		err = ImageDescribe(wrapPrompt)
	default:
		err = errors.New("invalid type")
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sse.SSE(c.Writer, c.Request, ch, body.Type)
}

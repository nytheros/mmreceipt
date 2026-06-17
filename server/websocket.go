package main

import "github.com/mattermost/mattermost/server/public/model"

const (
	WSEventDelivered = "read_receipt_delivered"
	WSEventRead      = "read_receipt_read"
)

func (p *Plugin) publishReceiptEvent(event string, rec ReceiptRecord, senderID, channelID string) {
	p.API.PublishWebSocketEvent(event, map[string]interface{}{
		"post_id":    rec.PostID,
		"user_id":    rec.UserID,
		"status":     rec.Status,
		"updated_at": rec.UpdatedAt,
	}, &model.WebsocketBroadcast{UserId: senderID, ChannelId: channelID})
}

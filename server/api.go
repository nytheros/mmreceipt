package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type receiptRequest struct {
	PostID string `json:"post_id"`
}

type statusResponse struct {
	PostID   string          `json:"post_id"`
	Receipts []ReceiptRecord `json:"receipts"`
}

func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/api/v1/") {
		http.NotFound(w, r)
		return
	}
	if !p.getConfiguration().EnableReadReceipts {
		http.Error(w, "read receipts disabled", http.StatusForbidden)
		return
	}
	userID := r.Header.Get("Mattermost-User-Id")
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/delivered":
		p.handleReceipt(w, r, userID, StatusDelivered)
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/read":
		p.handleReceipt(w, r, userID, StatusRead)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v1/status/"):
		p.handleStatus(w, r, userID, strings.TrimPrefix(r.URL.Path, "/api/v1/status/"))
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleReceipt(w http.ResponseWriter, r *http.Request, userID, status string) {
	var req receiptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PostID == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	post, appErr := p.API.GetPost(req.PostID)
	if appErr != nil || post == nil {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}
	if post.UserId == userID {
		http.Error(w, "sender cannot acknowledge own post", http.StatusBadRequest)
		return
	}
	channel, appErr := p.API.GetChannel(post.ChannelId)
	if appErr != nil || channel == nil {
		http.Error(w, "channel not found", http.StatusNotFound)
		return
	}
	if !p.channelAllowed(channel.Type) {
		http.Error(w, "unsupported channel", http.StatusForbidden)
		return
	}
	if cm, appErr := p.API.GetChannelMember(post.ChannelId, userID); appErr != nil || cm == nil {
		http.Error(w, "not a channel member", http.StatusForbidden)
		return
	}
	rec, appErr := p.store.Upsert(req.PostID, userID, status, model.GetMillis())
	if appErr != nil {
		http.Error(w, appErr.Error(), http.StatusInternalServerError)
		return
	}
	if rec.Status == status {
		event := WSEventDelivered
		if status == StatusRead {
			event = WSEventRead
		}
		p.publishReceiptEvent(event, *rec, post.UserId, post.ChannelId)
	}
	writeJSON(w, rec)
}

func (p *Plugin) handleStatus(w http.ResponseWriter, r *http.Request, userID, postID string) {
	post, appErr := p.API.GetPost(postID)
	if appErr != nil || post == nil {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}
	if post.UserId != userID {
		http.Error(w, "only the sender can view receipts", http.StatusForbidden)
		return
	}
	channel, appErr := p.API.GetChannel(post.ChannelId)
	if appErr != nil || channel == nil || !p.channelAllowed(channel.Type) {
		http.Error(w, "unsupported channel", http.StatusForbidden)
		return
	}
	receipts, appErr := p.store.GetForPost(postID)
	if appErr != nil {
		http.Error(w, appErr.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, statusResponse{PostID: postID, Receipts: receipts})
}

func (p *Plugin) channelAllowed(t model.ChannelType) bool {
	cfg := p.getConfiguration()
	if t == model.ChannelTypeDirect {
		return cfg.EnableFor == "DM" || cfg.EnableFor == "DM_AND_GM"
	}
	if t == model.ChannelTypeGroup {
		return cfg.EnableFor == "GM" || cfg.EnableFor == "DM_AND_GM"
	}
	return false
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

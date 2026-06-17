package main

import (
	"sync"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

const manifestID = "readreceipt"

type configuration struct {
	EnableReadReceipts bool   `json:"EnableReadReceipts"`
	EnableFor          string `json:"EnableFor"`
}

type Plugin struct {
	plugin.MattermostPlugin
	store             *ReceiptStore
	configurationLock sync.RWMutex
	configuration     *configuration
}

func (p *Plugin) OnActivate() error { p.store = NewReceiptStore(p.API); return nil }

func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()
	if p.configuration == nil {
		return &configuration{EnableReadReceipts: false, EnableFor: "DM_AND_GM"}
	}
	return p.configuration
}

func (p *Plugin) OnConfigurationChange() error {
	cfg := &configuration{EnableReadReceipts: false, EnableFor: "DM_AND_GM"}
	if err := p.API.LoadPluginConfiguration(cfg); err != nil {
		return err
	}
	if cfg.EnableFor == "" {
		cfg.EnableFor = "DM_AND_GM"
	}
	p.configurationLock.Lock()
	p.configuration = cfg
	p.configurationLock.Unlock()
	return nil
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	cfg := p.getConfiguration()
	if !cfg.EnableReadReceipts {
		return
	}
	channel, appErr := p.API.GetChannel(post.ChannelId)
	if appErr != nil || channel == nil || !p.channelAllowed(channel.Type) {
		return
	}
	members, appErr := p.API.GetChannelMembers(post.ChannelId, 0, 100)
	if appErr != nil {
		return
	}
	now := model.GetMillis()
	for _, member := range members {
		if member.UserId == post.UserId {
			continue
		}
		rec, err := p.store.Upsert(post.Id, member.UserId, StatusDelivered, now)
		if err != nil {
			continue
		}
		if rec != nil && rec.Status == StatusDelivered {
			p.publishReceiptEvent(WSEventDelivered, *rec, post.UserId, post.ChannelId)
		}
	}
}

func main() { plugin.ClientMain(&Plugin{}) }

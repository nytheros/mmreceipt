package main

import (
	"sync"

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

func main() { plugin.ClientMain(&Plugin{}) }

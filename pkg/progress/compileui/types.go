package compileui

import (
	"time"

	"github.com/tensorchord/MIDI/pkg/editor/vscode"
)

type Action string

const (
	ActionStart Action = "start"
	ActionEnd   Action = "end"
)

type Result struct {
	plugins map[string]*PluginInfo
	ZSHInfo *ZSHInfo
}

type PluginInfo struct {
	vscode.Plugin
	startTime *time.Time
	endTime   *time.Time
	cached    bool
}

type ZSHInfo struct {
	OHMYZSH   string
	startTime *time.Time
	endTime   *time.Time
	cached    bool
}

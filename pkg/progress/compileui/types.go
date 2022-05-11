package compileui

import (
	"time"

	"github.com/tensorchord/envd/pkg/editor/vscode"
)

type Action string

const (
	ActionStart Action = "start"
	ActionEnd   Action = "end"
)

type Result struct {
	plugins []*PluginInfo
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

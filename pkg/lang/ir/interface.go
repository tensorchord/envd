package ir

import (
	"errors"

	"github.com/tensorchord/MIDI/pkg/vscode"
)

func Base(os, language string) {
	DefaultGraph.Language = language
	DefaultGraph.OS = os
}

func PyPIPackage(deps []string) {
	DefaultGraph.PyPIPackages = append(DefaultGraph.PyPIPackages, deps...)
}

func SystemPackage(deps []string) {
	DefaultGraph.SystemPackages = append(DefaultGraph.SystemPackages, deps...)
}

func CUDA(version, cudnn string) {
	DefaultGraph.CUDA = &version
	DefaultGraph.CUDNN = &cudnn
}

func VSCodePlugins(plugins []string) error {
	for _, p := range plugins {
		plugin, err := vscode.ParsePlugin(p)
		if err != nil {
			return err
		}
		DefaultGraph.VSCodePlugins = append(DefaultGraph.VSCodePlugins, plugin)
	}
	return nil
}

func UbuntuAPT(mode, source string) error {
	if source == "" {
		if mode == ubuntuAPTModeAuto {
			return errors.New("auto-mode not implemented")
		}
		return errors.New("source is required")
	}

	DefaultGraph.UbuntuAPTSource = &source
	return nil
}

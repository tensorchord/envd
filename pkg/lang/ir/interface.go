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

// UbuntuAPT updates the Ubuntu apt source.list in the image.
func UbuntuAPT(mode, source string) error {
	if source == "" {
		if mode == mirrorModeAuto {
			// If the mode is set to `auto`, MIDI detects the location of the run
			// then set to the nearest mirror
			return errors.New("auto-mode not implemented")
		}
		return errors.New("source is required")
	}

	DefaultGraph.UbuntuAPTSource = &source
	return nil
}

func PyPIMirror(mode, mirror string) error {
	if mirror == "" {
		if mode == mirrorModeAuto {
			// If the mode is set to `auto`, MIDI detects the location of the run
			// then set to the nearest mirror.
			return errors.New("auto-mode not implemented")
		}
		return errors.New("mirror is required")
	}

	DefaultGraph.PyPIMirror = &mirror
	return nil
}

func Shell(shell string) error {
	DefaultGraph.Shell = shell
	return nil
}

func Jupyter(pwd string, port int64) error {
	DefaultGraph.JupyterConfig = &JupyterConfig{
		Password: pwd,
		Port:     port,
	}
	return nil
}

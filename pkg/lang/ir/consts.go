package ir

const (
	osDefault       = "ubuntu20.04"
	languageDefault = "python3.8"
	mirrorModeAuto  = "auto"

	aptSourceFilePath  = "/etc/apt/sources.list"
	pypiMirrorFilePath = "/etc/pip.conf"

	pypiConfigTemplate = `
[global]
index-url=%s`
)

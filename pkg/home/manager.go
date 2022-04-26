package home

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	cacheDirName = "cache"
)

type Manager interface {
	HomeDir() string
	CacheDir() string
	ConfigFile() string
}

type generalManager struct {
	homeDir    string
	cacheDir   string
	configFile string

	logger *logrus.Entry
}

var (
	defaultManager *generalManager
	once           sync.Once
)

func Intialize(homeDir, configFile string) error {
	once.Do(func() {
		defaultManager = &generalManager{}
	})
	if err := defaultManager.init(homeDir, configFile); err != nil {
		return err
	}
	return nil
}

func GetManager() Manager {
	return defaultManager
}

func (m generalManager) CacheDir() string {
	return m.cacheDir
}

func (m generalManager) ConfigFile() string {
	return m.configFile
}

func (m generalManager) HomeDir() string {
	return m.homeDir
}

func (m *generalManager) init(homeDir, configFile string) error {
	expandedDir, err := expandHome(homeDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(expandedDir, 0755); err != nil {
		return err
	}
	m.homeDir = expandedDir

	m.cacheDir = filepath.Join(expandedDir, cacheDirName)

	expandedFilePath, err := expandHome(configFile)
	if err != nil {
		return err
	}

	_, err = os.Stat(expandedFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(configFile); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	m.configFile = expandedFilePath

	m.logger = logrus.WithFields(logrus.Fields{
		"homeDir":  m.homeDir,
		"cacheDir": m.cacheDir,
		"config":   m.configFile,
	})

	m.logger.Debug("home manager initialized")
	return nil
}

func expandHome(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, path[2:])
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

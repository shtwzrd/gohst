package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const Dir = "gohst"
const Hist = "hist"
const Cfg = "gohst.json"
const Key = "gohst.key"

const EnvUser = "GOHST_USER"
const EnvPassword = "GOHST_PASSWORD"
const EnvDomain = "GOHST_DOMAIN"
const EnvKey = "GOHST_KEY"

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
	Key      []byte `json:"key"`
}

func LoadConfig(dir string) Config {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Could not build directory hierarchy: %s", dir))
	}
	env := LoadEnv()
	cfg := LoadConfigFile(dir)
	key := LoadKeyFile(dir)
	dst := Config{}
	dst = env
	if cfg.Domain != "" {
		dst.Domain = cfg.Domain
	}
	if cfg.Username != "" {
		dst.Username = cfg.Username
	}
	if cfg.Password != "" {
		dst.Password = cfg.Password
	}
	if len(cfg.Key) > 0 {
		dst.Key = cfg.Key
	}
	if len(key) > 0 {
		dst.Key = key
	}
	return dst
}

func LoadKeyFile(dir string) []byte {
	// skip if not found
	if _, err := os.Stat(filepath.Join(dir, Key)); err != nil {
		return make([]byte, 0)
	}

	content, err := ioutil.ReadFile(filepath.Join(dir, Key))
	if err != nil {
		panic(fmt.Sprintf("Cannot read %s file in %s: %s", Key, dir, err))
	}
	return content
}

func LoadConfigFile(dir string) Config {
	// create if not found
	file := filepath.Join(dir, Cfg)
	if _, err := os.Stat(file); err != nil {
		empty, err := json.MarshalIndent(Config{}, "", "\t")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err2 := ioutil.WriteFile(file, empty, os.ModePerm)
		if err2 != nil {
			panic(fmt.Sprintf("Failed to create the missing config file: %s", err2))
		}
	}

	content, err := ioutil.ReadFile(filepath.Join(dir, Cfg))
	if err != nil {
		panic(fmt.Sprintf("Cannot read %s file in %s: %s", Cfg, dir, err))
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		panic(fmt.Sprintf("%s file contains invalid JSON: %s", Cfg, err))
	}
	return conf
}

func LoadEnv() Config {
	dst := Config{}
	dst.Username = os.Getenv(EnvUser)
	dst.Password = os.Getenv(EnvPassword)
	dst.Domain = os.Getenv(EnvDomain)
	dst.Key = []byte(os.Getenv(EnvKey))
	return dst
}

func DefaultDir() string {
	if runtime.GOOS == "linux" || runtime.GOOS == "freebsd" {
		return DefaultDirLinux()
	}
	if runtime.GOOS == "windows" {
		return DefaultDirWindows()
	}
	if runtime.GOOS == "darwin" {
		return DefaultDirWindows()
	}
	return DefaultDirLinux()
}

func DefaultDirLinux() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		return filepath.Join(xdg, Dir)
	}
	home := os.Getenv("HOME")
	return filepath.Join(home, ".config", Dir)
}

func DefaultDirWindows() string {
	appdata := os.Getenv("%APPDATA%")
	return filepath.Join(appdata, Dir)
}

func DefaultDirOSX() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, "Library", "Preferences", Dir)
}

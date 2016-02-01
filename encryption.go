package main

import (
	"github.com/jacobsa/crypto/siv"
	g "github.com/warreq/gohstd/common"
)

func EncryptInvocations(cmds g.Invocations, key []byte) (g.Invocations, error) {
	dst := make([]byte, 0)
	assoc := make([][]byte, 0)
	for i, c := range cmds {
		com, err := siv.Encrypt(dst, key, []byte(c.Command), assoc)
		dir, err := siv.Encrypt(dst, key, []byte(c.Directory), assoc)
		host, err := siv.Encrypt(dst, key, []byte(c.Host), assoc)
		shell, err := siv.Encrypt(dst, key, []byte(c.Shell), assoc)
		user, err := siv.Encrypt(dst, key, []byte(c.User), assoc)
		tags := make([]string, 0)
		for _, tag := range c.Tags {
			t, err := siv.Encrypt(dst, key, []byte(tag), assoc)
			if err != nil {
				return nil, err
			}
			tags = append(tags, string(t))
		}
		if err != nil {
			return nil, err
		}
		cmds[i].Command = string(com)
		cmds[i].Directory = string(dir)
		cmds[i].Host = string(host)
		cmds[i].Shell = string(shell)
		cmds[i].User = string(user)
		cmds[i].Tags = tags
	}
	return cmds, nil
}

func DecryptInvocations(cmds g.Invocations, key []byte) (g.Invocations, error) {
	assoc := make([][]byte, 0)
	for i, c := range cmds {
		com, err := siv.Decrypt(key, []byte(c.Command), assoc)
		dir, err := siv.Decrypt(key, []byte(c.Directory), assoc)
		host, err := siv.Decrypt(key, []byte(c.Host), assoc)
		shell, err := siv.Decrypt(key, []byte(c.Shell), assoc)
		user, err := siv.Decrypt(key, []byte(c.User), assoc)
		tags := make([]string, 0)
		for _, tag := range c.Tags {
			if len(tag) > 0 {
				t, err := siv.Decrypt(key, []byte(tag), assoc)
				if err != nil {
					return nil, err
				}
				tags = append(tags, string(t))
			}
		}
		if err != nil {
			return nil, err
		}
		cmds[i].Command = string(com)
		cmds[i].Directory = string(dir)
		cmds[i].Host = string(host)
		cmds[i].Shell = string(shell)
		cmds[i].User = string(user)
		cmds[i].Tags = tags
	}
	return cmds, nil
}

func DecryptCommands(cmds g.Commands, key []byte) (g.Commands, error) {
	assoc := make([][]byte, 0)
	for i, c := range cmds {
		com, err := siv.Decrypt(key, []byte(c), assoc)
		if err != nil {
			return nil, err
		}
		cmds[i] = g.Command(com)
	}
	return cmds, nil
}

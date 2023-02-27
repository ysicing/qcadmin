// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ssh

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
)

func (s *SSH) getClientConfig() *ssh.ClientConfig {
	if s.clientConfig == nil {
		auth := s.sshAuthMethod(s.Password, s.PkFile, s.PkData, s.PkPassword)
		config := ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		}
		defaultTimeout := time.Duration(15) * time.Second
		if s.Timeout <= 0 {
			s.Timeout = defaultTimeout
		}
		s.clientConfig = &ssh.ClientConfig{
			User:            s.User,
			Auth:            auth,
			Timeout:         s.Timeout,
			Config:          config,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint:gosec
		}
	}
	return s.clientConfig
}

// SSH connection operation
func (s *SSH) connect(host string) (*ssh.Client, error) {
	clientConfig := s.getClientConfig()
	ip, port := getSSHHostIPAndPort(host)
	addr := s.addrReformat(ip, port)
	return ssh.Dial("tcp", addr, clientConfig)
}

func newSession(client *ssh.Client) (*ssh.Session, error) {
	session, err := client.NewSession()
	if err != nil {
		_ = client.Close()
		return nil, err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     //disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		_ = session.Close()
		_ = client.Close()
		return nil, err
	}
	return session, nil
}

func (s *SSH) Connect(host string) (sshClient *ssh.Client, session *ssh.Session, err error) {
	try := 0
	if err := wait.ExponentialBackoff(defaultBackoff, func() (bool, error) {
		try++
		s.log.Debugf("the %d/%d time tring to ssh to %s with user %s", try, defaultBackoff.Steps, host, s.User)
		sshClient, session, err = s.newClientAndSession(host)
		if err != nil {
			return false, nil
		}
		return true, nil
	}); err != nil {
		return nil, nil, fmt.Errorf("ssh init dialer failed, %s %w", host, err)
	}
	return
}

func (s *SSH) newClientAndSession(host string) (*ssh.Client, *ssh.Session, error) {
	sshClient, err := s.connect(host)
	if err != nil {
		return nil, nil, err
	}
	session, err := newSession(sshClient)
	return sshClient, session, err
}

func (s *SSH) isLocalAction(host string) bool {
	return s.localAddress != nil && isLocalIP(host, s.localAddress)
}

func (s *SSH) sshAuthMethod(password, pkFile, pkData, pkPasswd string) (auth []ssh.AuthMethod) {
	if pkData != "" {
		signer, err := parsePrivateKey([]byte(pkData), []byte(pkPasswd))
		if err == nil {
			auth = append(auth, ssh.PublicKeys(signer))
		}
	}
	if password != "" {
		auth = append(auth, ssh.Password(password))
	} else {
		if pkFile == "" {
			pkFile = fmt.Sprintf("%s/.ssh/id_rsa", zos.GetHomeDir())
		}
		if file.CheckFileExists(pkFile) {
			signer, err := parsePrivateKeyFile(pkFile, pkPasswd)
			if err == nil {
				auth = append(auth, ssh.PublicKeys(signer))
			}
		}
	}
	return auth
}

func parsePrivateKey(pemBytes []byte, password []byte) (ssh.Signer, error) {
	if len(password) == 0 {
		return ssh.ParsePrivateKey(pemBytes)
	}
	return ssh.ParsePrivateKeyWithPassphrase(pemBytes, password)
}

func parsePrivateKeyFile(filename string, password string) (ssh.Signer, error) {
	pemBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file %v", err)
	}
	return parsePrivateKey(pemBytes, []byte(password))
}

func (s *SSH) addrReformat(host, port string) string {
	if !strings.Contains(host, ":") {
		host = fmt.Sprintf("%s:%s", host, port)
	}
	return host
}

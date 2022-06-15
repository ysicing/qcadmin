// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package kutil

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/ztime"
)

const (
	NodeToken = "/var/lib/rancher/k3s/server/node-token"
)

func GetNodeToken() string {
	b, err := ioutil.ReadFile(NodeToken)
	if err != nil {
		return ""
	}
	return string(b)
}

// NeedCacheHelmFile helm repo update
func NeedCacheHelmFile() bool {
	cachefile := fmt.Sprintf("%s/.566964cd0285e57cd088caa251ae863a.lock", common.DefaultCacheDir)
	if file.CheckFileExists(cachefile) {
		data, err := file.ReadAll(cachefile)
		if err != nil {
			return true
		}
		old, _ := ztime.TimeParse("200601021504", string(data))
		now := time.Now()
		if now.Sub(old) > time.Hour {
			return true
		}
		return false
	}
	file.Writefile(cachefile, ztime.NowUnixString())
	return true
}

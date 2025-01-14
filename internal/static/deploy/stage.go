// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package deploy

import (
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/common"
)

func Stage(dataDir string) error {
	for _, name := range AssetNames() {
		content, err := Asset(name)
		if err != nil {
			return err
		}
		p := filepath.Join(dataDir, name)
		os.MkdirAll(filepath.Dir(p), common.FileMode0755)
		if err := os.WriteFile(p, content, common.FileMode0755); err != nil {
			return errors.Wrapf(err, "failed to write to %s", name)
		}
	}
	return nil
}

// func StageFunc(dataDir string, templateVars map[string]string) error {
// 	log := log.GetInstance()
// 	log.Debugf("writing manifest: %s/hack/manifests/plugins", dataDir)
// 	for _, name := range AssetNames() {
// 		content, err := Asset(name)
// 		if err != nil {
// 			return err
// 		}
// 		for k, v := range templateVars {
// 			content = bytes.Replace(content, []byte(k), []byte(v), -1)
// 		}
// 		p := filepath.Join(dataDir, name)
// 		os.MkdirAll(filepath.Dir(p), 0700)

// 		if err := os.WriteFile(p, content, 0600); err != nil {
// 			return errors.Wrapf(err, "failed to write to %s", name)
// 		}
// 	}
// 	return nil
// }

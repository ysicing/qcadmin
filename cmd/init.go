// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"os"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/providers"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"

	// default provider
	_ "github.com/easysoft/qcadmin/internal/pkg/providers/incluster"
	_ "github.com/easysoft/qcadmin/internal/pkg/providers/native"

	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Run this command in order to set up the QuCheng control plane",
	}
	skip bool
)

func init() {
	initCmd.PersistentFlags().BoolVar(&skip, "skip-precheck", false, "skip precheck")
}

func newCmdInit(f factory.Factory) *cobra.Command {
	log := f.GetLog()

	name := "native"
	if file.CheckFileExists(common.GetDefaultKubeConfig()) {
		name = "incluster"
	}

	initCmd.PreRun = func(cmd *cobra.Command, args []string) {
		defaultArgs := os.Args
		if file.CheckFileExists(common.GetCustomConfig(common.InitFileName)) {
			log.Donef("quickon is already initialized, just run %s get cluster status", color.SGreen("%s status", defaultArgs[0]))
			os.Exit(0)
		}
		if name == "incluster" {
			// TODO Check k8s ready
		}
	}
	initCmd.Run = func(cmd *cobra.Command, args []string) {
		if name == "native" {

		}
		q
	}
	return initCmd
}

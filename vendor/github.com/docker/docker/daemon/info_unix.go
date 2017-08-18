// +build !windows

package daemon

import (
	"context"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/pkg/sysinfo"
)

// FillPlatformInfo fills the platform related info.
func (daemon *Daemon) FillPlatformInfo(v *types.Info, sysInfo *sysinfo.SysInfo) {
	v.MemoryLimit = sysInfo.MemoryLimit
	v.SwapLimit = sysInfo.SwapLimit
	v.KernelMemory = sysInfo.KernelMemory
	v.OomKillDisable = sysInfo.OomKillDisable
	v.CPUCfsPeriod = sysInfo.CPUCfsPeriod
	v.CPUCfsQuota = sysInfo.CPUCfsQuota
	v.CPUShares = sysInfo.CPUShares
	v.CPUSet = sysInfo.Cpuset
	v.Runtimes = daemon.configStore.GetAllRuntimes()
	v.DefaultRuntime = daemon.configStore.GetDefaultRuntimeName()
	v.InitBinary = daemon.configStore.GetInitPath()

	v.ContainerdCommit.Expected = dockerversion.ContainerdCommitID
	if sv, err := daemon.containerd.GetServerVersion(context.Background()); err == nil {
		v.ContainerdCommit.ID = sv.Revision
	} else {
		logrus.Warnf("failed to retrieve containerd version: %v", err)
		v.ContainerdCommit.ID = "N/A"
	}

	v.RuncCommit.Expected = dockerversion.RuncCommitID
	if rv, err := exec.Command(DefaultRuntimeBinary, "--version").Output(); err == nil {
		parts := strings.Split(strings.TrimSpace(string(rv)), "\n")
		if len(parts) == 3 {
			parts = strings.Split(parts[1], ": ")
			if len(parts) == 2 {
				v.RuncCommit.ID = strings.TrimSpace(parts[1])
			}
		}

		if v.RuncCommit.ID == "" {
			logrus.Warnf("failed to retrieve %s version: unknown output format: %s", DefaultRuntimeBinary, string(rv))
			v.RuncCommit.ID = "N/A"
		}
	} else {
		logrus.Warnf("failed to retrieve %s version: %v", DefaultRuntimeBinary, err)
		v.RuncCommit.ID = "N/A"
	}

	v.InitCommit.Expected = dockerversion.InitCommitID
	if rv, err := exec.Command(DefaultInitBinary, "--version").Output(); err == nil {
		// examples of how Tini outputs version info:
		//   "tini version 0.13.0 - git.949e6fa"
		//   "tini version 0.13.2"
		parts := strings.Split(strings.TrimSpace(string(rv)), " - ")

		v.InitCommit.ID = ""
		if v.InitCommit.ID == "" && len(parts) >= 2 {
			gitParts := strings.Split(parts[1], ".")
			if len(gitParts) == 2 && gitParts[0] == "git" {
				v.InitCommit.ID = gitParts[1]
				v.InitCommit.Expected = dockerversion.InitCommitID[0:len(v.InitCommit.ID)]
			}
		}
		if v.InitCommit.ID == "" && len(parts) >= 1 {
			vs := strings.TrimPrefix(parts[0], "tini version ")
			v.InitCommit.ID = "v" + vs
		}

		if v.InitCommit.ID == "" {
			logrus.Warnf("failed to retrieve %s version: unknown output format: %s", DefaultInitBinary, string(rv))
			v.InitCommit.ID = "N/A"
		}
	} else {
		logrus.Warnf("failed to retrieve %s version", DefaultInitBinary)
		v.InitCommit.ID = "N/A"
	}
}

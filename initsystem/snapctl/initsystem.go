package snapctl

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"sigs.k8s.io/etcdadm/apis"
	"sigs.k8s.io/etcdadm/service"
)

// InitSystem implements initsystem.InitSystem for snapctl.
type InitSystem struct {
	etcdAdmConfig *apis.EtcdAdmConfig
}

func New(etcdAdmConfig *apis.EtcdAdmConfig) *InitSystem {
	return &InitSystem{etcdAdmConfig}
}

func (s *InitSystem) Install() error {
	// No installation for snaps. Etcd is already installed as part of the snap package.
	return nil
}

func (s *InitSystem) Configure() error {
	if err := service.WriteEnvironmentFile(s.etcdAdmConfig); err != nil {
		return fmt.Errorf("failed writing environment file for etcd service: %w", err)
	}
	return nil
}

func (s *InitSystem) IsActive() (bool, error) {
	stdout, err := exec.Command("snapctl", "services", "microk8s.daemon-etcd").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve service status from snapctl: %w", err)
	}
	return bytes.Contains(stdout, []byte("active")), nil
}

func (s *InitSystem) EnableAndStartService() error {
	if err := exec.Command("snapctl", "start", "--enable", "microk8s.daemon-etcd").Run(); err != nil {
		return fmt.Errorf("failed to execute snapctl: %w", err)
	}
	return nil
}

func (s *InitSystem) DisableAndStopService() error {
	if err := exec.Command("snapctl", "stop", "--disable", "microk8s.daemon-etcd").Run(); err != nil {
		return fmt.Errorf("failed to execute snapctl: %w", err)
	}
	return nil
}

func (s *InitSystem) StartupTimeout() time.Duration {
	return 10 * time.Second
}

func (s *InitSystem) SetDefaults(cfg *apis.EtcdAdmConfig) {
	cfg.EnvironmentFile = filepath.Join(cfg.SnapDataDir, "etcd.env")
	cfg.EtcdctlEnvFile = filepath.Join(cfg.SnapDataDir, "etcdctl.env")
	cfg.EtcdctlShellWrapper = filepath.Join(cfg.SnapDataDir, "etcdctl.sh")
	cfg.CacheDir = filepath.Join(cfg.SnapDataDir, "cache", "etcd", fmt.Sprintf("v%s", cfg.Version))
}

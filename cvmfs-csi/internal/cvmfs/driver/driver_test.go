package driver

import (
	"fmt"
	"os"
	"testing"

	"github.com/kubernetes-csi/csi-test/v5/pkg/sanity"
	"golang.org/x/sync/errgroup"
)

func TestDriverSuite(t *testing.T) {
	socket := fmt.Sprintf("/var/lib/kubelet/plugins/%s/csi.sock", DefaultName)
	endpoint := "unix://" + socket
	if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove unix domain socket file %s, error: %s", socket, err)
	}

	driver := Driver{&Opts{
		DriverName:  DefaultName,
		CSIEndpoint: endpoint,
	}}

	var eg errgroup.Group
	eg.Go(func() error {
		return driver.Run()
	})

	cfg := sanity.NewTestConfig()
	if err := os.RemoveAll(cfg.TargetPath); err != nil {
		t.Fatalf("failed to delete target path %s: %s", cfg.TargetPath, err)
	}
	if err := os.RemoveAll(cfg.StagingPath); err != nil {
		t.Fatalf("failed to delete staging path %s: %s", cfg.StagingPath, err)
	}
	cfg.Address = endpoint
	cfg.IdempotentCount = 0
	sanity.Test(t, cfg)

	if err := eg.Wait(); err != nil {
		t.Errorf("driver run failed: %s", err)
	}
}

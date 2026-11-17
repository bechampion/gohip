package others

import (
	"os"
	"path/filepath"
	"testing"
)

type blockDev struct {
	uuid   string
	slaves []string
}

// fakeBlock builds a fake sysfs block tree and points sysBlock at it for the test.
func fakeBlock(t *testing.T, devs map[string]blockDev) {
	t.Helper()
	root := t.TempDir()
	for name, d := range devs {
		if d.uuid != "" {
			dir := filepath.Join(root, name, "dm")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "uuid"), []byte(d.uuid+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		for _, s := range d.slaves {
			if err := os.MkdirAll(filepath.Join(root, name, "slaves", s), 0o755); err != nil {
				t.Fatal(err)
			}
		}
	}
	old := sysBlock
	sysBlock = root + "/"
	t.Cleanup(func() { sysBlock = old })
}

func TestIsDmCrypt(t *testing.T) {
	fakeBlock(t, map[string]blockDev{
		"dm-0": {uuid: "CRYPT-LUKS2-abc-luks-abc", slaves: []string{"nvme0n1p3"}},
		"dm-1": {uuid: "LVM-xyz", slaves: []string{"dm-0"}},
		"dm-2": {uuid: "LVM-plain", slaves: []string{"sda1"}},
		"dm-3": {uuid: "LVM-mixed", slaves: []string{"dm-0", "sda1"}},
		"dm-4": {uuid: "CRYPT-LUKS2-def-luks-def", slaves: []string{"nvme0n1p4"}},
		"dm-5": {uuid: "LVM-spanned", slaves: []string{"dm-0", "dm-4"}},
	})

	for name, want := range map[string]bool{
		"dm-0": true,  // LUKS mapping itself
		"dm-1": true,  // LVM on top of LUKS
		"dm-2": false, // LVM on a plain partition
		"dm-3": false, // LVM spanning a LUKS and a plain device
		"dm-5": true,  // LVM spanning two LUKS devices
		"sda1": false, // plain partition
		"nope": false, // unknown device
	} {
		if got := isDmCrypt(name, 0); got != want {
			t.Errorf("isDmCrypt(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestEncStateMissingDevice(t *testing.T) {
	if got := encState("/dev/does-not-exist"); got != "unencrypted" {
		t.Errorf("encState of missing device = %q, want unencrypted", got)
	}
}

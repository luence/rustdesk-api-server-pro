package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSafeZipDestinationRejectsTraversal(t *testing.T) {
	dst := t.TempDir()

	if _, err := safeZipDestination(dst, "../evil.txt"); err == nil {
		t.Fatal("expected traversal path to be rejected")
	}
}

func TestSafeZipDestinationRejectsAbsolutePath(t *testing.T) {
	dst := t.TempDir()

	if _, err := safeZipDestination(dst, filepath.Join(string(os.PathSeparator), "tmp", "evil.txt")); err == nil {
		t.Fatal("expected absolute path to be rejected")
	}
}

func TestSafeZipDestinationAllowsNestedPath(t *testing.T) {
	dst := t.TempDir()

	got, err := safeZipDestination(dst, "dir/file.txt")
	if err != nil {
		t.Fatalf("expected nested path to be allowed: %v", err)
	}

	want := filepath.Join(dst, "dir", "file.txt")
	if got != want {
		t.Fatalf("unexpected destination: got %q want %q", got, want)
	}
}

func TestSafeModesRemoveGroupAndWorldWrite(t *testing.T) {
	if got := safeDirMode(0777); got != 0755 {
		t.Fatalf("unexpected dir mode: got %o want 0755", got)
	}
	if got := safeFileMode(0777); got != 0755 {
		t.Fatalf("unexpected executable file mode: got %o want 0755", got)
	}
	if got := safeFileMode(0666); got != 0644 {
		t.Fatalf("unexpected regular file mode: got %o want 0644", got)
	}
}

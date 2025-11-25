package manager

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMapArch(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"amd64": "amd64",
		"386":   "i386",
		"arm64": "aarch64",
		"arm":   "armv7l",
	}
	for input, expected := range tests {
		got, err := mapArch(input)
		if err != nil {
			t.Fatalf("mapArch(%s) error = %v", input, err)
		}
		if got != expected {
			t.Fatalf("mapArch(%s) = %s, want %s", input, got, expected)
		}
	}
	if _, err := mapArch("mips"); err == nil {
		t.Fatalf("expected error for unsupported arch")
	}
}

func TestUnzipArchiveAndFindBinary(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	archive := filepath.Join(tmpDir, "snell.zip")
	if err := createTestZip(archive); err != nil {
		t.Fatalf("create zip: %v", err)
	}

	extractDir := filepath.Join(tmpDir, "extract")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := unzipArchive(archive, extractDir); err != nil {
		t.Fatalf("unzipArchive() error = %v", err)
	}
	binary, err := findSnellBinary(extractDir)
	if err != nil {
		t.Fatalf("findSnellBinary() error = %v", err)
	}
	data, err := os.ReadFile(binary)
	if err != nil {
		t.Fatalf("read binary: %v", err)
	}
	if string(data) != "test" {
		t.Fatalf("unexpected binary content: %s", data)
	}
}

func TestSnellInstallerGetVersion(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip on windows")
	}
	dir := t.TempDir()
	binary := filepath.Join(dir, "snell-server")
	script := "#!/bin/sh\necho 'snell-server 5.0.1'\n"
	if err := os.WriteFile(binary, []byte(script), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	installer := NewSnellInstaller(binary, nil)
	version, err := installer.GetVersion(context.Background())
	if err != nil {
		t.Fatalf("GetVersion() error = %v", err)
	}
	if version != "snell-server 5.0.1" {
		t.Fatalf("unexpected version: %s", version)
	}
}

func createTestZip(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := zip.NewWriter(file)
	defer writer.Close()

	entry, err := writer.Create("snell-server-test")
	if err != nil {
		return err
	}
	if _, err := entry.Write([]byte("test")); err != nil {
		return err
	}
	return nil
}

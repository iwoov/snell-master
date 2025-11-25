package manager

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/iwoov/snell-master/backend/agent/internal/client"
	"github.com/iwoov/snell-master/backend/pkg/logger"
)

const (
	defaultDownloadTimeout = 5 * time.Minute
)

// SnellInstaller 负责 Snell Server 的下载与安装。
type SnellInstaller struct {
	BinaryPath   string
	masterClient *client.MasterClient
	httpClient   *http.Client
}

// NewSnellInstaller 创建 SnellInstaller。
func NewSnellInstaller(binaryPath string, masterClient *client.MasterClient) *SnellInstaller {
	return &SnellInstaller{
		BinaryPath:   binaryPath,
		masterClient: masterClient,
		httpClient: &http.Client{
			Timeout: defaultDownloadTimeout,
		},
	}
}

// IsInstalled 判断 Snell 二进制是否已存在。
func (s *SnellInstaller) IsInstalled() bool {
	info, err := os.Stat(s.BinaryPath)
	return err == nil && !info.IsDir()
}

// GetVersion 读取 Snell Server 版本。
func (s *SnellInstaller) GetVersion(ctx context.Context) (string, error) {
	if s.BinaryPath == "" {
		return "", fmt.Errorf("binary path not configured")
	}
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, s.BinaryPath, "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("snell version: %w - %s", err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

// DetectArch 返回当前运行环境对应的 Snell 架构名称。
func (s *SnellInstaller) DetectArch() (string, error) {
	return mapArch(runtime.GOARCH)
}

// Install 根据 Master 返回的配置下载安装 Snell。
func (s *SnellInstaller) Install(ctx context.Context, cfg *client.SnellConfig) error {
	if cfg == nil {
		return fmt.Errorf("snell config is nil")
	}
	arch, err := s.DetectArch()
	if err != nil {
		return err
	}

	downloadURL := cfg.DownloadURLs[arch]
	if downloadURL == "" {
		return fmt.Errorf("snell download url not found for arch %s", arch)
	}

	log := logger.WithModule("installer")
	log.Infof("Downloading Snell Server version %s for arch %s", cfg.Version, arch)

	tmpDir, err := os.MkdirTemp("", "snell-install-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, "snell.zip")
	if err := s.downloadFile(ctx, downloadURL, archivePath); err != nil {
		return fmt.Errorf("download snell: %w", err)
	}

	extractDir := filepath.Join(tmpDir, "extract")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("create extract dir: %w", err)
	}

	if err := unzipArchive(archivePath, extractDir); err != nil {
		return fmt.Errorf("unzip snell: %w", err)
	}

	binary, err := findSnellBinary(extractDir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.BinaryPath), 0o755); err != nil {
		return fmt.Errorf("create binary dir: %w", err)
	}

	if err := moveFile(binary, s.BinaryPath); err != nil {
		return fmt.Errorf("install binary: %w", err)
	}

	if err := os.Chmod(s.BinaryPath, 0o755); err != nil {
		return fmt.Errorf("chmod binary: %w", err)
	}

	version, err := s.GetVersion(ctx)
	if err != nil {
		log.Warnf("Snell installed but failed to read version: %v", err)
	} else {
		log.Infof("Snell installed successfully: %s", version)
	}
	return nil
}

func (s *SnellInstaller) downloadFile(ctx context.Context, url, dest string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), defaultDownloadTimeout)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("download request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: status %d", resp.StatusCode)
	}
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()
	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func unzipArchive(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer reader.Close()
	for _, file := range reader.File {
		targetPath := filepath.Join(dest, file.Name)
		if !strings.HasPrefix(targetPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid zip entry: %s", file.Name)
		}
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, file.Mode()); err != nil {
				return fmt.Errorf("mkdir entry: %w", err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return fmt.Errorf("mkdir file dir: %w", err)
		}
		dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("create entry file: %w", err)
		}
		rc, err := file.Open()
		if err != nil {
			dst.Close()
			return fmt.Errorf("open entry: %w", err)
		}
		if _, err := io.Copy(dst, rc); err != nil {
			dst.Close()
			rc.Close()
			return fmt.Errorf("copy entry: %w", err)
		}
		dst.Close()
		rc.Close()
	}
	return nil
}

var errBinaryLocated = errors.New("snell-binary-located")

func findSnellBinary(root string) (string, error) {
	var found string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), "snell-server") {
			found = path
			return errBinaryLocated
		}
		return nil
	})
	if err != nil && !errors.Is(err, errBinaryLocated) {
		return "", err
	}
	if found == "" {
		return "", fmt.Errorf("snell binary not found in archive")
	}
	return found, nil
}

func moveFile(src, dest string) error {
	if err := os.Rename(src, dest); err == nil {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return os.Remove(src)
}

func mapArch(goarch string) (string, error) {
	switch goarch {
	case "amd64":
		return "amd64", nil
	case "386":
		return "i386", nil
	case "arm64":
		return "aarch64", nil
	case "arm":
		return "armv7l", nil
	default:
		return "", fmt.Errorf("unsupported architecture: %s", goarch)
	}
}

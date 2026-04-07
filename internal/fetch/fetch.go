package fetch

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const repoAPI = "https://api.github.com/repos/Anthony-Maxwell1/BST-Core/releases/latest"

type Release struct {
	Assets []struct {
		BrowserDownloadURL string `json:"browser_download_url"`
		Name               string `json:"name"`
	} `json:"assets"`
}

func FetchLatest() error {
	// Determine OS zip name
	var osZip string
	switch runtime.GOOS {
	case "linux":
		osZip = "app-linux-x64.zip"
	case "windows":
		osZip = "app-win-x64.zip"
	case "darwin":
		osZip = "app-osx-x64.zip"
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	// Get latest release
	req, _ := http.NewRequest("GET", repoAPI, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return err
	}

	// Find the correct asset
	var url string
	for _, asset := range release.Assets {
		if asset.Name == osZip {
			url = asset.BrowserDownloadURL
			break
		}
	}

	if url == "" {
		return fmt.Errorf("%s not found in latest release", osZip)
	}

	// Download
	fmt.Println("Downloading:", url)
	assetResp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer assetResp.Body.Close()

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	coreDir := filepath.Join(exeDir, "core")
	os.MkdirAll(coreDir, 0755)
	zipPath := filepath.Join(coreDir, osZip)

	// Download file
	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, assetResp.Body)
	if err != nil {
		return err
	}

	// Extract
	fmt.Println("Extracting...")
	if err := unzip(zipPath, coreDir); err != nil {
		return err
	}
	os.Remove(zipPath)
	fmt.Println("Extraction complete")

	if runtime.GOOS != "windows" {
		fmt.Println("Setting permissions...")
		// Set permissions for non-Windows OS
		files, err := os.ReadDir(coreDir)
		if err != nil {
			return err
		}
		for _, file := range files {
			if !file.IsDir() {
				if err := perms(filepath.Join(coreDir, file.Name())); err != nil {
					return err
				}
			}
		}
		fmt.Println("Permissions set")
	}

	return nil
}

// permissions helper
func perms(path string) error {
	return os.Chmod(path, 0755)
}

// unzip helper
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.Create(fpath)
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

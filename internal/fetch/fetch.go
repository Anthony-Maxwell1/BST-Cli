package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const repoAPI = "https://api.github.com/repos/Anthony-Maxwell1/BST-Core/releases/latest"

type Release struct {
	Assets []struct {
		BrowserDownloadURL string `json:"browser_download_url"`
		Name               string `json:"name"`
	} `json:"assets"`
}

func FetchLatest() error {
	resp, err := http.Get(repoAPI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return err
	}

	if len(release.Assets) == 0 {
		return fmt.Errorf("no assets found")
	}

	url := release.Assets[0].BrowserDownloadURL
	fmt.Println("Downloading:", url)

	assetResp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer assetResp.Body.Close()

	os.MkdirAll("core", 0755)
	filePath := filepath.Join("core", release.Assets[0].Name)

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, assetResp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Saved to", filePath)
	return nil
}

package github

import (
	"encoding/json"
	"fmt"
	"rustdesk-api-server-pro/util"
)

type Release struct {
	ID              int     `json:"id"`
	NodeID          string  `json:"node_id"`
	TagName         string  `json:"tag_name"`
	TargetCommitish string  `json:"target_commitish"`
	Name            string  `json:"name"`
	Draft           bool    `json:"draft"`
	Prerelease      bool    `json:"prerelease"`
	CreatedAt       string  `json:"created_at"`
	PublishedAt     string  `json:"published_at"`
	Assets          []Asset `json:"assets"`
}

type Asset struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Label              string `json:"label"`
	ContentType        string `json:"content_type"`
	State              string `json:"state"`
	Size               int    `json:"size"`
	DownloadCount      int    `json:"download_count"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func GetReleases(repo string) (*[]Release, error) {
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/releases", repo)
	resp, err := util.HttpGetString(endpoint)
	if err != nil {
		return nil, fmt.Errorf("get releases request failed: %w", err)
	}
	releases := &[]Release{}
	if err = json.Unmarshal([]byte(resp), releases); err != nil {
		return nil, fmt.Errorf("decode releases response failed: %w", err)
	}
	return releases, nil
}

func GetLatestRelease(repo string) (*Release, error) {
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := util.HttpGetString(endpoint)
	if err != nil {
		return nil, fmt.Errorf("get latest release request failed: %w", err)
	}
	release := &Release{}
	if err = json.Unmarshal([]byte(resp), release); err != nil {
		return nil, fmt.Errorf("decode latest release response failed: %w", err)
	}
	return release, nil
}

func GetReleaseByTag(repo, tag string) (*Release, error) {
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", repo, tag)
	resp, err := util.HttpGetString(endpoint)
	if err != nil {
		return nil, fmt.Errorf("get release by tag request failed: %w", err)
	}
	release := &Release{}
	if err = json.Unmarshal([]byte(resp), release); err != nil {
		return nil, fmt.Errorf("decode release by tag response failed: %w", err)
	}
	return release, nil
}

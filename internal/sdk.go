package internal

import (
	"os"
	"path/filepath"
)

type SDKInfo struct {
	Name    string
	Version string
	Path    string
}

// DiscoverSDKs scans the sdkRoot for available SDKs and versions
func DiscoverSDKs(sdkRoot string) (map[string][]SDKInfo, error) {
	sdks := make(map[string][]SDKInfo)
	entries, err := os.ReadDir(sdkRoot)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			lang := entry.Name()
			langPath := filepath.Join(sdkRoot, lang)
			vers, _ := os.ReadDir(langPath)
			for _, v := range vers {
				if v.IsDir() {
					sdks[lang] = append(sdks[lang], SDKInfo{
						Name:    lang,
						Version: v.Name(),
						Path:    filepath.Join(langPath, v.Name()),
					})
				}
			}
		}
	}
	return sdks, nil
}

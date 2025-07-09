package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

// GenerateActivateBat creates activate.bat for CMD
func GenerateActivateBat(cfg map[string]string, sdkRoot, outDir string) error {
	bat := "@echo off\n"
	for lang, ver := range cfg {
		bat += fmt.Sprintf("set %s_HOME=%s\n", lang, filepath.Join(sdkRoot, lang, ver))
		bat += fmt.Sprintf("set PATH=%s;%s\n", filepath.Join(sdkRoot, lang, ver, "bin"), "%PATH%")
	}
	return os.WriteFile(filepath.Join(outDir, "activate.bat"), []byte(bat), 0644)
}

// GenerateActivatePs1 creates activate.ps1 for PowerShell
func GenerateActivatePs1(cfg map[string]string, sdkRoot, outDir string) error {
	ps1 := ""
	for lang, ver := range cfg {
		ps1 += fmt.Sprintf("$env:%s_HOME = '%s'\n", lang, filepath.Join(sdkRoot, lang, ver))
		ps1 += fmt.Sprintf("$env:PATH = '%s;' + $env:PATH\n", filepath.Join(sdkRoot, lang, ver, "bin"))
	}
	return os.WriteFile(filepath.Join(outDir, "activate.ps1"), []byte(ps1), 0644)
}

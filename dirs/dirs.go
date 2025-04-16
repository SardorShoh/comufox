package dirs

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

func RegistryDirectory() string {
	if envPath := os.Getenv("PLAYWRIGHT_BROWSERS_PATH"); envPath != "" {
		return envPath
	}
	if runtime.GOOS == "linux" {
		if envPath := os.Getenv("XDG_CACHE_HOME"); envPath != "" {
			return path.Join(envPath, "ms-playwright")
		} else {
			return path.Join(os.Getenv("HOME"), ".cache", "ms-playwright")
		}
	} else if runtime.GOOS == "darwin" {
		return path.Join(os.Getenv("HOME"), "Library", "Caches", "ms-playwright")
	} else if runtime.GOOS == "windows" {
		if envPath := os.Getenv("LOCALAPPDATA"); envPath != "" {
			return path.Join(envPath, "ms-playwright")
		} else {
			return path.Join(os.Getenv("HOME"), "AppData", "Local")
		}
	} else {
		panic(fmt.Sprintf("unsupported operating system: %s", runtime.GOOS))
	}
}

func GetExecutableName() string {
	ents, err := os.ReadDir(RegistryDirectory())
	if err != nil {
		panic(err)
	}
	var exePath string
	for _, ent := range ents {
		if ent.IsDir() {
			exePath = path.Join(RegistryDirectory(), ent.Name())
			break
		}
	}

	switch NormalizeOS(runtime.GOOS) {
	case "linux":
		return exePath + "/camoufox-bin"
	case "macos":
		return exePath + "/Camoufox.app/Contents/MacOS/camoufox"
	case "windows":
		return exePath + "\\camoufox.exe"
	default:
		// This should never be reached due to the check in normalizeOS
		return ""
	}
}

func RemoveOtherVersions(version string) error {
	ents, err := os.ReadDir(RegistryDirectory())
	if err != nil {
		return err
	}
	for _, ent := range ents {
		if ent.IsDir() && ent.Name() != version {
			return os.RemoveAll(path.Join(RegistryDirectory(), ent.Name()))
		}
	}
	return nil
}

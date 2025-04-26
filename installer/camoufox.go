package installer

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/SardorShoh/comufox/dirs"

	"github.com/artdarek/go-unzip"
	"github.com/cavaliergopher/grab/v3"
	"github.com/k0kubun/pp/v3"
	"github.com/tidwall/gjson"
	"resty.dev/v3"
)

var ExecPath string = ""

func convertByte(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func downloadCamoufox(dir, url string) error {
	client := grab.NewClient()
	req, err := grab.NewRequest(dir, url)
	if err != nil {
		return err
	}
	pp.Printf("Downloading %v...\n", url)
	resp := client.Do(req)
	pp.Printf("  %v\n", resp.HTTPResponse.Status)
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				convertByte(resp.BytesComplete()),
				convertByte(resp.Size()),
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		return err
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)
	return nil
}

// install Camoufox if not already installed
func InstallCamoufox() {
	client := resty.New()
	defer client.Close()
	resp, err := client.R().Get("https://api.github.com/repos/daijro/camoufox/releases/latest")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	result := gjson.ParseBytes(resp.Bytes())
	dir := dirs.RegistryDirectory()
	ExecPath = path.Join(dir, "camoufox-"+result.Get("tag_name").String())
	if _, err := os.Stat(ExecPath); os.IsNotExist(err) {
		if err := dirs.RemoveOtherVersions(result.Get("tag_name").String()); err != nil {
			panic(err)
		}
		if err := os.MkdirAll(ExecPath, 0750); err != nil {
			panic(fmt.Sprintf("could not create directory: %v", err))
		}

		var camoufoxZipFilename, downloadableURL string
		var platform, arch string
		switch runtime.GOOS {
		case "darwin":
			platform = "mac"
		case "linux":
			platform = "lin"
		case "windows":
			platform = "win"
		default:
			panic(fmt.Sprintf("unsupported operating system: %s", runtime.GOOS))
		}
		switch runtime.GOARCH {
		case "amd64":
			arch = "x86_64"
		case "arm64":
			arch = "arm64"
		case "386":
			arch = "i686"
		default:
			panic(fmt.Sprintf("unsupported architecture: %s", runtime.GOARCH))
		}
		camoufoxZipFilename = fmt.Sprintf(`camoufox-%s-%s.%s.zip`, strings.TrimPrefix(result.Get("tag_name").String(), "v"), platform, arch)
		for _, asset := range result.Get("assets").Array() {
			if asset.Get("name").String() == camoufoxZipFilename {
				downloadableURL = asset.Get("browser_download_url").String()
			}
		}
		log.Println("Installing camoufox from " + downloadableURL)
		log.Println("Into " + ExecPath)
		if err = downloadCamoufox(ExecPath, downloadableURL); err != nil {
			panic(fmt.Sprintf("could not download camoufox: %v", err))
		}
		uz := unzip.New(path.Join(ExecPath, camoufoxZipFilename), ExecPath)
		if err = uz.Extract(); err != nil {
			panic(fmt.Sprintf("could not unzip camoufox: %v", err))
		}
		os.Remove(path.Join(ExecPath, camoufoxZipFilename))

		// url = "https://github.com/plord12/webscrapers/releases/download/" + launchVer + "/" + launchZipFilename
		// log.Println("Installing launch from " + url)
		// log.Println("Into " + browserDirectory)
		// _, err = grab.Get(browserDirectory, url)
		// if err != nil {
		// 	panic(fmt.Sprintf("could not download launch: %v", err))
		// }
		// uz = unzip.New(path.Join(browserDirectory, launchZipFilename), browserDirectory)
		// err = uz.Extract()
		// if err != nil {
		// 	panic(fmt.Sprintf("could not unzip launch: %v", err))
		// }
		// os.Chmod(path.Join(browserDirectory, launchZipFilename), 0755)
		// os.Remove(path.Join(browserDirectory, launchZipFilename))
	}
}

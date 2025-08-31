package main

import "installer-bundler/core"

var items = map[string]string{
	"Chrome":  "https://dl.google.com/chrome/install/375.126/chrome_installer.exe",
	"Firefox": "https://download-installer.cdn.mozilla.net/pub/firefox/releases/113.0/win64/en-US/Firefox%20Setup%20113.0.exe",
	"VLC":     "https://get.videolan.org/vlc/3.0.18/win64/vlc-3.0.18-win64.exe",
	"7-Zip":   "https://www.7-zip.org/a/7z1900-x64.exe",
}

func main() {
	var coreItems []core.Item
	for title, link := range items {
		coreItems = append(coreItems, core.Item{
			Title: title,
			Link:  link,
		})
	}
	bundler := core.NewBundler(coreItems, "installer-runtime")
	projectDir, err := bundler.GenerateProject()
	if err != nil {
		panic(err)
	}

	installerPath, err := core.BuildProject(projectDir)
	if err != nil {
		panic(err)
	}

	println("Installer built at:", installerPath)
}

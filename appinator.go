package main

import (
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

const InfoPlistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>CFBundleDevelopmentRegion</key>
		<string>English</string>
		<key>CFBundleExecutable</key>
		<string>{{.Name}}</string>
		<key>CFBundleIconFile</key>
		<string>Icon</string>
		<key>CFBundleIdentifier</key>
		<string>{{.Package}}.{{.Name}}</string>
		<key>CFBundleInfoDictionaryVersion</key>
		<string>6.0</string>
		<key>CFBundlePackageType</key>
		<string>APPL</string>
		<key>CFBundleSignature</key>
		<string>{{.Contents}}</string>
		<key>CFBundleVersion</key>
		<string>1.0</string>
	</dict>
</plist>`

type appinator struct {
	Name       string `json:"-"`
	Contents   string `json:"-"`
	App        string `json:"app,omitempty"`
	Package    string `json:"package,omitempty"`
	Icon       string `json:"icon,omitempty"`
	Resources  string `json:"resources,omitempty"`
	Frameworks string `json:"frameworks,omitempty"`
	Docs       string `json:"docs,omitempty"`
	Debug      bool   `json:"debug,omitempty"`
}

func (a *appinator) Build() error {
	if a.Name == "" {
		a.Name = filepath.Base(a.App)
	}
	a.Contents = a.Name + "#"

	if a.Debug && runtime.GOOS == "darwin" {
		if err := os.MkdirAll(a.Name+".app/Contents/MacOS/", 0755); err != nil {
			return err
		} else if err = os.Symlink(a.App, a.Name+".app/Contents/MacOS/"+a.Name); err != nil {
			return err
		}
	} else if err := copy(a.App, a.Name+".app/Contents/MacOS/"+a.Name); err != nil {
		return err
	}

	if err := copy(a.Resources, a.Name+".app/Contents/Resources"); err != nil {
		return err
	} else if err = copy(a.Icon, a.Name+".app/Contents/Resources/Icon"); err != nil {
		return err
	} else if err = copy(a.Docs, a.Name+".app/Contents/Resources/Docs"); err != nil {
		return err
	} else if err = copy(a.Frameworks, a.Name+".app/Contents/Frameworks"); err != nil {
		return err
	}

	if InfoPlist, err := os.Create(a.Name + ".app/Contents/Info.plist"); err != nil {
		return err
	} else if InfoTemplate, err := template.New("Info.plist").Parse(InfoPlistTemplate); err != nil {
		return err
	} else if err := InfoTemplate.Execute(InfoPlist, a); err != nil {
		return err
	} else if err := InfoPlist.Close(); err != nil {
		return err
	}

	if PkgInfo, err := os.Create(a.Name + ".app/Contents/PkgInfo"); err != nil {
		return err
	} else if _, err := PkgInfo.WriteString(a.Contents); err != nil {
		return err
	} else if err := PkgInfo.Close(); err != nil {
		return err
	}

	return nil
}

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

func (self *appinator) Build() error {
	if self.Name == "" {
		self.Name = filepath.Base(self.App)
	}
	self.Contents = self.Name + "#"

	if self.Debug && runtime.GOOS == "darwin" {
		if err := os.MkdirAll(self.Name+".app/Contents/MacOS/", 0755); err != nil {
			return err
		} else if err = os.Symlink(self.App, self.Name+".app/Contents/MacOS/"+self.Name); err != nil {
			return err
		}
	} else if err := copy(self.App, self.Name+".app/Contents/MacOS/"+self.Name); err != nil {
		return err
	}

	if err := copy(self.Resources, self.Name+".app/Contents/Resources"); err != nil {
		return err
	} else if err = copy(self.Icon, self.Name+".app/Contents/Resources/Icon"); err != nil {
		return err
	} else if err = copy(self.Docs, self.Name+".app/Contents/Resources/Docs"); err != nil {
		return err
	} else if err = copy(self.Frameworks, self.Name+".app/Contents/Frameworks"); err != nil {
		return err
	}

	if InfoPlist, err := os.Create(self.Name + ".app/Contents/Info.plist"); err != nil {
		return err
	} else if InfoTemplate, err := template.New("Info.plist").Parse(InfoPlistTemplate); err != nil {
		return err
	} else if err := InfoTemplate.Execute(InfoPlist, self); err != nil {
		return err
	} else if err := InfoPlist.Close(); err != nil {
		return err
	}

	if PkgInfo, err := os.Create(self.Name + ".app/Contents/PkgInfo"); err != nil {
		return err
	} else if _, err := PkgInfo.WriteString(self.Contents); err != nil {
		return err
	} else if err := PkgInfo.Close(); err != nil {
		return err
	}

	return nil
}

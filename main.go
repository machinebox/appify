package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(2)
	}
}

func run() error {
	var (
		name       = flag.String("name", "My Go Application", "app name")
		author     = flag.String("author", "Appify by Machine Box", "author")
		version    = flag.String("version", "1.0", "app version")
		identifier = flag.String("id", "", "bundle identifier")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return errors.New("missing executable argument")
	}
	bin := args[0]
	appname := *name + ".app"
	contentspath := filepath.Join(appname, "Contents")
	apppath := filepath.Join(contentspath, "MacOS")
	binpath := filepath.Join(apppath, appname)
	if err := os.MkdirAll(apppath, 0777); err != nil {
		return errors.Wrap(err, "os.MkdirAll")
	}
	fdst, err := os.Create(binpath)
	if err != nil {
		return errors.Wrap(err, "create bin")
	}
	defer fdst.Close()
	fsrc, err := os.Open(bin)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New(bin + " not found")
		}
		return errors.Wrap(err, "os.Open")
	}
	defer fsrc.Close()
	if _, err := io.Copy(fdst, fsrc); err != nil {
		return errors.Wrap(err, "copy bin")
	}
	if err := exec.Command("chmod", "+x", apppath).Run(); err != nil {
		return errors.Wrap(err, "chmod: "+apppath)
	}
	if err := exec.Command("chmod", "+x", binpath).Run(); err != nil {
		return errors.Wrap(err, "chmod: "+binpath)
	}
	id := *identifier
	if id == "" {
		id = *author + "." + *name
	}
	info := infoListData{
		Name:               *name,
		Executable:         filepath.Join("MacOS", appname),
		Identifier:         id,
		Version:            *version,
		InfoString:         *name + " by " + *author,
		ShortVersionString: *version,
	}
	tpl, err := template.New("template").Parse(infoPlistTemplate)
	if err != nil {
		return errors.Wrap(err, "infoPlistTemplate")
	}
	fplist, err := os.Create(filepath.Join(contentspath, "Info.plist"))
	if err != nil {
		return errors.Wrap(err, "create Info.plist")
	}
	defer fplist.Close()
	if err := tpl.Execute(fplist, info); err != nil {
		return errors.Wrap(err, "execute Info.plist template")
	}
	if err := ioutil.WriteFile(filepath.Join(contentspath, "README"), []byte(readme), 0666); err != nil {
		return errors.Wrap(err, "ioutil.WriteFile")
	}
	return nil
}

type infoListData struct {
	Name               string
	Executable         string
	Identifier         string
	Version            string
	InfoString         string
	ShortVersionString string
}

const infoPlistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>CFBundlePackageType</key>
		<string>APPL</string>
		<key>CFBundleInfoDictionaryVersion</key>
		<string>6.0</string>
		<key>CFBundleName</key>
		<string>{{ .Name }}</string>
		<key>CFBundleExecutable</key>
		<string>{{ .Executable }}</string>
		<key>CFBundleIdentifier</key>
		<string>{{ .Identifier }}</string>
		<key>CFBundleVersion</key>
		<string>{{ .Version }}</string>
		<key>CFBundleGetInfoString</key>
		<string>{{ .InfoString }}</string>
		<key>CFBundleShortVersionString</key>
		<string>{{ .ShortVersionString }}</string>
	</dict>
</plist>
`

// readme goes into a README file inside the package for
// future reference.
const readme = `Made with Appify by Machine Box
https://github.com/machinebox/appify

Inspired by https://gist.github.com/anmoljagetia/d37da67b9d408b35ac753ce51e420132 
`

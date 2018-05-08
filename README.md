# appify

Create a macOS Application from an executable (like a Go binary)

![](appify-icon-small.png)

## Install

To install `appify`:

```bash
go get github.com/machinebox/appify
```

## Usage

```
appify -name "My Go Application" /path/to/bin
```

It will create a macOS Application:

![Output of appify is a mac application](preview.png)

## What next?

If you want to build a Go backed web based desktop application, check out our [machinebox/desktop](https://github.com/machinebox/desktop) project.

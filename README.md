# Switchboard
![GitHub Actions Status](https://github.com/Cian911/switchboard/workflows/Release/badge.svg) ![GitHub Actions Status](https://github.com/Cian911/switchboard/workflows/Test%20Suite/badge.svg)  [![Go Report Card](https://goreportcard.com/badge/github.com/cian911/switchboard)](https://goreportcard.com/report/github.com/cian911/switchboard) ![Downloads](https://img.shields.io/github/downloads/Cian911/switchboard/total.svg) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/Cian911/switchboard.svg)](https://github.com/Cian911/switchboard) [![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/cian911/switchboard) [![GitHub stars](https://badgen.net/github/stars/Cian911/switchboard)](https://GitHub.com/Cian911/switchboard/starazers/) [![GitHub forks](https://badgen.net/github/forks/Cian911/switchboard/)](https://GitHub.com/Cian911/switchboard/network/)

<p align="center">
  <img style="float: right;" src="examples/logo.png" alt="Gomerge logo"/>
</p>

### Description
Do you ever get annoyed that your Downloads folder gets cluttered with all types of files? Do you wish you could automatically organise them into seperate, organised folders? Switchboard is a tool to help simplfy file organization on your machine/s. 

Switchboard works by monitoring a directory you provide (or list of directories), and uses file system notifications to move a matched file to the appropriate destination directory of your choosing.

See the video below as example. Here, I give switchboard a path to watch, a destination where I want matched files to move to, and the file extension of the type of files I want to move.

[![asciicast](https://asciinema.org/a/OwbnYltbn0jcSAGzfdmujwklJ.svg)](https://asciinema.org/a/OwbnYltbn0jcSAGzfdmujwklJ)

### Installation

You can install switchboard pre-compiled binary in a number of ways.

##### Homebrew

```sh
brew tap Cian911/switchboard
brew install switchboard

// Check everything is working as it should be
switchboard -h
```

##### Go Install

```sh
go install github.com/Cian911/switchboard@latest
```

##### Manually

You can download the pre-compiled binary for your specific OS type from the [OSS releases page](https://github.com/Cian911/switchboard/releases). You will need to copy these and extract the binary, then move it to you local bin directory. See the example below.

```sh
wget https://github.com/Cian911/switchboard/releases/download/${VERSION}/${PACKAGE_NAME}
sudo tar -xvf ${PACKAGE_NAME} -C /usr/local/bin/
sudo chmod +x /usr/local/bin/switchboard
```

### Quick Start

Using switchboard is pretty easy. Below lists the set of commands and flags you can pass in.

```sh
Run the switchboard application passing in the path, destination, and file type you'd like to watch for.

Usage:
   watch [flags]

Flags:
      --config string        Pass an optional config file containing multiple paths to watch.
  -d, --destination string   Path you want files to be relocated.
  -e, --ext string           File type you want to watch for.
  -h, --help                 help for watch
  -p, --path string          Path you want to watch.
      --poll int             Specify a polling time in seconds. (default 3)
```

To get started quickly, you can run the following command, passing in the path, destination, and file extenstion you want to watch for. See the example below.

```sh
switchboard watch -p /home/user/Downloads -d /home/user/Movies -e .mp4
```

> We highly recommend using absolute file paths over relative file paths. Always include the `.` when passing the file extension to switchboard.

And that's it! Once ran, switchboard will start observing the user downloads folder for mp4 files added. Once it receives a new create event with the correct file extension, it will move the file to the users movies folder.

### Important Notes

##### Absolute File Path

As you might have noticed in the example above, we passed in the absolute file path. While relative file paths will work too, they have not been tested in all OS systems. Therefore we strongly recommend you use absolute file paths when running switchboard.

##### File Extenstion

You may have also noticed in the above example, we used `.mp4` including the prefixed `.`. This is important, as switchboard will not match file extenstions correctly if the given `--ext` flag does not contain the `.`.

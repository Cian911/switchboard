# Switchboard
![GitHub Actions Status](https://github.com/Cian911/switchboard/workflows/Release/badge.svg) ![GitHub Actions Status](https://github.com/Cian911/switchboard/workflows/Test%20Suite/badge.svg)  [![Go Report Card](https://goreportcard.com/badge/github.com/cian911/switchboard)](https://goreportcard.com/report/github.com/cian911/switchboard) ![Homebrew Downloads](https://img.shields.io/badge/dynamic/json?color=success&label=Downloads&query=count&url=https://github.com/Cian911/switchboard/blob/master/count.json?raw=True&logo=homebrew) ![Downloads](https://img.shields.io/github/downloads/Cian911/switchboard/total.svg) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/Cian911/switchboard.svg)](https://github.com/Cian911/switchboard) [![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/cian911/switchboard) [![GitHub stars](https://badgen.net/github/stars/Cian911/switchboard)](https://GitHub.com/Cian911/switchboard/starazers/) [![GitHub forks](https://badgen.net/github/forks/Cian911/switchboard/)](https://GitHub.com/Cian911/switchboard/network/)

<p align="center">
  <img style="float: right;width:400px;height:400px;" src="examples/logo.png" alt="Gomerge logo"/>
</p>

### Description
Do you ever get annoyed that your Downloads folder gets cluttered with all types of files? Do you wish you could automatically organise them into seperate, organised folders? Switchboard is a tool to help simplfy file organization on your machine/s. 

Switchboard works by monitoring a directory you provide (or list of directories), and uses file system notifications to move a matched file to the appropriate destination directory of your choosing.

See the video below as example. Here, I give switchboard a path to watch, a destination where I want matched files to move to, and the file extension of the type of files I want to move.

### Pro

As of version `v1.0.0` we have released a pro version which has a ton more features and functionality. Head over to [https://goswitchboard.io/pro](https://goswitchboard.io/pro) for more info.

**Switchboard Pro** gives you extra features and support over the free open-source version.

Purchasing a **pro** or **enterprise** license for **Switchboard Pro** helps us to continue working on both the pro and free version of the software, and bring more features to **_YOU_**!

- [x] Support for **prioritising specific file events** over others.
- [x] **Regex support** so you can watch for any file name or type you choose.
- [x] Support for archival file extractions, **.zip/.rar et al**.
- [x] Support for **optional file removal**.
- [x] Product support should you run into any issues.
- [x] Access to product roadmap.
- [x] Priority feature requests.

---

[![asciicast](https://asciinema.org/a/OwbnYltbn0jcSAGzfdmujwklJ.svg)](https://asciinema.org/a/OwbnYltbn0jcSAGzfdmujwklJ)

You can also visit https://goswitchboard.io/ for all your documentation needs, news, and updates!


### Installation

You can install switchboard pre-compiled binary in a number of ways.

##### Homebrew

```sh
brew tap Cian911/switchboard
brew install switchboard

// Check everything is working as it should be
switchboard -h
```

You can also upgrade the version of `switchboard` you already have installed by doing the following.

```sh
brew upgrade switchboard
```

##### Docker

```sh
docker pull ghcr.io/cian911/switchboard:${VERSION}

docker run -d -v ${SRC} -v ${DEST} ghcr.io/cian911/switchboard:${VERSION} watch -h
```

##### Go Install

```sh
go install github.com/Cian911/switchboard@${VERSION}
```

##### Manually

You can download the pre-compiled binary for your specific OS type from the [OSS releases page](https://github.com/Cian911/switchboard/releases). You will need to copy these and extract the binary, then move it to you local bin directory. See the example below for extracting a zipped version.

```sh
curl https://github.com/Cian911/switchboard/releases/download/${VERSION}/${PACKAGE_NAME} -o ${PACKAGE_NAME}
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
      --poll int             Specify a polling time in seconds. (default 60)
  -r, --regex-pattern string Pass a regex pattern to watch for any files matching this pattern.
```

To get started quickly, you can run the following command, passing in the path, destination, and file extenstion you want to watch for. See the example below.

```sh
switchboard watch -p /home/user/Downloads -d /home/user/Movies -e .mp4
```

> We highly recommend using absolute file paths over relative file paths. Always include the `.` when passing the file extension to switchboard.

And that's it! Once ran, switchboard will start observing the user downloads folder for mp4 files added. Once it receives a new create event with the correct file extension, it will move the file to the users movies folder.

### Important Notes

##### Polling

We set a high polling time on switchboard as in some operating systems we don't get file closed notifications. Therefore switchboard implements a polling solution to check for when a file was last written to. If the file falls outside the time since last polled, the file is assumed to be closed and will be moved to the destination directory. This obviously is not ideal, as we can't guarentee that a file is _actually_ closed. Therefore the option is there to set the polling interval yourself. In some cases, a higher polling time might be necessary.

##### Polling & Linux
As of release `v1.0.0` we now support `IN_CLOSE_WRITE` events in _linux_ systems. For context, this event tells us when a process has finished writing to a file (something we don't get on OSX & Windows). This means we do not need to use polling for linux systems (though we do for _some_ circumstances) however the functionaity still exists should you wish to use it.  

##### Absolute File Path

As you might have noticed in the example above, we passed in the absolute file path. While relative file paths will work too, they have not been tested in all OS systems. Therefore we strongly recommend you use absolute file paths when running switchboard.

##### File Extenstion

You may have also noticed in the above example, we used `.mp4` including the prefixed `.`. This is important, as switchboard will not match file extenstions correctly if the given `--ext` flag does not contain the `.`.

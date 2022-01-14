# launchctl

Below outlines an example for MacOS users using `launchctl`. launchctl will ensure switchboard is started and left running in the background, so you can get on with what you need to do.

Create the following file `sudo touch /Library/LaunchDaemons/switchboard.plist` and copy and paste the contents below. Note, you should make any changes necessary to the arguments list, depending on how you want to use the tool.

```bash
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN"
    "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>switchboard</string>
    <key>ServiceDescription</key>
    <string>File system watcher</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/switchboard</string>
        <string>watch</string>
        <string>--config</string>
        <string>/path_to_your_config_file/config.yaml</string>
    </array>
    <key>RunAtLoad</key>
    <false/>
</dict>
</plist>
```

Then, run the following commands to load, start, and list the running service and ensuring it has started.

```bash
sudo launchctl load /Library/LaunchDaemons/switchboard.plist
sudo launchctl start switchboard

sudo launchctl list | grep switchboard
```

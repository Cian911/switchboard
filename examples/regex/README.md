### Regex

As of `v1.0.0` switchboard now supports regex patterns.

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

Below is an example of using regex patterns with switchboard in your `config.yaml` file.

```yaml
pollingInterval: 10
watchers:
  - path: "/home/cian/Downloads"
    destination: "/home/cian/Documents"
    ext: ".txt"
  - path: "/home/cian/Downloads"
    destination: "/home/cian/Documents/Reports"
    ext: ".txt"
    pattern: "(?i)(financial-report-[a-z]+-[0-9]+.txt)"
  - path: "/home/cian/Downloads"
    destination: "/home/cian/Videos"
    ext: ".mp4"
```

Or you can pass a regex pattern via the cli like so.

```bash
switchboard watch -p /home/john/Downloads -d /home/john/Documents -e .csv -r "(?i)(financial-report-[a-z]+-[0-9]+.txt)"
```



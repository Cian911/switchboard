# Yaml Config

Switchboard can also with a config.yaml file. See below for an example configuration.

```bash
watchers:
  - path: "/home/user/input"
    destination: "/home/user/output"
    ext: ".txt"
  - path: "/home/user/downloads"
    destination: "/home/user/movies"
    ext: ".mp4"
```

With the content above, create a `config.yaml` file. You can then pass it as a flag to switchboard like so:

```bash
switchboard watch --config config.yaml

###
Using config file: yaml/config.yaml
2022/01/04 22:53:15 Observing
```

# Counter-Strike 2 Plugin Manager (cs2pm)
Counter-Strike 2 Plugin Manager is a tool for easily installing and uninstalling CS2 Server Metamod (and CounterStrikeSharp) Plugins.

> **NOTE: cs2pm is in very early stages of development. Once I am finished making it fit my own use case, I plan to expand on it to have functionality similar to that of [Scoop](https://scoop.sh/), but for CS2 server plugins!**

## Usage

### Create a cs2pm.yaml manifest file
```yaml
plugins:
  - name: CS2-SimpleAdmin
    description: Simple Admin Plugin
    downloadUrl: https://github.com/daffyyyy/CS2-SimpleAdmin/releases/download/build-230/CS2-SimpleAdmin.zip
    extractPrefix: addons/counterstrikesharp/plugins
    uninstall:
      directories:
        - addons/counterstrikesharp/plugins/CS2-SimpleAdmin
```

### Install plugins from the manifest file
```bash
cs2pm install
```

### Uninstall plugins from the manifest file
```bash
cs2pm uninstall
```

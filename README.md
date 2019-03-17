# ghstatus

> A simple tool to retrieve and update your GitHub status from your terminal

```bash
# Get your current status
$ ghstatus
Current status: 🐋  Containerizing and shipping all the things
```

```bash
# Update your status
$ ghstatus set :whale: Containerizing all the things!
🎉 Updated your status!

# Or set yourself as busy
$ ghstatus --busy set :sleeping: Taking a nap
🎉 Updated your status!
```

## installation

Head over to the [releases page](https://github.com/BrunoScheufler/ghstatus/releases) and get your ghstatus binary!

Once you've downloaded a binary and made it executable, run it once and head over to the newly-generated config file (located in `~/.config/ghstatus/config.json` by default, you can change this behaviour by setting the `--config` flag to point to your file location of choice) and enter your GitHub access token. The token has to be permitted to access the following scopes: `user`, `read:org`.

## additional notes

All of the included functionality is achieved by using a few great packages, such as

- [mapstructure](https://github.com/mitchellh/mapstructure) for handling JSON -> map -> struct workflows
- [aurora](https://github.com/logrusorgru/aurora) for all of the terminal colors
- [emoji](https://github.com/kyokomi/emoji) for the :sparkles: emoji magic


## license

This project is licensed under the [MIT License](LICENSE.md)

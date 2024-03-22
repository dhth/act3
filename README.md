# act3

✨ Overview
---

Glance at the results of the last 3 runs of your Github Actions.

<p align="center">
  <img src="./act3.png" alt="Usage" />
</p>

💾 Installation
---

**go**:

```sh
go install github.com/dhth/act3@latest
```

⚡️ Usage
---

Create a config file that looks like the following (the default location `act3`)
will look for this file is `~/.config/act3/act3.yml.`

```yaml
workflows:
- id: ABC
  repo: dhth/outtasync
  name: release
- id: XYZ
  repo: dhth/ecsv
  name: release
  key: key-will-supersede repo/name in the output
- id: EFG
  repo: dhth/cueitup
  name: release
```

You can find the ID for your workflow as follows:

```bash
curl -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <YOUR_GH_TOKEN>" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/OWNER/REPO/actions/workflows/<WORKFLOW_FILE>

# use node_id from the response
```

```bash
ACT3_GH_ACCESS_TOKEN="<YOUR_GH_TOKEN> \
act3"
```

Acknowledgements
---

`act3` is built using the TUI framework [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea

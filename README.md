# act3

✨ Overview
---

Glance at the results of the last 3 runs of your Github Actions.

<p align="center">
  <img src="./act3.png" alt="Usage" />
</p>

💾 Installation
---

**homebrew**:

```sh
brew install dhth/tap/act3
```

**go**:

```sh
go install github.com/dhth/act3@latest
```

🛠️ Configuration
---

Create a config file that looks like the following (`act3` will look for this
file at `~/.config/act3/act3.yml.` by default).

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

⚡️ Usage

### CLI output

```bash
ACT3_GH_ACCESS_TOKEN="<YOUR_GH_TOKEN> \
act3"
```

### HTML output

`act3` can also output the results in HTML format.


```bash
ACT3_GH_ACCESS_TOKEN="<YOUR_GH_TOKEN> \
    act3" \
    -config-file=./examples/html/act3.yml \
    -format=html \
    -html-template-file=./examples/html/template.html
```

A sample page generated via `act3` is running at [https://dhth.github.io/act3](https://dhth.github.io/act3), the source code for which is in the [examples/html](./examples/html) directory.

Acknowledgements
---

`act3` is built using the TUI framework [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea

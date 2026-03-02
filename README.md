<p align="center">
  <h1 align="center">act3</h1>
  <p align="center">
    <a href="https://github.com/dhth/act3/actions/workflows/main.yml"><img alt="Build Status" src="https://img.shields.io/github/actions/workflow/status/dhth/act3/main.yml?style=flat-square"></a>
    <a href="https://github.com/dhth/act3/actions/workflows/vulncheck.yml"><img alt="Vulnerability Check" src="https://img.shields.io/github/actions/workflow/status/dhth/act3/vulncheck.yml?style=flat-square&label=vulncheck"></a>
    <a href="https://github.com/dhth/act3/releases/latest"><img alt="Latest release" src="https://img.shields.io/github/release/dhth/act3.svg?style=flat-square"></a>
    <a href="https://github.com/dhth/act3/releases/latest"><img alt="Commits since latest release" src="https://img.shields.io/github/commits-since/dhth/act3/latest?style=flat-square"></a>
  </p>
</p>

Glance at the results of the last 3 runs of your Github Actions.

> Want to see a demo before you read the rest of the documentation?
> View `act3` in action [here][3].

[![usage](https://asciinema.org/a/BS5LAEWzxV81uXJG3tiFcGULW.svg)](https://asciinema.org/a/BS5LAEWzxV81uXJG3tiFcGULW)

💾 Installation
---

**homebrew**:

```sh
brew install dhth/tap/act3
```

**[X-CMD](https://x-cmd.com)**

```sh
x install act3
```

**go**:

```sh
go install github.com/dhth/act3@latest
```

Or get the binary directly from a [release][2]. Read more about verifying the
authenticity of released artifacts [here](#-verifying-release-artifacts).

🔑 Authentication
---

You can have `act3` make authenticated calls to GitHub on your behalf in either
of two ways:

- Have an authenticated instance of [gh](https://github.com/cli/cli) available
    (recommended).
- Provide a valid Github token via `$GH_TOKEN`, which has the following
    permissions for the repos you want to query data for.
    - `actions:read`
    - `checks:read`

⚡️ Usage
---

### Basic Usage

```text
Usage:
  act3 [flags]
  act3 [command]

Available Commands:
  config      Interact with act3's config
  help        Help about any command

Flags:
  -c, --config-path string            location of act3's config file (default "/Users/user/Library/Application Support/act3/act3.yml")
  -f, --format string                 output format to use; possible values: default, table, html (default "default")
  -g, --global                        whether to use workflows defined globally via the config file
  -h, --help                          help for act3
      --html-template-path string     path of the HTML template file to use
      --html-title string             title to use in the HTML output (default "act3")
  -o, --open-failed                   whether to open failed workflows
  -r, --repos strings                 repos to fetch workflows for, in the format "owner/repo"
  -n, --workflow-name-filter string   regex expression to filter workflows by name
```

By default, `act3` will show results for the repository associated with the
current directory. Simply run `act3` from the project root.

You can also specify a list of repositories to fetch results for using the `-r`
flag.

```bash
act3 -r dhth/act3,dhth/bmm
```

[![run-on-configured-workflows](https://asciinema.org/a/IXj8L58ILEDawB7NhF5A2uKlJ.svg)](https://asciinema.org/a/IXj8L58ILEDawB7NhF5A2uKlJ)

### Specific Workflows

You can also fetch results for specific workflows using a config file, that
looks like the following. Run `act3 -h` to view the default location where
`act3` looks for this config file.

```yaml
workflows:
- id: W_kwDOLkC0eM4FaKV_
  repo: dhth/act3
  name: build
- id: W_kwDOLb3Pms4FRxjX
  repo: dhth/cueitup
  name: build
  key: cueitup-release  # optional
- id: W_kwDOLb3Pms4FRxjY
  repo: dhth/cueitup
  name: release
  url: https://asampleurl.com/{{runNumber}} # optional
```

`{{runNumber}}` gets replaced with the actual run number of the workflow.

You can generate this configuration using `act3` itself.

```bash
Usage:
  act3 config gen [flags]

Flags:
  -h, --help                          help for gen
  -r, --repos strings                 repos to generate the config for, in the format "owner/repo"
  -n, --workflow-name-filter string   regex expression to filter workflows by name
```

[![generate-own-config](https://asciinema.org/a/FoiwDy42w4HjV94E5xRa9Kav9.svg)](https://asciinema.org/a/FoiwDy42w4HjV94E5xRa9Kav9)

### Tabular output

`act3` can also output results in a tabular format.

```bash
act3 -f table
```

<p align="center">
  <img src="https://tools.dhruvs.space/images/act3/act3-2.png" alt="Usage" />
</p>

### HTML output

`act3` can also output results in HTML format. You can also specify a template
using the `-t` flag (refer to
[./internal/ui/assets/template.html](./internal/ui/assets/template.html) for the
default template.)

```bash
act3 -f html
```

The resultant HTML page looks like this.

<p align="center">
  <img src="https://tools.dhruvs.space/images/act3/act3-html-1.png" alt="Usage" />
</p>

> You can see this in action [here][3].

🔐 Verifying release artifacts
---

In case you get the `act3` binary directly from a [release][2], you may want to
verify its authenticity. Checksums are applied to all released artifacts, and
the resulting checksum file is signed using
[cosign](https://docs.sigstore.dev/cosign/installation/).

Steps to verify (replace `x.y.z` in the commands listed below with the version
you want):

1. Download the following files from the release:

   - act3_x.y.z_checksums.txt
   - act3_x.y.z_checksums.txt.pem
   - act3_x.y.z_checksums.txt.sig

2. Verify the signature:

   ```shell
   cosign verify-blob act3_x.y.z_checksums.txt \
       --certificate act3_x.y.z_checksums.txt.pem \
       --signature act3_x.y.z_checksums.txt.sig \
       --certificate-identity-regexp 'https://github\.com/dhth/act3/\.github/workflows/.+' \
       --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
   ```

3. Download the compressed archive you want, and validate its checksum:

   ```shell
   curl -sSLO https://github.com/dhth/act3/releases/download/vx.y.z/act3_x.y.z_linux_amd64.tar.gz
   sha256sum --ignore-missing -c act3_x.y.z_checksums.txt
   ```

3. If checksum validation goes through, uncompress the archive:

   ```shell
   tar -xzf act3_x.y.z_linux_amd64.tar.gz
   ./act3
   # profit!
   ```

[2]: https://github.com/dhth/act3/releases
[3]: https://dhth.github.io/act3-runner

# act3

[![Build Workflow Status](https://img.shields.io/github/actions/workflow/status/dhth/act3/build.yml?style=flat-square)](https://github.com/dhth/act3/actions/workflows/build.yml)
[![Vulncheck Workflow Status](https://img.shields.io/github/actions/workflow/status/dhth/act3/vulncheck.yml?style=flat-square&label=vulncheck)](https://github.com/dhth/act3/actions/workflows/vulncheck.yml)
[![Latest Release](https://img.shields.io/github/release/dhth/act3.svg?style=flat-square)](https://github.com/dhth/act3/releases/latest)
[![Commits Since Latest Release](https://img.shields.io/github/commits-since/dhth/act3/latest?style=flat-square)](https://github.com/dhth/act3/releases)

Glance at the results of the last 3 runs of your Github Actions.

<p align="center">
  <img src="https://tools.dhruvs.space/images/act3/act3-1.png" alt="Usage" />
</p>

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

Flags:
  -c string                       path of the config file (default "/Users/user/.config/act3/act3.yml")
  -f string                       output format to use; possible values: default, table, html (default "default")
  -t string                       path of the HTML template file to use
  -r string                       repo to fetch workflows for, in the format "owner/repo"
  -g bool                         whether to use workflows defined globally via the config file (default false)
  -o bool                         whether to open failed workflows (via your OS's "open" command) (default false)
  -h, --help                      help for act3
```

By default, `act3` will show results for the repository associated with the
current directory. Simply run `act3` from the project root.

You can also specify a repository to fetch results for using the `-r` flag.

```bash
act3 -r neovim/neovim
```

### Specific Workflows

You can also fetch results for specific workflows using a config file, that
looks like the following.

```yaml
workflows:

- id: W_kwDOLkC0eM4FaKV_
  repo: dhth/act3
  name: build
  url: https://asampleurl.com/{{runNumber}}
- id: W_kwDOLkC0eM4FaKWA
  repo: dhth/act3
  name: release
  url: https://asampleurl.com/{{runNumber}}

- id: W_kwDOLb3Pms4FRxjX
  repo: dhth/cueitup
  name: build
  url: https://dhth.github.io/cueitup
- id: W_kwDOLb3Pms4FRxjY
  repo: dhth/cueitup
  name: release
  url: https://dhth.github.io/cueitup

- id: W_kwDOLghtl84FWTlZ
  repo: dhth/ecsv
  name: build
- id: W_kwDOLghtl84FWTla
  repo: dhth/ecsv
  name: release
```

`{{runNumber}}` gets replaced with the actual run number of the workflow.

You can find the ID of a workflow as follows:

```bash
curl -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <GH_TOKEN>" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/<OWNER>/<REPO>/actions/workflows

# or

gh api repos/<OWNER>/<REPO>/actions/workflows

# use node_id from the response
```

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

> A sample page generated via `act3` is running at
> [https://dhth.github.io/act3](https://dhth.github.io/act3).

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

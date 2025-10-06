<!-- template:define:options
{
  "nodescription": true
}
-->
![logo](https://liam.sh/-/gh/svg/lrstanley/clix?icon=material-symbols%3Aterminal&icon.height=100&layout=left&icon.color=rgba%2839%2C+132%2C+85%2C+1%29&font=1.2)

<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<p align="center">
  <a href="https://github.com/lrstanley/clix/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/clix?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/clix/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/clix?style=flat-square">
  </a>



  <a href="https://github.com/lrstanley/clix/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/actions/workflow/status/lrstanley/clix/test.yml?branch=master&label=test&style=flat-square">
  </a>




  <a href="https://codecov.io/gh/lrstanley/clix">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/lrstanley/clix/master?style=flat-square">
  </a>

  <a href="https://pkg.go.dev/github.com/lrstanley/clix">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/lrstanley/clix?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/clix">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/lrstanley/clix?style=flat-square">
  </a>
</p>
<p align="center">
  <a href="https://github.com/lrstanley/clix/issues?q=is:open+is:issue+label:bug">
    <img title="Bug reports" src="https://img.shields.io/github/issues/lrstanley/clix/bug?label=issues&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/clix/issues?q=is:open+is:issue+label:enhancement">
    <img title="Feature requests" src="https://img.shields.io/github/issues/lrstanley/clix/enhancement?label=feature%20requests&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/clix/pulls">
    <img title="Open Pull Requests" src="https://img.shields.io/github/issues-pr/lrstanley/clix?label=prs&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/clix/discussions/new?category=q-a">
    <img title="Ask a Question" src="https://img.shields.io/badge/support-ask_a_question!-blue?style=flat-square">
  </a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/discord-bytecord-blue.svg?style=flat-square" title="Discord Chat"></a>
</p>
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [Features](#sparkles-features)
  - [Usage](#gear-usage)
  - [Example help output](#example-help-output)
  - [Generate Markdown](#generate-markdown)
    - [Example output](#example-output)
      - [Application Options](#application-options)
      - [Example Group](#example-group)
      - [Logging Options](#logging-options)
  - [Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [Contributing](#handshake-contributing)
  - [License](#balance_scale-license)
<!-- template:end:toc -->

## :sparkles: Features

**clix** is a [kong](https://github.com/alecthomas/kong) wrapper, that handles parsing of struct
tags, unmarshalling config files, DI, and more. Includes builtin plugins for:

- Logging (JSON, pretty printed, logging to a file, or discarding) using `log/slog`.
  - Exposed handler and logger which you can use as a base for any additional
    logging configuration.
  - Change logging levels easily (and automatically when using `--debug`).
- Versioning (`--version`, `--version-json`) which aids in printing version
  information.
  - Exposes Go 1.18's build metadata, the ability to use that as the version
    info, automatically using VCS information if available.
  - Printing dependencies and build flags.
  - Embedding useful links (support, repo, homepage, etc) in both version
    output, and help output.
- Markdown (generate markdown from the CLI's help information). See [example 1](./examples/simple/README.md)
  and [example 2](./examples/multiple-commands/README.md). See [below](#generate-markdown)
  for more details.
- Built-in `--debug` flag.
- [godotenv](github.com/joho/godotenv) integration to auto-load environment variables
  from `.env` files, before parsing flags.

**clix** is configurable, so all of the above can be turned on/off, with a reasonable
default configuration that should work for most basic apps.

## :gear: Usage

<!-- template:begin:goget -->
<!-- do not edit anything in this "template" block, its auto-generated -->
```console
go get -u github.com/lrstanley/clix/v2@latest
```
<!-- template:end:goget -->

Example:

```go
package main

import (
	"github.com/lrstanley/clix/v2"
)

type Flags struct {
	Name string `name:"name" default:"world" help:"name to print"`
}

var cli = clix.NewWithDefaults[Flags]()

func main() {
	logger := cli.GetLogger()

	if cli.Debug {
		logger.Debug("thinking really hard...")
	}

	logger.Info("hello", "name", cli.Flags.Name)
}
```

## Example help output

Below shows an example of a non-tagged revision of `myproject`. When using
Git tags, the appropriate tag should be applied as the version, rather than
`(devel)`.

```console
$ ./myproject --help
Usage: myproject <command> [flags]

github.com/user/project :: (devel)

    build commit :: 063bc199778e1ae617c35b70608e4ef104d6ce31
      build date :: 2025-10-06T04:35:02Z
      go version :: go1.25.1 linux/amd64

Flags:
  -h, --help            Show context-sensitive help.
  -v, --version         prints version information and exits
      --version-json    prints version information in JSON format and exits
      --name="world"    name to print
  -D, --debug           enables debug mode

Logging flags
  --log.level="info"    logging level (none: disables logging) ($LOG_LEVEL)
  --log.json            output logs in JSON format ($LOG_JSON)
  --log.path=STRING     path to log file (disables stderr logging) ($LOG_PATH)
```

And version output:

```console
$ ./myproject --version
github.com/user/project :: (devel)
|  build commit :: 063bc199778e1ae617c35b70608e4ef104d6ce31
|    build date :: 2025-10-06T04:35:02Z
|    go version :: go1.25.1 linux/amd64

build options:
|    -buildmode :: exe
|     -compiler :: gc
|   CGO_ENABLED :: 1
|    CGO_CFLAGS ::
|  CGO_CPPFLAGS ::
|  CGO_CXXFLAGS ::
|   CGO_LDFLAGS ::
|        GOARCH :: amd64
|          GOOS :: linux
|       GOAMD64 :: v1
|           vcs :: git
|  vcs.revision :: 063bc199778e1ae617c35b70608e4ef104d6ce31
|      vcs.time :: 2025-10-06T04:35:02Z
|  vcs.modified :: false

dependencies:
  h1:iq6aMJDcFYP9uFrLdsiZQ2ZMmcshduyGv4Pek0MQPW0= :: github.com/alecthomas/kong :: v1.12.1
  h1:7eLL/+HRGLY0ldzfGMeQkb7vMd0as4CfYvUVzLqw0N0= :: github.com/joho/godotenv :: v1.5.1
  h1:2CQzrL6rslrsyjqLDwD11bZ5OpLBPU+g3G/r5LSfS8w= :: github.com/lmittmann/tint :: v1.1.2
                                          unknown :: github.com/lrstanley/clix/v2 :: (devel)
                                          unknown :: github.com/lrstanley/x/logging/handlers :: (devel)
```

You can also use `--version-json` for a more programmatic approach to the above
information.

## Generate Markdown

When using **clix**, you can generate markdown for your commands by passing the
command `generate-markdown` (hidden when using `--help`). This will print the markdown
to stdout.

```console
./<your-project> generate-markdown > USAGE.md
```

This functionality is configurable using environment variables:

| Environment Variable | Description | Default |
|----------------------|-------------|---------|
| `CLIX_TEMPLATE_PATH` | Path to a directory containing template files to use for the markdown. These inherit from the built-in templates, so you can simply override a specific sub-template to override only a specific section of the markdown. | `<built-in templates>` |
| `CLIX_OUTPUT_PATH` | Path to write the markdown to, or `-` to write to stdout. | `-` |

See [example 1](./examples/simple/README.md) and [example 2](./examples/multiple-commands/README.md)
for more examples on what this can look like.

## Migrating to v2

v2 is an overhaul of the project, changing the underlying parser, logger, and more.
The goal is to use more stdlib functionality where possible (like `log/slog`),
making Markdown generation more flexible, and reduce external dependencies. Should
be quite a bit less now!

Start by fetching the new version:

```console
go get -u github.com/lrstanley/clix/v2@latest
```

High-level changes:

- Moved from [go-flags](https://github.com/jessevdk/go-flags) to
  [kong](https://github.com/alecthomas/kong). Kong is more maintained, powerful,
  flexible, and has less issues with parse ordering. Kong is also configured
  through struct tags, though they are slightly different. Some examples of the
  naming differences:
  - `long` -> `name`
  - `description` -> `help`
  - `env-delim` -> `sep`
  - `required:"true"` -> `required:""`
  - `optional:""` for optional commands.
  - `command` -> `cmd`
  - See [the full list here](https://github.com/alecthomas/kong#supported-tags).
- Kong clears all struct values before parsing, whereas go-flags does not.
- Kong supports pre-validation hooks, resolvers for loading config files, etc.
- Kong doesn't create an env var for flags automatically. In v1, this caused a bunch
  of bugs, as some auto-created env vars would conflict with commonly used env vars.
  An example of this being `--path`, which conflicts with the `PATH` env var used
  by most operating systems.
- struct-based flags no longer support `--flag foo --flag bar` (go-flags style),
  and now they should be passed in as a single flag, like `--flag "foo,bar"` with
  a struct-tag separator defined like `sep:","` (the default).
- Version output no longer includes colors, as this required multiple dependencies.
  I also don't think it added much value.
- Logging, versioning, markdown, and Github runner auto-debug mode are now no longer
  loaded as plugins when using `clix.New`, but are still loaded by default when
  using `clix.NewWithDefaults`.
- Logging has migrated from `github.com/apex/log` to `log/slog`. apex also had
  pretty bad performance under high load, so you may notice a performance boost.
- Both prettified and JSON logging fields and output look different. JSON logging
  uses the standard `log/slog` JSON handler, with source lines added.
- Markdown generation now no longer requires "required" flags to be set, which was
  quite annoying before. It has moved from `--generate-markdown` (a flag) to
  `generate-markdown` (a command).
- Runner functionality has moved to an external package, [github.com/lrstanley/x/scheduler](https://pkg.go.dev/github.com/lrstanley/x/scheduler). It now supports crontab-style intervals,
  and a builder-style interface for creating jobs. See [examples/with-scheduler](./examples/with-scheduler)
  for an example.

See the [examples](./examples) for more details on on what your app should look like.

---

<!-- template:begin:support -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :raising_hand_man: Support & Assistance

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for
     guidelines on ensuring everyone has the best experience interacting with
     the community.
* :raising_hand_man: Take a look at the [support](.github/SUPPORT.md) document on
     guidelines for tips on how to ask the right questions.
* :lady_beetle: For all features/bugs/issues/questions/etc, [head over here](https://github.com/lrstanley/clix/issues/new/choose).
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :handshake: Contributing

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for guidelines
     on ensuring everyone has the best experience interacting with the
    community.
* :clipboard: Please review the [contributing](.github/CONTRIBUTING.md) doc for submitting
     issues/a guide on submitting pull requests and helping out.
* :old_key: For anything security related, please review this repositories [security policy](https://github.com/lrstanley/clix/security/policy).
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :balance_scale: License

```
MIT License

Copyright (c) 2021 Liam Stanley <liam@liam.sh>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

_Also located [here](LICENSE)_
<!-- template:end:license -->

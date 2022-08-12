<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<p align="center">clix -- go-flags wrapper with useful helpers</p>
<p align="center">
  <a href="https://github.com/lrstanley/clix/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/clix?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/clix/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/clix?style=flat-square">
  </a>


  <a href="https://github.com/lrstanley/clix/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/workflow/status/lrstanley/clix/test/master?label=test&style=flat-square&event=push">
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

  - []()
<!-- template:end:toc -->

## :sparkles: Features

- go-flags wrapper, that handles parsing and decoding, with additional
  helpers. Uses Go 1.18's generics to allow embedding your own custom struct.
- Auto-generated logger, and auto-provided logger flags that:
  - allows users to switch between plain-text, colored/pretty, JSON output
    (and quiet for no output).
  - Allows configuring the level dynamically.
  - Uses the built-in debug flag to automatically change the logging level.
- Built-in `--debug` flag.
- Built-in `--version` flag that provides a lot of useful features:
  - Using Go 1.18's build metadata, the ability to use that as the version
    info, automatically using VCS information if available.
  - Printing dependencies and build flags.
  - Embedding useful links (support, repo, homepage, etc) in both version
    output, and help output.
  - Colored output!
- `--generate-markdown` flag (hidden) that allows generating markdown from
  the CLI's help information (see below!).
- Uses [godotenv](github.com/joho/godotenv) to auto-load environment variables
  from `.env` files, before parsing flags.
- Many flags to enable/disable functionality to suit your needs.

## :ballot_box_with_check: TODO

- [ ] Generate commands/sub-commands/etc ([example](https://github.com/jessevdk/go-flags/pull/364))
- [ ] Custom markdown formatting

## :gear: Usage

<!-- template:begin:goget -->
<!-- do not edit anything in this "template" block, its auto-generated -->
```console
$ go get -u github.com/lrstanley/clix@latest
```
<!-- template:end:goget -->

Example:

```go
package main

import (
 "github.com/apex/log"
 clix "github.com/lrstanley/clix"
)

var (
  cli = &clix.CLI[Flags]{
  Links: clix.GithubLinks("github.com/lrstanley/myproject", "master", "https://mysite.com"),
 }
 logger log.Interface
)

type Flags struct {
 EnableHTTP bool   `short:"e" long:"enable-http" description:"enable the http server"`
 File       string `env:"FILE" short:"f" long:"file" description:"some file that does something"`

 SubFlags struct {
  Username string `env:"USERNAME" short:"u" long:"username" default:"admin" description:"example username"`
  Password string `env:"PASSWORD" short:"p" long:"password" description:"example password"`
 } `group:"Example Group" namespace:"example" env-namespace:"EXAMPLE"`
}

func main() {
 // Initializes cli flags, and a pre-configured logger, based off user-provided
 // flags (e.g. --log.json, --log.level, etc). Also automatically handles
 // --version, --help, etc.
 cli.Parse()
 logger = cli.Logger

 logger.WithFields(log.Fields{
  "debug":     cli.Debug,
  "file_path": cli.Flags.File,
 }).Info("hello world")
}
```

## Example help output

Below shows an example of a non-tagged revision of `myproject`. When using
Git tags, the appropriate tag should be applied as the version, rather than
`(devel)`.

```console
$ ./myproject --version
github.com/lrstanley/myproject :: (devel)
  build commit :: 2a00b14d2ff16b79ecbb2afc54f480c2c1e28172
    build date :: unknown
    go version :: go1.18.1 linux/amd64

helpful links:
      homepage :: https://myproject
        github :: https://github.com/lrstanley/myproject
        issues :: https://github.com/lrstanley/myproject/issues/new/choose
       support :: https://github.com/lrstanley/myproject/blob/master/.github/SUPPORT.md
  contributing :: https://github.com/lrstanley/myproject/blob/master/.github/CONTRIBUTING.md
      security :: https://github.com/lrstanley/myproject/security/policy

build options:
     -compiler :: gc
   CGO_ENABLED :: 1
    CGO_CFLAGS ::
  CGO_CPPFLAGS ::
  CGO_CXXFLAGS ::
   CGO_LDFLAGS ::
        GOARCH :: amd64
          GOOS :: linux
       GOAMD64 :: v1
           vcs :: git
  vcs.revision :: 2a00b14d2ff16b79ecbb2afc54f480c2c1e28172
      vcs.time :: 2022-04-26T00:08:11Z
  vcs.modified :: true

dependencies:
  h1:IVj9dxSeAC0CRxjM+AYLIKbNdCzAjUnsUjAp/td7kYo= :: ariga.io/atlas :: v0.3.8-0.20220424181913-f64001131c0e
  h1:Jlkg6X37VI/k5U02yTBB3MjKGniiBAmUGfA1TC1+dtU= :: ariga.io/entcache :: v0.0.0-20211014200019-283c566a429b
  h1:XE5df6hfIlK/YeRxY6ynRxoMKCqn2mIOcOjId8JrQN8= :: entgo.io/ent :: v0.10.2-0.20220424193633-04e0dc936be9
  h1:YB2fHEn0UJagG8T1rrWknE3ZQzWM06O8AMAatNn7lmo= :: github.com/agext/levenshtein :: v1.2.3
  h1:FHtw/xuaM8AgmvDDTI9fiwoAL25Sq2cxojnZICUU8l0= :: github.com/apex/log :: v1.9.0
  [...]
```

You can also use `./myproject --version-json` for a more programmatic approach
to the above information.

## Generate Markdown

When using **clix**, you can generate markdown for your commands by passing the
flag `--generate-markdown`, which is a hidden flag. This will print the markdown
to stdout.

```console
./<my-script> --generate-markdown > USAGE.md
```

### Example output

Raw:

```markdown
#### Application Options

| Environment vars | Flags               | Type   | Description                          |
| ---------------- | ------------------- | ------ | ------------------------------------ |
| -                | `-e, --enable-http` | bool   | enable the http server               |
| `FILE`           | `-f, --file`        | string | some file that does something        |
| -                | `-v, --version`     | bool   | prints version information and exits |
| `DEBUG`          | `-D, --debug`       | bool   | enables debug mode                   |

#### Example Group

| Environment vars   | Flags                    | Type   | Description                           |
| ------------------ | ------------------------ | ------ | ------------------------------------- |
| `EXAMPLE_USERNAME` | `-u, --example.username` | string | example username [**default: admin**] |
| `EXAMPLE_PASSWORD` | `-p, --example.password` | string | example password                      |

#### Logging Options

| Environment vars | Flags          | Type   | Description                                                                      |
| ---------------- | -------------- | ------ | -------------------------------------------------------------------------------- |
| `LOG_QUIET`      | `--log.quiet`  | bool   | disable logging to stdout (also: see levels)                                     |
| `LOG_LEVEL`      | `--log.level`  | string | logging level [**default: info**] [**choices: debug, info, warn, error, fatal**] |
| `LOG_JSON`       | `--log.json`   | bool   | output logs in JSON format                                                       |
| `LOG_PRETTY`     | `--log.pretty` | bool   | output logs in a pretty colored format (cannot be easily parsed)                 |
```

Generated:

---

#### Application Options

| Environment vars | Flags               | Type   | Description                          |
| ---------------- | ------------------- | ------ | ------------------------------------ |
| -                | `-e, --enable-http` | bool   | enable the http server               |
| `FILE`           | `-f, --file`        | string | some file that does something        |
| -                | `-v, --version`     | bool   | prints version information and exits |
| `DEBUG`          | `-D, --debug`       | bool   | enables debug mode                   |

#### Example Group

| Environment vars   | Flags                    | Type   | Description                           |
| ------------------ | ------------------------ | ------ | ------------------------------------- |
| `EXAMPLE_USERNAME` | `-u, --example.username` | string | example username [**default: admin**] |
| `EXAMPLE_PASSWORD` | `-p, --example.password` | string | example password                      |

#### Logging Options

| Environment vars | Flags          | Type   | Description                                                                      |
| ---------------- | -------------- | ------ | -------------------------------------------------------------------------------- |
| `LOG_QUIET`      | `--log.quiet`  | bool   | disable logging to stdout (also: see levels)                                     |
| `LOG_LEVEL`      | `--log.level`  | string | logging level [**default: info**] [**choices: debug, info, warn, error, fatal**] |
| `LOG_JSON`       | `--log.json`   | bool   | output logs in JSON format                                                       |
| `LOG_PRETTY`     | `--log.pretty` | bool   | output logs in a pretty colored format (cannot be easily parsed)                 |

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

Copyright (c) 2021 Liam Stanley <me@liamstanley.io>

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

# cocopilot 💩🧑‍✈️

> fetches API tokens to use GitHub Copilot with any tool

[![GitHub Issues](https://img.shields.io/github/issues/wwmoraes/cocopilot.svg)](https://github.com/wwmoraes/cocopilot/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/wwmoraes/cocopilot.svg)](https://github.com/wwmoraes/cocopilot/pulls)

![Codecov](https://img.shields.io/codecov/c/github/wwmoraes/cocopilot)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

______________________________________________________________________

## 📝 Table of Contents

- [About](#-about)
- [Getting Started](#-getting-started)
- [Usage](#-usage)

<!-- - [TODO](./TODO.md) -->

- [Contributing](./CONTRIBUTING.md)

## 🧐 About

> Once upon a time, some "CISO wizards" (believe me, self-entitled) in a
> bureaucratic company decided to allow GitHub Copilot - but only through Visual
> Studio Code. Such bullshit decision smelled so bad that an engineer set out to
> find a way to get Copilot tokens to use with other tools such as aichat. This
> is the end result.

`cocopilot` is a CLI that [does one thing, and does it well][unix-philosophy]:
retrieves and refreshes GitHub Copilot tokens for tools to interact with its
inference APIs. This allows you to feed any tool that supports GitHub Copilot,
from CLIs to agents and anything in-between.

The name means two things: first, that it is a co-copilot, as in a third-tier
pilot. Second, "cocô" means "poop" in Portuguese, which perfectly represents
such bullshit policy that made me create this unnecessary tool in the first
place. I bet such "security decision" is for the [greater good][greater-good].

## 🏁 Getting Started

Using nix: `nix run github:wwmoraes/cocopilot`. Its also available in my NUR.

## 🎈 Usage

Run `cocopilot`, it'll automagically retrieve and return a valid token on the
standard output; this may be an existing token, a refreshed token or a new one
after asking you to authenticate.

The command is save to use directly as an assignment to environment variables,
such as `COPILOT_API_KEY` for aichat:

```shell
env COPILOT_API_KEY=$(cocopilot) aichat --model copilot:gpt-4.1 hi
```

> [!NOTE] aichat in specific requires an openai-compatible client configuration
> for GitHub Copilot to make the example command above work ;)

[greater-good]: https://www.youtube.com/watch?v=5u8vd_YNbTw
[unix-philosophy]: https://www.linfo.org/unix_philosophy.html

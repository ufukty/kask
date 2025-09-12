# Docs

<img src=".assets/card-og.png" style="width:min(100%, 640px);border-radius:8px">

👋 Welcome to Kask documentation.

## Repository

Repository is on [GitHub](https://github.com/ufukty/kask)

## Install

To embed the version information into binary (and produced files for reproducibility):

```sh
cd "$(mktemp -d)"
git clone https://github.com/ufukty/kask
cd kask
git fetch --tags --quiet
git checkout "$(git tag --list 'v*' | sort -Vr | head -n 1)"
make install
```

Make sure you have the `make`, `git` and `go` commands accessible by your `PATH`. Users don't care about the version number stamping can as well use the `go install` command directlty.

## License

See [LICENSE](https://github.com/ufukty/kask/blob/main/LICENSE).

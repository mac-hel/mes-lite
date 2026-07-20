# Go Modules [see](https://go.dev/ref/mod#modules-overview)
**module** - collection of **packages** released, versioned, and distributed together.
    - **module path** is defined in `go.mod`; describes where to find it and what it does, e.g.:
        `golang.org/x/net` (module in repo's root)
        `golang.org/x/tools/gopls` (mosule in repo's `gopls` subdir)
        `golang.org/x/tools/gopls/v2` (module released under major version 2; for 1 is omitted)
    - **module root directory** is dir containing `go.mod`
    - **main module** is the one containing dir (with `go.mod`) where currently go command is executed (active project)
**package** - collection of source files in the same directory (all packages in module are compiled together)
    - **package path** = `MODULE_PATH/PACKAGE_SUBDIR`, e.g.: `golang.org/x/net/html` (`html` is package)

git init --initial-branch=main --object-format=sha1

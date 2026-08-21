# GO language crib

module - versioned(MVS) packages, go.mod{path, go.sum,
    `cmd/ api/main.go web/main.go`      - multiple binaries
    `internal/`                         - compiler/toolchain enforces hidden packages
    `pkg/`                              - public packages (e.g. lib) or in root dir
package - compil/encaps unit, dir{*.go          ;cohesive responsibilities (business, not technical)
'main' - exec prog, func main() (compos. root)

graceful shutdown: signal, stop accepting, mark app unready, shutdown http serv, finish pending reqs,
                   cancel app/root ctx, wait goroutines, flush observab., close resources, exit

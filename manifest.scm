;; Development environment for git-cone.
;; Usage:
;;   guix shell -m manifest.scm
;;
;; For cgo-based tools (for example govulncheck), use Zig's C compiler:
;;   CC="zig cc" CGO_ENABLED=1 guix shell -m manifest.scm -- <tool> <args>

(specifications->manifest
 (list "go"
       "gopls"
       "govulncheck"
       "go-staticcheck"
       "go-tools"
       "gofumpt"
       "zig"
       "podman"
       "podman-compose"
       "make"
       "git"
       "openssh"
       "sqlite"))

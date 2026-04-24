;; Development environment for git-cone.
;; Usage:
;;   guix shell -m manifest.scm

(specifications->manifest
 (list "go"
       "gopls"
       "govulncheck"
       "go-staticcheck"
       "go-golangci-lint"
       "podman"
       "podman-compose"
       "make"
       "git"
       "openssh"
       "gcc-toolchain"
       "sqlite"))

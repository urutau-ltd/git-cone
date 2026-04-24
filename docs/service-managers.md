# Running git-cone Without systemd

This fork is Guix-first and does not ship systemd packaging.
Examples below assume:

- binary: `/usr/local/bin/cone`
- data path: `/var/lib/git-cone`
- config/data env file: `/etc/git-cone/env`
- runtime user: `git`

Example env file:

```sh
GIT_CONE_DATA_PATH=/var/lib/git-cone
GIT_CONE_NAME=Git Cone
GIT_CONE_SSH_PUBLIC_URL=ssh://git.example.com
GIT_CONE_HTTP_PUBLIC_URL=https://git.example.com
GIT_CONE_SECURITY_STRICT=true
```

Create state directories first:

```sh
install -d -m 0750 -o git -g git /var/lib/git-cone
install -d -m 0755 /etc/git-cone
```

## SysVinit

Example `/etc/init.d/git-cone`:

```sh
#!/bin/sh
### BEGIN INIT INFO
# Provides:          git-cone
# Required-Start:    $network
# Required-Stop:     $network
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
### END INIT INFO

DAEMON=/usr/local/bin/cone
NAME=git-cone
PIDFILE=/run/$NAME.pid
RUNDIR=/run
USER=git
WORKDIR=/var/lib/git-cone
ENVFILE=/etc/git-cone/env

[ -r "$ENVFILE" ] && . "$ENVFILE"

start() {
  mkdir -p "$RUNDIR"
  chown "$USER" "$RUNDIR"
  start-stop-daemon --start --background \
    --make-pidfile --pidfile "$PIDFILE" \
    --chuid "$USER" --chdir "$WORKDIR" \
    --exec "$DAEMON" -- serve
}

stop() {
  start-stop-daemon --stop --pidfile "$PIDFILE" --retry TERM/30/KILL/5
}

case "$1" in
  start) start ;;
  stop) stop ;;
  restart) stop; start ;;
  *) echo "usage: $0 {start|stop|restart}"; exit 1 ;;
esac
```

## OpenRC

Example `/etc/init.d/git-cone`:

```sh
#!/sbin/openrc-run

name="git-cone"
description="git-cone Git server"
command="/usr/local/bin/cone"
command_args="serve"
command_user="git:git"
directory="/var/lib/git-cone"
pidfile="/run/${RC_SVCNAME}.pid"
command_background="yes"
output_log="/var/log/git-cone.log"
error_log="/var/log/git-cone.log"
env_file="/etc/git-cone/env"

depend() {
  need net
}
```

Enable it:

```sh
rc-update add git-cone default
rc-service git-cone start
```

## Runit

Create `/etc/sv/git-cone/run`:

```sh
#!/bin/sh
set -eu
cd /var/lib/git-cone
exec 2>&1
[ -r /etc/git-cone/env ] && . /etc/git-cone/env
exec chpst -u git:git /usr/local/bin/cone serve
```

Optional log service `/etc/sv/git-cone/log/run`:

```sh
#!/bin/sh
exec svlogd -tt /var/log/git-cone
```

Enable it:

```sh
ln -s /etc/sv/git-cone /var/service/git-cone
sv up git-cone
```

## GNU Shepherd

The official Shepherd manual documents `service`,
`make-forkexec-constructor`, and `make-kill-destructor`.

Example service definition:

```scheme
(use-modules (shepherd service)
             (srfi srfi-1))

(define git-cone-env
  '("GIT_CONE_DATA_PATH=/var/lib/git-cone"
    "GIT_CONE_NAME=Git Cone"
    "GIT_CONE_SSH_PUBLIC_URL=ssh://git.example.com"
    "GIT_CONE_HTTP_PUBLIC_URL=https://git.example.com"
    "GIT_CONE_SECURITY_STRICT=true"))

(register-services
 (list
  (service
   '(git-cone)
   #:requirement '(networking)
   #:start (make-forkexec-constructor
            '("/usr/local/bin/cone" "serve")
            #:directory "/var/lib/git-cone"
            #:user "git"
            #:group "git"
            #:log-file "/var/log/git-cone.log"
            #:environment-variables git-cone-env)
   #:stop (make-kill-destructor)
   #:respawn? #t)))
```

Then:

```sh
herd start git-cone
herd status git-cone
```

Official references:

- https://www.gnu.org/software/shepherd/manual/shepherd.html
- https://www.gnu.org/software/shepherd/manual/html_node/Defining-Services.html
- https://www.gnu.org/software/shepherd/manual/html_node/Service-De_002d-and-Constructors.html

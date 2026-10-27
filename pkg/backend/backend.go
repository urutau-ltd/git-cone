package backend

import (
	"context"
	"fmt"
	"sync"

	"charm.land/log/v2"
	"github.com/urutau-ltd/git-cone/git"
	"github.com/urutau-ltd/git-cone/pkg/config"
	"github.com/urutau-ltd/git-cone/pkg/db"
	"github.com/urutau-ltd/git-cone/pkg/notify"
	"github.com/urutau-ltd/git-cone/pkg/store"
)

// Backend is the Soft Serve backend that handles users, repositories, and
// server settings management and operations.
type Backend struct {
	ctx    context.Context
	cfg    *config.Config
	db     *db.DB
	store  store.Store
	logger *log.Logger

	// importGate reports why imports and mirror syncs are off, or "" when
	// they are allowed. It runs the git version check on first use, so
	// nothing has to happen in a particular order at startup for the gate
	// to be in place.
	importGate func() string

	notifier        notify.Notifier
	webhookFailures sync.Map // key: int64 (webhook ID), value: *webhookCounter
}

// New returns a new Soft Serve backend.
func New(ctx context.Context, cfg *config.Config, db *db.DB, st store.Store) *Backend {
	logger := log.FromContext(ctx).WithPrefix("backend")
	b := &Backend{
		ctx:      ctx,
		cfg:      cfg,
		db:       db,
		store:    st,
		logger:   logger,
		notifier: notify.Noop{},
	}
	b.importGate = sync.OnceValue(b.checkGitVersion)

	return b
}

// SetNotifier sets the security event notifier.
func (b *Backend) SetNotifier(n notify.Notifier) {
	b.notifier = n
}

// TrackAuthFailure increments the auth failure counter for the given IP.
// Fires a notification when 5 failures occur within 60 seconds.
func (b *Backend) TrackAuthFailure(ip string) {
	notify.TrackAuthFailure(b.notifier, ip)
}

// ImportsDisabledReason returns why repository imports and mirror syncs are
// off, or "" when they are allowed. The reason is for the server log; users
// get proto.ErrImportsDisabled.
func (d *Backend) ImportsDisabledReason() string {
	return d.importGate()
}

// minImportGitVersion is the oldest git that honours the settings the import
// and mirror guards rely on: http.curloptResolve (2.33) for pinning the
// validated address and http.followRedirects via GIT_CONFIG_* (2.31). Older
// git silently ignores both, so the guard degrades to nothing.
var minImportGitVersion = git.Version{Major: 2, Minor: 33}

// binVersion is how the gate learns the installed git version. Tests swap it
// out to exercise the disabled path without an old git on hand.
var binVersion = git.BinVersion

// checkGitVersion returns why imports must stay off, or "" if the installed
// git is new enough. Running a clone under a git that ignores the guards
// would be worse than refusing it, so an unreadable version also closes
// the gate.
func (d *Backend) checkGitVersion() string {
	v, err := binVersion()
	if err != nil {
		return fmt.Sprintf("cannot determine git version: %v", err)
	}
	if v.Less(minImportGitVersion) {
		return fmt.Sprintf("git %s is older than %s, which imports and mirror syncs need to be safe", v, minImportGitVersion)
	}
	return ""
}

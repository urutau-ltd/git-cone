package backend

import (
	"context"
	"sync"

	"charm.land/log/v2"
	"github.com/urutau-ltd/git-cone/pkg/config"
	"github.com/urutau-ltd/git-cone/pkg/db"
	"github.com/urutau-ltd/git-cone/pkg/notify"
	"github.com/urutau-ltd/git-cone/pkg/store"
	"github.com/urutau-ltd/git-cone/pkg/task"
)

// Backend is the Soft Serve backend that handles users, repositories, and
// server settings management and operations.
type Backend struct {
	ctx             context.Context
	cfg             *config.Config
	db              *db.DB
	store           store.Store
	logger          *log.Logger
	cache           *cache
	manager         *task.Manager
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
		manager:  task.NewManager(ctx),
		notifier: notify.Noop{},
	}

	// TODO: implement a proper caching interface
	cache := newCache(b, 1000)
	b.cache = cache

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

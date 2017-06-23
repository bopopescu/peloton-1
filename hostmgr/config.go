package hostmgr

import (
	"time"

	"code.uber.internal/infra/peloton/hostmgr/reconcile"
)

// Config is Host Manager specific configuration
type Config struct {
	Port int `yaml:"port"`
	// Time to hold offer for in seconds
	OfferHoldTimeSec int `yaml:"offer_hold_time_sec"`
	// Frequency of running offer pruner
	OfferPruningPeriodSec int `yaml:"offer_pruning_period_sec"`

	// Minimum backoff duration for retrying any Mesos connection.
	MesosBackoffMin time.Duration `yaml:"mesos_backoff_min"`

	// Maximum backoff duration for retrying any Mesos connection.
	MesosBackoffMax time.Duration `yaml:"mesos_backoff_max"`

	// FIXME(gabe): this isnt really the DB write concurrency. This is
	// only used for processing task updates and should be moved into
	// the storage namespace, and made clearer what this controls
	// (threads? rows? statements?)
	DbWriteConcurrency int `yaml:"db_write_concurrency"`

	// Number of go routines that will ack for status updates to mesos
	TaskUpdateAckConcurrency int `yaml:"taskupdate_ack_concurrency"`

	// Size of the channel buffer of the status updates
	TaskUpdateBufferSize int `yaml:"taskupdate_buffer_size"`

	TaskReconcilerConfig *reconcile.TaskReconcilerConfig `yaml:"task_reconciler"`

	HostmapRefreshInterval time.Duration `yaml:"hostmap_refresh_interval"`
}

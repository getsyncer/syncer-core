package drift

import (
	"github.com/getsyncer/syncer-core/config"
)

type RunData struct {
	// The root config file
	RootConfig *config.Root
	// The subsection of the config file that is relevant to this run
	RunConfig config.Dynamic
	// The registries of this run
	Registry Registry
	// Where we want to copy destination files to
	DestinationWorkingDir string
}

func (r *RunData) AutogenMsg() string {
	// Separate the string so that it doesn't get picked up by the autogen script.
	return MagicTrackedString
}

const (
	MagicTrackedString                = "THIS FILE IS AUTOGENERATED BY SYNCER" + " DO NOT EDIT"
	DefaultSyncerGeneratedGoDirectory = "internal"
	DefaultSyncerGeneratedGoFilename  = "syncer.go"
	DefaultSyncerConfigFileName       = ".syncer.yaml"
)

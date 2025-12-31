package client_adapter

import (
	"context"
	"io"
	"time"

	"github.com/restic/restic/internal/archiver"
	"github.com/restic/restic/internal/checker"
	"github.com/restic/restic/internal/crypto"
	"github.com/restic/restic/internal/fs"
	"github.com/restic/restic/internal/repository"
	"github.com/restic/restic/internal/restic"
	"github.com/restic/restic/internal/restorer"
	"github.com/restic/restic/internal/ui/progress"
)

// ============================================================================
// EXPORT TYPES
// ============================================================================

type Repository = repository.Repository
type Backend = restic.Backend
type Snapshot = restic.Snapshot
type Options = repository.Options
type ID = restic.ID
type FileType = restic.FileType
type Handle = restic.Handle
type FileInfo = restic.FileInfo
type RewindReader = restic.RewindReader

// ============================================================================
// REPOSITORY
// ============================================================================

func NewRepository(be Backend, opts Options) (*Repository, error) {
	return repository.New(be, opts)
}

// ============================================================================
// ARCHIVER
// ============================================================================

func NewArchiver(repo restic.Repository, fs fs.FS, opts archiver.Options) *archiver.Archiver {
	return archiver.New(repo, fs, opts)
}

type ArchiverOptions = archiver.Options
type SnapshotOptions = archiver.SnapshotOptions

// ============================================================================
// RESTORER
// ============================================================================

func NewRestorer(repo restic.Repository, sn *restic.Snapshot, sparse bool, p *progress.Counter) *restorer.Restorer {
	return restorer.NewRestorer(repo, sn, sparse, p)
}

// ============================================================================
// CHECKER
// ============================================================================

func NewChecker(repo restic.Repository, checkUnused bool) *checker.Checker {
	return checker.New(repo, checkUnused)
}

// ============================================================================
// HELPER WRAPPERS
// (Because some methods are attached to internal structs we can't fully alias)
// ============================================================================

func LoadSnapshot(ctx context.Context, repo restic.Repository, id ID) (*Snapshot, error) {
	return restic.LoadSnapshot(ctx, repo, id)
}

func FindUsedBlobs(ctx context.Context, repo restic.Repository, treeIDs []ID, blobs restic.BlobSet, p *progress.Counter) error {
	return restic.FindUsedBlobs(ctx, repo, treeIDs, blobs, p)
}

func ParseID(s string) (ID, error) {
	return restic.ParseID(s)
}

// Constants
var SnapshotFile = restic.SnapshotFile
var PackFile = restic.PackFile
var InvalidBlob = restic.InvalidBlob
var NumBlobTypes = restic.NumBlobTypes

// Interface Aliases
type MasterIndex = restic.MasterIndex
type PackedBlob = restic.PackedBlob
type BlobSet = restic.BlobSet
type CountedBlobSet = restic.CountedBlobSet

func NewIDSet() restic.IDSet {
	return restic.NewIDSet()
}

func NewCountedBlobSet() restic.CountedBlobSet {
	return restic.NewCountedBlobSet()
}

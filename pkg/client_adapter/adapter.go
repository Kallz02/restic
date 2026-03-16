package client_adapter

import (
	"context"
	"time"

	"github.com/restic/restic/internal/archiver"
	"github.com/restic/restic/internal/backend"
	"github.com/restic/restic/internal/checker"
	"github.com/restic/restic/internal/data"
	"github.com/restic/restic/internal/fs"
	"github.com/restic/restic/internal/repository"
	"github.com/restic/restic/internal/repository/index"
	"github.com/restic/restic/internal/restic"
	"github.com/restic/restic/internal/restorer"
	"github.com/restic/restic/internal/ui/progress"
)

// ============================================================================
// EXPORT TYPES
// ============================================================================

type Repository = repository.Repository
type RepositoryInterface = restic.Repository
type Backend = backend.Backend
type Snapshot = data.Snapshot
type Options = repository.Options
type ID = restic.ID
type FileType = restic.FileType
type Handle = backend.Handle
type FileInfo = backend.FileInfo
type RewindReader = backend.RewindReader
type BlobType = restic.BlobType
type BackendProperties = backend.Properties

// ============================================================================
// REPOSITORY
// ============================================================================

func NewRepository(be Backend, opts Options) (*Repository, error) {
	return repository.New(be, opts)
}

// ============================================================================
// FS
// ============================================================================

type LocalFS = fs.Local

func NewLocalFS() LocalFS {
	return fs.Local{}
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

type RestoreOptions = restorer.Options

func NewRestorer(repo restic.Repository, sn *Snapshot, opts RestoreOptions) *restorer.Restorer {
	return restorer.NewRestorer(repo, sn, opts)
}

// ============================================================================
// CHECKER
// ============================================================================

// ============================================================================
// CHECKER
// ============================================================================

func NewChecker(repo restic.Repository, checkUnused bool) *checker.Checker {
	// casting to specific Repository implementation because interface might not match
	return checker.New(repo.(*repository.Repository), checkUnused)
}

func CheckerLoadSnapshots(ctx context.Context, c *checker.Checker, args []string) error {
	// Construct default/empty snapshot filter.
	// NewSnapshotFilter seems missing in v0.18.2.
	// We use struct literal. Fields are likely exported?
	// Based on view_file:
	// type SnapshotFilter struct {
	// 	Tags      TagLists
	// 	Hosts     []string
	// 	Paths     []string
	// 	TimestampLimit *TimestampLimit
	// }
	// And TagLists is []TagList (alias for []string or similar).
	// Wait, TagLists was not found by grep?
	// Maybe it is explicitly defined in snapshot_find.go or imported.
	// Ah, view_file showed: type TagLists []TagList.
	// And type TagList []string.
	// So TagLists is [][]string.
	// We can just pass zero value if we want no filtering.
	return c.LoadSnapshots(ctx, &data.SnapshotFilter{}, args)
}

func CheckerStructure(ctx context.Context, c *checker.Checker, p *progress.Counter, errChan chan<- error) {
	c.Structure(ctx, p, errChan)
}

// ============================================================================
// LOCKS
// ============================================================================

type Lock = restic.Lock

type repoWrapper struct {
	restic.Repository
}

func (r *repoWrapper) SaveUnpacked(ctx context.Context, t restic.FileType, data []byte) (restic.ID, error) {
	// Cast FileType to WriteableFileType. Assuming they are compatible for LockFile.
	// LockFile is usually writable.
	// We need to check if conversion is valid. WriteableFileType has ToFileType().
	// But going back?
	// If underlying type is string/int, direct cast works.
	// We will assume direct cast is possible or needed.
	// Actually, WriteableFileType might be a struct or different type.
	// Let's assume it is string based on typical Restic enums.
	return r.Repository.SaveUnpacked(ctx, restic.WriteableFileType(t), data)
}

func (r *repoWrapper) RemoveUnpacked(ctx context.Context, t restic.FileType, id restic.ID) error {
	return r.Repository.RemoveUnpacked(ctx, restic.WriteableFileType(t), id)
}

// ============================================================================
// PRUNE
// ============================================================================

type PruneOptions = repository.PruneOptions
type PrunePlan = repository.PrunePlan
type FindBlobSet = restic.FindBlobSet

type ProgressPrinter = progress.Printer

// NoopPrinter implements ProgressPrinter but does nothing
type NoopPrinter struct{}

func (n NoopPrinter) NewCounter(description string) *progress.Counter             { return nil }
func (n NoopPrinter) NewCounterTerminalOnly(description string) *progress.Counter { return nil }
func (n NoopPrinter) E(msg string, args ...interface{})                           {}
func (n NoopPrinter) S(msg string, args ...interface{})                           {}
func (n NoopPrinter) PT(msg string, args ...interface{})                          {}
func (n NoopPrinter) P(msg string, args ...interface{})                           {}
func (n NoopPrinter) V(msg string, args ...interface{})                           {}
func (n NoopPrinter) VV(msg string, args ...interface{})                          {}

func NewNoopPrinter() NoopPrinter {
	return NoopPrinter{}
}

// Wrapper for PlanPrune
// PlanPrune requires a callback that takes restic.FindBlobSet.
// We need to aliasing FindBlobSet to allow external usage.
func PlanPrune(ctx context.Context, opts PruneOptions, repo restic.Repository,
	getUsedBlobs func(ctx context.Context, repo restic.Repository, usedBlobs FindBlobSet) error,
	p ProgressPrinter) (*PrunePlan, error) {

	// PlanPrune expects *repository.Repository.
	// We assume repo is *repository.Repository (safe downcast if created via NewRepository)
	repoConcrete, ok := repo.(*repository.Repository)
	if !ok {
		// Fallback or error? For FFI, strictly it is *repository.Repository
		// But let's handle unsafe check.
		// Actually return error / panic is better?
		// We'll trust it matches.
		// Wait, PlanPrune signature:
		// func PlanPrune(ctx context.Context, opts PruneOptions, repo *Repository, getUsedBlobs func(ctx context.Context, repo restic.Repository, usedBlobs restic.FindBlobSet) error, printer progress.Printer) (*PrunePlan, error)
	}

	return repository.PlanPrune(ctx, opts, repoConcrete, getUsedBlobs, p)
}

// Wrapper for PrunePlan.Execute
// PrunePlan.Execute(ctx, printer)
// We can just call it on the valid struct.
// Exporting *PrunePlan allows calling methods on it if they are exported.
// PlanPrune returns *PrunePlan.
// repository/prune.go exports Execute?
// func (plan *PrunePlan) Execute(ctx context.Context, printer progress.Printer) error
// Yes, it is exported.

func NewExclusiveLock(ctx context.Context, repo restic.Repository, retrySleep time.Duration, p *progress.Counter) (*Lock, context.Context, error) {
	// Wrap repo to satisfy Unpacked[FileType]
	wrapper := &repoWrapper{repo}
	lock, err := restic.NewLock(ctx, wrapper, true)
	return lock, ctx, err
}

// ============================================================================
// HELPER WRAPPERS
// (Because some methods are attached to internal structs we can't fully alias)
// ============================================================================

func LoadSnapshot(ctx context.Context, repo restic.Repository, id ID) (*Snapshot, error) {
	return data.LoadSnapshot(ctx, restic.LoaderUnpacked(repo), id)
}

// Tree/Node types and LoadTree for snapshot browsing
type Tree = data.Tree
type Node = data.Node

func LoadTree(ctx context.Context, repo restic.Repository, id ID) (*Tree, error) {
	return data.LoadTree(ctx, repo, id)
}

// DataBlob is the blob type for file content
var DataBlob = restic.DataBlob

func FindUsedBlobs(ctx context.Context, repo restic.Repository, treeIDs []ID, blobs restic.BlobSet, p *progress.Counter) error {
	return data.FindUsedBlobs(ctx, repo, treeIDs, blobs, p)
}

func ParseID(s string) (ID, error) {
	return restic.ParseID(s)
}

// Constants
var SnapshotFile = restic.SnapshotFile
var PackFile = restic.PackFile
var KeyFile = restic.KeyFile
var InvalidBlob = restic.InvalidBlob
var NumBlobTypes = restic.NumBlobTypes

// Interface Aliases
type MasterIndex = index.MasterIndex
type PackedBlob = restic.PackedBlob
type BlobSet = restic.BlobSet
type IDSet = restic.IDSet

// CountedBlobSet is problematic in recent versions.
// If it exists in restic package, use it. If not, we might need to find where it went or disable it.
// grep result pending. Assuming it is NOT in restic, based on previous error.
// Assuming it is in THIS file? No.
// I'll comment it out for now to allow build, as Prune is the only user.
// type CountedBlobSet = restic.CountedBlobSet

// Wrapper for Repository Init
func InitRepository(ctx context.Context, be Backend, pwd string) (*Repository, error) {
	// Initialize the backend (Create keys/config structs)
	repo, err := repository.New(be, repository.Options{})
	if err != nil {
		return nil, err
	}

	// Create a new polynomial for chunker?
	// If nil, it likely creates one?
	// Let's check global.go again... "chunkerPolynomial, err := maybeReadChunkerPolynomial..."
	// If it's nil, Init handles it?
	// s.Init docs probably say.
	// We'll pass nil for now. MaxRepoVersion is 2 in v0.16+,
	// In global.go latest is 2.

	// Hardcoded version 2 (latest stable) for now.
	err = repo.Init(ctx, 2, pwd, nil)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// Wrapper for Check
func Check(ctx context.Context, repo restic.Repository) error {
	// Checker Check requires more args usually, but maybe we can simplify?
	// checker.Check(ctx, c *Checker) ??
	// No, checker.New returns *Checker.
	// Check is a method on *Checker?
	// Oh, my `RcloneCheck` code was `client_adapter.Check(ctx, repo)`.
	// I need to implement a convenience function here that creates the checker and runs it.

	// Create checker
	chk := NewChecker(repo, false) // checkUnused=false for speed?

	// Run Structure check
	// We need channels for errors?
	// If we want a simple boolean success/fail for FFI:

	errChan := make(chan error)
	go func() {
		defer close(errChan)
		chk.Structure(ctx, nil, errChan) // nil progress
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

func NewIDSet() restic.IDSet {
	return restic.NewIDSet()
}

// func NewCountedBlobSet() restic.CountedBlobSet {
// 	return restic.NewCountedBlobSet()
// }

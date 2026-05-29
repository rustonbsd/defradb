package corekvbadger

import (
	"context"
	"errors"

	badgerds "github.com/dgraph-io/badger/v4"
	rustcorekv "github.com/rustonbsd/corekv"
	rustbadger "github.com/rustonbsd/corekv/badger_ffi"
	compatcorekv "github.com/sourcenetwork/corekv"
)

// This adapter preserves DefraDB's existing sourcenetwork/corekv interfaces and
// context transaction plumbing while routing every actual Badger operation
// through rustonbsd/corekv/badger_ffi.

type Options = badgerds.Options

type Datastore struct {
	inner *rustbadger.Datastore
}

var _ compatcorekv.TxnStore = (*Datastore)(nil)
var _ compatcorekv.Dropable = (*Datastore)(nil)

func NewDatastore(path string, opts Options) (*Datastore, error) {
	inner, err := rustbadger.NewDatastore(path, opts)
	if err != nil {
		return nil, translateError(err)
	}

	return &Datastore{inner: inner}, nil
}

// DefraDB shares transactions between subsystems via sourcenetwork/corekv's
// context key. Looking up the txn here does not invoke sourcenetwork's Badger
// implementation; it only lets this adapter re-use the FFI-backed transaction
// that DefraDB previously stored on the context.
func txnFromContext(ctx context.Context) (*Txn, bool) {
	return compatcorekv.TryGetCtxTxnG[*Txn](ctx)
}

func (d *Datastore) Get(ctx context.Context, key []byte) ([]byte, error) {
	if txn, ok := txnFromContext(ctx); ok {
		return txn.Get(ctx, key)
	}

	value, err := d.inner.Get(ctx, key)
	return value, translateError(err)
}

func (d *Datastore) Has(ctx context.Context, key []byte) (bool, error) {
	if txn, ok := txnFromContext(ctx); ok {
		return txn.Has(ctx, key)
	}

	has, err := d.inner.Has(ctx, key)
	return has, translateError(err)
}

func (d *Datastore) Set(ctx context.Context, key, value []byte) error {
	if txn, ok := txnFromContext(ctx); ok {
		return txn.Set(ctx, key, value)
	}

	return translateError(d.inner.Set(ctx, key, value))
}

func (d *Datastore) Delete(ctx context.Context, key []byte) error {
	if txn, ok := txnFromContext(ctx); ok {
		return txn.Delete(ctx, key)
	}

	return translateError(d.inner.Delete(ctx, key))
}

func (d *Datastore) Iterator(ctx context.Context, opts compatcorekv.IterOptions) (compatcorekv.Iterator, error) {
	if txn, ok := txnFromContext(ctx); ok {
		return txn.Iterator(ctx, opts)
	}

	it, err := d.inner.Iterator(ctx, toRustIterOptions(opts))
	if err != nil {
		return nil, translateError(err)
	}

	return &Iterator{inner: it}, nil
}

func (d *Datastore) Close() error {
	return translateError(d.inner.Close())
}

func (d *Datastore) DropAll() error {
	return translateError(d.inner.DropAll())
}

func (d *Datastore) NewTxn(readonly bool) compatcorekv.Txn {
	return &Txn{inner: d.inner.NewTxn(readonly)}
}

type Txn struct {
	inner rustcorekv.Txn
}

var _ compatcorekv.Txn = (*Txn)(nil)

func (t *Txn) Get(ctx context.Context, key []byte) ([]byte, error) {
	value, err := t.inner.Get(ctx, key)
	return value, translateError(err)
}

func (t *Txn) Has(ctx context.Context, key []byte) (bool, error) {
	has, err := t.inner.Has(ctx, key)
	return has, translateError(err)
}

func (t *Txn) Set(ctx context.Context, key, value []byte) error {
	return translateError(t.inner.Set(ctx, key, value))
}

func (t *Txn) Delete(ctx context.Context, key []byte) error {
	return translateError(t.inner.Delete(ctx, key))
}

func (t *Txn) Iterator(ctx context.Context, opts compatcorekv.IterOptions) (compatcorekv.Iterator, error) {
	it, err := t.inner.Iterator(ctx, toRustIterOptions(opts))
	if err != nil {
		return nil, translateError(err)
	}

	return &Iterator{inner: it}, nil
}

func (t *Txn) Commit() error {
	return translateError(t.inner.Commit())
}

func (t *Txn) Discard() {
	t.inner.Discard()
}

type Iterator struct {
	inner rustcorekv.Iterator
}

var _ compatcorekv.Iterator = (*Iterator)(nil)

func (i *Iterator) Next() (bool, error) {
	hasNext, err := i.inner.Next()
	return hasNext, translateError(err)
}

func (i *Iterator) Key() []byte {
	return i.inner.Key()
}

func (i *Iterator) Value() ([]byte, error) {
	value, err := i.inner.Value()
	return value, translateError(err)
}

func (i *Iterator) Seek(key []byte) (bool, error) {
	found, err := i.inner.Seek(key)
	return found, translateError(err)
}

func (i *Iterator) Reset() {
	i.inner.Reset()
}

func (i *Iterator) Close() error {
	return translateError(i.inner.Close())
}

func toRustIterOptions(opts compatcorekv.IterOptions) rustcorekv.IterOptions {
	return rustcorekv.IterOptions{
		Prefix:   opts.Prefix,
		Start:    opts.Start,
		End:      opts.End,
		Reverse:  opts.Reverse,
		KeysOnly: opts.KeysOnly,
	}
}

func translateError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, rustcorekv.ErrNotFound):
		return compatcorekv.ErrNotFound
	case errors.Is(err, rustcorekv.ErrEmptyKey):
		return compatcorekv.ErrEmptyKey
	case errors.Is(err, rustcorekv.ErrValueNil):
		return compatcorekv.ErrValueNil
	case errors.Is(err, rustcorekv.ErrDiscardedTxn):
		return compatcorekv.ErrDiscardedTxn
	case errors.Is(err, rustcorekv.ErrDBClosed):
		return compatcorekv.ErrDBClosed
	case errors.Is(err, rustcorekv.ErrTxnConflict):
		return compatcorekv.ErrTxnConflict
	case errors.Is(err, rustcorekv.ErrReadOnlyTxn):
		return compatcorekv.ErrReadOnlyTxn
	default:
		return err
	}
}

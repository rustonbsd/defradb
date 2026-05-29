// Copyright 2022 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	badger "github.com/sourcenetwork/defradb/internal/corekvbadger"
	acpDB "github.com/sourcenetwork/defradb/internal/db/acp"
)

func newBadgerDB(ctx context.Context) (*DB, error) {
	rootstore, err := badger.NewDatastore("", badger.Options{InMemory: true})
	if err != nil {
		return nil, err
	}

	adminInfo, err := acpDB.NewNACInfo(ctx, "", false)
	if err != nil {
		return nil, err
	}
	return newDB(ctx, rootstore, adminInfo)
}

func TestNewDB(t *testing.T) {
	ctx := context.Background()
	rootstore, err := badger.NewDatastore("", badger.Options{InMemory: true})
	require.NoError(t, err)

	adminInfo, err := acpDB.NewNACInfo(ctx, "", false)
	require.NoError(t, err)

	_, err = NewDB(ctx, rootstore, adminInfo)
	require.NoError(t, err)
}

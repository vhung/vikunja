// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package handler

import (
	"bytes"
	"testing"

	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/events"
	"code.vikunja.io/api/pkg/models"
	"code.vikunja.io/api/pkg/modules/migration"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartMigrationRequiresProjectCreator(t *testing.T) {
	config.ServiceProjectCreators.Set([]string{"user2"})
	t.Cleanup(func() {
		config.ServiceProjectCreators.Set([]string{})
	})

	t.Run("oauth migrator", func(t *testing.T) {
		clearMigrationStatus(t)
		events.ClearDispatchedEvents()
		u := getTestUser(t)

		err := StartMigration(&stubMigrator{NameValue: "stub"}, u)
		require.Error(t, err)
		assert.True(t, models.IsErrGenericForbidden(err))
		assert.Empty(t, events.GetDispatchedEvents((&MigrationRequestedEvent{}).Name()))

		status, err := migration.GetMigrationStatus(&stubMigrator{NameValue: "stub"}, u)
		require.NoError(t, err)
		assert.Zero(t, status.ID, "no claim may be left behind")
	})
	t.Run("file migrator", func(t *testing.T) {
		clearMigrationStatus(t)
		events.ClearDispatchedEvents()
		u := getTestUser(t)

		err := StartFileMigration(&stubFileMigrator{name: "file-stub"}, u, bytes.NewReader([]byte("x")), 1, nil)
		require.Error(t, err)
		assert.True(t, models.IsErrGenericForbidden(err))
		assert.Empty(t, events.GetDispatchedEvents((&FileMigrationRequestedEvent{}).Name()))
	})
	t.Run("listed user may still migrate", func(t *testing.T) {
		clearMigrationStatus(t)
		events.ClearDispatchedEvents()

		require.NoError(t, StartMigration(&stubMigrator{NameValue: "stub"}, getTestUserByID(t, 2)))
	})
}

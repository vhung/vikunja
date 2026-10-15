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

package webtests

import (
	"net/http"
	"testing"

	"code.vikunja.io/api/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCreatorsAPI(t *testing.T) {
	e, err := setupTestEnv()
	require.NoError(t, err)
	config.ServiceProjectCreators.Set("user2")
	t.Cleanup(func() {
		config.ServiceProjectCreators.Set([]string{})
	})
	unlisted := humaTokenFor(t, &testuser1)
	listed := humaTokenFor(t, &testuser2)

	t.Run("user info reports the permission", func(t *testing.T) {
		rec := humaRequest(t, e, http.MethodGet, "/api/v2/user", "", unlisted, "")
		require.Equal(t, http.StatusOK, rec.Code, "body: %s", rec.Body.String())
		assert.Contains(t, rec.Body.String(), `"can_create_projects":false`)

		rec = humaRequest(t, e, http.MethodGet, "/api/v2/user", "", listed, "")
		require.Equal(t, http.StatusOK, rec.Code, "body: %s", rec.Body.String())
		assert.Contains(t, rec.Body.String(), `"can_create_projects":true`)
	})
	t.Run("unlisted user cannot create a root project", func(t *testing.T) {
		rec := humaRequest(t, e, http.MethodPost, "/api/v2/projects", `{"title":"root by user1"}`, unlisted, "application/json")
		assert.Equal(t, http.StatusForbidden, rec.Code, "body: %s", rec.Body.String())
	})
	t.Run("unlisted user can still create a child under a writable parent", func(t *testing.T) {
		rec := humaRequest(t, e, http.MethodPost, "/api/v2/projects", `{"title":"child by user1","parent_project_id":10}`, unlisted, "application/json")
		assert.Equal(t, http.StatusCreated, rec.Code, "body: %s", rec.Body.String())
	})
	t.Run("listed user can create a root project", func(t *testing.T) {
		rec := humaRequest(t, e, http.MethodPost, "/api/v2/projects", `{"title":"root by user2"}`, listed, "application/json")
		assert.Equal(t, http.StatusCreated, rec.Code, "body: %s", rec.Body.String())
	})
	t.Run("unlisted user cannot duplicate to the top level", func(t *testing.T) {
		rec := humaRequest(t, e, http.MethodPost, "/api/v2/projects/1/duplicate", `{}`, unlisted, "application/json")
		assert.Equal(t, http.StatusForbidden, rec.Code, "body: %s", rec.Body.String())
	})
}

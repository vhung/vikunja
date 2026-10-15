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

package models

import (
	"testing"

	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/files"
	"code.vikunja.io/api/pkg/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setProjectCreators(t *testing.T, value interface{}) {
	t.Helper()
	config.ServiceProjectCreators.Set(value)
	t.Cleanup(func() {
		config.ServiceProjectCreators.Set([]string{})
	})
}

func TestProjectCreators(t *testing.T) {
	t.Run("create root denied for unlisted user", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user2"})

		can, err := (&Project{Title: "root"}).CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.False(t, can)
	})
	t.Run("create root allowed for listed user", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user2", "user1"})

		can, err := (&Project{Title: "root"}).CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.True(t, can)
	})
	t.Run("username match is case sensitive", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"USER1"})

		can, err := (&Project{Title: "root"}).CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.False(t, can)
	})
	t.Run("empty list keeps upstream behaviour", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{})

		can, err := (&Project{Title: "root"}).CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.True(t, can)
	})
	t.Run("username comes from the database, not the auth claims", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user2"})

		can, err := (&Project{Title: "root"}).CanCreate(s, &user.User{
			ID:       1,
			Username: "user2",
		})
		require.NoError(t, err)
		assert.False(t, can)
	})
	t.Run("env style comma and space separated string", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()

		for _, value := range []string{"user2,user1", "user2 user1", " user2 , user1 "} {
			setProjectCreators(t, value)
			can, err := IsProjectCreator(s, &user.User{ID: 1})
			require.NoError(t, err)
			assert.True(t, can, value)
		}

		setProjectCreators(t, "user2,user3")
		can, err := IsProjectCreator(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.False(t, can)
	})
	t.Run("link share is never a creator", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user1"})

		can, err := IsProjectCreator(s, &LinkSharing{ID: 1})
		require.NoError(t, err)
		assert.False(t, can)
	})
	t.Run("create child under writable parent still allowed", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user2"})

		// User 1 has write access to project 10, which user 1 does not own.
		p := &Project{
			Title:           "child",
			ParentProjectID: Ptr(int64(10)),
		}
		can, err := p.CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.True(t, can)
	})
	t.Run("duplicate to root denied for unlisted user", func(t *testing.T) {
		files.InitTestFileFixtures(t)
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user1"})

		can, err := (&ProjectDuplicate{ProjectID: 3}).CanCreate(s, &user.User{ID: 3})
		require.NoError(t, err)
		assert.False(t, can)
	})
	t.Run("duplicate to root allowed for listed user", func(t *testing.T) {
		files.InitTestFileFixtures(t)
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user3"})

		can, err := (&ProjectDuplicate{ProjectID: 3}).CanCreate(s, &user.User{ID: 3})
		require.NoError(t, err)
		assert.True(t, can)
	})
	t.Run("duplicate under writable parent still allowed", func(t *testing.T) {
		files.InitTestFileFixtures(t)
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user2"})

		pd := &ProjectDuplicate{
			ProjectID:       1,
			ParentProjectID: 10,
		}
		can, err := pd.CanCreate(s, &user.User{ID: 1})
		require.NoError(t, err)
		assert.True(t, can)
	})
	t.Run("detach to root denied for unlisted user", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user2"})

		project := Project{
			ID:              21,
			Title:           "Test21",
			ParentProjectID: Ptr(int64(0)),
		}
		err := project.Update(s, &user.User{ID: 1})
		require.Error(t, err)
		assert.True(t, IsErrGenericForbidden(err))
	})
	t.Run("detach to root allowed for listed user", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user1"})

		project := Project{
			ID:              21,
			Title:           "Test21",
			ParentProjectID: Ptr(int64(0)),
		}
		require.NoError(t, project.Update(s, &user.User{ID: 1}))
	})
	t.Run("inbox project on user creation is not restricted", func(t *testing.T) {
		db.LoadAndAssertFixtures(t)
		s := db.NewSession()
		defer s.Close()
		setProjectCreators(t, []string{"user2"})

		require.NoError(t, CreateNewProjectForUser(s, &user.User{ID: 1}))
	})
}

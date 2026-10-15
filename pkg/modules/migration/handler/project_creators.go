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
	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/models"
	user2 "code.vikunja.io/api/pkg/user"
)

// requireProjectCreator rejects an import up front: every migrator creates
// top-level projects, and the queued job would otherwise only fail later.
func requireProjectCreator(u *user2.User) error {
	s := db.NewSession()
	defer s.Close()

	can, err := models.IsProjectCreator(s, u)
	if err != nil {
		return err
	}
	if !can {
		return models.ErrGenericForbidden{}
	}
	return nil
}

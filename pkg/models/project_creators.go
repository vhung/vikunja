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
	"strings"

	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/user"
	"code.vikunja.io/api/pkg/web"

	"xorm.io/xorm"
)

// projectCreators returns the configured service.projectcreators usernames.
// An env value arrives as one string and viper only splits it on whitespace,
// so commas are split here as well.
func projectCreators() []string {
	var names []string
	for _, entry := range config.ServiceProjectCreators.GetStringSlice() {
		for _, name := range strings.FieldsFunc(entry, func(r rune) bool {
			return r == ',' || r == ' ' || r == '\t' || r == '\n'
		}) {
			names = append(names, name)
		}
	}
	return names
}

// IsProjectCreator reports whether a may create top-level projects.
// An empty service.projectcreators keeps the upstream behaviour (everyone may).
func IsProjectCreator(s *xorm.Session, a web.Auth) (bool, error) {
	names := projectCreators()
	if len(names) == 0 {
		return true, nil
	}

	authUser, is := a.(*user.User)
	if !is {
		return false, nil
	}

	// JWT claims may carry a stale username after a rename.
	u, err := user.GetUserByID(s, authUser.ID)
	if err != nil {
		return false, err
	}

	for _, name := range names {
		// Exact match: usernames are unique case-sensitively, so "HungLV" can be
		// registered next to "hunglv" and must not inherit its permission.
		if name == u.Username {
			return true, nil
		}
	}
	return false, nil
}

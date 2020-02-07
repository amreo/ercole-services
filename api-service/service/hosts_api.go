// Copyright (c) 2019 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package service is a package that provides methods for querying data
package service

import (
	"strings"
	"time"

	"github.com/amreo/ercole-services/utils"
)

// SearchHosts search hosts
func (as *APIService) SearchHosts(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.SearchHosts(full, strings.Split(search, " "), sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}

// GetHost return the host specified in the hostname param
func (as *APIService) GetHost(hostname string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetHost(hostname, olderThan)
}
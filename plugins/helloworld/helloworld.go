/**
    This file is part of gomapper.

    Gomapper is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Gomapper is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with gomapper.  If not, see <http://www.gnu.org/licenses/>.
**/

package main

import (
	"fmt"
	"github.com/tinyzimmer/gomapper/gomapperdb"
	"github.com/tinyzimmer/gomapper/nmapresult"
)

type DatabaseConnector interface{}

func init() {
	fmt.Println("I was loaded")
}

func OnScanComplete(nmapRun *nmapresult.NmapRun, db interface{}) error {
	fmt.Println("I saw that")
	_, ok := db.(gomapperdb.MemoryDatabase)
	if !ok {
		fmt.Println("Invalid database")
	}
	return nil
}

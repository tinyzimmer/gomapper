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

package gomapperdb

import (
	"net"

	"github.com/hashicorp/go-memdb"
	"github.com/tinyzimmer/gomapper/formats"
	"github.com/tinyzimmer/gomapper/logging"
	"github.com/tinyzimmer/gomapper/netutils"
	"github.com/tinyzimmer/gomapper/plugininterface"
)

type dbnetwork formats.DbNetwork

type MemoryDatabase struct {
	LocalAddr net.IP
	Session   *memdb.MemDB
	Plugins   plugininterface.LoadedPlugins
}

func GetMemoryDatabase(plugins plugininterface.LoadedPlugins, addr net.IP) (MemoryDatabase, error) {
	database := MemoryDatabase{}
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"dbnetwork": &memdb.TableSchema{
				Name: "dbnetwork",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Subnet"},
					},
				},
			},
		},
	}
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return database, err
	}
	database.Session = db
	database.Plugins = plugins
	return database, nil
}

func (d MemoryDatabase) GetNetwork(network string) (result *formats.DbNetwork, err error) {
	txn := d.Session.Txn(false)
	defer txn.Abort()
	raw, err := txn.First("dbnetwork", "id", network)
	if err != nil {
		logging.LogError(err.Error())
		return
	}
	if raw == nil {
		return
	}
	result = raw.(*formats.DbNetwork)
	return
}

func (d MemoryDatabase) GetAllNetworks() (networks []formats.DbNetwork) {
	txn := d.Session.Txn(false)
	defer txn.Abort()
	iter, err := txn.Get("dbnetwork", "id")
	checkError(err)
	for {
		raw := iter.Next()
		if raw == nil {
			break
		} else {
			network := raw.(formats.DbNetwork)
			networks = append(networks, network)
		}
	}
	return
}

func (d MemoryDatabase) AddNetwork(network formats.DbNetwork) {
	txn := d.Session.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("dbnetwork", network); err != nil {
		logging.LogError(err.Error())
	} else {
		txn.Commit()
	}
}

func checkError(err error) {
	if err != nil {
		logging.LogError(err.Error())
	}
}

func checkDbNetwork(network string, ip string) (subnet string) {
	if network != "" {
		subnet = network
	} else {
		subnet = netutils.FormatDefaultNetmask(ip)
	}
	return
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

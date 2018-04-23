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
	"github.com/hashicorp/go-memdb"
)

type Network struct {
	Subnet string
	Hosts  []*DbHost
}

type DbHost struct {
	IP       string
	MAC      string
	Services []*DbService
}

type DbService struct {
	Port string
	Name string
}

type MemoryDatabase struct {
	Session *memdb.MemDB
}

func getMemoryDatabase() (MemoryDatabase, error) {
	database := MemoryDatabase{}
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"network": &memdb.TableSchema{
				Name: "network",
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
	return database, nil
}

func (d MemoryDatabase) AddScanResultsByNetwork(network string, results *NmapRun) {
	txn := d.Session.Txn(true)
	defer txn.Abort()
	net, err := d.GetNetwork(network)
	if net == nil || err != nil {
		net = &Network{}
		net.Subnet = network
	}
	for _, host := range results.Hosts {
		ip, mac := getScanAddrs(host.Addresses)
		dbHost := &DbHost{}
		dbHost.IP = ip
		dbHost.MAC = mac
		for _, port := range host.Ports.Ports {
			service := &DbService{}
			service.Port = fmt.Sprintf("%s/%s", port.PortId, port.Protocol)
			service.Name = port.Service.Name
			dbHost.Services = append(dbHost.Services, service)
		}
		net.Hosts = append(net.Hosts, dbHost)
	}
	if err := txn.Insert("network", net); err != nil {
		logError(err.Error())
	} else {
		txn.Commit()
	}
}

func (d MemoryDatabase) GetNetwork(network string) (result *Network, err error) {
	txn := d.Session.Txn(false)
	defer txn.Abort()
	raw, err := txn.First("network", "id", network)
	if err != nil {
		logError(err.Error())
		return
	}
	if raw == nil {
		return
	}
	result = raw.(*Network)
	return
}

func (d MemoryDatabase) GetAllNetworks() (networks []*Network) {
	txn := d.Session.Txn(false)
	defer txn.Abort()
	iter, err := txn.Get("network", "id")
	checkError(err)
	for {
		raw := iter.Next()
		if raw == nil {
			break
		} else {
			network := raw.(*Network)
			networks = append(networks, network)
		}
	}
	return
}

func (d MemoryDatabase) AddNetwork(network string) {
	txn := d.Session.Txn(true)
	defer txn.Abort()
	n := &Network{}
	n.Subnet = network
	if err := txn.Insert("network", n); err != nil {
		logError(err.Error())
	} else {
		txn.Commit()
	}
}

func checkError(err error) {
	if err != nil {
		logError(err.Error())
	}
}

func checkDbNetwork(network string, ip string) (subnet string) {
	if network != "" {
		subnet = network
	} else {
		subnet = formatDefaultNetmask(ip)
	}
	return
}

func getScanAddrs(addrs []Address) (ip string, mac string) {
	for _, addr := range addrs {
		if addr.AddrType == "ipv4" {
			ip = addr.Addr
		} else if addr.AddrType == "mac" {
			mac = addr.Addr
		}
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

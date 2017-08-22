// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodbadapter

import (
	"errors"
	"runtime"

	"github.com/casbin/casbin/model"
	"gopkg.in/mgo.v2"
)

// CasbinRule represents a rule in Casbin.
type CasbinRule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

// Adapter represents the MongoDB adapter for policy storage.
type Adapter struct {
	url        string
	session    *mgo.Session
	collection *mgo.Collection
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
	a.close()
}

// NewAdapter is the constructor for Adapter. If database name is not provided
// in the Mongo URL, 'casbin' will be used as database name.
func NewAdapter(url string) *Adapter {
	a := &Adapter{}
	a.url = url

	// Open the DB, create it if not existed.
	a.open()

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a
}

func (a *Adapter) createIndice() {
	var err error

	index := mgo.Index{
		Key: []string{"ptype"},
	}
	err = a.collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key: []string{"v0"},
	}
	err = a.collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key: []string{"v1"},
	}
	err = a.collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key: []string{"v2"},
	}
	err = a.collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key: []string{"v3"},
	}
	err = a.collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key: []string{"v4"},
	}
	err = a.collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key: []string{"v5"},
	}
	err = a.collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func (a *Adapter) open() {
	dI, err := mgo.ParseURL(a.url)
	if err != nil {
		panic(err)
	}

	if dI.Database == "" {
		dI.Database = "casbin"
	}

	session, err := mgo.DialWithInfo(dI)
	if err != nil {
		panic(err)
	}

	db := session.DB(dI.Database)
	collection := db.C("casbin_rule")

	a.session = session
	a.collection = collection

	a.createIndice()
}

func (a *Adapter) close() {
	a.session.Close()
}

func (a *Adapter) createTable() {
}

func (a *Adapter) dropTable() {
	if a.collection == nil {
		return
	}

	err := a.collection.DropCollection()
	if err != nil {
		if err.Error() != "ns not found" {
			panic(err)
		}
	}
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	key := line.PType
	sec := key[:1]

	tokens := []string{}
	if line.V0 != "" {
		tokens = append(tokens, line.V0)
	} else {
		goto LineEnd
	}

	if line.V1 != "" {
		tokens = append(tokens, line.V1)
	} else {
		goto LineEnd
	}

	if line.V2 != "" {
		tokens = append(tokens, line.V2)
	} else {
		goto LineEnd
	}

	if line.V3 != "" {
		tokens = append(tokens, line.V3)
	} else {
		goto LineEnd
	}

	if line.V4 != "" {
		tokens = append(tokens, line.V4)
	} else {
		goto LineEnd
	}

	if line.V5 != "" {
		tokens = append(tokens, line.V5)
	} else {
		goto LineEnd
	}

LineEnd:
	model[sec][key].Policy = append(model[sec][key].Policy, tokens)
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	line := CasbinRule{}
	iter := a.collection.Find(nil).Iter()
	for iter.Next(&line) {
		loadPolicyLine(line, model)
	}

	if err := iter.Close(); err != nil {
		return err
	}
	return nil
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{}

	line.PType = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) error {
	a.dropTable()
	a.createTable()

	var lines []interface{}

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, &line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, &line)
		}
	}

	err := a.collection.Insert(lines...)
	return err
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}

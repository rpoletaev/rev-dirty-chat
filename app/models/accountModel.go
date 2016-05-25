// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

/*
	BuoyModels contains the models for the buoy service
*/
package models

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	Account struct {
		ID             bson.ObjectId `bson:"_id,omitempty"`
		Email          string        `bson:"email"`
		HashedPassword string        `bson:"password"`
		Login          string        `bson:"login"`
		IsAdmin        bool          `bson:"isadmin"`
	}
)

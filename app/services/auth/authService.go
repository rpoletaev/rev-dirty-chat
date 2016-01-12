// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

/*
	Buoy implements the service for the buoy functionality
*/
package auth

import (
	"github.com/rpoletaev/rev-dirty-chat/app/models"
	"github.com/rpoletaev/rev-dirty-chat/app/services"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//** PUBLIC FUNCTIONS

func FindUserByEmail(service *services.Service, email string) (user *models.User, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindUserByEmail")

	tracelog.STARTED(service.UserId, "FindUserByEmail")

	queryMap := bson.M{"email": email}
	tracelog.TRACE(helper.MAIN_GO_ROUTINE, "FindUserByEmail", "Query : %s", mongo.ToString(queryMap))

	// Execute the query
	user = &models.User{}
	err = service.DBAction("users",
		func(collection *mgo.Collection) error {
			return collection.Find(queryMap).One(user)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindUserByEmail")
		return user, err
	}

	tracelog.COMPLETED(service.UserId, "FindUserByEmail")
	return user, err
}

func VerifyPassword(password string, user *models.User) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err
}

func SetUserPassword(user *models.User) (err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
	user.HashedPassword = string(hash)
	return err
}

func InsertUser(service *services.Service, user *models.User) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "InsertUser")

	SetUserPassword(user)

	err = service.DBAction("users",
		func(collection *mgo.Collection) error {
			return collection.Insert(user)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "InsertUser")
		return err
	}

	tracelog.COMPLETED(service.UserId, "InsertUser")
	return nil
}

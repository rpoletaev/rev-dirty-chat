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

func FindAccountByEmail(service *services.Service, email string) (account *models.Account, err error) {
	defer helper.CatchPanic(&err, service.UserId, "FindAccountByEmail")

	tracelog.STARTED(service.UserId, "FindAccountByEmail")

	queryMap := bson.M{"email": email}
	tracelog.TRACE(helper.MAIN_GO_ROUTINE, "FindAccountByEmail", "Query : %s", mongo.ToString(queryMap))

	// Execute the query
	account = &models.Account{}
	err = service.DBAction("accounts",
		func(collection *mgo.Collection) error {
			return collection.Find(queryMap).One(account)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "FindAccountByEmail")
		return account, err
	}

	tracelog.COMPLETED(service.UserId, "FindAccountByEmail")
	return account, err
}

func VerifyPassword(password string, account *models.Account) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(account.HashedPassword), []byte(password))
	return err
}

func SetAccountPassword(account *models.Account) (err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(account.HashedPassword), bcrypt.DefaultCost)
	account.HashedPassword = string(hash)
	return err
}

func InsertAccount(service *services.Service, account *models.Account) (err error) {
	defer helper.CatchPanic(&err, service.UserId, "InsertAccount")

	SetAccountPassword(account)

	err = service.DBAction("accounts",
		func(collection *mgo.Collection) error {
			return collection.Insert(account)
		})

	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "InsertAccount")
		return err
	}

	tracelog.COMPLETED(service.UserId, "InsertAccount")
	return nil
}

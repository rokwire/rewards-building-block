/*
 *   Copyright (c) 2020 Board of Trustees of the University of Illinois.
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package core

import (
	"log"
	cacheadapter "rewards/driven/cache"
	"rewards/driven/storage"
)

//Application represents the core application code based on hexagonal architecture
type Application struct {
	version string
	build   string

	Services Services //expose to the drivers adapters

	storage      Storage
	cacheAdapter *cacheadapter.CacheAdapter
}

// Start starts the core part of the application
func (app *Application) Start() {
	err := app.storeMultiTenancyData()
	if err != nil {
		log.Fatalf("error initializing multi-tenancy data: %s", err.Error())
	}

	app.storage.SetListener(app)
}

//as the service starts supporting multi-tenancy we need to add the needed multi-tenancy fields for the existing data,
func (app *Application) storeMultiTenancyData() error {
	log.Println("storeMultiTenancyData...")
	//in transaction
	transaction := func(context storage.TransactionContext) error {

		//check if we need to apply multi-tenancy data
		var applyData bool
		items, err := app.storage.FindAllRewardTypeItems()
		if err != nil {
			return err
		}
		for _, current := range items {
			if len(current.AppID) > 0 {
				log.Printf("\thas already app_id:%s", current.AppID)
				applyData = false
				break
			} else {
				log.Print("\tno app_id")
				applyData = true
				break
			}
		}

		//apply data if necessary
		if applyData {
			log.Print("\tapplying multi-tenancy data..")
			//TODO
		} else {
			log.Print("\tno need to apply multi-tenancy data, so do nothing")
		}

		return nil

	}

	err := app.storage.PerformTransaction(transaction)
	if err != nil {
		log.Printf("error performing transaction for multi tenancy")
		return err
	}
	return nil
}

// NewApplication creates new Application
func NewApplication(version string, build string, storage Storage, cacheadapter *cacheadapter.CacheAdapter) *Application {
	application := Application{
		version:      version,
		build:        build,
		storage:      storage,
		cacheAdapter: cacheadapter}

	// add the drivers ports/interfaces
	application.Services = &servicesImpl{app: &application}

	return &application
}

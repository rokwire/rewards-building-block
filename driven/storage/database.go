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

package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	mongoDBAuth  string
	mongoDBName  string
	mongoTimeout time.Duration

	db       *mongo.Database
	dbClient *mongo.Client

	rewardTypes *collectionWrapper
	rewardPools *collectionWrapper
}

func (m *database) start() error {

	log.Println("database -> start")

	//connect to the database
	clientOptions := options.Client().ApplyURI(m.mongoDBAuth)
	connectContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	client, err := mongo.Connect(connectContext, clientOptions)
	cancel()
	if err != nil {
		return err
	}

	//ping the database
	pingContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	err = client.Ping(pingContext, nil)
	cancel()
	if err != nil {
		return err
	}

	//apply checks
	db := client.Database(m.mongoDBName)

	rewardTypes := &collectionWrapper{database: m, coll: db.Collection("reward_types")}
	err = m.applyRewardTypesChecks(rewardTypes)
	if err != nil {
		return err
	}

	rewardPools := &collectionWrapper{database: m, coll: db.Collection("reward_pools")}
	err = m.applyRewardPoolsChecks(rewardPools)
	if err != nil {
		return err
	}

	//asign the db, db client and the collections
	m.db = db
	m.dbClient = client

	m.rewardTypes = rewardTypes
	m.rewardPools = rewardPools

	return nil
}

func (m *database) applyRewardTypesChecks(posts *collectionWrapper) error {
	log.Println("apply reward_types checks.....")

	indexes, _ := posts.ListIndexes()
	indexMapping := map[string]interface{}{}
	if indexes != nil {

		for _, index := range indexes {
			name := index["name"].(string)
			indexMapping[name] = index
		}
	}

	if indexMapping["name_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "name", Value: 1},
			}, true)
		if err != nil {
			return err
		}
	}

	if indexMapping["building_block_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "building_block", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	log.Println("reward_types checks passed")
	return nil
}

func (m *database) applyRewardPoolsChecks(posts *collectionWrapper) error {
	log.Println("apply reward_pools checks.....")

	indexes, _ := posts.ListIndexes()
	indexMapping := map[string]interface{}{}
	if indexes != nil {

		for _, index := range indexes {
			name := index["name"].(string)
			indexMapping[name] = index
		}
	}

	if indexMapping["code_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "code_1", Value: 1},
			}, true)
		if err != nil {
			return err
		}
	}

	if indexMapping["building_block_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "building_block", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	if indexMapping["active_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "active", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	if indexMapping["date_created_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "date_created", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	log.Println("reward_pools checks passed")
	return nil
}

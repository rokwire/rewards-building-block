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

// Listener listens for storage updates
type Listener interface {
	OnRewardTypesChanged()
}

type database struct {
	listener Listener

	mongoDBAuth  string
	mongoDBName  string
	mongoTimeout time.Duration

	db       *mongo.Database
	dbClient *mongo.Client

	rewardTypes       *collectionWrapper
	rewardInventories *collectionWrapper
	rewardHistory     *collectionWrapper
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
	go rewardTypes.Watch(nil)

	rewardInventories := &collectionWrapper{database: m, coll: db.Collection("reward_inventories")}
	err = m.applyRewardInventoriesChecks(rewardInventories)
	if err != nil {
		return err
	}

	rewardHistory := &collectionWrapper{database: m, coll: db.Collection("reward_history")}
	err = m.applyRewardHistoryChecks(rewardHistory)
	if err != nil {
		return err
	}

	//asign the db, db client and the collections
	m.db = db
	m.dbClient = client

	m.rewardTypes = rewardTypes
	m.rewardInventories = rewardInventories
	m.rewardHistory = rewardHistory

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

	if indexMapping["building_block_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "building_block", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	if indexMapping["reward_type_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "reward_type", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	log.Println("reward_types checks passed")
	return nil
}

func (m *database) applyRewardInventoriesChecks(posts *collectionWrapper) error {
	log.Println("apply reward_inventories checks.....")

	indexes, _ := posts.ListIndexes()
	indexMapping := map[string]interface{}{}
	if indexes != nil {
		for _, index := range indexes {
			name := index["name"].(string)
			indexMapping[name] = index
		}
	}

	if indexMapping["org_id_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "org_id", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	if indexMapping["reward_type_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "reward_type", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	if indexMapping["depleted_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "depleted", Value: 1},
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

	log.Println("reward_inventories checks passed")
	return nil
}

func (m *database) applyRewardHistoryChecks(posts *collectionWrapper) error {
	log.Println("apply reward_history checks.....")

	indexes, _ := posts.ListIndexes()
	indexMapping := map[string]interface{}{}
	if indexes != nil {

		for _, index := range indexes {
			name := index["name"].(string)
			indexMapping[name] = index
		}
	}

	if indexMapping["user_id_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "user_id", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	if indexMapping["pool_id_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "pool_id", Value: 1},
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

	log.Println("reward_history checks passed")
	return nil
}

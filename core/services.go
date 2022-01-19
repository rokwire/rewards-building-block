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

import "rewards/core/model"

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getRewardTypes(ids []string) ([]model.RewardType, error) {
	return app.storage.GetRewardTypes(ids)
}

func (app *Application) getRewardType(id string) (*model.RewardType, error) {
	return app.storage.GetRewardType(id)
}

func (app *Application) createRewardType(item model.RewardType) (*model.RewardType, error) {
	return app.storage.CreateRewardType(item)
}

func (app *Application) updateRewardType(id string, item model.RewardType) (*model.RewardType, error) {
	return app.storage.UpdateRewardType(id, item)
}

func (app *Application) deleteGetRewardTypes(id string) error {
	return app.storage.DeleteRewardTypes(id)
}

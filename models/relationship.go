// Copyright 2016 zm@huantucorp.com
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
/*
                   _ooOoo_
                  o8888888o
                  88" . "88
                  (| -_- |)
                  O\  =  /O
               ____/`---'\____
             .'  \\|     |//  `.
            /  \\|||  :  |||//  \
           /  _||||| -:- |||||-  \
           |   | \\\  -  /// |   |
           | \_|  ''\---/''  |   |
           \  .-\__  `-`  ___/-. /
         ___`. .'  /--.--\  `. . __
      ."" '<  `.___\_<|>_/___.'  >'"".
     | | :  `- \`.;`\ _ /`;.`/ - ` : | |
     \  \ `-.   \_ __\ /__ _/   .-` /  /
======`-.____`-.___\_____/___.-`____.-'======
                   `=---='
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
         佛祖保佑       永无BUG
*/
package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

const (
	UserLike    = "liked"
	UserMatch   = "matched"
	UserDislike = "disliked"
)

type UserRelationshipAdapter interface {
	NewRelationship(userId, otherUserId, state string) (*UserRelationship, error)
	ListUserRelationships(userId string) ([]*UserRelationship, error)
}

func (this *userDefault) NewRelationship(userId, otherUserId, state string) (*UserRelationship, error) {
	if state != UserLike && state != UserDislike {
		return nil, errors.New("invalid params")
	}

	relationship := &UserRelationship{
		UserId:      userId,
		OtherUserId: otherUserId,
		Type:        "relationship",
		State:       state,
	}

	if state == UserLike {
		otherRelationship, err := this.getUserRelationship(otherUserId, userId)
		if err == nil {
			if otherRelationship.State == state {
				relationship.State = UserMatch
			}
		}
	}

	if err := this.db.Create(relationship).Error; err != nil {
		return nil, err
	}

	return relationship, nil
}

func (this *userDefault) ListUserRelationships(userId string) ([]*UserRelationship, error) {
	if userId == "" {
		return nil, errors.New("invalid params")
	}

	relationships := []*UserRelationship{}
	if err := this.db.Where(&UserRelationship{UserId: userId}).Find(&relationships).Error; err != nil {
		return nil, err
	}

	return relationships, nil
}

func (this *userDefault) getUserRelationship(userId, otherUserId string) (*UserRelationship, error) {
	if userId == "" || otherUserId == "" {
		return nil, errors.New("invalid params")
	}

	relationship := &UserRelationship{}
	if err := this.db.Where(&UserRelationship{UserId: userId, OtherUserId: otherUserId}).First(&relationship).Error; err != nil {
		return nil, err
	}

	return relationship, nil
}

type UserRelationship struct {
	gorm.Model
	UserId      string `gorm:"index" json:"user_id"`
	OtherUserId string `gorm:"index" json:"other_user_id"`
	Type        string `gorm:"index" json:"type"`
	State       string `gorm:"index" json:"state"`
}

func (this *UserRelationship) TableName() string {
	return "user_relationship"
}

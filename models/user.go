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
	"github.com/aiwuTech/devKit/random"
	"github.com/jinzhu/gorm"
)

type UserAdapter interface {
	NewUser(name string) (*User, error)
	ListUsers() ([]*User, error)
	DBAdapter
	UserRelationshipAdapter
}

type userDefault struct {
	db *gorm.DB
}

func NewUserAdapter(db *gorm.DB) (UserAdapter, error) {
	user := &userDefault{
		db: db,
	}

	if err := user.SyncDB(); err != nil {
		return nil, err
	}

	return user, nil
}

func (this *userDefault) SyncDB() error {
	tx := this.db.Begin()
	{
		user := &User{}
		if tx.HasTable(user) {
			if err := tx.AutoMigrate(user).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.CreateTable(user).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		userRelation := &UserRelationship{}
		if tx.HasTable(userRelation) {
			if err := tx.AutoMigrate(userRelation).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.CreateTable(userRelation).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (this *userDefault) NewUser(name string) (*User, error) {
	if name == "" {
		return nil, errors.New("invalid params")
	}

	user := &User{
		Id:   this.getUniqueUserId(),
		Name: name,
		Type: "user",
	}

	if err := this.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// this func used for generate unique user id, at this moment, just for testing
func (this *userDefault) getUniqueUserId() string {
	return random.RandomNumeric(11)
}

func (this *userDefault) ListUsers() ([]*User, error) {
	users := []*User{}
	if err := this.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

type User struct {
	Id   string `gorm:"primary_key" json:"id"`
	Name string `gorm:"unique" json:"name"`
	Type string `gorm:"index" json:"type"`
}

func (this *User) TableName() string {
	return "user"
}

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
package handlers

import (
	"errors"
	"github.com/aiwuTech/tantanBackendTest/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

func GetUserRelationship(param martini.Params, r render.Render) {
	userId := param["user_id"]

	relationships, err := userAdapter.ListUserRelationships(userId)
	if err != nil {
		r.JSON(http.StatusInternalServerError, NewRspError(err))
		return
	}

	data := []map[string]interface{}{}
	for _, relationship := range relationships {
		data = append(data, map[string]interface{}{
			"user_id": relationship.OtherUserId,
			"state":   relationship.State,
			"type":    relationship.Type,
		})
	}

	r.JSON(http.StatusOK, data)
}

type PutUserRelationshipReqParam struct {
	State string `form:"state" json:"state" binding:"required"`
}

func PutUserRelationship(param martini.Params, reqParam PutUserRelationshipReqParam, r render.Render) {
	if reqParam.State != models.UserLike && reqParam.State != models.UserDislike {
		r.JSON(http.StatusBadRequest, NewRspError(errors.New("state = 'liked'|'disliked'")))
		return
	}

	userId := param["user_id"]
	otherUserId := param["other_user_id"]
	relationship, err := userAdapter.NewRelationship(userId, otherUserId, reqParam.State)
	if err != nil {
		r.JSON(http.StatusInternalServerError, NewRspError(err))
		return
	}

	r.JSON(http.StatusOK, map[string]interface{}{
		"user_id": otherUserId,
		"state":   relationship.State,
		"type":    relationship.Type,
	})
}

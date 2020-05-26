// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package issueapi

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikehelmick/tek-verification-server/pkg/api"
	"github.com/mikehelmick/tek-verification-server/pkg/controller"
	"github.com/mikehelmick/tek-verification-server/pkg/database"
)

type IssueAPI struct {
	database database.Database
}

func New(db database.Database) controller.Controller {
	return &IssueAPI{db}
}

func (iapi *IssueAPI) Execute(c *gin.Context) {
	var request api.IssuePINRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := api.IssuePINResponse{}

	// Generate PIN
	source := make([]byte, 6)
	_, err := rand.Read(source)
	if err != nil {
		response.Error = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	pinCode := base64.RawStdEncoding.EncodeToString(source)

	claims := map[string]string{
		"lab":   "test r us",
		"batch": "test batch number",
	}
	_, err = iapi.database.InsertPIN(pinCode, request.Risks, claims, request.ValidFor)
	if err != nil {
		response.Error = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.PIN = pinCode
	c.JSON(http.StatusOK, response)
}
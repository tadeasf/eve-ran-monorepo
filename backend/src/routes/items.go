// Copyright (C) 2024 Tadeáš Fořt
// 
// This file is part of EVE Ran Services.
// 
// EVE Ran Services is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// EVE Ran Services is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with EVE Ran Services.  If not, see <https://www.gnu.org/licenses/>.

package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db/queries"
	"github.com/tadeasf/eve-ran/src/services"
)

func FetchAndStoreItems(c *gin.Context) {
	items, err := services.FetchAllItems(50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, item := range items {
		err = queries.UpsertESIItem(item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Items fetched and stored successfully", "count": len(items)})
}

func GetAllItems(c *gin.Context) {
	items, err := queries.GetAllESIItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func GetItemByTypeID(c *gin.Context) {
	typeID := c.Param("typeID")
	itemTypeID, err := strconv.Atoi(typeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type ID"})
		return
	}

	item, err := queries.GetESIItemByTypeID(itemTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

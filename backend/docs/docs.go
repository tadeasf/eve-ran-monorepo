// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/characters": {
            "get": {
                "description": "Fetch all characters from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "characters"
                ],
                "summary": "Get all characters",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Character"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a new character ID to the database and fetch all kills\nAdd a new character ID to the database and fetch all kills",
                "consumes": [
                    "application/json",
                    "application/json"
                ],
                "produces": [
                    "application/json",
                    "application/json"
                ],
                "tags": [
                    "characters",
                    "characters"
                ],
                "summary": "Add a new character ID",
                "parameters": [
                    {
                        "description": "Character ID",
                        "name": "character",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Character"
                        }
                    },
                    {
                        "description": "Character ID",
                        "name": "character",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Character"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Character"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/characters/stats": {
            "get": {
                "description": "Fetch stats for all characters from the database with optional filters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "characters"
                ],
                "summary": "Get all character stats",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "collectionFormat": "csv",
                        "description": "Region IDs",
                        "name": "regionID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Start date (YYYY-MM-DD)",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date (YYYY-MM-DD)",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.CharacterStats"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/characters/{id}": {
            "delete": {
                "description": "Remove a character from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "characters"
                ],
                "summary": "Remove a character",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Character ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/characters/{id}/kills": {
            "get": {
                "description": "Fetch and store kills for a character from zKillboard",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "characters"
                ],
                "summary": "Get character kills",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Character ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Kill"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/characters/{id}/kills/db": {
            "get": {
                "description": "Fetch kills for a character from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "characters"
                ],
                "summary": "Get character kills from database",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Character ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page size",
                        "name": "pageSize",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PaginatedResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/kills": {
            "get": {
                "description": "Fetch all kills from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "kills"
                ],
                "summary": "Get all kills",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Kill"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/kills/region/{regionID}": {
            "get": {
                "description": "Fetch kills for a region from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "kills"
                ],
                "summary": "Get kills by region",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Region ID",
                        "name": "regionID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page size",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Start date (YYYY-MM-DD)",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date (YYYY-MM-DD)",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PaginatedResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/regions": {
            "get": {
                "description": "Fetch all regions from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "regions"
                ],
                "summary": "Get all regions",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Region"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.CharacterStats": {
            "type": "object",
            "properties": {
                "character_id": {
                    "type": "integer"
                },
                "kill_count": {
                    "type": "integer"
                },
                "total_isk": {
                    "type": "number"
                }
            }
        },
        "models.Attacker": {
            "type": "object",
            "properties": {
                "alliance_id": {
                    "type": "integer"
                },
                "character_id": {
                    "type": "integer"
                },
                "corporation_id": {
                    "type": "integer"
                },
                "damage_done": {
                    "type": "integer"
                },
                "faction_id": {
                    "type": "integer"
                },
                "final_blow": {
                    "type": "boolean"
                },
                "security_status": {
                    "type": "number"
                },
                "ship_type_id": {
                    "type": "integer"
                },
                "weapon_type_id": {
                    "type": "integer"
                }
            }
        },
        "models.Character": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "race_id": {
                    "type": "integer"
                },
                "security_status": {
                    "type": "number"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "models.Item": {
            "type": "object",
            "properties": {
                "flag": {
                    "type": "integer"
                },
                "item_type_id": {
                    "type": "integer"
                },
                "quantity_destroyed": {
                    "type": "integer"
                },
                "quantity_dropped": {
                    "type": "integer"
                },
                "singleton": {
                    "type": "integer"
                }
            }
        },
        "models.Kill": {
            "type": "object",
            "properties": {
                "attackers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Attacker"
                    }
                },
                "awox": {
                    "type": "boolean"
                },
                "character_id": {
                    "type": "integer"
                },
                "destroyed_value": {
                    "type": "number"
                },
                "dropped_value": {
                    "type": "number"
                },
                "fitted_value": {
                    "type": "number"
                },
                "hash": {
                    "type": "string"
                },
                "killmail_id": {
                    "type": "integer"
                },
                "killmail_time": {
                    "type": "string"
                },
                "locationID": {
                    "type": "integer"
                },
                "npc": {
                    "type": "boolean"
                },
                "points": {
                    "type": "integer"
                },
                "solar_system_id": {
                    "type": "integer"
                },
                "solo": {
                    "type": "boolean"
                },
                "total_value": {
                    "type": "number"
                },
                "victim": {
                    "$ref": "#/definitions/models.Victim"
                }
            }
        },
        "models.PaginatedResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "page": {
                    "type": "integer"
                },
                "pageSize": {
                    "type": "integer"
                },
                "totalItems": {
                    "type": "integer"
                },
                "totalPages": {
                    "type": "integer"
                }
            }
        },
        "models.Position": {
            "type": "object",
            "properties": {
                "x": {
                    "type": "number"
                },
                "y": {
                    "type": "number"
                },
                "z": {
                    "type": "number"
                }
            }
        },
        "models.Region": {
            "type": "object",
            "properties": {
                "constellations": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "region_id": {
                    "type": "integer"
                }
            }
        },
        "models.Victim": {
            "type": "object",
            "properties": {
                "alliance_id": {
                    "type": "integer"
                },
                "character_id": {
                    "type": "integer"
                },
                "corporation_id": {
                    "type": "integer"
                },
                "damage_taken": {
                    "type": "integer"
                },
                "faction_id": {
                    "type": "integer"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Item"
                    }
                },
                "position": {
                    "$ref": "#/definitions/models.Position"
                },
                "ship_type_id": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "EVE Ran API",
	Description:      "This is the API for EVE Ran application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

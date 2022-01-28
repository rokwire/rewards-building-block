// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/reward_pool": {
            "post": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Create a new reward pool",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminCreateRewardPool",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/RewardPool"
                        }
                    }
                }
            }
        },
        "/admin/reward_pool/{id}": {
            "put": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Updates a reward pool with the specified id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminUpdateRewardPool",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/RewardPool"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Deletes a reward pool with the specified id",
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminDeleteRewardPool",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/admin/reward_pools": {
            "get": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Retrieves  all reward types",
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminGetRewardPools",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Coma separated IDs of the desired records",
                        "name": "ids",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/RewardPool"
                            }
                        }
                    }
                }
            }
        },
        "/admin/reward_pools/{id}": {
            "get": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Retrieves a reward pool by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminRewardPool",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/RewardPool"
                        }
                    }
                }
            }
        },
        "/admin/reward_types": {
            "get": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Retrieves  all reward types",
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminGetRewardTypes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Coma separated IDs of the desired records",
                        "name": "ids",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/RewardType"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Create a new reward type",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminCreateRewardType",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/RewardType"
                        }
                    }
                }
            }
        },
        "/admin/reward_types/{id}": {
            "get": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Retrieves a reward type by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminRewardTypes",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/RewardType"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Updates a reward type with the specified id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminUpdateRewardType",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/RewardType"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Deletes a reward type with the specified id",
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminDeleteRewardType",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/int/reward_history": {
            "post": {
                "security": [
                    {
                        "InternalApiAuth": []
                    }
                ],
                "description": "Create a new reward history entry from another BB",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "InternalCreateRewardHistoryEntry",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/RewardHistoryEntry"
                        }
                    }
                }
            }
        },
        "/user/balance": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves balance for each user's wallet",
                "tags": [
                    "Client"
                ],
                "operationId": "GetUserBalance",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/version": {
            "get": {
                "description": "Gives the service version.",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "Version",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/wallet/{code}/balance": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves  the user balance",
                "tags": [
                    "Client"
                ],
                "operationId": "GetWalletBalance",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/wallet/{code}/history": {
            "get": {
                "security": [
                    {
                        "UserAuth": []
                    }
                ],
                "description": "Retrieves the user history",
                "tags": [
                    "Client"
                ],
                "operationId": "GetWalletHistory",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "RewardHistoryEntry": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "date_created": {
                    "type": "string"
                },
                "date_updated": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "description": "Do we need it here?",
                    "type": "string"
                },
                "pool_id": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "RewardPool": {
            "type": "object",
            "properties": {
                "active": {
                    "type": "boolean"
                },
                "amount": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/model.JSONData"
                },
                "date_created": {
                    "type": "string"
                },
                "date_updated": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "reward_code": {
                    "description": "illini_cash",
                    "type": "string"
                }
            }
        },
        "RewardType": {
            "type": "object",
            "properties": {
                "active": {
                    "type": "boolean"
                },
                "amount": {
                    "description": "5",
                    "type": "integer"
                },
                "building_block": {
                    "description": "\"content\"",
                    "type": "string"
                },
                "date_created": {
                    "type": "string"
                },
                "date_updated": {
                    "type": "string"
                },
                "display_name": {
                    "description": "Win five point by five readings",
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "description": "\"win_five_point_by_five_readings\"",
                    "type": "string"
                },
                "reward_code": {
                    "description": "illini_cash",
                    "type": "string"
                }
            }
        },
        "model.JSONData": {
            "type": "object",
            "additionalProperties": true
        }
    },
    "securityDefinitions": {
        "AdminGroupAuth": {
            "type": "apiKey",
            "name": "GROUP",
            "in": "header"
        },
        "AdminUserAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header (add Bearer prefix to the Authorization value)"
        },
        "InternalApiAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header (add INTERNAL-API-KEY with correct value as a header)"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.0.1",
	Host:        "localhost",
	BasePath:    "/content",
	Schemes:     []string{"https"},
	Title:       "Rewards Building Block API",
	Description: "RoRewards Building Block API Documentation.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}

// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "https://community.taiko.xyz/",
            "email": "info@taiko.xyz"
        },
        "license": {
            "name": "MIT",
            "url": "hhttps://github.com/taikoxyz/taiko-client/blob/main/LICENSE.md"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/assignment": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Try to accept a block proof assignment",
                "operationId": "create-assignment",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.ProposeBlockResponse"
                        }
                    },
                    "422": {
                        "description": "prover does not have capacity",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get current prover server status",
                "operationId": "get-status",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.Status"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "server.ProposeBlockResponse": {
            "type": "object",
            "properties": {
                "prover": {
                    "type": "string"
                },
                "signedPayload": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "server.Status": {
            "type": "object",
            "properties": {
                "currentCapacity": {
                    "type": "integer"
                },
                "maxExpiry": {
                    "type": "integer"
                },
                "minProofFee": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "prover-api.test.taiko.xyz",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Taiko Prover API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

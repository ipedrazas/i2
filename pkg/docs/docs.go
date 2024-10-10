// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Ivan Pedrazas",
            "url": "https://i2.alacasa.uk",
            "email": "ipedrazas@gmail.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/dns/:zone/entries": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "List DNS entries",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Cloud Provider",
                        "name": "provider",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dns.DNSEntry"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/dns/:zone/records": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a DNS record",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dns.DNSRecord"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/dns/:zone/records/:id": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Read a DNS record",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Cloud Provider",
                        "name": "provider",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dns.DNSRecord"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update a DNS record",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dns.DNSRecord"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Delete a DNS record",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Cloud Provider",
                        "name": "provider",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dns.DNSRecord"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/dns/ip/:ip": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Check if an IP is in use",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dns.DNSRecord"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/healtz/ready": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Check if the service is ready",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/proxmox/nodes": {
            "get": {
                "description": "Get cluster nodes",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "prxmx"
                ],
                "summary": "Get cluster nodes",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/proxmox/vms": {
            "get": {
                "description": "Get virtual machines",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "prxmx"
                ],
                "summary": "Get virtual machines",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/prxmx.Node"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dns.DNSEntry": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "domain": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "provider": {
                    "type": "string"
                },
                "ttl": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "dns.DNSRecord": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "provider": {
                    "type": "string"
                },
                "ttl": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "prxmx.Node": {
            "type": "object",
            "properties": {
                "ip": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "running": {
                    "type": "boolean"
                },
                "uptime": {
                    "$ref": "#/definitions/prxmx.Uptime"
                }
            }
        },
        "prxmx.Uptime": {
            "type": "object",
            "properties": {
                "days": {
                    "type": "integer"
                },
                "hours": {
                    "type": "integer"
                },
                "minutes": {
                    "type": "integer"
                },
                "raw": {
                    "type": "integer"
                },
                "seconds": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

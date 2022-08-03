// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/excel_template": {
            "get": {
                "description": "download excel table template",
                "produces": [
                    "application/vnd.ms-excel"
                ],
                "tags": [
                    "excel"
                ],
                "summary": "ExcelTemplate",
                "operationId": "excel_template",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of user",
                        "name": "X-User-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of supplier",
                        "name": "X-Supplier-Id",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errinfo.errInf"
                        }
                    }
                }
            }
        },
        "/load_from_excel": {
            "post": {
                "description": "upload excel table containing products info",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "excel"
                ],
                "summary": "LoadFromExcel",
                "operationId": "load_from_excel",
                "parameters": [
                    {
                        "type": "file",
                        "description": "excel file",
                        "name": "excel_file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of user",
                        "name": "X-User-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of supplier",
                        "name": "X-Supplier-Id",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/product_history": {
            "post": {
                "description": "get product list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "GetProductsHistory",
                "operationId": "products_history",
                "parameters": [
                    {
                        "type": "string",
                        "description": "limit of contents",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "offset of contents",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of user",
                        "name": "X-User-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of supplier",
                        "name": "X-Supplier-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "The input body.\u003cbr /\u003e UploadID is ID of previously uploaded task.",
                        "name": "InputBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.GetProductsHistoryRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Product"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errinfo.errInf"
                        }
                    }
                }
            }
        },
        "/task_history": {
            "get": {
                "description": "get task list",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "GetSupplierTaskHistory",
                "operationId": "supplier_task_history",
                "parameters": [
                    {
                        "type": "string",
                        "description": "limit of contents",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "offset of contents",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of user",
                        "name": "X-User-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of supplier",
                        "name": "X-Supplier-Id",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Task"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errinfo.errInf"
                        }
                    }
                }
            }
        },
        "/tecdoc_articles": {
            "post": {
                "description": "get tecdoc articles",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tecdoc"
                ],
                "summary": "GetTecDocArticles",
                "operationId": "tecdoc_articles",
                "parameters": [
                    {
                        "description": "The input body.\u003cbr /\u003eBrand is brand name, ArticleNumber is article of product.",
                        "name": "InputBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.GetTecDocArticlesRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Article"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errinfo.errInf"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "errinfo.errInf": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "model.Article": {
            "type": "object",
            "properties": {
                "articleNumber": {
                    "type": "string"
                },
                "brand": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "ean": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "oemnumbers": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "packageHeight": {
                    "$ref": "#/definitions/model.ArticleCriteria"
                },
                "packageLength": {
                    "$ref": "#/definitions/model.ArticleCriteria"
                },
                "packageWidth": {
                    "$ref": "#/definitions/model.ArticleCriteria"
                },
                "panoramicImages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Image"
                    }
                },
                "pictures": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Image"
                    }
                },
                "productGroups": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "relatedVehicles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "replacedByArticles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "weight": {
                    "$ref": "#/definitions/model.ArticleCriteria"
                }
            }
        },
        "model.ArticleCriteria": {
            "type": "object",
            "properties": {
                "criteriaAbbrDescription": {
                    "type": "string"
                },
                "criteriaDescription": {
                    "type": "string"
                },
                "criteriaId": {
                    "type": "integer"
                },
                "criteriaType": {
                    "type": "string"
                },
                "criteriaUnitDescription": {
                    "type": "string"
                },
                "formattedValue": {
                    "type": "string"
                },
                "immediateDisplay": {
                    "type": "boolean"
                },
                "isInterval": {
                    "type": "boolean"
                },
                "isMandatory": {
                    "type": "boolean"
                },
                "rawValue": {
                    "type": "string"
                }
            }
        },
        "model.GetProductsHistoryRequest": {
            "type": "object",
            "properties": {
                "UploadID": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "model.GetTecDocArticlesRequest": {
            "type": "object",
            "properties": {
                "ArticleNumber": {
                    "type": "string",
                    "example": "0451103274XF"
                },
                "Brand": {
                    "type": "string",
                    "example": "BOSCH"
                }
            }
        },
        "model.Image": {
            "type": "object",
            "properties": {
                "fileName": {
                    "type": "string"
                },
                "headerDescription": {
                    "type": "string"
                },
                "headerKey": {
                    "type": "integer"
                },
                "imageURL100": {
                    "type": "string"
                },
                "imageURL1600": {
                    "type": "string"
                },
                "imageURL200": {
                    "type": "string"
                },
                "imageURL3200": {
                    "type": "string"
                },
                "imageURL400": {
                    "type": "string"
                },
                "imageURL50": {
                    "type": "string"
                },
                "imageURL800": {
                    "type": "string"
                },
                "typeDescription": {
                    "type": "string"
                },
                "typeKey": {
                    "type": "integer"
                }
            }
        },
        "model.Product": {
            "type": "object",
            "properties": {
                "article": {
                    "type": "string"
                },
                "brand": {
                    "type": "string"
                },
                "cardNumber": {
                    "type": "integer"
                },
                "category": {
                    "type": "string"
                },
                "errorResponse": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "manufacturerArticle": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "providerArticle": {
                    "type": "string"
                },
                "sku": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "updateDate": {
                    "type": "string"
                },
                "uploadDate": {
                    "type": "string"
                },
                "uploadId": {
                    "type": "integer"
                }
            }
        },
        "model.Task": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "ip": {
                    "type": "string"
                },
                "productsFailed": {
                    "type": "integer"
                },
                "productsProcessed": {
                    "type": "integer"
                },
                "productsTotal": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
                },
                "supplierID": {
                    "type": "integer"
                },
                "updateDate": {
                    "type": "string"
                },
                "uploadDate": {
                    "type": "string"
                },
                "userID": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8002",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Tec-Doc API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

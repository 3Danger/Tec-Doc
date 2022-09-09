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
        "/articles/enrichment": {
            "post": {
                "description": "get task list",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "GetTecDocArticles",
                "operationId": "articles_enrichment",
                "parameters": [
                    {
                        "description": "brand \u0026\u0026 article - about product",
                        "name": "request",
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
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/model.Article"
                                }
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
        "/excel": {
            "get": {
                "description": "download excel table template",
                "tags": [
                    "excel"
                ],
                "summary": "ExcelTemplate",
                "operationId": "excel_template",
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
            },
            "post": {
                "description": "upload excel table containing products info",
                "tags": [
                    "excel"
                ],
                "summary": "LoadFromExcel",
                "operationId": "load_from_excel",
                "parameters": [
                    {
                        "description": "binary excel file",
                        "name": "excel_file",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
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
        "/excel/enrichment": {
            "post": {
                "description": "Enrichment excel file, limit entiies in file = 10000",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "excel"
                ],
                "summary": "ProductsEnrichedExcel",
                "operationId": "enrich_excel",
                "parameters": [
                    {
                        "description": "binary excel file",
                        "name": "excel_file",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
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
                            "$ref": "#/definitions/model.UploadIdRequest"
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
                "articleCriteria": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ArticleCriteria"
                    }
                },
                "articleNumber": {
                    "type": "string"
                },
                "crossNumbers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.CrossNumbers"
                    }
                },
                "genericArticleDescription": {
                    "type": "string"
                },
                "images": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "linkageTargets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.LinkageTargets"
                    }
                },
                "mfrName": {
                    "type": "string"
                },
                "oemNumbers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.OEM"
                    }
                },
                "packageArticleCriteria": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ArticleCriteria"
                    }
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
                "criteriaType": {
                    "type": "string"
                },
                "criteriaUnitDescription": {
                    "type": "string"
                },
                "formattedValue": {
                    "type": "string"
                },
                "rawValue": {
                    "type": "string"
                }
            }
        },
        "model.CrossNumbers": {
            "type": "object",
            "properties": {
                "articleNumber": {
                    "type": "string"
                },
                "mfrName": {
                    "type": "string"
                }
            }
        },
        "model.GetTecDocArticlesRequest": {
            "type": "object",
            "properties": {
                "articleNumber": {
                    "type": "string"
                },
                "brand": {
                    "type": "string"
                }
            }
        },
        "model.LinkageTargets": {
            "type": "object",
            "properties": {
                "beginYearMonth": {
                    "type": "string"
                },
                "endYearMonth": {
                    "type": "string"
                },
                "linkageTargetId": {
                    "type": "integer"
                },
                "mfrName": {
                    "type": "string"
                },
                "vehicleModelSeriesName": {
                    "type": "string"
                }
            }
        },
        "model.OEM": {
            "type": "object",
            "properties": {
                "articleNumber": {
                    "type": "string"
                },
                "mfrName": {
                    "type": "string"
                }
            }
        },
        "model.Product": {
            "type": "object",
            "properties": {
                "article": {
                    "type": "string"
                },
                "articleSupplier": {
                    "type": "string"
                },
                "barcode": {
                    "type": "string"
                },
                "brand": {
                    "type": "string"
                },
                "errorResponse": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "price": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
                },
                "subject": {
                    "type": "string"
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
        },
        "model.UploadIdRequest": {
            "type": "object",
            "properties": {
                "uploadID": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8002",
	BasePath:         "/api/v1/",
	Schemes:          []string{"http"},
	Title:            "Tec-Doc API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

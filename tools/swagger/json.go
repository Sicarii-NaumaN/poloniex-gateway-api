package swagger

const (
	// SwaggerJSON raw
	SwaggerJSON = `{
		   "swagger":"2.0",
		   "info":{
			  "title": "%s",
			  "version":"version not set"
		   },
		   "schemes":[
			  "http"
		   ],
           "host": "localhost:3000",
		   "paths":{
			 %s
		   },
		   "definitions": 
			 %s
		}`

	// HandlerRaw raw
	HandlerRaw = `
		"%s":{
			 "%s":{
				"summary": "%s",
				%s
				"responses":{
				   "200":{
					  "description":"A successful response.",
						 %s
				   },
				   "default":{
					  "description":"An unexpected error response.",
						 %s
				   }
				},
				%s
				"tags":[
				   "%s"
				]
			 }
      	}`

	// HandlerBrakesRaw raw
	HandlerBrakesRaw = `
		"%s":{
			%s
      	}`

	HandlerRaw1 = `
		 "%s":{
				"summary": "%s",
				%s
				"responses":{
				   "200":{
					  "description":"A successful response.",
						 %s
				   },
				   "default":{
					  "description":"An unexpected error response.",
						 %s
				   }
				},
				%s
				"tags":[
				   "%s"
				]
			 }
	`

	// URL parameters raws
	parametersRaw = `
		"parameters": [
          %s
        ],`
	parameter = `
		{
			"name": "%s",
			"in": "%s",
			"required": %t,
			"description": "%s",
			"type": "%s"
			%s
		}`
	parameterArray = `
		,"items": {
              "type": "%s"
            },
            "collectionFormat": "multi"
		`
	parameterBody = `
		{
            "name": "body",
            "in": "body",
            "required": true,
            	%s
		}`
	parameterFormData = `
		{
			"name": "%s",
			"in": "formData",
			"required": %t,
			"description": "%s",
			"type": "%s"
			%s
		}`

	// EmptyObject case no response
	EmptyObject  = `"schema": { "type":"object" }`
	bodySchema   = `"schema": { "$ref": "#/definitions/%s" }`
	fileResponse = `"schema": { "type": "file" }`

	responseProduce = `
		"produces": [
          %s
        ]
	`
	requestConsumes = `
		"consumes": [
          "%s"
        ]
	`
)

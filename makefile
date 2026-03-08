.PHONY: up down swagger-ui-yaml swagger-ui-json

up:
	docker-compose up --build -d

down:
	docker-compose down

# Swagger UI from openapi.yaml at http://localhost:8081
swagger-ui-yaml:
	docker run --rm -p 8081:8080 \
		-e SWAGGER_JSON=/openapi.yaml \
		-v $(PWD)/openapi.yaml:/openapi.yaml \
		swaggerapi/swagger-ui

# Swagger UI from hospital_swagger.json at http://localhost:8082
swagger-ui-json:
	docker run --rm -p 8082:8080 \
		-e SWAGGER_JSON=/hospital_swagger.json \
		-v $(PWD)/hospital_swagger.json:/hospital_swagger.json \
		swaggerapi/swagger-ui
up:
	docker-compose up --build -d
down:
	docker-compose down
#http://localhost:8081
swagger-ui:
	docker run -p 8081:8080 \
	-v $(pwd)/openapi.yaml:/openapi.yaml \
	swaggerapi/swagger-ui
create-migrate::
	migrate create -seq -ext .sql -dir ./migrations create_event_table

add-migrate::
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5433/demo?sslmode=disable" up
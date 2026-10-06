create-migrate::
	migrate create -ext sql -dir migrations/platform -seq init

add-migrate::
	migrate -path migrations/platform -database "postgres://cloud-native-shop:cloud-native-shop-password@localhost:5433/cloud-native-shop-db?sslmode=disable" up
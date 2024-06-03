migrate:
	goose -dir .\postgres\schema postgres postgresql://sonnvt:sonnvt@localhost:5432/demo?sslmode=disable up

oath:
	cd oath-keeper && .\down_oath_keeper.bat
	cd oath-keeper && .\run_oath_keeper.bat
	cd ..

run-all:
	docker-compose up -d
	cd oath-keeper && .\run_oath_keeper.bat
	cd ..

down-all:
	docker-compose down
	cd oath-keeper && .\down_oath_keeper.bat
	cd ..

run-pg:
	docker exec -it datn-postgres psql -U sonnvt -d demo

run-scylla:
	docker exec -it datn-scylladb cqlsh
	use authsvc;


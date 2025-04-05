PG_CONTAINER := local_db
PG_IMAGE := postgres:13.5
PG_PASS := secret
PG_USER := daylear
PG_DB := daylear
PG_PORT := 5432

S3_CONTAINER := local_s3
S3_IMAGE := adobe/s3mock
S3_BUCKET := daylear

.PHONY: run_daylear
run_daylear: run_db | run_s3 
	@go run ./...

## Run a postgres database locally with docker (PHONY)
run_db:
	@docker rm -f $(PG_CONTAINER)
	@docker run --name $(PG_CONTAINER) -e POSTGRES_PASSWORD=$(PG_PASS) -e POSTGRES_USER=$(PG_USER) -e POSTGRES_DB=$(PG_DB) -d -p $(PG_PORT):5432 $(PG_IMAGE) -c log_statement=all -c logging_collector=on -c log_directory=/var/log/postgresql -c log_filename=postgressql.log

run_s3:
	@docker rm -f $(S3_CONTAINER)
	@docker run -d -p 9090:9090 -e initialBuckets=daylear -e debug=true -t --name $(S3_CONTAINER) $(S3_IMAGE) 
help:
	@echo "Usage: make <command>\n"
	@echo "Commands:"
	@printf "%3sstart-test-db: Starts the test SQL Docker container\n"
	@printf "%3sstop-test-db: Stops the test SQL Docker container\n"
	@exit 0

start-test-db:
	@bash ./tools/start_testdb.sh

stop-test-db:
	@bash ./tools/stop_testdb.sh
help:
	@echo "Usage: make <command>\n"
	@echo "Commands:"
	@printf "%3sstart-testdb: Starts the test SQL Docker container\n"
	@printf "%3sstop-testdb: Stops the test SQL Docker container\n"
	@exit 0

start-testdb:
	@bash ./tools/start_testdb.sh

stop-testdb:
	@bash ./tools/stop_testdb.sh
# Define your Go application name
APP_NAME := .

# Define the number of replicas
NUM_REPLICAS := 7

# Define targets for building and running the application
.PHONY: build run

build:
	cd node/cmd && go run $(APP_NAME)

#run:
#	zip -r archive.zip .
#	expect -c 'spawn zipcloak archive.zip; expect "Enter password:"; send "123\r"; interact'
#
#	cd node/cmd && go run $(APP_NAME)
#
#	@for i in $$(seq $(NUM_REPLICAS)); do \
#		echo "Running replica $$i"; \
#		cd node/cmd && go run $(APP_NAME) & \
#	done; \
#	echo "All replicas are running in the background"

run:
	@for i in $$(seq $(NUM_REPLICAS)); do \
		echo "Running replica $$i"; \
		cd node/cmd && go run $(APP_NAME) & \
	done; \
	echo "All replicas are running in the background"

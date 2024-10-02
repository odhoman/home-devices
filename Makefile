LAMBDA_DIR = lambdas
GO_BUILD = GOOS=linux GOARCH=amd64 go build -o 
GO_TEST = go test -v

test_build_and_deploy_stack_all_lambdas:
	@echo "Testing, Building and Deploying stack for lambdas: $(LAMBDA)"
	@$(MAKE) test_and_build_all
	npm test
	cdk deploy

test_build_and_deploy_stack_for_single_lambda:
	@echo "Testing, Building and Deploying stack for lambdas: $(LAMBDA)"
	@$(MAKE) test_and_build_single_lambda LAMBDA=$(LAMBDA)
	npm test
	cdk deploy

test_and_build_all:
	@echo "Testing and Building all lambdas..."
	@$(MAKE) test_all || { echo "Tests failed. Build aborted."; exit 1; }
	@$(MAKE) build_single_lambda LAMBDA=createDevice
	@$(MAKE) build_single_lambda LAMBDA=deleteDevice
	@$(MAKE) build_single_lambda LAMBDA=updateDevice
	@$(MAKE) build_single_lambda LAMBDA=getDevice
	@$(MAKE) build_single_lambda LAMBDA=homeDeviceListener
	@$(MAKE) build_single_lambda LAMBDA=kinesisListener
	@echo "Testing and Building all lambdas: Completed."
	
build_all:
	@echo "Testing and Building all lambdas..."
	@$(MAKE) build_single_lambda LAMBDA=createDevice
	@$(MAKE) build_single_lambda LAMBDA=deleteDevice
	@$(MAKE) build_single_lambda LAMBDA=updateDevice
	@$(MAKE) build_single_lambda LAMBDA=getDevice
	@$(MAKE) build_single_lambda LAMBDA=kinesisListener
	@echo "Testing and Building all lambdas: Completed."	

test_and_build_createDevice:
	@echo "Testing all and Building createDevice..."
	@$(MAKE) test_and_build_single_lambda LAMBDA=createDevice
	@echo "Build of createDevice completed."

test_and_build_deleteDevice:
	@echo "Testing all and Building deleteDevice..."
	@$(MAKE) test_and_build_single_lambda LAMBDA=deleteDevice
	@echo "Build of deleteDevice completed."

test_and_build_updateDevice:
	@echo "Testing all and Building updateDevice..."
	@$(MAKE) test_and_build_single_lambda LAMBDA=updateDevice
	@echo "Build of updateDevice completed."

test_and_build_getDevice:
	@echo "Testing all and Building getDevice..."
	@$(MAKE) test_and_build_single_lambda LAMBDA=getDevice
	@echo "Build of getDevice completed."

test_and_build_homeDeviceListener: 
	@echo "Testing all and Building homeDeviceListener..."
	@$(MAKE) test_and_build_single_lambda LAMBDA=homeDeviceListener
	@echo "Build of homeDeviceListener completed."

test_and_build_single_lambda:
	@$(MAKE) test_all || { echo "Tests failed. Build aborted."; exit 1; }
	@$(MAKE) build_single_lambda LAMBDA=$(LAMBDA)

build_single_lambda:
	@echo "Building lambda: $(LAMBDA)"
	@echo "cd $(LAMBDA_DIR)"
	cd "$(LAMBDA_DIR)"
	@echo "running cd $(LAMBDA_DIR)"
	cd $(LAMBDA_DIR) && $(GO_BUILD) ./cmd/$(LAMBDA)/bootstrap ./cmd/$(LAMBDA)/$(LAMBDA).go
	@echo "Build of $(LAMBDA) completed."

test_all:
	@echo "Running all tests in $(LAMBDA_DIR)..."
	@$(MAKE) run_tests_in_dir dir=$(LAMBDA_DIR) || { echo "Tests failed. Aborting."; exit 1; }

run_tests_in_dir:
	@for sub_dir in $(dir)/*; do \
		if [ -d $$sub_dir ]; then \
			$(MAKE) run_tests_in_dir dir=$$sub_dir || exit 1; \
		fi \
	done
	@echo "Checking for test files in $(dir)..."
	@test_files=`find $(dir) -maxdepth 1 -type f -name "*_test.go"`; \
	if [ -n "$$test_files" ]; then \
		echo "Running tests in $(dir)..."; \
		cd $(dir) && $(GO_TEST); \
		if [ $$? -ne 0 ]; then \
			echo "Tests failed in $(dir). Aborting."; \
			exit 1; \
		fi; \
	fi

.PHONY: test_build_and_deploy_stack_all_lambdas \
        test_build_and_deploy_stack_for_single_lambda \
        test_and_build_all \
        test_and_build_createDevice \
        test_and_build_deleteDevice \
        test_and_build_updateDevice \
        test_and_build_getDevice \
        test_and_build_homeDeviceListener \
        test_and_build_single_lambda \
        build_single_lambda \
        test_all \
        run_tests_in_dir \
        createDevice \
        deleteDevice \
        updateDevice \
        getDevice \
        homeDeviceListener


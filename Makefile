.PHONY: help
help:
	@echo 'make <target>'
	@echo ''
	@echo 'Targets:'
	@echo ''
	@echo '    help          Show this help'
	@echo '    pre-commit    Run pre-commit checks'
	@echo ''
	@echo '    gcp_clean     Clean up state files'
	@echo '    gcp_configure Configure the deployment'
	@echo '    gcp_deploy    Deploy configured resources'
	@echo '    gcp_init      Initialize modules, providers'
	@echo '    gcp_install   Install Terraform'
	@echo '    gcp_lint      Run linters'
	@echo '    gcp_plan      Show deployment plan'
	@echo '    gcp_test      Run tests'
	@echo ''

.PHONY: pre-commit
pre-commit:
	@pre-commit run -a

.PHONY: gcp_clean
gcp_clean:
	@cd gcp/project && rm -rf .terraform *.tfstate* .terraform.lock.hcl
	@cd gcp/project/test && rm -f go.mod go.sum

.PHONY: gcp_configure
gcp_configure:
	@cd gcp && ./scripts/configure.sh -e example -m US -o kpeder -p us-east1 -s us-central1 -t devops

.PHONY: gcp_deploy
gcp_deploy: gcp_configure gcp_init
	@cd gcp/project/test && go test -v

.PHONY: gcp_init
gcp_init: gcp_configure
	@cd gcp/project && terraform init
	@cd gcp/project/test && go mod init project_test.go; go mod tidy

.PHONY: gcp_install
gcp_install:
	@chmod +x ./scripts/install_terraform.sh
	@sudo ./scripts/install_terraform.sh -v ./gcp/versions.yaml

.PHONY: gcp_lint
gcp_lint: gcp_configure gcp_init
	@cd gcp/project/test && golangci-lint run --print-linter-name --verbose project_test.go

.PHONY: gcp_plan
gcp_plan: gcp_configure gcp_init
	@cd gcp/project && terraform plan

.PHONY: gcp_test
gcp_test: gcp_configure gcp_init
	@cd gcp/project/test && go test -v -destroy

## Terraform Deployment Example
Example Terraform deployment that includes Go test routines using Terratest.

### Decision Records
This repository uses architecture decision records to record design decisions about important elements of the solution.

The ADR index is available [here](./docs/decisions/index.md).

### Requirements
Tested on Go version 1.21 on Ubuntu Linux.

Uses installed packages:
```
gcloud
golangci-lint
make
pre-commit
terraform
terraform-docs
terragrunt
```

### Configuration
1. Install the packages listed above.
1. Make a copy of the gcp/gcp.yaml file, named local.gcp.yaml, and fill in the fields with configuration values for the target platform.
1. Use gcloud to log into the platform. Terraform uses application default credentials:
    ```
    $ gcloud auth application-default login
    ```

### Deployment
Automated installation configuration, and deployment steps are managed using Makefile targets. Use ```make help``` for a list of configured targets:
```
$ make help 
make <target>

Targets:

    help          Show this help
    pre-commit    Run pre-commit checks

    gcp_clean     Clean up state files
    gcp_configure Configure the deployment
    gcp_deploy    Deploy configured resources
    gcp_init      Initialize modules, providers
    gcp_install   Install Terraform
    gcp_lint      Run linters
    gcp_plan      Show deployment plan
    gcp_test      Run tests
```

Note that additional targets can be added in order to configure multiple environments, for example to create development and production environments.

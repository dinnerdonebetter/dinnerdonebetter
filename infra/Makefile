TERRAFORM_CMD := terraform

.PHONY: terraformat
terraformat:
	@(cd deploy/environments/dev/terraform && $(TERRAFORM_CMD) fmt)

.PHONY: lint_terraform
lint_terraform: terraformat
	@(cd deploy/environments/dev/terraform && $(TERRAFORM_CMD) init -upgrade && $(TERRAFORM_CMD) validate && $(TERRAFORM_CMD) fmt && $(TERRAFORM_CMD) fmt -check)

# project

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5.0 |
| <a name="requirement_random"></a> [random](#requirement\_random) | ~> 4.84.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | n/a |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_folder"></a> [folder](#module\_folder) | terraform-google-modules/folders/google | 4.0.0 |
| <a name="module_iam_audit_config"></a> [iam\_audit\_config](#module\_iam\_audit\_config) | terraform-google-modules/iam/google//modules/audit_config | 7.7.0 |
| <a name="module_project_factory"></a> [project\_factory](#module\_project\_factory) | terraform-google-modules/project-factory/google | 14.3.0 |
| <a name="module_project_name"></a> [project\_name](#module\_project\_name) | git@github.com:kpeder/terraform-module.git | 0.1.0 |

## Resources

| Name | Type |
|------|------|
| [google_project.main](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/project) | data source |
| [google_storage_bucket.state](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/storage_bucket) | data source |

## Inputs

No inputs.

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_audit_log_config"></a> [audit\_log\_config](#output\_audit\_log\_config) | Audit configuration on the configured project |
| <a name="output_billing_account"></a> [billing\_account](#output\_billing\_account) | Billing account of the configured project |
| <a name="output_enabled_apis"></a> [enabled\_apis](#output\_enabled\_apis) | APIs enabled on the configured project |
| <a name="output_folder_id"></a> [folder\_id](#output\_folder\_id) | ID of the configured project's parent folder |
| <a name="output_folder_name"></a> [folder\_name](#output\_folder\_name) | Name of the configured project's parent folder |
| <a name="output_project_id"></a> [project\_id](#output\_project\_id) | ID of the configured project |
| <a name="output_project_labels"></a> [project\_labels](#output\_project\_labels) | Labels configured on the project |
| <a name="output_project_name"></a> [project\_name](#output\_project\_name) | Name of the configured project |
| <a name="output_state_bucket_labels"></a> [state\_bucket\_labels](#output\_state\_bucket\_labels) | Labels configured on the state bucket for the configured project |
| <a name="output_state_bucket_name"></a> [state\_bucket\_name](#output\_state\_bucket\_name) | Name of the Terraform state bucket for the configured project |
| <a name="output_state_bucket_project"></a> [state\_bucket\_project](#output\_state\_bucket\_project) | Parent project of the Terraform state bucket for the configured project |
| <a name="output_state_bucket_versioning"></a> [state\_bucket\_versioning](#output\_state\_bucket\_versioning) | Versioning configuration on the state bucket for the configured project |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->

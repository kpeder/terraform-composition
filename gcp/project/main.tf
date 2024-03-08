locals {
  env    = yamldecode(file("../env.yaml"))
  inputs = yamldecode(file("./inputs.yaml"))
  gcp    = fileexists("../local.gcp.yaml") ? yamldecode(file("../local.gcp.yaml")) : yamldecode(file("../gcp.yaml"))
}

terraform {
  required_version = "~> 1.5.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.84.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 4.84.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5.0"
    }
  }
}

data "google_storage_bucket" "state" {
  name = split("://", module.project_factory.project_bucket_url[0])[1]
}

data "google_project" "main" {
  project_id = module.project_factory.project_id
}

module "folder" {
  source  = "terraform-google-modules/folders/google"
  version = "4.0.0"

  names  = local.inputs.folder.names
  parent = format("%s/%s", local.gcp.parent_type, local.gcp.parent_id)
  prefix = try(local.gcp.prefix, null)
}

module "iam_audit_config" {
  source  = "terraform-google-modules/iam/google//modules/audit_config"
  version = "7.7.0"

  audit_log_config = local.inputs.audit_log_config
  project          = module.project_factory.project_id
}

module "project_factory" {
  source  = "terraform-google-modules/project-factory/google"
  version = "14.3.0"

  activate_apis       = local.inputs.project.activate_apis
  auto_create_network = local.inputs.project.auto_create_network
  billing_account     = coalesce(local.inputs.project.billing_account_override, local.gcp.billing_account)

  # State bucket configuration
  bucket_force_destroy = local.inputs.project.state_bucket.force_destroy
  bucket_labels        = local.env.labels
  bucket_location      = local.env.locations.multiregion
  bucket_name          = format("%s-state", coalesce(local.inputs.project.project_id_override, module.project_name.random_pet))
  bucket_project       = local.gcp.build_project
  bucket_ula           = local.inputs.project.state_bucket.uniform_access
  bucket_versioning    = local.inputs.project.state_bucket.versioning

  default_service_account = local.inputs.project.default_service_account
  folder_id               = module.folder.id
  labels                  = local.env.labels
  name                    = coalesce(local.inputs.project.project_id_override, module.project_name.random_pet)
  org_id                  = local.gcp.organization_id
  random_project_id       = local.inputs.project.random_project_id
}

module "project_name" {
  source = "git@github.com:kpeder/terraform-module.git?ref=0.1.0"

  length = 1
  prefix = try(local.gcp.prefix, null)
}

output "audit_log_config" {
  description = "Audit configuration on the configured project"
  value       = module.iam_audit_config.audit_log_config
}

output "billing_account" {
  description = "Billing account of the configured project"
  value       = data.google_project.main.billing_account
}

output "enabled_apis" {
  description = "APIs enabled on the configured project"
  value       = module.project_factory.enabled_apis
}

output "folder_id" {
  description = "ID of the configured project's parent folder"
  value       = module.folder.id
}

output "folder_name" {
  description = "Name of the configured project's parent folder"
  value       = module.folder.name
}

output "project_id" {
  description = "ID of the configured project"
  value       = module.project_factory.project_id
}

output "project_labels" {
  description = "Labels configured on the project"
  value       = data.google_project.main.labels
}

output "project_name" {
  description = "Name of the configured project"
  value       = module.project_factory.project_name
}

output "state_bucket_labels" {
  description = "Labels configured on the state bucket for the configured project"
  value       = data.google_storage_bucket.state.labels
}

output "state_bucket_name" {
  description = "Name of the Terraform state bucket for the configured project"
  value       = data.google_storage_bucket.state.name
}

output "state_bucket_project" {
  description = "Parent project of the Terraform state bucket for the configured project"
  value       = data.google_storage_bucket.state.project
}

output "state_bucket_versioning" {
  description = "Versioning configuration on the state bucket for the configured project"
  value       = data.google_storage_bucket.state.versioning
}

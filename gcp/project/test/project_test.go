package test

import (
	"flag"
	"os"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// Flag to destroy the target environment after tests
var destroy = flag.Bool("destroy", false, "destroy environment after tests")

func TestGCPProject(t *testing.T) {
	// Set execution directory
	terraformOptions := &terraform.Options{
		TerraformDir: "../.",
	}

	// Check for versions file
	if !assert.FileExists(t, terraformOptions.TerraformDir+"/../versions.yaml") {
		t.Fail()
	}

	// Read and store the versions.yaml
	yfile, err := os.ReadFile(terraformOptions.TerraformDir + "/../versions.yaml")
	if err != nil {
		t.Fail()
	}

	versions := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &versions)
	if err != nil {
		t.Fail()
	}

	// Read the version output and verify the configured version
	goversion := runtime.Version()

	if assert.GreaterOrEqual(t, goversion, "go"+versions["golang_runtime_version"].(string)) {
		t.Logf("Go runtime version check PASSED, expected version >= '%s', got '%s'", "go"+versions["golang_runtime_version"].(string), goversion)
	} else {
		t.Errorf("Go runtime version check FAILED, expected version >= '%s', got '%s'", "go"+versions["golang_runtime_version"].(string), goversion)
	}

	// Check for env.yaml file
	if !assert.FileExists(t, terraformOptions.TerraformDir+"/../env.yaml") {
		t.Fail()
	}

	// Read and store the env.yaml
	yfile, err = os.ReadFile(terraformOptions.TerraformDir + "/../env.yaml")
	if err != nil {
		t.Fail()
	}

	env := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &env)
	if err != nil {
		t.Fail()
	}

	// Check for gcp.yaml file or a local override
	if !assert.FileExists(t, terraformOptions.TerraformDir+"/../local.gcp.yaml") {
		if !assert.FileExists(t, terraformOptions.TerraformDir+"/../gcp.yaml") {
			t.Fail()
		}
	}

	// Read and store the gcp.yaml or a local override
	if assert.FileExists(t, terraformOptions.TerraformDir+"/../local.gcp.yaml") {
		yfile, err = os.ReadFile(terraformOptions.TerraformDir + "/../local.gcp.yaml")
		if err != nil {
			t.Fail()
		}
	} else {
		yfile, err = os.ReadFile(terraformOptions.TerraformDir + "/../gcp.yaml")
		if err != nil {
			t.Fail()
		}
	}

	gcp := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &gcp)
	if err != nil {
		t.Fail()
	}

	// Check for inputs.yaml file
	if !assert.FileExists(t, terraformOptions.TerraformDir+"/inputs.yaml") {
		t.Fail()
	}

	// Read and store the inputs.yaml
	yfile, err = os.ReadFile(terraformOptions.TerraformDir + "/inputs.yaml")
	if err != nil {
		t.Fail()
	}

	inputs := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &inputs)
	if err != nil {
		t.Fail()
	}

	// Sanity test
	terraform.Validate(t, terraformOptions)

	// Initialize the deployment
	terraform.Init(t, terraformOptions)

	// Read the version command output
	version := terraform.RunTerraformCommand(t, terraformOptions, terraform.FormatArgs(terraformOptions, "version")...)

	// Verify configured Terraform version
	if assert.Contains(t, version, "Terraform v"+versions["terraform_binary_version"].(string)) {
		t.Logf("Terraform version check PASSED, expected version '~> %s', got \n%s", versions["terraform_binary_version"].(string), version)
	} else {
		t.Errorf("Terraform version check FAILED, expected version '~> %s', got \n%s", versions["terraform_binary_version"].(string), version)
	}

	// Verify configured Google provider version
	if assert.Contains(t, version, "provider registry.terraform.io/hashicorp/google v"+versions["google_provider_version"].(string)) {
		t.Logf("Provider version check PASSED, expected hashicorp/google version '~> %s', got \n%s", versions["google_provider_version"].(string), version)
	} else {
		t.Errorf("Provider version check FAILED, expected hashicorp/google version '~> %s', got \n%s", versions["google_provider_version"].(string), version)
	}

	// Defer Terraform destroy only if flag is set
	if *destroy {
		defer terraform.Destroy(t, terraformOptions)
	}

	// Create resources
	terraform.Apply(t, terraformOptions)

	// Store outputs
	outputs := terraform.OutputAll(t, terraformOptions)

	// Test project id
	if assert.NotNil(t, outputs["project_id"]) {
		t.Logf("Output test PASSED. Expected output to be string, got %s.", outputs["project_id"].(string))
	} else {
		t.Error("Output test FAILED. Expected output to be string, got nil.")
	}
	if assert.Equal(t, strings.Split(outputs["project_id"].(string), "-")[0], gcp["prefix"]) {
		t.Logf("Prefix test PASSED. Expected project name to start with %s, got %s.", gcp["prefix"], strings.Split(outputs["project_id"].(string), "-")[0])
	} else {
		t.Errorf("Prefix test FAILED. Expected project name to start with %s, got %s.", gcp["prefix"], strings.Split(outputs["project_id"].(string), "-")[0])
	}
	if inputs["project"].(map[string]interface{})["random_project_id"].(bool) {
		if (assert.Len(t, strings.Split(outputs["project_id"].(string), "-"), len(strings.Split(outputs["project_name"].(string), "-"))+1)) &&
			(assert.Contains(t, outputs["project_id"].(string), outputs["project_name"].(string))) {
			t.Logf("Suffix test PASSED. Expected random suffix, got %s.", strings.Split(outputs["project_id"].(string), "-")[len(strings.Split(outputs["project_name"].(string), "-"))])
		} else {
			t.Error("Suffix test FAILED. Expected random suffix, got nil.")
		}
	}

	// Test enabled APIs
	api_failed := false
	apis := []string{}
	for _, api := range inputs["project"].(map[string]interface{})["activate_apis"].([]interface{}) {
		if !assert.Contains(t, outputs["enabled_apis"].([]interface{}), api) {
			t.Errorf("APIs test FAILED. Expected API %s to be enabled, not found.", api)
			api_failed = true
		} else {
			apis = append(apis, api.(string))
		}
	}
	if !api_failed {
		t.Logf("APIs test PASSED. Expected project APIs are enabled: %v", apis)
	}

}

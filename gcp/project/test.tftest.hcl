variables {
  apis = yamldecode(file("./inputs.yaml"))
}

run "plan" {

  command = plan

  assert {
    condition     = length(module.project_factory.enabled_apis) == length(var.apis.project.activate_apis)
    error_message = format("Expected %d APIs to be enabled, got %d", length(var.apis.project.activate_apis), length(module.project_factory.enabled_apis))
  }
}

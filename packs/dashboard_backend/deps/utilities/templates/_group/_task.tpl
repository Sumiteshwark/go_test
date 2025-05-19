# see https://developer.hashicorp.com/nomad/docs/job-specification/env
[[- define "task_environment_variables" ]]
[[- $env_vars := var "nomad_task_env_vars" . ]]
env {
    [[- range $var := $env_vars ]]
        [[- $clean_name := $var.key | upper ]]
        [[ $clean_name ]] = [[ $var.value | quote ]]
    [[- end ]]
}
[[- end ]]



# see https://developer.hashicorp.com/nomad/docs/job-specification/service
[[ define "task_service" -]]
[[ $service := var "nomad_task_service" . ]]
      service {
        name = [[ $service.service_name | quote ]]
        port = [[ $service.service_port_label | quote ]]
        tags = [[ $service.service_tags | toJson ]]
        
        # TODO: add sidecar car for implementing the mTls

        check {
          type     = [[ $service.check_type | quote ]]
          [[- if $service.check_path]]
          path     = [[ $service.check_path | quote ]]
          [[- end]]
          interval = [[ $service.check_interval | quote ]]
          timeout  = [[ $service.check_timeout | quote ]]
        }
      }
[[- end ]]



# see https://developer.hashicorp.com/nomad/docs/job-specification/resources
[[- define "task_resources" ]]
resources {
  [[- $resources := var "nomad_task_resources" . ]]
  cpu        = [[ $resources.cpu ]]
  cores      = [[ $resources.cores | default "null" ]]
  memory     = [[ $resources.memory ]]

  # TODO: add support for memory oversubscription
  # memory_max = [[ $resources.memory_max ]]
}
[[- end ]]
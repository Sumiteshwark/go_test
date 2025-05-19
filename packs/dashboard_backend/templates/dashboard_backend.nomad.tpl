job "[[ var "nomad_job_name" . ]]" { 

  [[ template "util_job_meta" . ]]

  group "[[ var "nomad_group_name" . ]]" {
    
    network {
      mode = "[[ var "nomad_group_network_mode" . ]]"
      [[ template "util_job_group_network_dns" . ]]
      [[ template "util_job_group_network_ports" . ]]
    }
    constraint {
          attribute = "${meta.env}"
          operator = "="
          value     = "[[ var "env" . ]]"     
    }
    [[ template "util_job_group_awaiting_task" . ]]

    task "[[ var "nomad_task_name" . ]]" {
      # see https://developer.hashicorp.com/nomad/docs/drivers
      driver = "[[ var "nomad_task_driver" . ]]"

      config {
        # image = "ghcr.io/covlant/dashboard_backend:[[var "nomad_task_image_tag" .]]"
        image = "ghcr.io/covlant/client_dashboard_backend:[[var "nomad_task_image_tag" .]]"
        
        [[- $ports := list -]]
        [[ range var "nomad_task_ports" . ]]
          [[- $ports = append $ports .name ]]
        [[- end ]]
        ports = [[ toJson $ports ]]
        auth {
          username = "${GITHUB_USERNAME}"
          password = "${GITHUB_TOKEN}"
          server_address = "ghcr.io"
        }
      }

      [[- $env_vars := var "nomad_task_env_vars" . ]]
      env {
        # From nomad_task_env_vars
        [[- range $var := var "nomad_task_env_vars" . ]]
          [[- $clean_name := $var.key | upper ]]
          [[ $clean_name ]] = [[ $var.value | quote ]]
        [[- end ]]

        # Vault-rendered or template values
        DB_PASSWORD         = [[ var "DB_PASSWORD" . | quote ]]
        NOMAD_BASE_URL      = [[ var "NOMAD_BASE_URL" . | quote ]]

        # Directly pulled from environment (i.e., GitHub Actions secrets or Nomad env injection)
        AUTH0_CLIENT_ID     = "${AUTH0_CLIENT_ID}"
        AUTH0_CLIENT_SECRET = "${AUTH0_CLIENT_SECRET}"
      }


      [[ template "vault_kv" . ]]

      [[ template "task_service" .]]

      [[ template "task_resources" .]]

    }
  } 
}

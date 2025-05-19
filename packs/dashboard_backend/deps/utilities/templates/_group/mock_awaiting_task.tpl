[[- define "util_job_group_awaiting_task" -]]
[[- /* Get the list of services to await */ -]]
[[- $services := var "nomad_awaiting_task_services" . ]]

task "await-services" {
  driver = "docker"

  config {
    image   = "busybox:1.28"
    command = "sh"
    args    = [
      "-c", <<-EOF
        services="[[ join " " $services ]]"
        for svc in $services; do
          echo "Waiting for $svc..."
          until nslookup "$svc.service.consul" >/dev/null 2>&1; do
            echo -n "."
            sleep 2
          done
          echo "$svc is available."
        done
      EOF
    ]
    network_mode = "host"
  }

  resources {
    cpu    = 200
    memory = 128
  }

  lifecycle {
    hook    = "prestart"
    sidecar = false
  }
}
[[- end ]]
[[- define "util_job_group_network_ports" ]]
# see https://developer.hashicorp.com/nomad/docs/job-specification/network#port
[[- $ports := var "nomad_task_ports" . ]]

    [[- range $port := $ports ]]
        port [[ $port.name | quote ]] {
            to = [[ $port.to ]]
            [[- if $port.static ]]
            static = [[ $port.static ]]
            [[- end ]]
        }
    [[- end ]]
[[- end ]]


[[- define "util_job_group_network_dns" ]]
# see https://developer.hashicorp.com/nomad/docs/job-specification/network#dns-parameters
[[- $dns := var "nomad_group_network_dns" . ]]
dns {
    servers = [[ $dns.servers | toJson ]]
}    
[[- end ]]
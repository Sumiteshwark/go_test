[[- define "util_job_meta" ]]
region      = "[[ var "nomad_job_region" . | default "global" ]]"
datacenters = [[ var "nomad_job_datacenters" . | default (list "dc1") | toJson ]]
type        = "[[ var "nomad_job_type" . | default "service" ]]"
namespace   = "[[ var "nomad_job_namespace" . | default "dev" ]]"
priority    = [[ var "nomad_job_priority" . | default 50 ]]
[[- end ]]

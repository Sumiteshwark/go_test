[[- define "vault_kv" -]]
[[- /* Get Vault kv secret */ -]]

template {
  destination = "secrets/secret.env"
  env         = true
  change_mode = "restart"
  data        = <<EOH
{{ with secret "secret/data/github" }}
   GITHUB_USERNAME={{ .Data.data.username }}
   GITHUB_TOKEN={{ .Data.data.token }}
{{ end }}
{{ with secret "secret/data/nomad" }}
	NOMAD_TOKEN={{ .Data.data.nomad_token }}
	NOMAD_ADDR={{ .Data.data.nomad_addr }}
	NOMAD_PACK_VERSION={{ .Data.data.nomad_pack_version }}
{{ end }}
{{ with secret "secret/data/auth0" }}
    AUTH0_CLIENT_ID={{ .Data.data.client_id }}
    AUTH0_CLIENT_SECRET={{ .Data.data.client_secret }}
{{ end }}
EOH
}

vault {
 policies = ["nomad-workloads"]
}

[[- end ]]
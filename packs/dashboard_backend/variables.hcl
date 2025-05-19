variable "env" {
  type = string
  description = "Nomad client env"
  default = "dev"
}

variable "nomad_job_name" {
  type        = string
  description = "Name for the Job."
  default = "dashboard-backend"
}

variable "nomad_group_name" {
  type        = string
  description = "Name for the Group."
  default     = "dashboard-backend"
}

variable "nomad_group_network_mode" {
  type        = string
  description = "Network Mode for the Group."
  default     = "host"
}

variable "nomad_group_network_dns" {
  type        = object({
    servers = list(string)
  })
  description = "Network dns for the Group."
  default     = {
    servers = ["172.17.0.1"]
  }
}

variable "nomad_task_ports" {
  type = list(object({
      name   = string
      to     = number
      static = number
    }))
  description = "Network ports for the group"
  default = [
      {
        "name" = "http",
        "to" = 8080,
        "static" = 8080,
      }
    ]
}

variable "nomad_awaiting_task_services" {
  type        = list(string)
  description = "Service on which dashboard backend in waiting for."
  default     = ["postgres"]
}


variable "nomad_task_driver" {
  type        = string
  description = "Driver to use for the Task."
  default     = "docker"
}

variable "nomad_task_name" {
  type        = string
  description = "Name for the Task."
  default     = "dashboard-backend"
}

variable "nomad_task_image_tag" {
  type        = string
  description = "Tag of the docker iamge"
  default     = "latest"
}

variable "nomad_task_env_vars" {
  type = list(object({
    key   = string
    value = string
  }))
  description = "Environment variables for the codeparser"
  default = [
    {
      key   = "DB_HOST"
      value = "postgres.service.consul"
    },
    {
      key   = "DB_USER"
      value = "ne_user"
    },
    {
      key   = "DB_NAME"
      value = "neural_engine_db"
    },
    {
      key   = "DB_PORT"
      value = "5432"
    },
    {
      key   = "AUTH0_DOMAIN"
      value = "covlant.us.auth0.com"
    },
    {
      key   = "AUTH0_AUDIENCE"
      value = "https://covlant.us.auth0.com/api/v2/"
    }
  ]
}

variable "NOMAD_BASE_URL" {
  type        = string
  description = "Nomad url"
  default     = "http://74.50.126.215:4646/"
}

variable "DB_PASSWORD" {
  type        = string
  description = "DB Password"
  default     = "ne_pass"
}

variable "AUTH0_CLIENT_ID" {
  type        = string
  description = "Auth0 client id"
  default     = "client_id"
}

variable "AUTH0_CLIENT_SECRET" {
  type        = string
  description = "Client secret of Auth0"
  default     = "client_secret"
}


variable "nomad_task_service" {
  type = object({
    service_port_label = string
    service_name       = string
    service_tags       = list(string)
    check_enabled      = bool
    check_type         = string
    check_path         = string
    check_interval     = string
    check_timeout      = string
  })
  description = ""
  default = {
    service_port_label = "http",
    service_name       = "dashboard-backend",
    service_tags       = ["dashboard-backend"],
    check_enabled      = true,
    check_type         = "tcp",
    check_path         = "/api/health",
    check_interval     = "30s",
    check_timeout      = "2s",
  }
}

variable "nomad_task_resources" {
  type = object({
    cpu        = number
    memory     = number
  })
  description = "Resource Limits for the Task."
  default = {
    # value in MHz
    cpu = 1000
    # value in MB
    memory = 1024
  }
}
variable "github_password" {
  type        = string
  description = "Github password for accesing the ghcr.io registry."
  default = "password"
}

variable "github_username" {
  type        = string
  description = "Github username for accesing the ghcr.io registry."
  default = "username"
}

variable "os_auth_url"     { type = string }
variable "os_region"       { type = string }
variable "os_username"     { type = string }
variable "os_password"     { type = string }
variable "os_project_name" { type = string }

variable "image_name"  {
  type    = string
  default = "debian"
}

variable "flavor_name" {
  type    = string
  default = "m1.large"
}

variable "external_network_id"   { type = string }
variable "external_network_name" { type = string }

variable "ssh_public_key" {
  type    = string
  default = "/opt/stack/.ssh/id_rsa.pub"
}

variable "ssh_private_key" {
  type    = string
  default = "/opt/stack/.ssh/id_rsa"

}
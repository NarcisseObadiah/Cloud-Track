terraform {
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.52.1"
    }
  }
}

provider "openstack" {
  auth_url    = var.os_auth_url
  region      = var.os_region
  user_name   = var.os_username
  password    = var.os_password
  tenant_name = var.os_project_name
}

# SSH Key
resource "openstack_compute_keypair_v2" "keypair" {
  name       = "terraform-key"
  public_key = file(var.ssh_public_key)
}

# Network & Subnet
resource "openstack_networking_network_v2" "private_net" {
  name = "terraform-net"
}

resource "openstack_networking_subnet_v2" "private_subnet" {
  name            = "terraform-subnet"
  network_id      = openstack_networking_network_v2.private_net.id
  cidr            = "192.168.100.0/24"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8"]
}

# Router
resource "openstack_networking_router_v2" "router" {
  name                = "terraform-router"
  external_network_id = var.external_network_id
}

resource "openstack_networking_router_interface_v2" "router_interface" {
  router_id = openstack_networking_router_v2.router.id
  subnet_id = openstack_networking_subnet_v2.private_subnet.id
}

# Security Group
resource "openstack_compute_secgroup_v2" "secgroup" {
  name        = "terraform-secgroup"
  description = "Allow SSH, ICMP, Kubernetes"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 6443
    to_port     = 6443
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 10250
    to_port     = 10250
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

# Master Node
resource "openstack_compute_instance_v2" "master" {
  name            = "k8s-master"
  image_name      = var.image_name
  flavor_name     = var.flavor_name
  key_pair        = openstack_compute_keypair_v2.keypair.name
  security_groups = [openstack_compute_secgroup_v2.secgroup.name]

  network {
    uuid = openstack_networking_network_v2.private_net.id
  }

  depends_on = [openstack_networking_router_interface_v2.router_interface]
}

resource "openstack_networking_floatingip_v2" "master_fip" {
  pool = var.external_network_name
}

resource "openstack_compute_floatingip_associate_v2" "master_fip_assoc" {
  floating_ip = openstack_networking_floatingip_v2.master_fip.address
  instance_id = openstack_compute_instance_v2.master.id
}

# Worker Nodes
resource "openstack_compute_instance_v2" "workers" {
  count           = 2
  name            = "k8s-worker-${count.index + 1}"
  image_name      = var.image_name
  flavor_name     = var.flavor_name
  key_pair        = openstack_compute_keypair_v2.keypair.name
  security_groups = [openstack_compute_secgroup_v2.secgroup.name]

  network {
    uuid = openstack_networking_network_v2.private_net.id
  }

  depends_on = [openstack_compute_instance_v2.master]
}

resource "openstack_networking_floatingip_v2" "worker_fips" {
  count = 2
  pool  = var.external_network_name
}
                                              
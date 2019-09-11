// a source represents a reusable setting for a system boot/start.
source "virtualbox-iso" "vb-ubuntu-12.04" {
  # I think TF concept of data is very good,
  iso = data.iso.ubuntu1604_server_amd64

  boot_wait      = "10s"
  http_directory = "xxx"
  boot_command   = ["..."]

  shutdown_command = "echo 'vagrant' | sudo -S shutdown -P now"
}

data "iso" "ubuntu1604_server_amd64" {
  url           = "http://releases.ubuntu.com/12.04/ubuntu-12.04.5-server-amd64.iso"
  checksum      = "769474248a3897f4865817446f9a4a53"
  checksum_type = "md5"
}

# If connection is it's own abstraction it could resolve one of Packers main issues.
connection "ssh" "vagrant" {
  password     = "vagrant"
  username     = "vagrant"
  wait_timeout = "10000s"
}

source "amazon-ebs" "aws-ubuntu-16.04" {
  instance_type = "t2.micro"
  encrypt_boot  = true
  region        = "eu-west-3"
  # I think TF concept of data is very good,
  # so instead of the below source_ami_filter:
  ami = data.amazon_ami.ubuntu1604

  #ssh_username = "ubuntu"
}

data "amazon_ami" "ubuntu1604" {
  filters {
    virtualization-type = "hvm"
    name                = "ubuntu/images/*ubuntu-xenial-{16.04}-amd64-server-*"
    root-device-type    = "ebs"
  }
  owners      = ["099720109477"]
  most_recent = true
}

import_sources {
  from = "packer.io/some/library/defining/sources"
}

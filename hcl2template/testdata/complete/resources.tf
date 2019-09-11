// a source represents a reusable setting for a system boot/start.
source "virtualbox-iso" "vb-ubuntu-12.04" {
    iso_url = "http://releases.ubuntu.com/12.04/ubuntu-12.04.5-server-amd64.iso"
    iso_checksum = "769474248a3897f4865817446f9a4a53"
    iso_checksum_type = "md5"

    ssh_password = "vagrant"
    ssh_username = "vagrant"
    ssh_wait_timeout = "10000s"

    boot_wait = "10s"
    http_directory = "xxx"
    boot_command = ["..."]

    shutdown_command = "echo 'vagrant' | sudo -S shutdown -P now"
}

source "amazon-ebs" "aws-ubuntu-16.04" {
    instance_type = "t2.micro"
    encrypt_boot = true
    region = "eu-west-3"
    source_ami_filter {
        filters {
            virtualization-type = "hvm"
            name = "ubuntu/images/*ubuntu-xenial-{16.04}-amd64-server-*"
            root-device-type = "ebs"
        }
        owners = [
            "099720109477"
        ]
        most_recent = true
    },

    ssh_username = "ubuntu"
}

import_sources {
    from = "packer.io/some/library/defining/sources"
}

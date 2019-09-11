
source "virtualbox-iso" "vb-ubuntu-1204" {
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


source "virtualbox-iso" "vb-ubuntu-1604" {
    iso_url = "http://releases.ubuntu.com/16.04/ubuntu-16.04.5-server-amd64.iso"
    iso_checksum = "769474248a3897f4865817446f9WRONG"
    iso_checksum_type = "md5"

    ssh_password = "vagrant"
    ssh_username = "vagrant"
    ssh_wait_timeout = "10000s"

    boot_wait = "10s"
    http_directory = "xxx"
    boot_command = ["..."]

    shutdown_command = "echo 'vagrant' | sudo -S shutdown -P now"
}
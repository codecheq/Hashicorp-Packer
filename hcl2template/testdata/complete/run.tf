
// starts resources to provision them.
build {
    from src.amazon-ebs.ubuntu-1604 {
        ami_name = "ami_name_blih_blah"
        // this creates a new resource with settings inherited from the source  
    }

    from src.virtualbox-iso.ubuntu-1204 {
        outout_dir = "path/"
    }

    provision {
        communicator = comm.ssh.vagrant

        shell {
            inline = [
                "echo '{{user `my_secret`}}' :D"
            ]
        }

        shell {
            script = [
                "script-1.sh",
                "script-2.sh",
            ]
            override "vmware-iso" {
                execute_command = "echo 'password' | sudo -S bash {{.Path}}"
            }
        }

        upload "log.go" "/tmp" {
            timeout = "5s"
        }

    }
}

build {
    // build an ami using the ami from the previous build block.
    from "src.amazon.{{user `image_name`}-ubuntu-1.0" {
        ami_name = "fooooobaaaar"
    }

    provision {
        communicator = comm.ssh.vagrant

        shell {
            inline = [
                "echo HOLY GUACAMOLE !"
            ]
        }
    }
}
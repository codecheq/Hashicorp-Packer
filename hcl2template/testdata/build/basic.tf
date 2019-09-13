
// starts resources to provision them.
build {
    output "aws_ami" "{{user `image_name`}}-aws-ubuntu-16.04" {
        from =  "aws-ubuntu-16.04" 
        // this creates a new resource with settings inherited from the source  
    }

    output "aws_ami" "{{user `image_name`}}-vb-ubuntu-12.04" {
        from =  "vb-ubuntu-12.04" 

        override_source_settings {
            // resulting source will get settings from source + this setting :
            ssh_username = "ubuntu"
        }
    }

    output "aws_ami" "{{user `image_name`}}-vmw-ubuntu-16.04" {
        from =  "packer-vmw-ubuntu-16.04"
    }

    provision {
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
    output "aws_ami" "fooooobaaaar" {
        from = "{{user `image_name`}}-aws-ubuntu-16.04"
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
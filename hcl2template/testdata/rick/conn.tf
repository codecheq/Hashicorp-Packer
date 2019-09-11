# If connection is it's own abstraction it could resolve one of Packers main issues.
connection "ssh" "vagrant" {
  password     = "vagrant"
  username     = "vagrant"
  wait_timeout = "10000s"
}

connection "ssh" "password" {
  password     = "s3cr3t"
  username     = "root"
  wait_timeout = "10000s"
}

build {
  # ...

  provisioners {
    connection = connection.ssh.vagrant

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

  # continue provisioning with another connection
  provisioners {
    connection = connection.ssh.root

    shell {
      inline = [
        "echo '{{user `my_secret`}}' :D"
      ]
    }
  }
}


# For AWS the temp key can be something like an TF attribute
source "amazon-ebs" "aws-ubuntu-16.04" {
  # ...
  create_temp_keypair = true
}

connection "ssh" "aws_temp_key" {
  ssh_private_key = source.amazon-ebs.aws-ubuntu-16.04.ssh_keypair
  username        = "ec2-user"
}

# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.box = "ubuntu/bionic64"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # NOTE: This will enable public access to the opened port
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine and only allow access
  # via 127.0.0.1 to disable public access
  config.vm.network "forwarded_port", guest: 3000, host: 3000, host_ip: "127.0.0.1"
  config.vm.network "forwarded_port", guest: 8080, host: 8080, host_ip: "127.0.0.1"
  
  # This one is so that we can use `elm reactor` to rapidly test out ui changes
  # without conflicting with the actual echoserver
  config.vm.network "forwarded_port", guest: 8000, host: 8000, host_ip: "127.0.0.1"
  config.vm.network "forwarded_port", guest: 19999, host: 19999, host_ip: "127.0.0.1"

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "./configs", "/config_files"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  config.vm.provider "virtualbox" do |vb|
    # Display the VirtualBox GUI when booting the machine
    vb.gui = false

    # Customize the amount of memory on the VM:
    vb.memory = "1024"
  end
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.
  $script = <<-SHELL
    # First get all updates and 3rd party artifacts (container images)
    cp /home/vagrant
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    sudo apt-get update
    sudo apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io
    sudo docker pull consul
    sudo docker pull vault
    sudo docker pull fluent/fluent-bit
    sudo docker pull telegraf
    sudo docker pull influxdb
    sudo docker pull netdata/netdata:v1.19.0
    sudo docker pull tecnativa/docker-socket-proxy

    # Install golang
    wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz
    mkdir /home/vagrant/go
    sudo chown vagrant:vagrant /home/vagrant/go
    echo 'PATH=$PATH:/usr/local/go/bin:/home/vagrant/go/bin' >> ~/.bashrc
    echo 'GOPATH=~/go' >> ~/.bashrc
    PATH=$PATH:/usr/local/go/bin:/home/vagrant/go/bin

    # Fetch the sd-notify proxy and install into path
    go get github.com/coreos/sdnotify-proxy && sudo cp ~/go/bin/sdnotify-proxy /usr/local/bin/

    # Install protobuf and grpc
    go get -u google.golang.org/grpc
    sudo apt install unzip
    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.11.2/protoc-3.11.2-linux-x86_64.zip 
    sudo mv include/* /usr/local/include/; sudo mv bin/* /usr/local/bin/
    rmdir include; rmdir bin
    go get -u github.com/golang/protobuf/protoc-gen-go
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
    go get -u github.com/golang/protobuf/protoc-gen-go
    go get -u github.com/caarlos0/env

    # Install elm
    curl -L -o elm.gz https://github.com/elm/compiler/releases/download/0.19.1/binary-for-linux-64-bit.gz
    gunzip elm.gz
    chmod +x elm
    sudo mv elm /usr/local/bin/

    # Get the actual project code and build the container
    go get -u github.com/spf13/cobra/cobra
    go get -u github.com/brnsampson/echopilot
    mkdir -p /etc/echopilot
    ln -s /home/vagrant/go/src/github.com/brnsampson/echopilot /home/vagrant/echopilot

    # Create self-signed certs for this development server and add them to trusted certs
    openssl req -x509 -newkey rsa:4096 -nodes -subj "/C=US/ST=California/L=Who knows/O=Shady Interprises/OU=self signers/CN=www.example.com" -sha256 -keyout /etc/echopilot/key.pem -out /etc/echopilot/cert.pem -days 365
    cp /etc/echopilot/cert.pem /usr/local/share/ca-certificates/
    update-ca-certificates

    # Make convenience link and set up service/config files for all side-cars
    sudo docker build -t echopilot /home/vagrant/go/src/github.com/brnsampson/echopilot/
    sudo cp echopilot/systemd/echopilot.service /etc/systemd/system/
    sudo cp echopilot/systemd/consul.service /etc/systemd/system/
    sudo cp echopilot/systemd/telegraf.service /etc/systemd/system/
    sudo cp echopilot/systemd/fluent-bit.service /etc/systemd/system/
    sudo cp echopilot/systemd/netdata.service /etc/systemd/system/
    sudo cp echopilot/systemd/docker-proxy.service /etc/systemd/system/
    sudo cp -r echopilot/etc/consul/* /etc/consul/
    sudo mkdir -p /etc/telegraf/
    sudo cp echopilot/etc/telegraf.conf /etc/telegraf/
    sudo mkdir -p /etc/fluent-bit
    sudo cp echopilot/etc/fluent-bit.conf /etc/fluent-bit/

    # Any development environment setup
    git clone https://github.com/fatih/vim-go.git ~/.vim/pack/plugins/start/vim-go

    # Now we will set up npm (currently not used)
    # wget -qO- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.0/install.sh | bash
    # nvm install node
  SHELL

  config.vm.provision "shell",
    inline: $script,
    privileged: false
end

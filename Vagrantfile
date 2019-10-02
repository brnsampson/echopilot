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
  # config.vm.network "forwarded_port", guest: 80, host: 8080, host_ip: "127.0.0.1"

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
  config.vm.provision "shell", inline: <<-SHELL
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    apt-get update
    apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
    apt-get install -y docker-ce docker-ce-cli containerd.io
    docker pull consul
    docker pull vault
    docker pull fluent/fluent-bit
    docker pull telegraf
    docker pull influxdb
    wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/vagrant/.bashrc
    go get -u go.uber.org/zap
    go get github.com/coreos/sdnotify-proxy && sudo cp ~/go/bin/sdnotify-proxy /usr/local/bin/
    go get github.com/brnsampson/echopilot.git
    PATH=$PATH:/usr/local/go/bin:/home/vagrant/go/bin
    ln -s /home/vagrant/go/src/github.com/brnsampson/echopilot /home/vagrant/echopilot
    docker build -t echopilot /home/vagrant/go/src/github.com/brnsampson/echopilot/
    cp echopilot/systemd/echopilot.service /etc/systemd/system/
    cp echopilot/systemd/telegraf.service /etc/systemd/system/
    cp echopilot/systemd/fluent-bit.service /etc/systemd/system/
    mkdir -p /etc/telegraf/
    cp echopilot/etc/telegraf.conf /etc/telegraf/
    mkdir -p /etc/fluent-bit
    cp echopilot/etc/fluent-bit.conf /etc/fluent-bit/
    git clone https://github.com/fatih/vim-go.git ~/.vim/pack/plugins/start/vim-go
  SHELL
end

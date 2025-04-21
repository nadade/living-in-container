# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu-server24/ubuntu_24.04.02"
  config.vm.box_version = "1.0.1"
  config.vm.provision "docker", images: ["alpine"]
  config.vm.provision "shell", inline: <<-SHELL
    apt-get update && apt-get install golang-go curl -y
    curl https://cdimage.ubuntu.com/ubuntu-base/noble/daily/current/noble-base-amd64.tar.gz -o /home/vagrant/ubuntufs.tar.gz
    mkdir /home/vagrant/ubuntufs
    tar xvfz /home/vagrant/ubuntufs.tar.gz -C /home/vagrant/ubuntufs
    chown -R vagrant:vagrant /home/vagrant/ubuntufs/
  SHELL
end

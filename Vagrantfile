# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
    config.vm.box = "ubuntu-trusty"
    config.vm.box_url = "https://cloud-images.ubuntu.com/vagrant/trusty/current/trusty-server-cloudimg-amd64-vagrant-disk1.box"
    config.vm.network :private_network, ip: "192.168.69.100"

    config.vm.provider "virtualbox" do |v|
        v.name = "elwyn"
        v.customize ["modifyvm", :id, "--memory", 2048]
    end

    config.vm.provision "shell",
        inline: $install_prereqs

    config.vm.provision "shell",
        inline: $install_go

    config.vm.provision "shell",
        inline: $build_elwyn
end

$install_prereqs = <<SCRIPT
apt-get install -y \
    nginx mongodb make git
SCRIPT

$install_go = <<SCRIPT
wget -q -nc https://go.googlecode.com/files/go1.2.1.linux-amd64.tar.gz
tar -xzf go1.2.1.linux-amd64.tar.gz
cp go/bin/* /usr/local/bin/
cp -r go /usr/lib
SCRIPT


$build_elwyn = <<SCRIPT
export GOROOT=/usr/lib/go
export GOPATH=/go

ln -sf /vagrant /srv/elwyn
cd /srv/elwyn

chown vagrant:vagrant -R .

make clean
make

cp -f deploy/upstart/local.conf /etc/init/elwyn.conf
cp -f deploy/nginx/elwyn.conf /etc/nginx/sites-enabled/elwyn
rm -f /etc/nginx/sites-enabled/default
service elwyn restart
service nginx restart
SCRIPT

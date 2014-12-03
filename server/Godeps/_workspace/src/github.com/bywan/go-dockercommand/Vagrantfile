require 'vagrant-openstack-provider'

Vagrant.configure('2') do |config|

  config.vm.box = 'ubuntu/trusty64'

  config.vm.provider :openstack do |os, override|
    override.ssh.username    = ENV['OS_SSH_USERNAME']
    os.openstack_auth_url    = ENV['OS_AUTH_URL']
    os.tenant_name           = ENV['OS_TENANT_NAME']
    os.username              = ENV['OS_USERNAME']
    os.password              = ENV['OS_PASSWORD']
    os.floating_ip_pool      = ENV['OS_FLOATING_IP_POOL']
    os.flavor                = ENV['OS_FLAVOR']
    os.image                 = ENV['OS_IMAGE']
  end

  config.vm.provision 'shell', inline: 'curl -sSL https://get.docker.com/ubuntu/ | sudo sh'
  config.vm.provision 'shell', inline: 'apt-get update -y && apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git mercurial bzr'
  config.vm.provision 'shell', inline: 'mkdir /goroot && curl https://storage.googleapis.com/golang/go1.3.1.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1'
  config.vm.provision 'shell', inline: 'mkdir /gopath'
  # export GOROOT=/goroot
  # export GOPATH=/gopath
  # export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
end

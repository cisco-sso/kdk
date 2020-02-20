# -*- mode: ruby -*-
# vi: set ft=ruby :

$memory = 4096  # In megabytes
$cpus   = 2

# Keybase
keybase_root = nil
if File.directory?("/keybase/team")
  keybase_root = "/keybase"
elsif File.directory?("/Volumes/Keybase/team")
  keybase_root = "/Volumes/Keybase"
elsif File.directory?("k:/team")
  keybase_root = "k:"
else
  puts "WARNING: Optional keybase.io VirtualFS is not present"
  puts "  On Linux and OSX, KeybaseFS would be mounted under /keybase"
  puts "  On Windows, KeybaseFS would be mounted under k:"
  puts "This Vagrant automation as well as virtual machine configuration"
  puts "  depend on the secrets in keybase.io"
end

Vagrant.configure("2") do |config|

  # References for `dcwangmit01/kdk`
  #   - https://github.com/cisco-sso/kdk/blob/master/files/provision.sh
  #   - https://github.com/cisco-sso/kdk/tree/master/packer
  config.vm.box = "dcwangmit01/kdk"

  # If needed, enable bridged networking.
  #   For example:
  #     If you want your machine to have its own IP on the network.
  #     Or, if you are using sshuttle which seems to require it.
  #   The vagrant base image has ssh password logins disabled.
  # config.vm.network "public_network",
  #   bridge: [  # Vagrant falls back to first match.
  #     "en8: Belkin USB-C LAN",
  #     "en9: USB 10/100/1000 LAN",
  #     "en0: Wi-Fi (AirPort)"
  #   ]

  config.vm.provider "virtualbox" do |vb|
    vb.memory = $memory
    vb.cpus = $cpus
  end

  config.vm.provider "hyperv" do |hv|
    hv.memory = $memory
    hv.cpus = $cpus
  end

  # Mount user-defined directories if they exist
  personal_dirs = {
    "~/Dev" => "/home/vagrant/Dev",
    "~/x" => "/home/vagrant/x",
    "~/git" => "/home/vagrant/git",
    "~/code/cisco" => "/home/vagrant/git",
    "~/Documents/GitHub" => "/home/vagrant/github",
  }
  personal_dirs.each { |host_path, guest_path|
    if File.directory?(File.expand_path(host_path))
      # parent dirs to be auto-created by synced_folder mount
      config.vm.synced_folder File.expand_path(host_path), guest_path
      puts "Enabling optional host directory mount"
      puts "  host_path: " + host_path
      puts "  guest_path:  " + guest_path
    end
  }

  ## Set up SSH forwarding
  config.ssh.forward_agent = true

  ## Place to store secrets.
  ## ATTENTION: Requires Keybase client activation/sign-in on host OS.
  require 'pathname'
  keybase_dirs = [ "private", "public", "team" ]
  unless keybase_root.nil?
    keybase_dirs.each { |dir|
      host_path = String(Pathname.new(keybase_root + "/" + dir).realpath)
      guest_path = "/keybase/" + dir

      if File.directory?(host_path)
        # parent dirs to be auto-created by synced_folder mount
        config.vm.synced_folder host_path, guest_path
      else
        puts "WARNING: Failed to mount keybase.io VirtualFS"
        puts "  host_path: " + host_path
        puts "  guest_path: " + guest_path
      end
    }
  end

  ## Fix Line Endings
  ##   For Windows Hyperv, fix carriage-return line endings.
  ##   For Unix-based systems dos2unix does not modify files.
  config.vm.provision "shell",
                        run: "always",
                        inline: <<-SHELL
    sudo apt-get install -y dos2unix
    pushd /vagrant
    dos2unix $(find . -type d \\( -name .git -o -name vendor -o -name pkg -o -name _dist -o -name .vagrant \\) -prune -o -type f -name '*.go' -prune -o -type f -print) 2>&1
  SHELL

  ## Provision the VM upon first boot.
  config.vm.provision "shell",
    path: "./files/vagrant-provision.sh",
    privileged: false

  if Vagrant.has_plugin?("vagrant-cachier")
    # Copy/Paste from: https://github.com/fgrehm/vagrant-cachier#quick-start
    #
    # Configure cached packages to be shared between instances of the same base box.
    # More info on http://fgrehm.viewdocs.io/vagrant-cachier/usage
    config.cache.scope = :box

    # For more information please check http://docs.vagrantup.com/v2/synced-folders/basic_usage.html
  end
end

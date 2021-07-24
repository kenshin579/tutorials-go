#!/usr/bin/env bash

# ruby install
apt-get update
apt-get install -y wget
wget https://raw.githubusercontent.com/redis/redis/4.0/src/redis-trib.rb
chmod 0755 ./redis-trib.rb

apt install gnupg2
apt install curl nodejs dirmngr gnupg2 build-essential libssl-dev git-core zlib1g-dev libreadline-dev libyaml-dev libsqlite3-dev sqlite3 libxml2-dev software-properties-common libxslt1-dev libcurl4-openssl-dev libffi-dev
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
apt update
apt install yarn
git clone https://github.com/rbenv/rbenv.git ~/.rbenv
echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.bashrc
echo 'eval "$(rbenv init -)"' >> ~/.bashrc
git clone https://github.com/rbenv/ruby-build.git ~/.rbenv/plugins/ruby-build
echo 'export PATH="$HOME/.rbenv/plugins/ruby-build/bin:$PATH"' >> ~/.bashrc
rbenv install 2.7.0
gem install redis

# echopilot
A repo for building a containerpilot based dumb echo server in order to test out library insertion and external service integrations


# Vagrant Build
Want to test? There is a bundled Vagrantfile which create a VM that already has the code cloned, the container built, and the service file placed. Convenient!

To vagrant: [sudo] vagrant up; [sudo] vagrant ssh

This VM provides three (initially stopped) service: telegraf, fluent-bit, and echopilot. It also starts an influxdb container for telegraf to write into when you start it.

Some things that would be nice for the future:

- start an instance of elasticsearch/kibana for fluent-bit to write logs into
- add some statsd output or advertise a service for telegraf to scrape in echopilot
- finish off watchdog part of echopilot app/systemd service


# Prerequisites:

- Install docker: https://docs.docker.com/install/
- Install go: https://golang.org/doc/install
- Install sdnotify-proxy: go get github.com/coreos/sdnotify-proxy && sudo cp ~/go/bin/sdnotify-proxy /usr/local/bin/
- Add everything in systemd/ to the /etc/systemd/system/ directory in order to add the new services to systemd
- Add the config files in etc/ to the relevant location in /etc/
- Build with docker build -t echopilot .

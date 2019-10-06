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


# Prerequisites (only needed if you really don't want to use vagrant for some reason):

- Install docker: https://docs.docker.com/install/
- Install go: https://golang.org/doc/install
- Install sdnotify-proxy: go get github.com/coreos/sdnotify-proxy && sudo cp ~/go/bin/sdnotify-proxy /usr/local/bin/
- Add everything in systemd/ to the /etc/systemd/system/ directory in order to add the new services to systemd
- Add the config files in etc/ to the relevant location in /etc/
- Build with docker build -t echopilot .



# Sources (in addition to logs of stack overflow and stuff I'm sure I lost track of):

general go patterns:
- https://github.com/thockin/go-build-template
- https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831

testing:
- https://blog.alexellis.io/golang-writing-unit-tests/
- https://github.com/benbjohnson/testing
- https://medium.com/@benbjohnson/structuring-tests-in-go-46ddee7a25c

signal handling:
- https://github.com/benbjohnson/testing
- https://www.openmymind.net/Golang-Hot-Configuration-Reload/
- https://gravitational.com/blog/golang-ssh-bastion-graceful-restarts/
- https://gist.github.com/peterhellberg/38117e546c217960747aacf689af3dc2

systemd hyjinks:
- https://github.com/coreos/go-systemd

docker stuff:
- https://hub.docker.com/r/phusion/baseimage/


Things to investigate for the future(?):
- https://gist.github.com/rivo/f96ad8710b54a49180a314ec4d68dbfb
- https://grisha.org/blog/2014/06/03/graceful-restart-in-golang/

If I used your resources, no matter how small, thank you! It would have been beyond my patience to figure all of these things out by myself.

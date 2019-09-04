# echopilot
A repo for building a containerpilot based dumb echo server in order to test out library insertion and external service integrations

Want to test? There is a bundled Vagrantfile which create a VM that already has the code cloned, the container built, and the service file placed. Convenient!

Prerequisites:

Install docker: https://docs.docker.com/install/
Install go: https://golang.org/doc/install
Install sdnotify-proxy: go get github.com/coreos/sdnotify-proxy && sudo cp ~/go/bin/sdnotify-proxy /usr/local/bin/

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
- Install elm: curl -L -o elm.gz https://github.com/elm/compiler/releases/download/0.19.1/binary-for-linux-64-bit.gz && gunzip elm.gz && sudo mv elm /usr/local/bin/
- Add everything in systemd/ to the /etc/systemd/system/ directory in order to add the new services to systemd
- Add the config files in etc/ to the relevant location in /etc/
- Build the elm code and populate /dist (see below)
- Build with docker build -t echopilot .

# Testing
run tests with
`go test github.com/brnsampson/echopilot/pkg/echoserver`

# Building
First build the protobuf definition, rest gateway, and swagger doc
 protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. --grpc-gateway_out=logtostderr=true:$GOPATH/src --go_out=plugins=grpc:$GOPATH/src proto/echo/echo.proto; cp proto/echo/echo.swagger.json ./dist/swagger.json

build and install with
`go install github.com/brnsampson/echopilot`

Then run by going to the echopilot directory and running
`echopilot serve`

To instead build the docker container simply go to the echopilot directory and run
`sudo docker build -t echopilot /home/vagrant/go/src/github.com/brnsampson/echopilot/`

Then run via
`sudo docker run --rm -p 3000:8080 echopilot`

# Working with the elm part
This is actually really simple. Elm just compiles down to either a single .js or a single html that we can serve as static concent.
- Working directory: ui/
- Build: elm make src/Main.elm && mv index.html ../dist/
- That's literally all there is to it

- the dockerfile copys in whatever is currently in ./dist/, so after updating the ui code and building you will have to rebuild the docker container.

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

elm:
- https://guide.elm-lang.org/
- https://package.elm-lang.org/
- https://elmprogramming.com/fetching-data-using-get.html

grpc:
- https://github.com/grpc/grpc-go/tree/master/examples
- https://github.com/grpc-ecosystem/grpc-gateway
- https://github.com/grpc-ecosystem/grpc-opentracing
- https://levelup.gitconnected.com/grpc-basics-part-2-rest-and-swagger-53ec2417b3c4
- https://github.com/scottyw/grpc-example
- https://github.com/swagger-api/swagger-ui
- https://medium.com/@ribice/serve-swaggerui-within-your-golang-application-5486748a5ed4

Things to investigate for the future(?):
- https://gist.github.com/rivo/f96ad8710b54a49180a314ec4d68dbfb
- https://grisha.org/blog/2014/06/03/graceful-restart-in-golang/

If I used your resources, no matter how small, thank you! It would have been beyond my patience to figure all of these things out by myself.

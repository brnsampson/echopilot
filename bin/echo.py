#!/usr/bin/env python

import socket
import sys
import logging
import signal

log = logging.getLogger(__name__)
out_hdlr = logging.StreamHandler(sys.stdout)
out_hdlr.setFormatter(logging.Formatter('%(asctime)s %(message)s'))
out_hdlr.setLevel(logging.INFO)
log.addHandler(out_hdlr)
log.setLevel(logging.INFO)


class GracefulKiller:
  kill_now = False
  def __init__(self):
    signal.signal(signal.SIGINT, self.exit_gracefully)
    signal.signal(signal.SIGTERM, self.exit_gracefully)

  def exit_gracefully(self,signum, frame):
    self.kill_now = True

if __name__ == '__main__':
    killer = GracefulKiller()

    # Create a TCP/IP socket
    socket.setdefaulttimeout(1)
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    
    # Bind the socket to the port
    server_address = ('localhost', 10000)
    log.info('starting up on {} port {}'.format(server_address[0], server_address[1]))
    sock.bind(server_address)
    
    # Listen for incoming connections
    sock.listen(1)
    
    log.info('ready for connection')
    while not killer.kill_now:
        # Wait for a connection
        try:
            connection, client_address = sock.accept()
        except socket.timeout:
            pass
        except:
            raise
        else:
            try:
                log.info('connection from {}'.format(client_address))
    
                # Receive the data in small chunks and retransmit it
                while not killer.kill_now:
                    try:
                        data = connection.recv(16)
                        log.info('received "{}"'.format(data))
                        if data:
                            log.info('sending data back to the client')
                            connection.sendall(data)
                        else:
                            log.info('no more data from {}'.format(client_address))
                            break
                    except socket.timeout:
                        pass
            finally:
                # Clean up the connection
                connection.close()
                log.info('ready for connection')
log.info('echo server stopping...')
sys.exit(0)

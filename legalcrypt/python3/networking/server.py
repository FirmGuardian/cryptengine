# echo client
import threading
from socket import *
from ssl import *

def echo_loop(connection, address):
    #while not finished
    finished = False

    while not finished:
        try:
            #send and receive data from the client socket
            data_in=connection.recv(1024)
            message=data_in.decode()
            print('client ' + str(client) + ' ' + str(address) + ' sent ' + message)

            if message=='quit':
                finished= True
                #close the connection
                connection.shutdown(SHUT_RDWR)
                connection.close()

                #close the server socket
                server_socket.shutdown(SHUT_RDWR)
                server_socket.close()
            else:
                data_out=message.encode()
                connection.send(data_out)
        except Exception as e:
            print("[-] Client " + str(address) + ' encountered ' + str(e))
            finished= True

            try:
                #close the connection
                connection.shutdown(SHUT_RDWR)
                connection.close()

                #close the server socket
                server_socket.shutdown(SHUT_RDWR)
                server_socket.close()
            except Exception as e2:
                print("SSL-Related Error Has Occurred (Violation of Protocol)")
                print("Quitting Now...")



#create socket
server_socket=socket(AF_INET, SOCK_STREAM)

#Bind to an unused port on the local machine
server_socket.bind(('localhost',6660))

#listen for connection
server_socket.listen(5)
tls_server = wrap_socket(server_socket, ssl_version=PROTOCOL_TLSv1, cert_reqs=CERT_NONE, server_side=True, keyfile='/Users/alexbujduveanu/Desktop/encryption-app-master/legalcrypt/python3/my.key', certfile='/Users/alexbujduveanu/Desktop/encryption-app-master/legalcrypt/python3/my.crt')
print('server started')

while True:
    #accept connection
    client, address = tls_server.accept()
    client.settimeout(4800)
    threading.Thread(target=echo_loop, args = (client,address)).start()
    print('started thread for connection from ' + str(address))

import SimpleHTTPServer
import SocketServer
import sys
# minimal web server.  serves files relative to the
# current directory.

PORT = int(sys.argv[1])

Handler = SimpleHTTPServer.SimpleHTTPRequestHandler

httpd = SocketServer.TCPServer(("", PORT), Handler)

print("serving at port", PORT)
httpd.serve_forever()
# Play with a simpel example 
1. Run the following command line to open a Chrome with a disable web security setting 
(Because Chrome blocks CORS, which blocks our http request. We are not familar with JS, and do not know how to handle it. So, we need to open browser with disable web security setting)
Fro mac:
open -n -a /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --args --user-data-dir="/tmp/chrome_dev_test" --disable-web-security

For windows (need to adjust the path to Chrome):
"C:\Program Files (x86)\Google\Chrome\Application\chrome.exe" --disable-web-security --disable-gpu --user-data-dir=~/chromeTemp

2. How to start the server: Go to the server folder, and type the following command line:  

go run server.go  

go run server_2.go  

So that we started two servers.  

3. Drag the two html file under the example_html_js folder into the browser. Then you can try the GUI for fun!

# How to configure more server and clients:
1. For server, copy the server.go file to a new file with another name, modify line 398 and 403. For line 398 add another port number into the string list. For line 403, specify myserver.index to the index of that server's port number in the line 398 list. 
2. For client, make a copy of the example_client.html file, modify the line in line 12, 13, and 14. For line 12, give that client a universally unique name, for line 13, append that server's address into the address array. For line 14 specify the index of the server that you want to connect to. 

For our original design, if the server died we would like the client to be able to find another server, but since we are 
not very fimiliar with JavaScript, in our current implementation if our server died, the client won't be able to find another
server by itself. 


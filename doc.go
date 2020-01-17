/*
Sample project:
Implements functionality for creating infinite number of nodes by means of go-libp2p
and its packages. Nodes can be created and run inside separate terminal windows
on same host or on different hosts.
Created network is p2p - independent from physical network topology, with an unstructured overlay
Each node has its own blockchain that contains blocks with data for each
remote peer connection
When a new connection to a remote peer is established - a new block is saved in the blockchain
Used databse is bolt.db - for simplicity

main goroutine expects two arguments:
-sp - source peer - default 0
-debug - defaut false

Source peer number is used for generating private key used for peer identification
Debug is used for chosing to generate same peer ID each time for debugging purposes

Database files are saved inside $HOME/.bolt/<peer-id>/blockchain.db

Running main() from command line as:
./libp2p-node -sp <port-number> -debug <boolean>
starts node creation and displays indications at each step as INFO or ERROR in the terminal

Start ./libp2p-node -sp <port-number> -debug <boolean> in each terminal window
*/
package main

# Edge-Matrix-Core

Edge-Matrix-Core was born from and built to improve the maintainability and modularity of the existing version of [Edge-Matrix-Computing](https://github.com/shawnstartup/edge-matrix-computing)

You can create a P2P network as shown below and connect the webapp to the edge of the network:
```
 +----------------+                 +----------------+                 +-----------------+              +-----------------+
 |                |                 |                |                 |                 |              |                 |
 |                | libp2p stream   |                | libp2p stream   |                 |    HTTP      |                 |
 |   Relay Node   <----------------->   Relay Node   <----------------->    Edge Node    <-------------->   Webapp SERVER |
 |                |   Discovery     |                | p2p Reservation |                 |              |                 |
 |  libp2p host   |                 |  libp2p host   |                 |   libp2p host   |              |   local host    |
 +------|---------+                 +-------|--------+                 +-----------------+              +-----------------+                   
    libp2p stream                     libp2p stream 
        |                                   |  
    Discovery                           Discovery
 +------|---------+                 +-------|--------+                 +-----------------+              +-----------------+
 |                |                 |                |                 |                 |              |                 |
 |                | libp2p stream   |                | libp2p stream   |                 |    HTTP      |                 |
 |   Relay Node   <----------------->   Relay Node   <----------------->    Edge Node    <-------------->   Webapp SERVER |
 |                |   Discovery     |                | p2p Reservation |                 |              |                 |
 |  libp2p host   |                 |  libp2p host   |                 |   libp2p host   |              |   local host    |
 +------|---------+                 +-------|--------+                 +-----------------+              +-----------------+                   
    libp2p stream                     libp2p stream 
        |                                   |  
    Discovery                           Discovery
 +------|---------+                 +-------|--------+                 +-----------------+              +-----------------+
 |                |                 |                |                 |                 |              |                 |
 |                | libp2p stream   |                | libp2p stream   |                 |    HTTP      |                 |
 |   Relay Node   <----------------->   Relay Node   <----------------->    Edge Node    <-------------->   Webapp SERVER |
 |                |   Discovery     |                | p2p Reservation |                 |              |                 |
 |  libp2p host   |                 |  libp2p host   |                 |   libp2p host   |              |   local host    |
 +------|---------+                 +-------|--------+                 +-----------------+              +-----------------+                   
    libp2p stream                     libp2p stream 
        |                                   |  
       ...                                 ...
```
Refer to the [Edge-Matrix-Computing](https://github.com/shawnstartup/edge-matrix-computing) project for usage


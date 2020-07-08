# Grover-Server - Fits Image Extractor
[![Build Status](https://ci.templetron.io/api/badges/Grover/grover-server/status.svg)](https://ci.templetron.io/Grover/grover-server)  
Grover allows you to visualize your fits image files as well as provide the metadata contained within each file: 
- Visualize and store fits files + images
- Convert extracted images to videos and mosiacs
- View metadata extracted from file

### Usage
***
**STARTING THE WEB SERVER**
***
**Make sure you have golang >v1.14**  
**If using [Grover-Client](https://github.com/Grover-Fits/grover-client) make sure to build client before running server configuration**
1. create .env file
```bash
    make configure
```
2. build & run the executable
```bash
    make start
```

You can access the site from http://localhost:8080
***
#### Known Issues
* Handling of fits file images with floating point bitpix value (-32, -64) renders a insufficient image. For now use fits files with only integer based bitpix values
* Large fits file handling crashes the client

***
## Contributing
[Contributions](https://github.com/Grover-Fits/grover-server/issues?q=is%3Aissue+is%3Aopen) are welcome. If interested, fork this repo and submit a PR.


#### TODO:
* GRPC TLS support
* implement support for tables / arrays being extracted from fits files
* ~~support single file, multi image fits data~~
* docker image
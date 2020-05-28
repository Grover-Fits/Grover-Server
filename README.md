# Grover-Server - Fits Image Extractor

Grover allows you to visualize your fits image files as well as provide the metadata contained within each file: 
- Visualize and store fits files + images
- View metadata extracted from file
- Convert extracted images to a video file(.mkv)

### Usage
***
**STARTING THE GRPC SERVER**
***
**Make sure you have golang >v1.14**  
build & run the executable
```bash
    make start
```

You can access the site from http://localhost:8085
***
#### TODO:
* TLS support
* implement support for tables / arrays being extracted from fits files
* support single file, multi image fits data
* docker image
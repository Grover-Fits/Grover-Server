#!/bin/bash
echo Grover -- Fits File Extractor -- Server Configuration
echo Where is your client located? ex. /opt/www/grover/client
read clientPath

echo CLIENT_PATH=$clientPath > .env
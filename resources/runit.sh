#!/bin/bash
set -x

curl -v https://localhost:8443/blank.html?baseUrl=https://www.google.com
#curl -v https://localhost:8443/caddy-vars-regex?baseUrl=https://github.com:443/amalto

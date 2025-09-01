#!/bin/bash

docker run --rm -p 0.0.0.0:8080:80 -v ./:/usr/share/nginx/html nginx


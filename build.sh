#!/bin/bash

docker build -t ${2:-"rocklviv/cloudflare-dns-updater"} .
docker login -u rocklviv -p dckr_pat_is5Nh4nnnM5zpadKOnneyMt9d58

docker tag rocklviv/cloudflare-dns-updater:latest rocklviv/cloudflare-dns-updater:${1}
docker push rocklviv/cloudflare-dns-updater:${1}
#!/usr/bin/env bash
function get_cuda_image_name {
    output=`curl -s https://hub.docker.com/v2/repositories/nvidia/cuda/tags/\?page_size\=100\&page\=1\&name\=cudnn8-devel-ubuntu20.04 | jq ".results | .[] | .name"`
    echo $output
}

get_cuda_image_name
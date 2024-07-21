#!/bin/bash

if [[ $1 == "up" ]]; then
    cd ./auth-msv
    make migrate-up

    cd ..
    cd ./order-msv
    make migrate-up

    cd ..
    cd ./shopcart-msv
    make migrate-up

    cd ..
    cd ./warehouse-msv
    make migrate-up
elif [[ $1 == "down" ]]; then
    cd ./warehouse-msv
    make migrate-down

    cd ..
    cd ./shopcart-msv
    make migrate-down

    cd ..
    cd ./warehouse-msv
    make migrate-down

    cd ..
    cd ./order-msv
    make migrate-down

    cd ..
    cd ./auth-msv
    make migrate-down
else 
    echo "Invalid parameters: ./migrator.sh [up|down]. If agrs are empty, runs 'up'"
fi

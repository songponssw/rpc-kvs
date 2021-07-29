#!/bin/bash

set -m

./storage &

sleep 2; ./frontend

fg %1
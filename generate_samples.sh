#!/bin/bash

mkdir -p samples
head -c 1M /dev/urandom > samples/1M
head -c 5M /dev/urandom > samples/5M
head -c 10M /dev/urandom > samples/10M

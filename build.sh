#!/bin/bash

./release/user-package.sh nosource noconf codename=$(git describe --tags) buildname=docker-fly

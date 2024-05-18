#!/bin/bash
docker run --env-file .env -p 50051:50051 docker.io/anmho/noteservice:latest

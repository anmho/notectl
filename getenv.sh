#!/bin/bash

aws secretsmanager get-secret-value \
  --secret-id noteservice-env \
  --query SecretString \
  --output text | tee .env

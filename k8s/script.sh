#!/bin/bash

# Get the password from the secret
password=$(kubectl get secrets test-backend-secrets -o jsonpath='{.data.DB_PASSWORD}' | base64 -d)

# The value you want to compare with
expected_password="swRPxUNQLmjaDMUF"

echo $password
echo $expected_password

# Perform the equality check
if [ "$password" == "$expected_password" ]; then
    echo "The passwords match."
else
    echo "The passwords do not match."
fi
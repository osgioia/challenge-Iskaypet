#!/bin/bash

# Compile main.go for Lambda (setting GOOS=linux because Lambda uses a Linux environment)
GOOS=linux GOARCH=amd64 go build -o main .

# Install AWS SAM CLI if not already installed
if ! command -v sam &> /dev/null
then
    echo "AWS SAM CLI is not installed. Proceeding with installation..."
    # Download the installation script
    curl -Lo aws-sam-cli-linux-x86_64.zip https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-linux-x86_64.zip
    # unzip the file
    unzip aws-sam-cli-linux-x86_64.zip -d sam-installation
    # Grant execute permissions to the script
    chmod +x ./sam-installation/install
    # Run the installation script with superuser permissions
    sudo ./sam-installation/install
    # Remove the installation script after installation
    rm aws-sam-cli-linux-x86_64.zip
    # Ensure the PATH is updated (a terminal restart may be required)
    export PATH="$HOME/.aws-sam-cli/bin:$PATH"
else
    echo "AWS SAM CLI is already installed."
fi

# Verify if SAM CLI is now in the PATH
if command -v sam &> /dev/null
then
    # Enter the directory
    cd lambda
    # Build the project with SAM
    sam build

    # Deploy to AWS using SAM. The --guided flag provides an interactive deployment guide
    sam deploy --guided

    # Message when done
    echo "Deployment completed!"
else
    echo "AWS SAM CLI could not be installed correctly or is not in the PATH."
fi

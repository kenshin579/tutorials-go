#!/usr/bin/env bash

#https://artifacthub.io/packages/helm/jenkinsci/jenkins
# helm repo add jenkins https://charts.jenkins.io

kubectl apply -f volumes.yaml
helm install my-jenkins jenkins/jenkins -f values.yaml

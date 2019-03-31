# Rickshaw

This is a typical 3-tier architecture application with a web layer, app layer and a database. Will be deploying this application to a kubernetes cluster.

| Tier | Image | Image source | Notes |
| --- | --- | --- | --- |
| web | nginx 1.14 | docker hub | static content |
| app | api 0.1 | build from sources| golang app |
| db | percona 8.0 | docker hub | database |

## Kubernetes
Setup kubernetes cluster using [cornet](https://github.com/jostho/cornet)

## Build

Ssh into bastion

    ssh -A ec2-user@54.x.x.x

Get sources from github

    cd ~/src && git clone git@github.com:jostho/rickshaw.git

Build app container image

    cd ~/src/rickshaw/app && make image

Verify the built image

    buildah images
    buildah inspect jostho/api:v0.1.0

Push the image to the private registry

    buildah push --tls-verify=false jostho/api:v0.1.0 docker://registry:5000/api:v0.1.0

Verify the image in the registry

    skopeo inspect --tls-verify=false docker://registry:5000/api:v0.1.0

## Deploy

Label system namespaces

    kubectl label namespace/kube-system role=kube-system
    kubectl label namespace/monitoring role=monitoring

Create service monitor in monitoring namespace

    kubectl -n monitoring apply -f ~/src/rickshaw/k8s/servicemonitor.yaml

Create a new namespace to deploy the application

    kubectl create ns lane1

Setup role

    kubectl -n lane1 apply -f ~/src/rickshaw/k8s/role.yaml

Setup network policy

    kubectl -n lane1 apply -f ~/src/rickshaw/k8s/networkpolicy.yaml

Setup the config

    kubectl -n lane1 create configmap mysql-initdb-config --from-file ~/src/rickshaw/db/data.sql
    kubectl -n lane1 create configmap nginx-root-config --from-file ~/src/rickshaw/web/
    kubectl -n lane1 apply -f ~/src/rickshaw/k8s/configmap.yaml
    kubectl -n lane1 apply -f ~/src/rickshaw/k8s/secret.yaml

Setup persistent volume

    kubectl -n lane1 apply -f ~/src/rickshaw/k8s/persistentvolume.yaml

Setup db

    kubectl -n lane1 apply -f ~/src/rickshaw/k8s/db.yaml

Setup app

    kubectl -n lane1 apply -f ~/src/rickshaw/k8s/app.yaml

Setup web

    kubectl -n lane1 apply -f ~/src/rickshaw/k8s/web.yaml

Setup ingress

    kubectl -n lane1 apply -f ~/src/rickshaw/k8s/ingress.yaml

## Test
Add hosts entry for bastion, and then curl from localhost

    curl -i http://rickshaw-web.example.com/
    curl -i http://rickshaw-app.example.com/

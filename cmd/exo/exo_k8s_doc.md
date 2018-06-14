# Deploy Kubernetes with EXO

## Synopsis

- In this documentation we will show you how to deploy kubernetes cluster with a simple exo command `$ exo k8s`

## Prerequisites

- Make sure you have [`rke`](https://github.com/rancher/rke/releases) and `ssh` program to your `$PATH`

## Usage

### EXO k8s command

```
$ exo k8s -h
Manage kubernetes clusters

Usage:
  exo k8s [command]

Available Commands:
  create      Create k8s cluster
  delete      Remove kubernetes cluster
  list        List kubernetes cluster(s)

Flags:
  -h, --help   help for k8s

Global Flags:
      --config string   Specify an alternate config file [env CLOUDSTACK_CONFIG]
  -r, --region string   config ini file section name [env CLOUDSTACK_REGION] (default "cloudstack")

Use "exo k8s [command] --help" for more information about a command.
$
```

- You can `create`, `delete` and  `list` your cluster(s)

### EXO k8s create command

```
$ exo k8s create -h
Create k8s cluster

Usage:
  exo k8s create <name> [flags]

Flags:
  -f, --firewall-rules-add      Add firewall rules for kubernetes (if --node not set)
  -h, --help                    help for create
      --no-auto
  -n, --node string             Node can provision existing instances [vm name | id, vm name | id,...]
      --node-capacity string    Node(s) capacity (if --node not set) (micro|tiny|small|medium|large|extra-large|huge|mega|titan) (default "medium")
      --node-number int         Node number to create (if --node not set) (default 1)
  -s, --security-group string   Create node(s) in a security group <security group name | id> (if --node not set) (default "default")

Global Flags:
      --config string   Specify an alternate config file [env CLOUDSTACK_CONFIG]
  -r, --region string   config ini file section name [env CLOUDSTACK_REGION] (default "cloudstack")
$
```

## Let's deploy !


### Configure your Exoscale firewall environement

```
$ exo firewall create k8s-security-group --description " Security group for my new kubernetes cluster"
```
### Deploy K8s

```
$ exo k8S create example-cluster --node-number 3 --security-group k8s-security-group --firewall-rules-add
```

#### this folowing command deploy k8s cluster:

- `--node-number 3` create a cluster with 3 nodes
- `--security-group k8s-security-group` put all nodes in `k8s-security-group` security group
- `--firewall-rules-add` create all firewall rule for `rke` and `kubernetes` connections in `k8s-security-group` security group

#### See the output

```
Creating sshkey
Deploying node-1...
Creating sshkey
Deploying node-2.....
Creating sshkey
Deploying node-3..
Installing Docker on node(s).....................
RKE install:
INFO[0000] Building Kubernetes cluster

......

INFO[0118] [addons] Successfully Saved addon to Kubernetes ConfigMap: rke-ingress-controller
INFO[0118] [addons] Executing deploy job..
INFO[0124] [ingress] ingress controller nginx is successfully deployed
INFO[0124] [addons] Setting up user addons
INFO[0124] [addons] no user addons defined
INFO[0124] Finished building Kubernetes cluster successfully

Your kubernetes cluster "example-cluster" is up !

"cluster.yml" file was created to use kubectl program

If you want use kubectl program:
    juste copy "kube_config_cluster.yml" file in "$HOME/.kube/config"
    "$ cp ~/.exoscale/k8s/clusters/example-cluster/kube_config_cluster.yml ~/.kube/config"

OR
    "$ kubectl --kubeconfig ~/.exoscale/k8s/clusters/example-cluster/kube_config_cluster.yml [all kubectl commands]"

If you want to connect to your web kubernetes dashboard follow this link:
    WIP https://dfjskjskl.com
```

#### congratulation you cluster is deployed !

## Setting up kubectl

- You can folow instruction on the result output to use kubctl

```
$ cp ~/.exoscale/k8s/clusters/example-cluster/kube_config_cluster.yml ~/.kube/config

```

- Try a kubectl command

```
$ kubectl cluster-info
Kubernetes master is running at https://159.100.250.110:6443
KubeDNS is running at https://159.100.250.110:6443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
```

## Access to kubernetes web dashboard

### Setting up the Dashboard

[Web UI (Dashboard)](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/)

```
$ kubectl create -f https://raw.githubusercontent.com/kubernetes/dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml
```

### Admin locally

```
$ kubectl --kubeconfig ~/.kube/config proxy
Starting to server on 127.0.0.1:8001
```
or
```
$ kubectl --kubeconfig ./kube_config_cluster.yml proxy
Starting to server on 127.0.0.1:8001
```

<http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/>

### Login using a token

To log into the dashboard, you need to authenticate as somebody or something (Service Account), a clever hack[®](https://github.com/kubernetes/dashboard/issues/2474#issuecomment-365704926) is to pick the token from the _clusterrole-aggregation-controller_.

```
$ kubectl -n kube-system describe secrets \
   `kubectl -n kube-system get secrets | awk '/clusterrole-aggregation-controller/ {print $1}'` \
       | awk '/token:/ {print $2}'
```

then copy and paste it into the token login

Congratulation you are in your kubernetes dashboard !


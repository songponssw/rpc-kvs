# vHive Single Node Installation By P'Boss

Step by Step Instruction to install vHive (Serverless Platform) as Single Node cluster on bare metal

### System

- OS: Ubuntu 20.04

### Pre-requirements

- Support KVM
- Sufficient CPU and Memory

## Firewall-cmd

Turn on port that required by the system

```
sudo firewall-cmd --permanent --add-port=6643/tcp
sudo firewall-cmd --permanent --add-port=10250/tcp
sudo firewall-cmd --reload
```

## Clone vHive Repository

Clone to the workspace

```
git clone https://github.com/ease-lab/vhive.git
cd vhive
git checkout tags/v1.3​
```

## Setup all node

Install dependencies and configuration files

```
mkdir -p /tmp/vhive-logs
bash -x ./scripts/cloudlab/setup_node.sh > >(tee -a /tmp/vhive-logs/setup_node.stdout) 2> >(tee -a /tmp/vhive-logs/setup_node.stderr >&2)
```

## Start Background Process

Start `containerd` process

```
sudo screen -dmS containerd bash -c "containerd > >(tee -a /tmp/vhive-logs/containerd.stdout) 2> >(tee -a /tmp/vhive-logs/containerd.stderr >&2)"; sleep 5;
```

Start Firecracker process

```
sudo PATH=$PATH screen -dmS firecracker bash -c "/usr/local/bin/firecracker-containerd --config /etc/firecracker-containerd/config.toml > >(tee -a /tmp/vhive-logs/firecracker.stdout) 2> >(tee -a /tmp/vhive-logs/firecracker.stderr >&2)"; sleep 5;
```

Build `vhive` executable file and run

```
source /etc/profile && go build
sudo screen -dmS vhive bash -c "./vhive > >(tee -a /tmp/vhive-logs/vhive.stdout) 2> >(tee -a /tmp/vhive-logs/vhive.stderr >&2)"; sleep 5;
```

## !!! Edit pod-cidr, metallb to 172.16.0.0/16

In case, local network is the same subnet as the default setting provided from vhive script (192.168.0.0/16)
so, this config is required

```
# @kubeadm init
create_one_node_cluster.sh
# @ip pool
configs/metallb/metallb-configmap.yaml
# cidr_pod
canal.yaml
```

## Install kube

Start install the cluster, might take some time

```
./scripts/cluster/./scripts/cluster/create_one_node_cluster.sh
```

If no problem, the installation is successful

## (If Error) Update istio-key

If Error: Istio-network configmap not found
It might unsuccessful install the config map so apply this to install again

```
kubectl apply --filename https://github.com/knative/net-istio/releases/download/v0.19.0/release.yaml
```

# Cleanup

If installation gone wrong

```
# down the cluster, delete all config created by vhive script
./scripts/github_runner/clean_cri_runner.sh

# down bridge created by canal.yaml
sudo if link delete flannel.1

# purge everything out
sudo apt-get purge docker-ce docker-ce-cli containerd.io
sudo apt-get purge kubeadm kubectl kubelet kubernetes-cni kube*
```

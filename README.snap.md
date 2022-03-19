# Using etcdadm in snap environment

## Initialization

```
ubuntu@etcd-1:~$ /root/etcdadm init --init-system snapctl --data-dir /var/snap/microk8s/current/etcd --certs-dir /var/snap/microk8s/current/certs/etcdadm --snap-data-dir /var/snap/microk8s/current/args --install-dir /snap/microk8s/current
```

## Reset

```
ubuntu@etcd-1:~$ /root/etcdadm reset --init-system snapctl --data-dir /var/snap/microk8s/current/etcd --certs-dir /var/snap/microk8s/current/certs/etcdadm --snap-data-dir /var/snap/microk8s/current/args
```

## Join node

```bash
# copy ca.crt and ca.key from initial node
ubuntu@etcd-2:~$ mkdir /var/snap/microk8s/current/certs/etcdadm
ubuntu@etcd-1:~$ scp /var/snap/microk8s/current/certs/etcdadm/ca.crt etcd-2:/var/snap/microk8s/current/certs/etcdadm/ca.crt
ubuntu@etcd-1:~$ scp /var/snap/microk8s/current/certs/etcdadm/ca.key etcd-2:/var/snap/microk8s/current/certs/etcdadm/ca.key

# join second node
ubuntu@etcd-2:~$ /root/etcdadm join --init-system snapctl --data-dir /var/snap/microk8s/current/etcd --certs-dir /var/snap/microk8s/current/certs/etcdadm --snap-data-dir /var/snap/microk8s/current/args --install-dir /snap/microk8s/current https://etcd-1:2379
```

## kube-apiserver arguments

```
--etcd-servers="https://127.0.0.1:2379"
--etcd-cafile=$SNAP_DATA/certs/etcdadm/ca.crt
--etcd-certfile=$SNAP_DATA/certs/etcdadm/apiserver-etcd-client.crt
--etcd-keyfile=$SNAP_DATA/certs/etcdadm/apiserver-etcd-client.key
```

## Entrypoint for service

```
#!/bin/bash
set -e
export PATH="$SNAP/usr/sbin:$SNAP/usr/bin:$SNAP/sbin:$SNAP/bin:$PATH"
ARCH="$($SNAP/bin/uname -m)"
export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:$SNAP/lib:$SNAP/usr/lib:$SNAP/lib/$ARCH-linux-gnu:$SNAP/usr/lib/$ARCH-linux-gnu"
export LD_LIBRARY_PATH=$SNAP_LIBRARY_PATH:$LD_LIBRARY_PATH

set -o allexport
if [ -e $SNAP_DATA/args/etcd.env ]
then
  . $SNAP_DATA/args/etcd.env
else
  echo "Missing etcd configuration file!"
  exit 1
fi

$SNAP/etcd
```

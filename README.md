# Kubernetes GitHub Token Authenticator 
Authenticates Kubernetes users against a Github organization

### Create DaemonSet
```sh
kubectl create -f https://raw.githubusercontent.com/TwistoPayments/k8s-github-auth/master/manifests/github-auth.yaml
```
### Configure Kubernetes for Auth
There are two options how to modify kube-apiserver:
* A) Create kubeadm-config.yaml and apply new configuration
```yaml
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: vX.XX.X
apiServer:
    extraArgs:
        authentication-token-webhook-config-file: /srv/k8s/auth/token-webhook-config.json
        authentication-token-webhook-cache-ttl: 30m
    extraVolumes:
      - name: "github-auth"
        hostPath: /srv/k8s/auth
        mountPath: /srv/k8s/auth
        readOnly: true
        pathType: Directory
```
Check changes: `kubeadm upgrade diff --config kubeadm-config.yaml`

Apply changes to cluster: `kubeadm upgrade apply --config kubeadm-config.yaml`

* B) Direct changes in /etc/kubernetes/manifests/kube-apiserver.yaml. Insert these lines to specified sections. 
*spec.containers.command.kube-apiserver*
```yaml
    - --authentication-token-webhook-config-file=/srv/k8s/auth/token-webhook-config.json
    - --authentication-token-webhook-cache-ttl=30m
```
*spec.containers.command.volumeMounts*
```yaml
- mountPath: /srv/k8s/auth
      name: github-auth
      readOnly: true
```
*spec.volumes*
```yaml
  - hostPath:
      path: /srv/k8s/auth
      type: Directory
    name: github-auth
```
Kube-apiserver apply changes automatically.
### Group Role-based configuration (RBAC)
Create testing namespaces `kubectl create namespace testnamespace`

Assign user to specific team within Github and modify this group in RBAC.yaml  
```sh
kubectl create -f https://raw.githubusercontent.com/TwistoPayments/k8s-github-auth/master/manifests/RBAC.yaml
```
## Inspiration
This project was inspired by [oursky/kubernetes-github-authn](https://github.com/oursky/kubernetes-github-authn).
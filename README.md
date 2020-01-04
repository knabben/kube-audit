Dynamic Kubernetes Audit
---

To run the APIserver localy, setup the apiserver with:

```
$ make build-apiserver
$ make run-apiserver
docker run -it -p 8080:8080 apiserver:latest
``` 

More details in:

https://knabben.github.io/posts/kube-audit/

### Audit policy file

Create the Kind cluster using the config.yaml file, this will export the local volumes
and mount the *config/audit.yaml* file on */etc/kubernetes/audit.yaml*, after that
the audit-policy-file flag from APIServer uses this file. 

```
$ cd config
config$ kind create cluster -v=8 --config config.yaml
```

### Dynamic backend

Get the API address from the container APIServer and use on kube.yaml:

```  
  webhook:
    throttle:
      qps: 10
      burst: 15
    clientConfig:
      url: "http://172.17.0.3:8080/"
```

Start the audit registration with Kubectl 
```
kubectl create -f config/kube.yaml

```

Immutable Kubernetes Audit Log
---

Allows the consumption of an audit sink and later store in a blockchain test network. 
Creating a trusted, immutable and shared tree of events from a Kubernetes cluster.    

### Setup

Build the APIserver image with:

```
$ make build-apiserver
``` 

Setup the container with the docker-compose configuration, this brings the Sawtooth network + the generated server:
```
docker-compose -f config/sawtooth-default.yaml up
```

PS: Kind and docker network, is necessary to change the default bridge network for the user **sawtooth_default** for this lab. 

### Dynamic backend

Set the host name in the audit registration (config/kube.yaml):

```  
apiVersion: auditregistration.k8s.io/v1alpha1
kind: AuditSink
metadata:
  name: mysink
spec:
  policy:
    level: RequestResponse
    stages:
    - RequestReceived
  webhook:
    throttle:
      qps: 10
      burst: 15
    clientConfig:
      url: "http://apiserver:8080/"
```

Start the audit registration with Kubectl
 
```
kubectl create -f config/kube.yaml

```

This sink will send all the data to the API server, that will forward the traffic to a sawtooth test network.

### Check blockchain state

After setting up the TransactionProcessor, you can check the state and transactions list.

```
root@0033771beeee:/# sawtooth state list --url http://sawtooth-rest-api-default:8008
ADDRESS                                                                 SIZE  DATA
000000a87cb5eafdcca6a8c983c585ac3c40d9b1eb2ec8ac9f31ff5ca4f3850ccc331a  45    b'\n+\n$sawtooth.consensus.algorithm.version\x12\x030.1'
000000a87cb5eafdcca6a8c983c585ac3c40d9b1eb2ec8ac9f31ff82a3537ff0dbce7e  46    b'\n,\n!sawtooth.consensus.algorithm.name\x12\x07Devmode'
000000a87cb5eafdcca6a8cde0fb0dec1400c5ab274474a6aa82c12840f169a04216b7  110   b'\nl\n&sawtooth.settings.vote.authorized_keys\x12B03f14e484e7e65779ccb7c28e51e7e5ba0245e6e6a343d14154ccf2e4de6e11148'
3ad5f1005259702d8a02a566bead83a834b797dfe3c1162efbba35a1ea2236edb66f9d  312   b'{"RequestURI":"/api/v1/services?resourceVersion=203\\u0026timeout=8m27s\\u0026timeoutSeconds=507\\u0026watch=true","Verb"
3ad5f102b37cddf6a5b1ab6c4e077066235cc9b1fee542279c90c5e6cd7de107ff21d4  268   b'{"RequestURI":"/apis/storage.k8s.io/v1?timeout=32s","Verb":"get","Code":0,"User":{"Username":"system:serviceaccount:kube-
3ad5f10302b6c3011f89e22337fa2627365efb078c8e28fdb8e96bb0f1d58f3cac7c90  287   b'{"RequestURI":"/apis/admissionregistration.k8s.io/v1beta1?timeout=32s","Verb":"get","Code":0,"User":{"Username":"system:s
3ad5f10333c83c0685458c943fa2c01ba17e15b108a439fda5f646f960816599536840  274   b'{"RequestURI":"/apis/apiextensions.k8s.io/v1?timeout=32s","Verb":"get","Code":0,"User":{"Username":"system:serviceaccount
3ad5f109595ccd117ed6a5a9f3bcb6418237365e769a5eb77c252d928f475a7238ae72  265   b'{"RequestURI":"/api/v1/namespaces/kube-system/endpoints/kube-scheduler?timeout=10s","Verb":"update","Code":0,"User":{"Use
3ad5f10b3a2cd98529f629f5f2611b85eff5fbbd09f53708c60be6073acaafb1537eb8  232   b'{"RequestURI":"/api/v1/namespaces/default/pods?limit=500","Verb":"list","Code":0,"User":{"Username":"kubernetes-admin"},"
3ad5f10cfd350849aa6501644ca3f8b5d44c415025d697c890b06d3ade5d4f5561e37d  276   b'{"RequestURI":"/apis/networking.k8s.io/v1beta1?timeout=32s","Verb":"get","Code":0,"User":{"Username":"system:serviceaccou
3ad5f10d04d61e070a374600347d5870b0980dcecf619e7c33e2f876450c7981511990  285   b'{"RequestURI":"/api/v1/namespaces/kube-system/events/coredns-6955765f44-pw2wh.15e6c82ec0e1c07c","Verb":"patch","Code":0,"
3ad5f10f130e6e0ac2071adb19b12de05fde4bdb034bed214b788de08e6f6b14df7944  273   b'{"RequestURI":"/api/v1/namespaces/kube-system/pods/coredns-6955765f44-r66gx/status","Verb":"patch","Code":0,"User":{"User
3ad5f111ba137e45c6da8ee338e57da7bda99c8f48896ea15d57c7b675de712e39c71b  280   b'{"RequestURI":"/api/v1/namespaces/kube-system/endpoints/kube-controller-manager?timeout=10s","Verb":"get","Code":0,"User"
3ad5f11613367fed4cb7aadfc41a3565ed6792d6f39c83274c4bc773560c53bca36166  274   b'{"RequestURI":"/apis/discovery.k8s.io/v1beta1?timeout=32s","Verb":"get","Code":0,"User":{"Username":"system:serviceaccoun
3ad5f11bb118ddf4f749bcdf53e38781fa6b7f5ff26a68edfaecfeee92528516d1e5bb  220   b'{"RequestURI":"/api/v1/namespaces/kube-public","Verb":"get","Code":0,"User":{"Username":"system:apiserver"},"Resource":""
3ad5f11f856b2786daaf0ad8998ea7350aeaaa67b467af3b5b3a5d6f963e2a43c49c91  265   b'{"RequestURI":"/apis/autoscaling/v1?timeout=32s","Verb":"get","Code":0,"User":{"Username":"system:serviceaccount:kube-sys
3ad5f120ad3e19debc205ded6731a2a9f09c65994385bd7d2be1aca12b416ce4a2960f  264   b'{"RequestURI":"/api/v1/namespaces/kube-system/pods/coredns-6955765f44-r66gx","Verb":"get","Code":0,"User":{"Username":"sy
```

More details in:

https://knabben.github.io/posts/kube-audit/

# How to reproduce

Setup the cluster, argo and the operator that rejects all pods.
```
kind create cluster --name mytest --wait 5m --image "kindest/node:v1.16.15@sha256:c10a63a5bda231c0a379bf91aebf8ad3c79146daca59db816fb963f731852a99"

# Install argo
kubectl create ns argo
kubectl apply -n argo -f yamls/install-argo.yaml --validate=false

# Install webhook
docker build -t mywebhook:1.0.0 -f deny-all-webhook/Dockerfile deny-all-webhook
kind load docker-image --name mytest mywebhook:1.0.0
kubectl apply -f yamls/webhook.yaml
```

Run some workflows
```
kubectl apply -f yamls/workflow-namespace.yaml
for i in `seq 1 100`; do sed "s/myworkflow/myworkflow$i/g" yamls/workflow.yaml | kubectl apply -f -; done
```

Just check that no pod exists
```
kubectl get pods -n test-workflow
```

Use the following command to see if a workflow is still RUNNING. Wait a couple of minutes before running the following command (just be be sure). In my experiment I waited around 8 minutes:
```
kubectl get workflows -n test-workflow
```

Examine a RUNNING workflow. In my tests, e.g. myworkflow96 was still running.
```
kubectl get workflows myworkflow96 -n test-workflow -o yaml
```

You will see something like this:
```
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"argoproj.io/v1alpha1","kind":"Workflow","metadata":{"annotations":{},"name":"myworkflow96","namespace":"test-workflow"},"spec":{"entrypoint":"start","templates":[{"container":{"args":["echo 'it should not run'"],"command":["/bin/sh","-c"],"image":"alpine:3.7"},"name":"start"}]}}
  creationTimestamp: "2021-03-01T15:28:09Z"
  generation: 2
  labels:
    workflows.argoproj.io/phase: Running
  name: myworkflow96
  namespace: test-workflow
  resourceVersion: "1766"
  selfLink: /apis/argoproj.io/v1alpha1/namespaces/test-workflow/workflows/myworkflow96
  uid: 716a6b23-51d7-42f5-b1e1-253681c3bf58
spec:
  arguments: {}
  entrypoint: start
  templates:
  - arguments: {}
    container:
      args:
      - echo 'it should not run'
      command:
      - /bin/sh
      - -c
      image: alpine:3.7
      name: ""
      resources: {}
    inputs: {}
    metadata: {}
    name: start
    outputs: {}
status:
  finishedAt: null
  nodes:
    myworkflow96:
      displayName: myworkflow96
      finishedAt: "2021-03-01T15:28:09Z"
      id: myworkflow96
      message: 'admission webhook "pod.validation.webhook" denied the request: Pod
        is not allowed'
      name: myworkflow96
      phase: Error
      progress: 1/1
      startedAt: "2021-03-01T15:28:09Z"
      templateName: start
      templateScope: local/myworkflow96
      type: Pod
  phase: Running
  progress: 1/1
  startedAt: "2021-03-01T15:28:09Z"
```

Cleanup
```
kind delete cluster --name mytest
```
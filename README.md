# admission-controller
demo k8s admission controller

# Dev Setup
## One Time
- Create Cluster
    ```shell
      > brew install kind                                                  # install a local k8s cluster using kind
      > kind create cluster --config dev/cluster-config.yaml --name demo   # create k8s cluster 
    ```
- Verify that cluster is created and set kubectl context
    ```shell
    > kind get clusters
    demo
    
    > kubectl --context kind-demo cluster-info
    Kubernetes master is running at https://127.0.0.1:55000
    KubeDNS is running at https://127.0.0.1:55000/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
    
    To further debug and diagnose cluster problems, use 'kubectl --context kind-demo cluster-info dump'.
    ```
- Create TLS server certificate [and ca]
  ```shell
  > make certs 
  ```
  Certs are located in the dev/certs folder.

- Build and run the latest code
  ```shell
  > make run
  ```
- Create namespace and enable validating webhook
  ```shell
  > ./k8s-setup.sh
  ```
  
## Testing Admission Controller
- Run `main` under the debugger
- Deploy a manifest to the cluster to see if it is intercepted by the admission controller
  ```
  kubectl --context kind-demo apply -f dev/busybox.yaml
  ```
# Custom Kubernetes Controller

## Description
A simple Kubernetes custom controller using the operator SDK

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [Minikube](https://minikube.sigs.k8s.io/docs/start/) to get a local cluster for testing
### Pull the image
You can find the image in this repo: https://hub.docker.com/repositories/tokhiarh

### Running on the cluster
  Install the CRDs into the cluster:

```sh
make install
```

or trigger it manually:

```sh
 kubectl apply -f config/crd/bases/dummy.interview.com_dummies.yaml
 ```

#### Creating the dummy Objects

Run the controller:

```sh
make run
```
and then:

```sh
    kubectl apply -f config/samples/dummy_v1alpha1_dummy.yaml
    kubectl apply -f config/samples/dummy_v1alpha2_dummy.yaml
```

Get the pods:

```sh
 kubectl get pods
```

Describe the pod:

```sh
  kubectl describe po dummy-sample
```


Delete the pod manually

```sh
kubectl delete pod dummy-sample
```

Now exit and then `make run` again to let the controller recreate the deleted pod automatically.

To delete a dummy object
```sh
kubectl delete dummy dummy-sample
```
Now the controller will aslo delete the related pod.


 #### Deployment
 Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=tokhiarh/dummy-controller:v0.1.1
```
Apply the dummy sample:

```sh
kubectl apply -f config/samples/dummy_v1alpha1_dummy.yaml
```

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


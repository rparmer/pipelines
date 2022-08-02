# CRD/ConfigMap Pipeline POC

The purpose of this project is to provide a minimal POC for how pipelines could be defined using labels, crds, and configmaps.  It only looks at how the various stages could be defined and does not look into how to promote deployments between the stages.

## Getting Started
- Install Flux onto your cluster following the [Flux bootstrapping procedures](https://fluxcd.io/docs/installation/#bootstrap)
- Install [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize)
- Fork this repo and update the `demo/repository.yaml` file to include your repo url and branch
- Run `make install` to build the generated files and install the CRDs onto the cluster
- Run `make demo` to deploy the example pipeline definitions and demo apps onto the cluster

## Usage
The demo cli has support for querying pipeline information based on a hybrid label w/ crd approach, full crd or configmap definition.  Each can be reached by running `go run .` and passing the pipeline type you want to test as an argument.

```bash
go run . # default - crd/label hybrid approach
go run . crd # full crd approach
go run . cm # configmap approach
```

There are also Makefile steps you can run instead

```bash
make run # default - crd/label hybrid approach
make run-crd # full crd approach
make run-cm # configmap approach
```

### CRD/Label Hybrid Approach
The hybrid approach uses a CRD to define a pipeline's stages (`dev`, `staging`, `prod`), but the pipeline definition does not have any deployment information directly associated to it.  Instead the application teams will add a `pipelines.wego.weave.works/name` label to their Flux `Kustomization` or `HelmRelease` config where this label value is the associated pipeline name.

```yaml
apiVersion: wego.weave.works/v1alpha1
kind: Pipeline
metadata:
  name: example-pipeline
  namespace: flux-system
spec:
  stages:
    - name: dev      # name of stage
      namespace: dev # namespace associated to stage
      order: 1       # stage order (dev > staging > prod)
      # cluster: ''  # optional field to support multi-cluster (not in scope for POC)
    - name: staging
      namespace: staging
      order: 2
    - name: prod
      namespace: prod
      order: 3
---
apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
kind: Kustomization
metadata:
  name: dev # maps to namespace defined under `dev` stage or pipeline
  labels:
    pipelines.wego.weave.works/name: example-pipeline # matches pipeline name above
spec:
  ...
---
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: podinfo-pipeline-helm
  namespace: default # not mapped to stage in pipeline definition, would not show up when looking up pipeline info
  labels:
    pipelines.wego.weave.works/name: example-pipeline # matches pipeline name above
spec:
  ...
```

With this approach operators have the control to define what namespace and/or cluster should be used for each stage and application teams have the flexibility to add/remove deployments from the pipeline just by adding labels.  Any labels added to resources that are not defined in a pipeline stage are ignored.  This will help prevent people from inadvertently (or intentionally) injecting an un-authorized deployment into the pipeline.

### Full CRD Approach
The full CRD approach removes the label option and moves the entire pipeline definition into the `Pipeline` resource.  It builds upon the hybrid CRD and adds a new `releaseRefs` field that holds the resource definitions for a stage.  This removes some of the flexibility the application team had in defining what resources belonged to a stage and will put a greater dependency on the operators.  But it will also give a clearer view of what a pipeline consists of.  The extra work on the operators could be lessened with automation templates and/or breaking the releaseRefs object into a seperate CRD that application teams can manage.

```yaml
apiVersion: wego.weave.works/v1alpha2
kind: Pipeline
metadata:
  name: example-pipeline
  namespace: flux-system
spec:
  stages:
    - name: dev
      namespace: dev
      order: 1
      releaseRefs: # list of `Kustomization` or `HelmRelease` objects. They MUST be in the defined stage namespace
        - name: podinfo-pipeline-helm
          kind: HelmRelease
        - name: dev # name of kustomization object (doesn't have to match stage name)
          kind: Kustomization
    - name: staging
      namespace: staging
      order: 2
      releaseRefs:
        - name: podinfo-pipeline-helm
          kind: HelmRelease
        - name: staging
          kind: Kustomization
    - name: prod
      namespace: prod
      order: 3
      releaseRefs:
        - name: podinfo-pipeline-helm
          kind: HelmRelease
        - name: prod
          kind: Kustomization
```

### Configmap Approach
The configmap approach looks to use existing Kubernetes resources without needing to create any CRDs.  This example mimics the hybrid approach, but it could be changed to mimic the full CRD approach as well.  It does work, but it seems to be very brittle and may be prone to config errors by the consumer.  For the purpose of this POC the cli demo is hard-coded to the configmap example below.  Feel free to play around with the stage definitions, but a configmap with the name `example-pipeline` is required.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: example-pipeline
  namespace: flux-system
data:
  pipelines: |
    stages:
      - name: dev
        namespace: dev
        order: 1
      - name: staging
        namespace: staging
        order: 2
      - name: prod
        namespace: prod
        order: 3
---
apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
kind: Kustomization
metadata:
  name: dev # maps to namespace defined under `dev` stage or pipeline
  labels:
    pipelines.wego.weave.works/name: example-pipeline # matches pipeline name above
spec:
  ...
---
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: podinfo-pipeline-helm
  namespace: default # not mapped to stage in pipeline definition, would not show up when looking up pipeline info
  labels:
    pipelines.wego.weave.works/name: example-pipeline # matches pipeline name above
spec:
  ...
```

## Results
The CRD examples do not address the use of controllers.  If it is decided that a controller is needed for reconciliation or to give better visibility to a pipeline then that will need to be addressed in a future spike.

Overall the hybrid or full CRD approach seem to be best suited for this solution.  The configmap approach will work, but it loses the safe guards and defined structure that the CRDs help provide.  I found it very easy to mess up the configmap data structure and I don't believe it will evolve well as the pipeline project matures.  Below are some summarized tables listing the pros and cons of each approach.

One issue with all these solutions is that the `Kustomization` object may contain multiple deployments.  For example in this demo `podinfo.yaml` and `podinfo-2.yaml` are defined under the same Kustomization and will both show up under the pipeline results.  We can recommend that application teams have a unique Kustomization per deployment, but there isn't a way to enforce that.  On the ui we would need to support this scenario (and that is why I included it in this demo).

### Hybrid
| Pros | Cons |
| ---- | ---- |
| Operator control of stage definitions | Full pipeline not directly visible (needs to be aggregated) |
| Application team flexibility | |
| RBAC (operators control stages) | RBAC (could be an issue with with multiple pipelines in same namespace) |

### CRD
| Pros | Cons |
| ---- | ---- |
| Operator control of stage definitions | Not easy for application teams to change deployments |
| Clear pipeline definition | |
| RBAC | |

### Configmap
| Pros | Cons |
| ---- | ---- |
| No CRDs needed | Difficult to implement RBAC |
| | Brittle and prone to errors |
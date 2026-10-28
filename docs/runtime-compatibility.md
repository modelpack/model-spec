# Runtime Compatibility

There are two ways to get a ModelPack artifact into a Kubernetes pod:

1. **Kubernetes image volumes** ([KEP-4639](https://kubernetes.io/docs/tasks/configure-pod-container/image-volumes/)) — the container runtime pulls and mounts the artifact directly. Support for ModelPack artifacts varies by runtime; see the matrix below.
2. **[Model CSI Driver](https://github.com/modelpack/model-csi-driver)** — a runtime-independent path that works on any CRI runtime and older Kubernetes versions.

_Last reviewed: 2026-10-28. Status tracked in [#233](https://github.com/modelpack/model-spec/issues/233)._

## Kubernetes Requirements

[Image volumes](https://kubernetes.io/docs/tasks/configure-pod-container/image-volumes/) (the `volumes[].image` field) require:

| Kubernetes version | Status | Notes |
|---|---|---|
| 1.31 | Alpha | Feature gate `ImageVolume` must be enabled |
| 1.33 | Beta (off by default) | `subPath`/`subPathExpr` supported |
| ~1.35 | Beta (on by default) | Exact version unconfirmed; see [KEP-4639](https://github.com/kubernetes/enhancements/issues/4639) |
| 1.36+ | GA | Enabled by default |

## Runtime Matrix

| Runtime | Min version | ModelPack artifact mounts | Notes |
|---|---|---|---|
| CRI-O | v1.33.0 | Yes | Raw layers only (see [Layer format notes](#layer-format-notes)). Each layer is mounted as a single file, named by the `org.opencontainers.image.title` OCI annotation (standard, not project-specific) or `org.cncf.model.filepath` (see [annotations.md](./annotations.md)) if present; layers with neither annotation are skipped. Can be disabled via CRI-O config. ([cri-o#8953](https://github.com/cri-o/cri-o/pull/8953), [cri-o#9062](https://github.com/cri-o/cri-o/pull/9062), [cri-o#9131](https://github.com/cri-o/cri-o/pull/9131)) |
| containerd | — | No | Tracked in [containerd/containerd#11381](https://github.com/containerd/containerd/issues/11381). Use the Model CSI Driver in the meantime. |
| moby / Docker Engine | — | No | Tracked in [moby/moby#49591](https://github.com/moby/moby/issues/49591). Use the Model CSI Driver in the meantime. |
| Model CSI Driver | v0.1.x | Yes | Runtime-independent. Works with any CRI runtime and compatible with older Kubernetes versions. Pulls and unpacks layers itself; supports P2P distribution via Dragonfly. |

## Layer Format Notes

CRI-O mounts each layer blob as a single file — it does **not** extract `.tar`, `.tar+gzip`, or `.tar+zstd` layers. Only `*.raw` media types result in a usable file on the mount path.

For CRI-O image volumes, use `application/vnd.cncf.model.weight.v1.raw` (or other `*.raw` variants) and set the `org.cncf.model.filepath` annotation on each layer so CRI-O names the mounted file correctly. See [annotations.md](./annotations.md) for the full annotation reference.

When packaging a model, ensure layers use `*.raw` media types. Consult your build tool's documentation for how to configure the layer format.

## Example: Image Volume with CRI-O

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: model-server
spec:
  containers:
    - name: server
      image: myregistry.com/myserver:latest
      volumeMounts:
        - name: model
          mountPath: /model
          readOnly: true
  volumes:
    - name: model
      image:
        reference: myregistry.com/mymodel:v1.0
        pullPolicy: IfNotPresent
```

## Model CSI Driver

The Model CSI Driver is the recommended path on containerd- or Docker-based nodes, or on Kubernetes older than 1.31. It handles both raw and tar layers independently of the container runtime.

See the [Model CSI Driver repository](https://github.com/modelpack/model-csi-driver) for installation and usage, including Helm-based deployment.

For the CSI driver usage example, see [Use Models in Kubernetes](./getting-started.md#use-models-in-kubernetes).

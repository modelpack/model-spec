# Runtime Compatibility

There are two ways to get a ModelPack artifact into a Kubernetes pod:

1. **Kubernetes [image volumes](https://kubernetes.io/docs/tasks/configure-pod-container/image-volumes/)** ([KEP-4639](https://github.com/kubernetes/enhancements/issues/4639)) — the container runtime pulls and mounts the artifact directly. Support for ModelPack artifacts varies by runtime; see the matrix below.
2. **[Model CSI Driver](https://github.com/modelpack/model-csi-driver)** — a runtime-independent path that works on any CRI runtime and older Kubernetes versions.

_Last reviewed: 2026-11-26. Upstream status: [containerd/containerd#11381](https://github.com/containerd/containerd/issues/11381), [moby/moby#49591](https://github.com/moby/moby/issues/49591)._

## Kubernetes Requirements

[Image volumes](https://kubernetes.io/docs/tasks/configure-pod-container/image-volumes/) (the `volumes[].image` field) are controlled by the `ImageVolume` [feature gate](https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/):

| Kubernetes version | Stage | Notes |
|---|---|---|
| 1.31–1.32 | Alpha (off by default) | Enable the `ImageVolume` feature gate |
| 1.33–1.34 | Beta (off by default) | Enable the `ImageVolume` feature gate; `subPath`/`subPathExpr` supported |
| 1.35 | Beta (on by default) | |
| 1.36+ | GA | Always enabled |

CRI-O artifact mounts need CRI-O v1.33 or later; otherwise use the Model CSI Driver.

## Runtime Matrix

| Runtime | Versions | ModelPack artifact mounts | Notes |
|---|---|---|---|
| CRI-O | v1.33.x | Yes | Names each mounted file from `org.opencontainers.image.title`, otherwise from the legacy `org.cnai.model.filepath`. Does **not** read `org.cncf.model.filepath`. |
| CRI-O | v1.34.0+ | Yes | Names each mounted file from `org.opencontainers.image.title`, otherwise from `org.cncf.model.filepath`. Does **not** read `org.cnai.model.filepath`. |
| containerd | — | No | Tracked in [containerd/containerd#11381](https://github.com/containerd/containerd/issues/11381). Use the Model CSI Driver in the meantime. |
| moby / Docker Engine | — | No | Tracked in [moby/moby#49591](https://github.com/moby/moby/issues/49591). Use the Model CSI Driver in the meantime. |
| Model CSI Driver | v0.1.x | Yes | Runtime-independent; works with any CRI runtime and older Kubernetes versions. Extracts `*.raw` and uncompressed `*.tar` layers; `*.tar+gzip` and `*.tar+zstd` are not supported. Reads `org.cncf.model.filepath` or `org.cnai.model.filepath`. Supports P2P distribution via Dragonfly. |

On every CRI-O version, a layer with neither of the two name annotations that version reads (`org.opencontainers.image.title` and its version-specific filepath key) is skipped. CRI-O only logs a warning on the node, and the pod still starts with the file missing. Artifact mounts can be turned off with `oci_artifact_mount_support = false` in the CRI-O configuration. See [cri-o#8953](https://github.com/cri-o/cri-o/issues/8953), [cri-o#9062](https://github.com/cri-o/cri-o/pull/9062), [cri-o#9131](https://github.com/cri-o/cri-o/pull/9131).

## Layer Format Notes

CRI-O bind-mounts each layer blob as a single file. It does **not** extract `.tar`, `.tar+gzip`, or `.tar+zstd` layers. For a usable CRI-O image volume, every layer must:

1. Use a `*.raw` media type, such as `application/vnd.cncf.model.weight.v1.raw`.
2. Carry a name annotation that the running CRI-O version reads. `org.opencontainers.image.title` set to the file's path is the recommended, portable choice: every CRI-O version from v1.33.0 reads it. Otherwise, use `org.cnai.model.filepath` on v1.33.x or `org.cncf.model.filepath` on v1.34.0+.

Also set `org.cncf.model.filepath` on every layer so the Model CSI Driver and other ModelPack tools can name the file (see [annotations.md](./annotations.md)). It is not a CRI-O requirement when `org.opencontainers.image.title` is set.

```json
{
  "mediaType": "application/vnd.cncf.model.weight.v1.raw",
  "digest": "sha256:3f907c1a03bf20f20355fe449e18ff3f9de2e49570ffb536f1a32f20c7179808",
  "size": 30327160,
  "annotations": {
    "org.opencontainers.image.title": "model.safetensors",
    "org.cncf.model.filepath": "model.safetensors"
  }
}
```

The Model CSI Driver does not read `org.opencontainers.image.title`. It requires `org.cncf.model.filepath` (or `org.cnai.model.filepath`) on raw layers. When the driver is configured with a Dragonfly endpoint, it requires `org.cncf.model.filepath` (or `org.cnai.model.filepath`) on every layer, including tar layers, and fails with `missing annotation filepath` otherwise. To serve CRI-O and the CSI driver from one artifact, set both `org.opencontainers.image.title` and `org.cncf.model.filepath` on every layer.

The [example manifest in the specification](./spec.md#example-image-manifest-for-model-artifacts) uses `.tar` layers without annotations. CRI-O would skip every layer, so it produces no usable files.

### Build tools

| Tool | Releases | Layer name annotation | CRI-O v1.33.x | CRI-O v1.34.0+ |
|---|---|---|---|---|
| [AIKit](./aikit.md) `packager/modelpack` | v0.21.0+ | `org.opencontainers.image.title` and `org.cncf.model.filepath` | Yes | Yes |
| [modctl](./modctl.md) | v0.1.0–v0.1.2, v0.2.1-cnai, v0.2.2, v0.2.3 (tag only; no GitHub release as of 2026-11-26) | `org.cnai.model.filepath` only (with `vnd.cnai.*` media types) | Yes | No, layers are skipped |
| [modctl](./modctl.md) | v0.2.0, v0.2.1 (pre-release on GitHub) | `org.cncf.model.filepath` only | No, layers are skipped | Yes |

AIKit (v0.21.0+) and modctl (v0.1.0 and later) build raw layers by default. Keep AIKit's default `layer_packaging=raw` and modctl's default `--raw=true`. No modctl release sets `org.opencontainers.image.title`.

Do not use the upstream `main` installation that the [modctl guide](./modctl.md) links to, because it writes only `org.cncf.model.filepath`. Install a tagged release that matches the table instead. For CRI-O v1.33.x, run `go install github.com/modelpack/modctl@v0.2.2` or use the matching release binary. Otherwise, pick the modctl release that matches your CRI-O version from the table above, or use the Model CSI Driver, which reads `org.cncf.model.filepath` or `org.cnai.model.filepath`.

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

The Model CSI Driver is the recommended path on containerd- or Docker-based nodes, on CRI-O older than v1.33, or on Kubernetes older than 1.31. It extracts `*.raw` and uncompressed `*.tar` layers independently of the container runtime (`*.tar+gzip` and `*.tar+zstd` are not supported in v0.1.x).

See the [Model CSI Driver repository](https://github.com/modelpack/model-csi-driver) for installation and usage, including Helm-based deployment.

For the CSI driver usage example, see [Use Models in Kubernetes](./getting-started.md#use-models-in-kubernetes).

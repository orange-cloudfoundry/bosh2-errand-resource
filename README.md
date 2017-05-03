# BOSH Errand Resource

A resource that will run an errand using the [BOSH CLI v2](https://bosh.io/docs/cli-v2.html).

## Differences from original BOSH Deployment Resource

The original [BOSH Errand Resource](https://github.com/starkandwayne/bosh-errand-resource)
uses the Ruby CLI and does not support newer BOSH features (for example UAA auth).

### Breaking Changes

* This resource requires that the target director's SSL certificate is trusted. If the director's certificate is not
 already trusted by normal root authorities, a custom CA certificate must be provided.

## Adding to your pipeline

To use the BOSH Errand Resource, you must declare it in your pipeline as a resource type:

```
resource_types:
- name: bosh-errand
  type: docker-image
  source:
    repository: starkandwayne/bosh2-errand-resource
```

## Source Configuration

* `deployment`: *Required.* The name of the deployment.
* `target`: *Optional.* The address of the BOSH director which will be used for the deployment. If omitted, target_file
  must be specified via out parameters, as documented below.
* `client`: *Required.* The username or UAA client ID for the BOSH director.
* `client_secret`: *Required.* The password or UAA client secret for the BOSH director.
* `ca_cert`: *Optional.* CA certificate used to validate SSL connections to Director and UAA. If omitted, the director's
  certificate must be already trusted.

### Example

``` yaml
- name: staging
  type: bosh-deployment
  source:
    deployment: staging-deployment-name
    target: https://bosh.example.com:25555
    client: admin
    client_secret: admin
    ca_cert: "-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"
```

### Dynamic Source Configuration

Sometimes source configuration cannot be known ahead of time, such as when a BOSH director is created as part of your
pipeline. In these scenarios, it is helpful to be able to have a dynamic source configuration. In addition to the
normal parameters for `put`, the following parameters can be provided to redefine the source:

* `source_file`: *Optional.* Path to a file containing a YAML or JSON source config. This allows the target to be determined
  at runtime, e.g. by acquiring a BOSH lite instance using the
  [Pool resource](https://github.com/concourse/pool-resource). The content of the `source_file` should have the same
  structure as the source configuration for the resource itself. The `source_file` will be merged into the exist source
  configuration.

_Note_: This is only supported for a `put`.

#### Example

```
- put: staging
  params:
    source_file: path/to/sourcefile
```

Sample source file:

```
{
    "target": "dynamic-director.example.com",
    "client_secret": "generated-secret",
    "vars_store": {
        "config": {
            "bucket": "my-bucket"
        }
    }
}
```

## Behaviour

### `out`: Run a BOSH errand

This will run any given errand for the specified deployment.

#### Parameters

* `name`: *Required.* Name of the errand to run

* `keep_alive`: *Optional.* Use existing VM to run an errand and keep it after completion

* `when_changed`: *Optional.* Run errand only if errand configuration has changed or if the previous run was unsuccessful

* `target_file`: *Optional.* Path to a file containing a BOSH director address.
  This allows the target to be determined at runtime, e.g. by acquiring a BOSH
  lite instance using the [Pool
  resource](https://github.com/concourse/pool-resource).

  If both `target_file` and `target` are specified, `target_file` takes
  precedence.

``` yaml
- put: staging
  params:
    name: smoke-tests
```

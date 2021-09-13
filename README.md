# OSM Health
OSM Health is a tool for troubleshooting issues with [Open Service Mesh](https://github.com/openservicemesh/osm). It can
be used to gather diagnostic information on the state of the control plane and the mesh components to ensure that 
everything is running as expected.


## Get started

1. Clone the repo
1. From root, build the binary

    ```bash
    make build-osm-health
    ```

1. Add it to `$PATH` to get started
    ```bash
    sudo cp ./bin/osm-health /usr/local/bin/osm-health
    ```

You can now use the CLI with:
```bash
osm-health <command>
```

## Commands
Currently, osm-health checks the connectivity between two pods by running a series of diagnostic checks on the meshed namespaces and pods, 
Envoy, SMI policies and core OSM control plane components. To run these checks, use:

```bash
osm-health connectivity pod-to-pod <SOURCE_POD> <DESTINATION_POD>
```

For the most up to date list of available commands:
```bash
osm-health --help
```

## Outcomes
A command runs a series of checks associated with that command.

Each check can return one of 4 outcomes:
1. `Pass`: indicates the check was successful and its result was as expected
1. `Fail`: indicates the check failed and returns the error that could be causing the failure. Failed checks highlight
   components that could require further investigation
1. `Info`: this is returned when the check is not generally expected to pass or fail, but rather the purpose of the 
   check is to simply provide information to the user. An info check prints out general diagnostic information generated
   by the check.
    > For example, when SMI TrafficTarget checks are run, they may return an `info` outcome that says that permissive 
   > traffic policy mode is enabled, so SMI access policies do not apply. Such an outcome cannot be categorized as a 
   > pass or a fail outcome because it is not an unexpected behavior
1. `Unknown`: this indicates the check could not come to a clear conclusion
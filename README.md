## cf-networking-helpers

This repo contains various helper tools used in our cf-networking and silk
releases. Some helpers may be used in testing, as well.

The `ci/` directory contains pipelines for testing the code in this repo with
different database configurations. The jobs are in the [opensource
CI](https://networking.ci.cf-app.com/teams/ga/pipelines/cf-networking-helpers).


### Running tests

To run tests use `./scripts/docker-test`.

This will use `dep` to get the dependencies necessary using the `Gopkg.toml` and `Gopkg.lock`.
Since this is a library those two files are only used to get dependencies for testing. DO NOT
vendor packages here.

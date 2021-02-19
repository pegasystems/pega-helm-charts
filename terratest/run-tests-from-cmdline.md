# Running tests from the command line

The test can also be run from the command line which avoids needing to have a license needed to use the IntelliJ plugin.

## Prerequisites
- A local clone of this repo
- Install an appropriate version go (v1.13+)
- Install dep

## Setup
Add {path to repo}/terratest to GOPATH environment variable.

Then:
- cd to {path to repo}/terratest/src/test
- run `dep ensure` -- this may take a little while and log warnings.

## Running Tests

Assuming that you're in the {path to repo}/terratest/src/test, you can run the following commands. 

Note that you can drop the `-v` if you do not care about the output.

### Running all tests
`go test -v test`

### Run All of the Tests for the Pega chart
`go test -v test/pega`

### Run a Single Test
`go test -v test/pega -run TestPegaTierDeployment` where `TestPegaTierDeployment` is the name of the test.

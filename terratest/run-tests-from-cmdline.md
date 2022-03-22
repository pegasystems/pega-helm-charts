# Running tests from the command line

The test can also be run from the command line which allows you to avoid complying with the IntelliJ plugin license requirements.

## Prerequisites
- You have a local clone of this repo
- You have installed an appropriate go version (v1.13 or later)
- You have installed dep
- You have disabled go modules by setting the environment variable `GO111MODULE="off"` (v1.16 or later)

## Setup
Add {path to repo}/terratest to the GOPATH environment variable on your local system.

Then:
- cd to {path to repo}/terratest/src/test
- run `dep ensure` -- this may take a little while and log warnings.

## Running Tests

From {path to repo}/terratest/src/test, run the following commands. 

You can omit the `-v` flag to suppress the output.

### Running all tests
`go test -v test`

### Run All of the Tests for the Pega chart
`go test -v test/pega`

### Run a Single Test
`go test -v test/pega -run TestPegaTierDeployment` where `TestPegaTierDeployment` is the name of the test.

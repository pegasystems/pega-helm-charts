# Running tests from the command line

## Prerequisites
- You have a local clone of this repo
- You have installed an appropriate go version (v1.13 or later)
- You enabled go modules by setting the environment variable `GO111MODULE="on"`

## Running Tests

From {path to repo}/terratest/src/test, run the following commands. 

You can omit the `-v` flag to suppress the output.

### Running all tests
`go test -v ./pega ./addons/ ./backingservices/`

### Run All of the Tests for the Pega chart
`go test -v ./pega`

### Run a Single Test
`go test -v ./pega -run TestPegaTierDeployment` where `TestPegaTierDeployment` is the name of the test.

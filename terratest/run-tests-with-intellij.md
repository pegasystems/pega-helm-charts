- Prerequisites:
    - make sure you have the latest version of Intellij installed
    - install Go - version >=1.13 is required
    - install Dep 
    - make sure you have helm version <3.0 installed, as version 3.0. is conflicting with terratest framework for running test through Intellij
- Install Go plugin for IntelliJ (only available for IntelliJ IDEA Ultimate edition)
- Go to settings (Ctrl + Alt + S)
- Find option "Languages & Frameworks | Go | GOPATH"
- In project GOPATH add {path to project}/terratest
- Find option "Languages & Frameworks | Go | GOROOT and select the Go version
- Add Dep Executable in "Languages & Frameworks | Go | Dep {path}\bin\dep but make sure NOT to select option "Enable dep integration"
- Click "Ok"
- from the commandline, while being in {path to project}/terratest/src/test run "dep ensure" after you close the settings. 
  Run it and wait while it will be finished and IntelliJ reindex dependencies. (you should be able to see vendor directory created in in src/test directory)
- Find  "Makefile" in the terratest directory and run the "helm dependency update ./charts/pega/" and "helm dependency update ./charts/addons/" from the commandline, 
  while being in the {path to project}/terratest to update the necessary helm dependencies;
- In order to run all tests you can do the following
    - create "Go test" configuration with package test kind and use "test" as a package
    - Run any single test using standard IntelliJ interface
#!/usr/bin/env groovy
def labels = ""

node {
      stage("Init"){
          // if (env.CHANGE_ID) {
            pullRequest.labels.each{
            echo "label: $it"
            validateProviderLabel(it)
            labels += "$it,"
            }
            labels = labels.substring(0,labels.length()-1)
            echo "PR labels -> $labels"
            // Comment on PR about the text execution
            pullRequest.comment("Starting ${env.BRANCH_NAME} validation on -> $labels")
            sh "curl -fsSL -o helm-v3.2.4-linux-amd64.tar.gz https://get.helm.sh/helm-v3.2.4-linux-amd64.tar.gz"
            sh "tar -zxvf helm-v3.2.4-linux-amd64.tar.gz"
            sh "mv linux-amd64/helm /usr/local/bin/helm"
         //  } else {
         //    currentBuild.result = 'ABORTED'
         //    throw new Exception("Aborting as this is not a PR job")
         // }
      }
      stage ("Checkout and Package Charts") {

            // Checkout PR Code
            def scmVars = checkout scm
            branchName = "${scmVars.GIT_BRANCH}"
            packageName = currentBuild.displayName
            prNumber = "${env.BRANCH_NAME}".split("-")[1]
            chartVersion = "${prNumber}.${env.BUILD_NUMBER}"
            
            // Perform Chart packaging
            sh "helm dependency update ./charts/pega/"
            sh "helm dependency update ./charts/addons/"
            sh "helm dependency update ./charts/backingservices/"
            sh "helm package --version ${chartVersion} ./charts/pega/"
            sh "helm package --version ${chartVersion} ./charts/addons/"
            sh "helm package --version ${chartVersion} ./charts/backingservices/"
            
            // Publish helm charts to test-automation GitHub Pages
            withCredentials([usernamePassword(credentialsId: "helmautomation",
              passwordVariable: 'AUTOMATION_APIKEY', usernameVariable: 'AUTOMATION_USERNAME')]) {

                sh "git clone https://pegaautomationuser:${AUTOMATION_APIKEY}@github.com/pegaautomationuser/helmcharts.git --branch=gh-pages gh-pages"
                sh "mv pega-${chartVersion}.tgz gh-pages/"
                sh "mv addons-${chartVersion}.tgz gh-pages/"
                sh "mv backingservices-${chartVersion}.tgz gh-pages/"
                dir("gh-pages") {
                  sh "helm repo index --merge index.yaml --url https://pegaautomationuser.github.io/helmcharts/ ."
                  sh "cat index.yaml"    
                  sh "git config user.email pegaautomationuser@gmail.com"
                  sh "git config user.name ${AUTOMATION_USERNAME}"
                  sh "git add ."
                  sh "git commit -m \"Jenkins build to publish test artefacts of version ${chartVersion}\""
                  sh "git push -u origin gh-pages"                                  
                }

            } 
      }

      stage("Setup Cluster and Execute Tests") {
          
          jobMap = [:]
          jobMap["job"] = "../kubernetes-test-orchestrator/master"
          jobMap["parameters"] = [
                                  string(name: 'PROVIDERS', value: labels),
                                  string(name: 'WEB_READY_IMAGE_NAME', value: ""),
                                  string(name: 'HELM_CHART_VERSION', value: chartVersion),
                              ]
          jobMap["propagate"] = true
          jobMap["quietPeriod"] = 0 
          resultWrapper = build jobMap
          currentBuild.result = resultWrapper.result
      } 
  }

def validateProviderLabel(String provider){
    def validProviders = ["integ-all","integ-eks","integ-gke","integ-aks"]
    def failureMessage = "Invalid provider label - ${provider}. valid labels are ${validProviders}"
    if(!validProviders.contains(provider)){
        currentBuild.result = 'FAILURE'
        pullRequest.comment("${failureMessage}")
        throw new Exception("${failureMessage}")
    }
}

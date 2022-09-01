#!/usr/bin/env groovy
def labels = ""

node {
      stage("Init"){
          if (env.CHANGE_ID) {
            pullRequest.labels.each{
            echo "label: $it"
            validateLabels(it)
            labels += "$it,"
            }
            labels = labels.substring(0,labels.length()-1)
            echo "PR labels -> $labels"
            // Comment on PR about the text execution
            pullRequest.comment("Starting ${env.BRANCH_NAME} validation on -> $labels")
            sh "curl -fsSL -o helm-v3.2.4-linux-amd64.tar.gz https://get.helm.sh/helm-v3.2.4-linux-amd64.tar.gz"
            sh "tar -zxvf helm-v3.2.4-linux-amd64.tar.gz"
            sh "mv linux-amd64/helm /usr/local/bin/helm"
          } else {
            currentBuild.result = 'ABORTED'
            throw new Exception("Aborting as this is not a PR job")
         }
      }
      stage ("Checkout and Package Charts and ConfigFiles") {

            // Checkout PR Code
            def scmVars = checkout scm
            branchName = "${scmVars.GIT_BRANCH}"
            packageName = currentBuild.displayName
            prNumber = "${env.BRANCH_NAME}".split("-")[1]
            chartVersion = "${prNumber}.${env.BUILD_NUMBER}"
            deployConfigsFileName = "deploy-config-${chartVersion}.tgz"
            installerConfigsFileName = "installer-config-${chartVersion}.tgz"
            // Perform Chart packaging
            sh "helm dependency update ./charts/pega/"
            sh "helm dependency update ./charts/addons/"
            sh "helm dependency update ./charts/backingservices/"
            sh "helm package --version ${chartVersion} ./charts/pega/"
            sh "helm package --version ${chartVersion} ./charts/addons/"
            sh "helm package --version ${chartVersion} ./charts/backingservices/"
            sh "tar -czvf ${deployConfigsFileName} --directory=./charts/pega/config deploy/context.xml.tmpl deploy/server.xml.tmpl deploy/prconfig.xml deploy/prlog4j2.xml"
            sh "mkdir -p ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/migrateSystem.properties.tmpl ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/prbootstrap.properties.tmpl ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/prconfig.xml.tmpl ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/prlog4j2.xml ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/prpcUtils.properties.tmpl ./charts/pega/charts/installer/config/installer && cp ./charts/pega/charts/installer/config/setupDatabase.properties.tmpl ./charts/pega/charts/installer/config/installer && tar -czvf ${installerConfigsFileName} --directory=./charts/pega/charts/installer/config installer/migrateSystem.properties.tmpl installer/prbootstrap.properties.tmpl installer/prconfig.xml.tmpl installer/prlog4j2.xml installer/prpcUtils.properties.tmpl installer/setupDatabase.properties.tmpl"

            // Publish helm charts to test-automation GitHub Pages
            withCredentials([usernamePassword(credentialsId: "helmautomation",
              passwordVariable: 'AUTOMATION_APIKEY', usernameVariable: 'AUTOMATION_USERNAME')]) {

                sh "git clone https://${AUTOMATION_USERNAME}:${AUTOMATION_APIKEY}@github.com/pegaautomationuser/helmcharts.git --branch=gh-pages gh-pages"
                sh "mv pega-${chartVersion}.tgz gh-pages/"
                sh "mv addons-${chartVersion}.tgz gh-pages/"
                sh "mv backingservices-${chartVersion}.tgz gh-pages/"
                sh "mv ${deployConfigsFileName} gh-pages/"
                sh "mv ${installerConfigsFileName} gh-pages/"
                dir("gh-pages") {
                  sh "helm repo index --merge index.yaml --url https://pegaautomationuser.github.io/helmcharts/ ."   
                  sh "git config user.email pegaautomationuser@gmail.com"
                  sh "git config user.name ${AUTOMATION_USERNAME}"
                  sh "git add ."
                  sh "git commit -m \"Jenkins build to publish test artefacts of version ${chartVersion}\""
                  sh "git push -u origin gh-pages --force"                                  
                }

            } 
      }

      stage("Setup Cluster and Execute Tests") {
          prLabels = labels.toString().split(",")

          if(prLabels.contains("integ-all") || prLabels.contains("integ-eks") || prLabels.contains("integ-gke") || prLabels.contains("integ-aks")) {
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
          } else {
              echo "Skipping 'Setup Cluster and Execute Tests' stage based on PR labels: $prLabels"
          }
      } 
  }

def validateLabels(String label){
    def validLabels = ["integ-all","integ-eks","integ-gke","integ-aks","configs"]
    def failureMessage = "Invalid label - ${label}. valid labels are ${validLabels}"
    if(!validLabels.contains(label)){
        currentBuild.result = 'FAILURE'
        pullRequest.comment("${failureMessage}")
        throw new Exception("${failureMessage}")
    }
}

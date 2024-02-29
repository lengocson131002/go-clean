def artifactId = ""
def dockerImage = ""
def dockerTag = "latest"
def now = (new Date()).format("yyMMddHHmmss", TimeZone.getTimeZone('Asia/Ho_Chi_Minh'))
def mainPathSvc = 'integration/microservices/new-mcs'

pipeline {
    agent { label "server141" }
    environment {
        DOCKER_REGISTRY = 'DEV-ESB-Log:5000'
        // K8S_NAMESPACE = 'ocb-new-mcs'
        K8S_CONTEXT = 'kubernetes-admin@kubernetes'
        // JENKINS_MCS_CARD_NGHIA_SONAR_TOKEN = credentials('nghiant5-jenkins-pipeline-jenkins')
        // JENKINS_SONAR_HOST = credentials('jenkins-sonar-host')
        JENKINS_API_TOKEN = credentials('JENKINS_API_TOKEN')
    }

    stages {

        stage('Checkout and build') {
           steps {
               script {
                    artifactId =  "golang-clean-architecture"
                    dockerImage = DOCKER_REGISTRY + "/nghiant5/mcs-card-test/" + artifactId
                    dockerTag =  now
               }
           }
        }

        stage("Build & Push Docker Image") {
           steps {
                script {
                    sh "echo '\n>>>> Build docker image'"
                    sh "docker build -t ${dockerImage}:${dockerTag} ."

                    sh "echo '\n>>>> Push to private registry'"
                    sh "docker push ${dockerImage}:${dockerTag}"

                }
           }
       }
    }    
}
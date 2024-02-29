def artifactId = ""
def dockerImage = ""
def dockerTag = "latest"
def now = (new Date()).format("yyMMddHHmmss", TimeZone.getTimeZone('Asia/Ho_Chi_Minh'))
def mainPathSvc = 'integration/microservices/new-mcs'

pipeline {
    agent { label "server141" }
    environment {
        DOCKER_REGISTRY = 'DEV-ESB-Log:5000'
        K8S_CONTEXT = 'kubernetes-admin@kubernetes'

    }

    stages {

        stage('Checkout and build') {
           steps {
               script {
                    sh "ls"

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
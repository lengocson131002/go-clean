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
        stage("Cleanup Workspace") {
            steps {
                cleanWs()
            }
        }

        stage("Checkout from SCM") {
               steps {
                   git branch: 'refactor-t24', credentialsId: 'github-nghia', url: 'https://github.com/lengocson131002/go-clean.git'
               }

        }

        stage('Checkout and build') {
           steps {
               script {

                    artifactId =  golang-clean-architecture
                    dockerImage = DOCKER_REGISTRY + "/nghiant5/mcs-card-test/" + artifactId
                    dockerTag = "${GIT_COMMIT}".substring(0, 8) + '_' + now
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
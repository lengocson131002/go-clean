def artifactId = ""
def dockerImage = ""
def dockerTag = "latest"
def now = (new Date()).format("yyMMddHHmmss", TimeZone.getTimeZone('Asia/Ho_Chi_Minh'))
def mainPathSvc = 'integration/microservices/new-mcs'

pipeline {
    agent { label "server141" }

    tools{
        go 'go-1.21.6'
    }

    environment {
        DOCKER_REGISTRY = 'DEV-ESB-Log:5000'
        K8S_CONTEXT = 'kubernetes-admin@kubernetes'
        JENKINS_API_TOKEN = credentials('JENKINS_API_TOKEN')

    }


    stages {

        // stage('SonarQube Analysis') {
        //   steps {
        //       script {
        //           sh "mvn clean verify sonar:sonar \
        //                   -Dsonar.projectKey=N3zG0iYdVhAqTERQDXZZSDF7jK1vL5p99 \
        //                   -Dsonar.projectName='mcs-card-nghiant-test' \
        //                   -Dsonar.host.url=$JENKINS_SONAR_HOST \
        //                   -Dsonar.token=$JENKINS_MCS_CARD_NGHIA_SONAR_TOKEN"
        //       }
        //   }
        // }
        stage('Checkout and build') {
           steps {
               script {
                  // sh "go test ./..."

                    artifactId =  "golang-clean-architecture"
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
        stage("Cleanup"){
            steps {
                script{
                    sh "echo '\n>>>> Remove docker image'"
                    sh "docker rmi ${dockerImage}:${dockerTag}"
                }
            }
        }
        //stage("Trivy Scan") {
//            steps {
//                script {
//                    sh " trivy --skip-db-update --severity HIGH,MEDIUM,CRITICAL  image 10.96.24.141:5001/nghiant5/mcs-card-test/mcs-card:${dockerTag}"
//                }
//            }
//     }
       stage("Trigger CD Pipeline") {
            steps {
                script {
                    sh "curl -v -k --user user:${JENKINS_API_TOKEN} -X POST -H 'cache-control: no-cache' -H 'content-type: application/x-www-form-urlencoded' --data 'dockerTag=${dockerTag}' 'http://10.96.24.141:8080/job/PipelineHuyPham/job/nghiant5-pipeline/job/gitops-go-clean-architecture/buildWithParameters?token=gitops-token'"
                }
            }
        }
    }    
}
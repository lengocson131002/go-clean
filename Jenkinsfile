pipeline {
    agent { label "Jenkins-Agent1" }
    environment {
            APP_NAME = "petcare-frontend"
    }

    stages {
        stage("Cleanup Workspace") {
            steps {
                cleanWs()
            }
        }

        stage("Checkout from SCM") {
               steps {
                   git branch: 'frontend-CD', credentialsId: 'github', url: 'https://github.com/nguyentrongnghia1702/gitops-hcmiu-petcare-app'
               }
        }

        stage("Update the Deployment Tags") {
            steps {
                dir("frontend-cd"){
                    sh """
                        cat deployment.yaml
                        sed -i 's/${APP_NAME}.*/${APP_NAME}:${IMAGE_TAG}/g' deployment.yaml
                        cat deployment.yaml
                    """
                }
                
            }
        }

        // stage("Push the changed deployment file to Git") {
        //     steps {
        //         sh """
        //            git config --global user.name "Jenkins"
        //            git config --global user.email "Jenkins@gmail.com"
        //            git add .
        //            git commit -m "Updated Deployment Manifest"
        //         """
        //         withCredentials([gitUsernamePassword(credentialsId: 'github', gitToolName: 'Default')]) {
        //           sh "git push https://github.com/nguyentrongnghia1702/gitops-hcmiu-petcare-app frontend-CD"
        //         }
        //     }
        // }
      
    }
}
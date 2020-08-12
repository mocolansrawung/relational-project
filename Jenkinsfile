podTemplate(containers: [
    containerTemplate(name: 'docker', image: 'docker:19.03.6', command: 'cat', ttyEnabled: true)
  ],
  volumes: [
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock')
  ]
  ) {
    node(POD_LABEL) {
            stage('Checkout') {
            checkout scm
            }

            stage('Build Docker Image') {
                container('docker') {
                        docker.withRegistry('https://107126629234.dkr.ecr.ap-southeast-1.amazonaws.com', 'ecr:ap-southeast-1:49feb1c9-1719-4520-aa17-67695b347b0e	') {
                            script {
                                if (env.BRANCH_NAME == 'dev'){
                                    sh 'docker build --network=host -f "Dockerfile" -t 107126629234.dkr.ecr.ap-southeast-1.amazonaws.com/boilerplate-go:1.1.$BUILD_NUMBER-dev .'
                                }
                                else if (env.BRANCH_NAME == 'master'){
                                    sh 'docker build --network=host -f "Dockerfile" -t 107126629234.dkr.ecr.ap-southeast-1.amazonaws.com/boilerplate-go:1.1.$BUILD_NUMBER-prod .'
                                }
                                else{
                                    sh 'docker build --network=host -f "Dockerfile" -t 107126629234.dkr.ecr.ap-southeast-1.amazonaws.com/boilerplate-go:1.1.$BUILD_NUMBER-$BRANCH_NAME .'
                                }
                            }
                        }
                }
            }

            stage('Push Docker Image') {
                container('docker') {
                        docker.withRegistry('https://107126629234.dkr.ecr.ap-southeast-1.amazonaws.com', 'ecr:ap-southeast-1:49feb1c9-1719-4520-aa17-67695b347b0e	') {
                            script {
                                if (env.BRANCH_NAME == 'dev'){
                                    sh 'docker push 107126629234.dkr.ecr.ap-southeast-1.amazonaws.com/boilerplate-go:1.1.$BUILD_NUMBER-dev'
                                }
                                else if (env.BRANCH_NAME == 'master'){
                                    sh 'docker push 107126629234.dkr.ecr.ap-southeast-1.amazonaws.com/boilerplate-go:1.1.$BUILD_NUMBER-prod'
                                }
                                else{
                                    sh 'docker push 107126629234.dkr.ecr.ap-southeast-1.amazonaws.com/boilerplate-go:1.1.$BUILD_NUMBER-$BRANCH_NAME'
                                }
                            }
                        }
                }
            }

            stage('Notification') {
                // telegramSend 'Job success for $JOB_NAME - Docker Image tag : boilerplate-go:1.1.$BUILD_NUMBER-$BRANCH_NAME'
                discordSend description: "Jenkins Pipeline Build Success!", footer: "boilerplate-go:1.1.$BUILD_NUMBER-$BRANCH_NAME", result: currentBuild.currentResult, title: "$JOB_NAME", webhookURL: "$DISCORD_WEBHOOK"
            }
    }
}
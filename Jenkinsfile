podTemplate(
    nodeSelector: 'evermos.com/serviceClass=t3a-large-jenkins',
    containers: [
    containerTemplate(name: 'docker', image: 'docker:19.03.6', command: 'cat', ttyEnabled: true),
    containerTemplate(
        name: 'sonarqube',
        image: 'cloudbees/java-build-tools:2.5.1',
        command: 'cat',
        ttyEnabled: true
    ),
  ],
  volumes: [
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock')
  ]
  ) {
    node(POD_LABEL) {
        def appName = "boilerplate-go"
        def appFullName
        def revision
        def message
        def repoURL
        
        stage('Checkout') {
            def scmVars = checkout([
                $class: 'GitSCM',
                branches: scm.branches,
                extensions: scm.extensions + [
                    [
                        $class: 'AuthorInChangelog'
                    ],
                    [
                        $class: 'ChangelogToBranch',
                        options: [
                            compareRemote: 'origin',
                            compareTarget: 'master'
                        ]
                    ]
                ],
                userRemoteConfigs: scm.userRemoteConfigs
                ])
            appFullName = "${appName}:${scmVars.GIT_COMMIT}"
            revision = "${scmVars.GIT_COMMIT}"
            repoURL = "${scmVars.GIT_URL}"
            echo repoURL
            message = sh(returnStdout: true, script: "git log --oneline -1 ${revision}")
        }

        stage('SonarQube analysis') {
            container('sonarqube'){
                def scannerHome = tool 'SonarQube';
                withSonarQubeEnv('sonarqube') {
                    script {
                        if(env.BRANCH_NAME =~ 'PR-.*'){
                            sh "echo sonar.pullrequest.key=${env.CHANGE_ID} >> sonar-project.properties"
                            sh "echo sonar.pullrequest.base=${env.CHANGE_TARGET} >> sonar-project.properties"
                            sh "echo sonar.pullrequest.branch=${env.CHANGE_BRANCH} >> sonar-project.properties"
                            sh "echo sonar.scm.revision=${revision} >> sonar-project.properties"
                            sh "${scannerHome}/bin/sonar-scanner"
                        }
                        else{
                            sh "echo sonar.branch.name=${env.BRANCH_NAME} >> sonar-project.properties"
                            sh "${scannerHome}/bin/sonar-scanner"
                        }
                    }
                }
            }
        }

        stage("Quality Gate"){
            timeout(time: 10, unit: 'MINUTES') { // Just in case something goes wrong, pipeline will be killed after a timeout
                def qg = waitForQualityGate() // Reuse taskId previously collected by withSonarQubeEnv
                if (qg.status != 'OK') {
                    error "Pipeline aborted due to quality gate failure: ${qg.status}"
                }
            }
        }

        // Check version compatibility on PR branches only.
        if(env.BRANCH_NAME =~ 'PR-.*') {
            stage('Check Go Version Compatibility') {
                def compatibleGoVersion = [ '1.14', '1.15' ]
                def buildJob = [:]
                for (ver in compatibleGoVersion) {
                    def goVer = ver.trim()
                    buildJob["Run build with go version ${goVer}"] = {
                        stage("Test build compatibility for Go version ${goVer}"){
                            container('docker') {
                                script {
                                    sh "docker build --build-arg GO_VERSION=${goVer} --network=host -t be-boilerplate-go:test-${goVer} ."
                                }
                            }
                        }
                    }
                }
                parallel buildJob
            }
        }

        // Build and push the image and notify via Discord only on PR merge to master.
        if (env.BRANCH_NAME == 'master') {
            stage('Build Docker Image') {
                container('docker') {
                    docker.withRegistry('https://107126629234.dkr.ecr.ap-southeast-1.amazonaws.com', 'ecr:ap-southeast-1:49feb1c9-1719-4520-aa17-67695b347b0e') {
                        script {
                            sh """docker build --network=host -f "Dockerfile" -t 107126629234.dkr.ecr.ap-southeast-1.amazonaws.com/${appFullName} ."""
                        }
                    }
                }
            }

            stage('Push Docker Image') {

                container('docker') {
                    docker.withRegistry('https://107126629234.dkr.ecr.ap-southeast-1.amazonaws.com', 'ecr:ap-southeast-1:49feb1c9-1719-4520-aa17-67695b347b0e	') {
                        script {
                            sh """docker push 107126629234.dkr.ecr.ap-southeast-1.amazonaws.com/${appFullName}"""
                        }
                    }
                }
            }

            stage('Notification') {
                notify("slack","${currentBuild.currentResult}","${message}","${appFullName}","#jenkins-build")
            }
        }
    }
}

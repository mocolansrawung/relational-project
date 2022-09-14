podTemplate(
    yaml: readTrusted('pod.yaml')
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
                            compareTarget: 'main'
                        ]
                    ]
                ],
                userRemoteConfigs: scm.userRemoteConfigs
                ])
            appFullName = "${appName}:${getShortCommit()}-${getTimeStamp()}"
            revision = "${scmVars.GIT_COMMIT}"
            repoURL = "${scmVars.GIT_URL}"
            echo repoURL
            message = sh(returnStdout: true, script: "git log --oneline -1 ${revision}")
        }

        stage('SonarQube analysis') {
            catchError(buildResult: 'SUCCESS', stageResult: 'UNSTABLE') {
                container('sonarqube'){
                    def scannerHome = tool 'SonarQube';
                    withSonarQubeEnv('sonarqube') {
                        script {
                            if(env.BRANCH_NAME =~ 'PR-.*'){
                                sh "${scannerHome}/bin/sonar-scanner \
                                -Dsonar.projectKey=${appName} \
                                -Dsonar.projectName=${appName} \
                                -Dsonar.pullrequest.key=${env.CHANGE_ID} \
                                -Dsonar.pullrequest.base=${env.CHANGE_TARGET} \
                                -Dsonar.pullrequest.branch=${env.CHANGE_BRANCH} \
                                -Dsonar.sources=. \
                                -Dsonar.projectVersion=1.0 \
                                -Dsonar.host.url=http://sonarqube-sonarqube.evm-internal:9000"
                            }
                            else{
                                sh "${scannerHome}/bin/sonar-scanner -D sonar.branch.name='${env.BRANCH_NAME}'"
                            }
                        }
                    }
                }
            }
        }

        stage("Quality Gate"){
            try {
                timeout(time: 5, unit: 'MINUTES') { // Just in case something goes wrong, pipeline will be killed after a timeout
                    def qg = waitForQualityGate() // Reuse taskId previously collected by withSonarQubeEnv
                    if (qg.status != 'OK') {
                        error "Pipeline aborted due to quality gate failure: ${qg.status}"
                    }
                }
            } catch (err) {
                echo "Quality gate failed: ${err}"
                currentBuild.result = 'SUCCESS'
            }
        }

        // Check version compatibility on PR branches only.
        if(env.BRANCH_NAME =~ 'PR-.*') {
            parallel (
                "1.14": {
                    stage ("compability Go 1.14"){
                        container("kaniko") {
                                script {
                                    sh "/kaniko/executor --build-arg GO_VERSION=1.14 --context `pwd` --no-push --destination boilerplate:test-1.14"                                
                                    }
                            }
                    }
                },
                "1.15": {
                    stage ("compability Go 1.15") {
                            container("kaniko2") {
                                    script {
                                        sh "/kaniko/executor --build-arg GO_VERSION=1.15 --context `pwd` --no-push --destination boilerplate:test-1.15"                                
                                        }
                                }
                        }
                    }
                
            )
        }

        // Build and push the image and notify via Discord only on PR merge to master.
        if (env.BRANCH_NAME == 'main') {
            stage('Build Docker Image') {
                container('kaniko') {
                    script {
                        sh """
                        /kaniko/executor --context `pwd` --destination 107126629234.dkr.ecr.ap-southeast-1.amazonaws.com/${appFullName}
                        """
                    }
                }
            }

            stage('Notification') {
                notify("slack","${currentBuild.currentResult}","${message}","${appFullName}","#jenkins-build")
            }
        }
    }
}

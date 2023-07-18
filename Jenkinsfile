pipeline{
    agent any
    environment{
        SERVICE_NAME = "domain.symper.vn"
        HOST_NAME = "domain-forward.symper.vn"
        KAFKA_SUBCRIBE = false
        APP_NAME=sh (script: "echo $SERVICE_NAME | cut -d'.' -f1", returnStdout: true).trim()
        Author_Name=sh(script: "git show -s --pretty=%ae", returnStdout: true).trim()
        BRANCH_NAME = "${GIT_BRANCH.split('/').size() > 1 ? GIT_BRANCH.split('/')[1..-1].join('/') : GIT_BRANCH}"
        COMMIT_MESSAGE = sh(script: "git log --format=%B -n 1", returnStdout: true).trim()
        MSTEAMS_WEBHOOK=credentials('ms_teams_webhook')
    }
    stages{
        stage ("quality control") {
            when {
                branch "dev"    
            }
            environment {
                SERVICE_ENV = "test"
                SSH_HOST = "10.20.166.15"
                POSTGRES_HOST = "10.20.166.52"
                POSTGRES_DB = "file_service_symper_vn"
                CLICKHOUSE_HOST = "10.20.166.52"
                CACHE_HOST= "redis-server.redis-server.svc.cluster.local"
                KAFKA_PREFIX = "10.20.166.6"
            }
            stages {
                stage("build"){
                    steps{
                        withCredentials([usernamePassword(credentialsId: 'docker-hub', passwordVariable: 'DOCKER_REGISTRY_PWD', usernameVariable: 'DOCKER_REGISTRY_USER')]) {
                            sh 'echo $DOCKER_REGISTRY_PWD | docker login -u $DOCKER_REGISTRY_USER --password-stdin localhost:5000'
                        }
                        script {
                            env.BUILD_VERSION = 'latest'
                            sh "docker build -t localhost:5000/${SERVICE_NAME}:${env.BUILD_VERSION} ."
                            sh "docker push localhost:5000/${SERVICE_NAME}:${env.BUILD_VERSION}"
                            sh "docker image rm localhost:5000/${SERVICE_NAME}:${env.BUILD_VERSION}"
                        }
                    }
                }
                stage("checkin deployment state") {
                    steps{
                        withCredentials([
                            usernamePassword(credentialsId: 'ssh_qc_vps', passwordVariable: 'USER_PASS', usernameVariable: 'USER_NAME')
                        ]) {
                            script {
                                sshagent(['ssh_qc_key']) {
                                    try{
                                        env.CURRENT_ROLE=sh (returnStdout: true, 
                                                            script: "ssh -o StrictHostKeyChecking=no $USER_NAME@$SSH_HOST 'echo -e \'$USER_PASS\' | sudo -S kubectl get services --field-selector metadata.name=\"$APP_NAME\" -o jsonpath={.items[0].spec.selector.role}'").trim()
                                    } catch (Exception e) {
                                        echo "$e" }
                                }
                                sh "echo role ${env.CURRENT_ROLE}"
                                if("$env.CURRENT_ROLE" == "" || "$env.CURRENT_ROLE" == "green") {
                                    env.CURRENT_ROLE = "green"
                                    env.TARGET_ROLE = "blue"
                                } else {
                                    env.TARGET_ROLE = "green"
                                }
                            }
                        }
                    }
                }
                stage("deploy to k8s") {
                    // when {
                    //     expression {
                    //         "${Author_Name}"?.startsWith("release")
                    //     }
                    // }
                    steps{
                        withCredentials([
                            usernamePassword(credentialsId: 'dev_database', passwordVariable: 'POSTGRES_PASS', usernameVariable: 'POSTGRES_USER'),
                            usernamePassword(credentialsId: 'ssh_qc_vps', passwordVariable: 'USER_PASS', usernameVariable: 'USER_NAME'),
                            usernamePassword(credentialsId: 'cache', passwordVariable: 'CACHE_PASSWORD', usernameVariable: 'CACHE_USER_NAME')
                        ]) {
                            sh "chmod +x shellscripts/*"
                            sshagent(['ssh_qc_key']) {
                                sh '''
                                    echo $TARGET_ROLE
                                '''
                                sh "./shellscripts/updateManifests.sh"
                                sh "./shellscripts/deployService.sh"
                            }
                        }
                    }
                }
            }
        }
        stage ("production") {
            when {
                branch "tags/*"
            }
            environment {
                SERVICE_ENV = "prod"
                POSTGRES_HOST = "10.20.166.193"
                SSH_HOST = "10.20.166.246"
                POSTGRES_DB = "file_service_symper_vn"
                CLICKHOUSE_HOST = "10.20.166.166"
                CACHE_HOST= "redis-server.redis-server.svc.cluster.local"
                KAFKA_PREFIX = "10.20.166.6"
            }
            stages {
                stage("build"){
                    steps{
                        withCredentials([usernamePassword(credentialsId: 'docker-hub', passwordVariable: 'DOCKER_REGISTRY_PWD', usernameVariable: 'DOCKER_REGISTRY_USER')]) {
                            sh 'echo $DOCKER_REGISTRY_PWD | docker login -u $DOCKER_REGISTRY_USER --password-stdin localhost:5000'
                        }
                        script {
                            latestTag = sh(returnStdout:  true, script: "git tag --sort=-creatordate | head -n 1").trim()
                            env.BUILD_VERSION = latestTag
                            sh "docker build -t localhost:5000/${SERVICE_NAME}:${env.BUILD_VERSION} ."
                            sh "docker push localhost:5000/${SERVICE_NAME}:${env.BUILD_VERSION}"
                            sh "docker image rm localhost:5000/${SERVICE_NAME}:${env.BUILD_VERSION}"
                        }
                    }
                }
                stage("checkin deployment state") {
                    steps{
                        withCredentials([
                            usernamePassword(credentialsId: 'ssh_prod_vps', passwordVariable: 'USER_PASS', usernameVariable: 'USER_NAME')
                        ]) {
                            script {
                                sshagent(['prod_ssh_key']) {
                                    try{
                                        env.CURRENT_ROLE=sh (returnStdout: true, 
                                                            script: "ssh -o StrictHostKeyChecking=no $USER_NAME@$SSH_HOST 'echo -e \'$USER_PASS\' | sudo -S kubectl get services --field-selector metadata.name=\"$APP_NAME\" -o jsonpath={.items[0].spec.selector.role}'").trim()
                                    } catch (Exception e) {
                                        echo "$e" }
                                }
                                sh "echo role ${env.CURRENT_ROLE}"
                                if("$env.CURRENT_ROLE" == "" || "$env.CURRENT_ROLE" == "green") {
                                    env.CURRENT_ROLE = "green"
                                    env.TARGET_ROLE = "blue"
                                } else {
                                    env.TARGET_ROLE = "green"
                                }
                            }
                        }
                    }
                }
                stage("deploy to k8s"){
                    // when {
                    //     expression {
                    //         "${Author_Name}"?.startsWith("release")
                    //     }
                    // }
                    steps{
                        withCredentials([
                            usernamePassword(credentialsId: 'sv_database', passwordVariable: 'POSTGRES_PASS', usernameVariable: 'POSTGRES_USER'),
                            usernamePassword(credentialsId: 'ssh_prod_vps', passwordVariable: 'USER_PASS', usernameVariable: 'USER_NAME'),
                            usernamePassword(credentialsId: 'cache', passwordVariable: 'CACHE_PASSWORD', usernameVariable: 'CACHE_USER_NAME')
                        ]) {
                            sh "chmod +x shellscripts/*"
                            sshagent(['prod_ssh_key']) {
                                sh '''
                                    echo $TARGET_ROLE
                                    echo $CURRENT_ROLE
                                '''
                                sh "./shellscripts/updateManifests.sh"
                                sh "./shellscripts/deployService.sh"
                            }
                        }
                    }
                }
            }
        }
    }
    post{
        always{
            echo "${env.SERVICE_NAME} - Build # ${env.BUILD_NUMBER} - ${env.BUILD_STATUS}: \nCheck console output at ${env.BUILD_URL} to view the results."
        }
        success{
            office365ConnectorSend webhookUrl: MSTEAMS_WEBHOOK,
                message: "${env.SERVICE_NAME} - Build # ${env.BUILD_NUMBER} - <span style=\"color:green\">${currentBuild.currentResult}</span>: <br> <strong style=\"color: deeppink;\">${env.COMMIT_MESSAGE}</strong>  <br> Check console output at ${env.BUILD_URL} to view the results.",
                status: 'Success'  
            echo "========pipeline executed successfully ========"
        }
        failure{
            office365ConnectorSend webhookUrl: MSTEAMS_WEBHOOK,
                message: "${env.SERVICE_NAME} - Build # ${env.BUILD_NUMBER} - <span style=\"color:red\">${currentBuild.currentResult}</span>: <br> <strong style=\"color: deeppink;\">${env.COMMIT_MESSAGE}</strong> <br> Check console output at ${env.BUILD_URL} to view the results.",
                status: 'Failed'   
            echo "========pipeline execution failed========"
        }
    }
}
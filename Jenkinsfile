//#!groovy

pipeline {
    // 指定集群节点
    agent any
    // 选项
    options {
        timestamps() //日志会有时间
        skipDefaultCheckout() //删除隐式checkout scm语句
        disableConcurrentBuilds() //禁止并行
        timeout(time: 1, unit: "HOURS") //流水线超市设置1h
        buildDiscarder(logRotator(numToKeepStr: '20')) //保留build数量
    }
    // 声明全局变量
    environment {
        harborUsername = "admin"
        harborPassword = "Harbor12345"
        harborAddress = "10.192.0.5:9002"
        harborRepo = "go-maxms"
    }
    // 流水线阶段
    stages {
        // 拉取代码
        stage("Checkout") {
            steps {
                echo "--------------------- Checkout Start ---------------------"
                timeout(time: 5, unit: "MINUTES"){
                    checkout([$class: "GitSCM", branches: [[name: '**']], extensions: [], userRemoteConfigs: [[url: "https://github.com/liuzhaomax/go-maxms.git"]]])
                }
                echo "--------------------- Checkout End ---------------------"
            }
        }
//         // 本地github会用到
//         stage("Update GitHub") {
//             steps {
//                 echo "--------------------- Update GitHub Start ---------------------"
//                 script {
//                     timeout(time: 20, unit: "MINUTES"){
//                         sh """
//                             git config --get remote.origin.url
//                             tr -d
//                             git rev-parse HEAD
//                         """
//                         // setting commit status
//                     }
//                 }
//                 echo "--------------------- Update GitHub End ---------------------"
//             }
//         }
        // 检查App版本
        stage("Version") {
            steps {
                echo "--------------------- Version Start ---------------------"
                echo "Branch: ${JOB_NAME}"
                script {
                    goHome = tool "go"
                    sh """
                        export GO_HOME=${goHome}
                        export PATH=\$GO_HOME/bin:\$PATH
                        rm -rf bin
                        ${goHome}/bin/go version
                    """
                }
                echo "--------------------- Version End ---------------------"
            }
        }
        // 语法格式检查
        stage("Lint") {
            steps {
                echo "--------------------- Lint Start ---------------------"
                script {
                    timeout(time: 15, unit: "MINUTES") {
                        goHome = tool "go"
                        sh """
                            export GO_HOME=${goHome}
                            export PATH=\$GO_HOME/bin:\$PATH
                            ${goHome}/bin/golangci-lint run -v --timeout 5m -c ./.golangci.yml ./...
                        """
                    }
                }
                echo "--------------------- Lint End ---------------------"
            }
        }
        // 构建
        stage("Build") {
            steps {
                echo "--------------------- Build Start ---------------------"
                script {
                    timeout(time: 15, unit: "MINUTES"){
                        goHome = tool "go" //变量名go在jenkins全局工具里定义的
                        sh """
                            export GO_HOME=${goHome}
                            export PATH=\$GO_HOME/bin:\$PATH
                            export ENV=dev
                            ${goHome}/bin/go build -o bin/main main/main.go
                        """
                    }
                }
                echo "--------------------- Build End ---------------------"
            }
        }
        // 静态代码分析SonarQube
        stage("SonarQube") {
            steps {
                echo "--------------------- SonarQube Start ---------------------"
                script {
                    timeout(time: 20, unit: "MINUTES"){
                        sonarScannerHome = "/var/jenkins_home/sonar-scanner"
                        String[] strArr = JOB_NAME.split("/")
                        String projectKey = strArr[0]
                        for (int i = 1; i < strArr.size(); i++) {
                            projectKey += "_" + strArr[i]
                        }
                        projectKey = projectKey.replaceAll("%2F", "_")
                        echo "SonarQube Project Key: ${projectKey}"
                        export PROJECT_KEY=${projectKey}
                        sh """
                            ${sonarScannerHome}/bin/sonar-scanner
                        """
                    }
                }
                echo "--------------------- SonarQube End ---------------------"
            }
        }
//         // 安全漏洞扫描Checkmarx
//         stage("Checkmarx") {
//             steps {
//                 echo "--------------------- Checkmarx Start ---------------------"
//                 echo "Checkmarx - SUCCESS"
//                 echo "--------------------- Checkmarx End ---------------------"
//             }
//         }
        // 构建镜像
        stage("Build Image") {
            steps {
                echo "--------------------- Build Image Start ---------------------"
                timeout(time: 10, unit: "MINUTES"){
                    sh """
                        docker build -t ${JOB_NAME}:${tag} .
                    """
                }
                echo "--------------------- Build Image End ---------------------"
            }
        }
        // 推送镜像到Harbor
        stage("Harbor") {
            steps {
                echo "--------------------- Push to Harbor Start ---------------------"
                timeout(time: 10, unit: "MINUTES"){
                    sh """
                        docker login -u ${harborUsername} -p ${harborPassword} ${harborAddress}
                        docker tag ${JOB_NAME}:${tag} ${harborAddress}/${harborRepo}/${JOB_NAME}:${tag}
                        docker push ${harborAddress}/${harborRepo}/${JOB_NAME}:${tag}
                    """
                }
                echo "--------------------- Push to Harbor End ---------------------"
            }
        }
        // 部署容器
        stage("Deploy") {
            steps {
                echo "--------------------- Deploy Start ---------------------"
                timeout(time: 10, unit: "MINUTES"){
                    sshPublisher(publishers: [sshPublisherDesc(configName: "test", transfers: [sshTransfer(cleanRemote: false, excludes: "", execCommand: "sudo deploy.sh $harborAddress $harborRepo $JOB_NAME $tag $container_port $host_port", execTimeout: 120000, flatten: false, makeEmptyDirs: false, noDefaultExcludes: false, patternSeparator: "[, ]+", remoteDirectory: "", remoteDirectorySDF: false, removePrefix: "", sourceFiles: "")], usePromotionTimestamp: false, useWorkspaceInPromotion: false, verbose: false)])
                }
                echo "--------------------- Deploy End ---------------------"
            }
        }
    }
    // 构建后的操作
    post {
        always {
            echo "********************************************************************"
            echo "********************************************************************"
            echo "****************** CI Pipeline about to Finish *********************"
            echo "********************************************************************"
            echo "********************************************************************"
        }

        success {
            script {
                echo "SUCCESS 成功"
                keepBuilds()
            }
            sh "docker image prune -f"
        }

        failure {
            script {
                echo "FAILURE 失败"
                keepBuilds()
            }
            error "错误发生，流水线失败"
        }

        aborted {
            echo "ABORTED 取消"
            error "流水线被终止"
        }
    }
}

// 保留最近5个构建，其中必须包含至少一次成功构建
def keepBuilds() {
    def buildsToKeep = []
    def successFound = false

    for (int i = currentBuild.number - 1; i >= 1 && buildsToKeep.size() < 5; i--) {
        def build = Jenkins.instance.getItemByFullName(env.JOB_NAME).getBuildByNumber(i)

        if (build.result == 'SUCCESS') {
            buildsToKeep << build
            successFound = true
        } else if (build.result == 'FAILURE' || build.result == 'ABORTED') {
            // 如果前一次构建成功，则保留这一次构建
            if (successFound) {
                buildsToKeep << build
            }
        }
    }

    // 设置保留的构建
    currentBuild.rawBuild.getAction(hudson.tasks.LogRotator.class).setBuildKeepDependencies(buildsToKeep)
}

// 组合job name为蛇形
def rewriteJobNameInSnake() {
    String[] strArr = JOB_NAME.split("/")
    String projectName = strArr[0..-1].join("_")
    return projectName
}
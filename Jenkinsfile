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
                            ${goHome}/bin/golangci-lint run -v --fast --timeout 5m /var/jenkins_home/workspace/${rewriteJobNameInSnake()}
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
                        sonarScannerHome = tool "sonar-scanner"
                        String[] strArr
                        strArr = JOB_NAME.split("/")
                        String projectKey = strArr[0]
                        for (int i = 1; i < strArr.size(); i++) {
                            projectKey = projectKey + "_" + strArr[i]
                        }
                        strArr = projectKey.split("%2F")
                        projectKey = strArr[0]
                        for (int i = 1; i < strArr.size(); i++) {
                            projectKey = projectKey + "_" + strArr[i]
                        }
                        sh """
                            ${sonarScannerHome}/bin/sonar-scanner \
                                -Dsonar.sources=./ \
                                -Dsonar.projectname=${JOB_NAME} \
                                -Dsonar.login=5cbe5f7092c9a2b8168d610c8efee1dfe938a6ad \
                                -Dsonar.projectKey=${projectKey} \
                                -Dsonar.nodejs.executable=/usr/bin/go \
                                -Dsonar.inclusions=src/**/*.go \
                                -Dsonar.coverage.exclusions=internal/**/*,env/**/*,specs/**/*,src/pb/**/* \
                                -Dsonar.qualitygate.wait=true
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
//#!groovy

pipeline {
    // 指定集群节点
    agent any
    // 选项
    options {
        timestamps() //日志会有时间
        skipDefaultCheckout() //删除隐式checkout scm语句
        disableConcurrentBuilds() //禁止并行
        timeout(time: 1, unit: "HOURS") //流水线超时设置1h
        buildDiscarder(logRotator(numToKeepStr: "20")) //保留build数量
    }
    // 声明全局变量
    environment {
        ENV = "dev" // 根据Config Selection步骤的input而定
        TAG = "" // 根据Config Selection步骤的input而定
        ProjectKey = "" // 根据SonarQube步骤的input而定
        harborUsername = "admin"
        harborPassword = "Harbor12345"
        harborAddress = "172.16.96.97:9002"
        harborRepo = "go-maxms"
        Container_port = "9999" // 启用随机端口，会被赋值
        Host_port = "9999" // 启用随机端口，会被赋值
        JobName = "go-maxms/main"
        DeploymentServerIP = "172.16.96.98"
    }
    // 流水线阶段
    stages {
        // 拉取代码
        stage("Checkout") {
            steps {
                echo "--------------------- Checkout Start ---------------------"
                timeout(time: 3, unit: "MINUTES"){
                    checkout([$class: "GitSCM", branches: [[name: "**"]], extensions: [], userRemoteConfigs: [[url: "https://github.com/liuzhaomax/go-maxms.git"]]])
                }
                echo "--------------------- Checkout End ---------------------"
            }
        }
//         // 本地github会用到
//         stage("Update GitHub") {
//             steps {
//                 echo "--------------------- Update GitHub Start ---------------------"
//                 script {
//                     timeout(time: 5, unit: "MINUTES"){
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
        stage("Config Selection") {
            steps {
                script {
                    timeout(time: 2, unit: "MINUTES"){
                        echo "Getting GitHub Tags..."
                        def tags = getGitHubTags()
                        def userInput
                        try {
                            timeout(time: 1, unit: "MINUTES") {
                                userInput = input(
                                    id: "userInput",
                                    message: "Please select environment and tag:",
                                    parameters: [
                                        [$class: "ChoiceParameterDefinition", name: "ENVIRONMENT", choices: "st\nsit\npnv\nqa\nprod", description: "Select environment"],
                                        [$class: "ChoiceParameterDefinition", name: "TAG", choices: tags, description: "Select tag"]
                                    ]
                                )
                            }
                        }
                        catch (e) {
                            echo "Using default env and tag values due to no operation in 1 min, or Exception caught: ${e}"
                        }
                        // 如果用户没有选择，使用默认值
                        ENV = userInput ? userInput.ENVIRONMENT : "st"
                        echo "Selected Environment: $ENV"
                        TAG = userInput ? userInput.TAG : tags.first()
                        echo "Selected Tag: $TAG"
                    }
                }
            }
        }
        // 检查App版本
        stage("Version") {
            steps {
                echo "--------------------- Version Start ---------------------"
                echo "Branch: ${JOB_NAME}"
                echo "App Tag: ${TAG}"
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
                    timeout(time: 5, unit: "MINUTES") {
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
                    timeout(time: 5, unit: "MINUTES"){
                        goHome = tool "go" //变量名go在jenkins全局工具里定义的
                        sh """
                            export GO_HOME=${goHome}
                            export PATH=\$GO_HOME/bin:\$PATH
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
                    timeout(time: 5, unit: "MINUTES"){
                        ProjectKey = genSonarProjectKey()
                        echo "SonarQube Project Key: ${ProjectKey}"
                        sonarScannerHome = tool "sonar-scanner"
                        sh """
                            export PROJECT_KEY=${ProjectKey}
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
            when {
                expression { return JOB_NAME == JobName }
            }
            steps {
                echo "--------------------- Build Image Start ---------------------"
                script {
                    timeout(time: 5, unit: "MINUTES"){
                        goHome = tool "go"
                        // 随机端口部署，将查询到的随机空闲端口写入配置文件
                        randomPort = sh(script: "${goHome}/bin/go run ./script/get_random_idle_port/main.go -e ${ENV}", returnStdout: true).trim()
                        Container_port = randomPort
                        Host_port = randomPort
                        sh """
                            docker build -t ${ProjectKey}:${TAG} .
                        """
                    }
                }
                echo "--------------------- Build Image End ---------------------"
            }
        }
        // 推送镜像到Harbor
        stage("Harbor") {
            when {
                expression { return JOB_NAME == JobName }
            }
            steps {
                echo "--------------------- Push to Harbor Start ---------------------"
                timeout(time: 5, unit: "MINUTES"){
                    sh """
                        docker login -u ${harborUsername} -p ${harborPassword} ${harborAddress}
                        docker tag ${ProjectKey}:${TAG} ${harborAddress}/${harborRepo}/${ProjectKey}:${TAG}
                        docker push ${harborAddress}/${harborRepo}/${ProjectKey}:${TAG}
                    """
                }
                echo "--------------------- Push to Harbor End ---------------------"
            }
        }
        // 部署容器
        stage("Deploy") {
            when {
                expression { return JOB_NAME == JobName }
            }
            steps {
                echo "--------------------- Deploy Start ---------------------"
                script {
                    timeout(time: 2, unit: "MINUTES") {
                        echo "ENV: ${ENV}"
                        echo "Port: ${Host_port}"
                        sh """
                            chmod +x ./deploy.sh
                            ./deploy.sh $harborAddress $harborRepo $ProjectKey $TAG $Container_port $Host_port $ENV $DeploymentServerIP
                        """
                    }
                }
                echo "--------------------- Deploy End ---------------------"
            }
        }
    }
    // 构建后的操作
    post {
        always {
            echo """
            ********************************************************************
            ********************************************************************
            ******************* Pipeline about to Finish ***********************
            ********************************************************************
            ********************************************************************
            """
        }
        cleanup {
            sh "docker image prune -f"
            sh "rm -rf ${WORKSPACE}/*"
        }
    }
}

// 保留最近5个构建，其中必须包含至少一次成功构建
def keepBuilds() {
    def buildsToKeep = []
    def successFound = false

    for (int i = currentBuild.number - 1; i >= 1 && buildsToKeep.size() < 5; i--) {
        def build = Jenkins.instance.getItemByFullName(env.JOB_NAME).getBuildByNumber(i)

        if (build.result == "SUCCESS") {
            buildsToKeep << build
            successFound = true
        } else if (build.result == "FAILURE" || build.result == "ABORTED") {
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

// 生成sonar的project key
def genSonarProjectKey() {
    String[] strArr = JOB_NAME.split("/")
    String pk = strArr[0]
    for (int i = 1; i < strArr.size(); i++) {
        pk += "_" + strArr[i]
    }
    pk = pk.replaceAll("%2F", "_")
    return pk.toLowerCase()
}

// 获取github tags
def getGitHubTags() {
    // 获取GitHub仓库的标签列表
    def tagsCommand = 'git ls-remote --tags origin'
    def tagsOutput = sh(script: tagsCommand, returnStdout: true).trim()
    // 处理输出，提取标签的名称
    def tagList = tagsOutput.readLines().collect { it.replaceAll(/.*refs\/tags\/(.*)(\^\{\})?/, '$1') }
    // 在 tagList 的 0 号索引位置，添加一个快照标签
    def timestamp = new Date().format("yyyyMMddHHmmss")
    tagList.add(0, "SNAPSHOT-$timestamp")
    return tagList
}

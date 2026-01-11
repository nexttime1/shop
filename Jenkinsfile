pipeline {
    agent any

    stages {
        // 第一个阶段：拉取代码
        stage('pull code') {
            steps {
                git branch: 'main', credentialsId: 'github-ed25519', url: 'git@github.com:nexttime1/GOMall_CloudWeGo.git'
            }
        }

        // 第二个阶段：代码构建 → 名称：构建项目
        stage('build project') {
            steps {
                sh '''
                    echo '开始构建'
                    echo '构建成功'
                '''
            }
        }
    }
}
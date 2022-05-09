pipeline{
    agen any
    stages{
        stage("verify tooling"){
            steps{
                sh '''
                    docker version
                    docker info
                    docker-compose version
                    curl --version
                    git --version
                '''
            }
        }
    }
}
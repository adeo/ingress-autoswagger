node('dockerhost') {
    env.TAG = '3.2'
    env.DOCKER_IMAGE = 'docker-devops.art.lmru.tech/bricks/ingress-autoswagger'
    env.DOCKER_REGISTRY_CREDS = 'lm-sa-devops'

    timestamps {
        ansiColor('xterm') {
            stage('Checkout') {
                checkout scm
            }

            stage('Build & Push Image') {
                if (env.CHANGE_ID) {
                    lint()
                } else {
                    image_build_and_push()
                }

            }

            stage('Wipe') {
                cleanWs()
            }
        }
    }
}

def lint() {
    // not needed here
}

def image_build_and_push() {
    def image = docker.build("${env.DOCKER_IMAGE}:${env.TAG}", ".")
    try {
        docker.withRegistry("https://$DOCKER_IMAGE", "$DOCKER_REGISTRY_CREDS") {
            image.push('$TAG')
        }
    }
    finally {
        sh "docker rmi $DOCKER_IMAGE:$TAG"
    }
}

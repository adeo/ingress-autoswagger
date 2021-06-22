node('dockerhost') {
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
                    helm_push()
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

def helm_push() {
    def currentPath = pwd()

    withCredentials([usernamePassword(credentialsId: 'helm-bricks-local-repo', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
        sh("docker run -v $currentPath:/pwd docker.art.lmru.tech/img-k8s-deployer:latest \
                /bin/bash -c 'helm repo add --username $USERNAME --password $PASSWORD bricks https://art.lmru.tech/helm-local-bricks; \
                helm push-artifactory pwd/helm bricks'")

    }
}

def image_build_and_push() {
    def image = docker.build("${env.DOCKER_IMAGE}:${env.GIT_TAG}", ".")
    try {
        docker.withRegistry("https://$DOCKER_IMAGE", "$DOCKER_REGISTRY_CREDS") {
            image.push('$GIT_TAG')
        }
    }
    finally {
        sh "docker rmi $DOCKER_IMAGE:$GIT_TAG"
    }
}

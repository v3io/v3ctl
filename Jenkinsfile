label = "${UUID.randomUUID().toString()}"
git_project = "v3ctl"
git_project_user = "v3io"
git_project_upstream_user = "v3io"
git_deploy_user = "iguazio-prod-git-user"
git_deploy_user_token = "iguazio-prod-git-user-token"
git_deploy_user_private_key = "iguazio-prod-git-user-private-key"

podTemplate(label: "${git_project}-${label}", inheritFrom: "jnlp-docker-golang1.14.2") {
    node("${git_project}-${label}") {
        pipelinex = library(identifier: 'pipelinex@development', retriever: modernSCM(
                [$class       : 'GitSCMSource',
                 credentialsId: git_deploy_user_private_key,
                 remote       : "git@github.com:iguazio/pipelinex.git"])).com.iguazio.pipelinex
        common.notify_slack {
            withCredentials([
                    string(credentialsId: git_deploy_user_token, variable: 'GIT_TOKEN')
            ]) {
                github.release(git_deploy_user, git_project, git_project_user, git_project_upstream_user, true, GIT_TOKEN) {
                    stage("get release") {
                        container('jnlp') {
                            RELEASE_ID = github.get_release_id(git_project, git_project_user, "${github.TAG_VERSION}", GIT_TOKEN)
                        }
                    }
                    stage('get dependencies') {
                        container('golang') {
                            dir("${github.BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                common.shellc("make get-dependencies")
                            }
                        }
                    }
                    parallel(
                        'build linux binaries': {
                            container('golang') {
                                stage('build linux binaries') {
                                    dir("${github.BUILD_FOLDER}/src/github.com/${git_project_upstream_user}/${git_project}") {
                                        common.shellc("V3CTL_SRC_PATH=\$(pwd) V3CTL_BIN_PATH=\$(pwd) V3CTL_TAG=${github.TAG_VERSION} GOARCH=amd64 GOOS=linux make v3ctl-bin")
                                    }
                                }
                            }
                        },
                    )
                    parallel(
                        'upload linux binaries artifactory': {
                            container('jnlp') {
                                stage('upload linux binaries artifactory') {
                                    withCredentials([
                                            string(credentialsId: pipelinex.PackagesRepo.ARTIFACTORY_IGUAZIO[2], variable: 'PACKAGES_ARTIFACTORY_PASSWORD')
                                    ]) {
                                        common.upload_file_to_artifactory(pipelinex.PackagesRepo.ARTIFACTORY_IGUAZIO[0], pipelinex.PackagesRepo.ARTIFACTORY_IGUAZIO[1], PACKAGES_ARTIFACTORY_PASSWORD, "iguazio-devops/k8s", "v3ctl-${github.TAG_VERSION}-linux-amd64")
                                    }
                                }
                            }
                        },
                    )
                }
            }
        }
    }
}

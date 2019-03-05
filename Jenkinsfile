label = "${UUID.randomUUID().toString()}"
git_project = "v3ctl"
git_project_user = "v3io"
git_deploy_user_token = "iguazio-prod-git-user-token"
git_deploy_user_private_key = "iguazio-prod-git-user-private-key"

podTemplate(label: "${git_project}-${label}", inheritFrom: "jnlp-docker-golang") {
    node("${git_project}-${label}") {
        pipelinex = library(identifier: 'pipelinex@reduction', retriever: modernSCM(
                [$class       : 'GitSCMSource',
                 credentialsId: git_deploy_user_private_key,
                 remote       : "git@github.com:iguazio/pipelinex.git"])).com.iguazio.pipelinex
        common.notify_slack {
            withCredentials([
                    string(credentialsId: git_deploy_user_token, variable: 'GIT_TOKEN')
            ]) {
                github.init_project(git_project, git_project_user, GIT_TOKEN) {
                    def RELEASE_ID
                    stage('prepare sources') {
                        container('jnlp') {
                            RELEASE_ID = github.get_release_id(git_project, git_project_user, "${github.TAG_VERSION}", GIT_TOKEN)

                            dir("${github.BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                git(changelog: false, credentialsId: git_deploy_user_private_key, poll: false, url: "git@github.com:${git_project_user}/${git_project}.git")
                                common.shellc("git checkout ${github.TAG_VERSION}")
                            }
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
                                dir("${github.BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                    common.shellc("V3CTL_TAG=${github.TAG_VERSION} GOARCH=amd64 GOOS=linux make v3ctl-bin")
                                }
                            }
                        },
                        'build darwin binaries': {
                            container('golang') {
                                dir("${github.BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                    common.shellc("V3CTL_TAG=${github.TAG_VERSION} GOARCH=amd64 GOOS=darwin make v3ctl-bin")
                                }
                            }
                        },
                        'build windows binaries': {
                            container('golang') {
                                dir("${github.BUILD_FOLDER}/src/github.com/v3io/${git_project}") {
                                    common.shellc("V3CTL_TAG=${github.TAG_VERSION} GOARCH=amd64 GOOS=windows make v3ctl-bin")
                                }
                            }
                        }
                    )

                    parallel(
                        'upload linux binaries': {
                            container('jnlp') {
                                github.upload_asset(git_project, git_project_user, "v3ctl-${github.TAG_VERSION}-linux-amd64", RELEASE_ID, GIT_TOKEN)
                            }
                        },
                        'upload linux binaries artifactory': {
                            container('jnlp') {
                                withCredentials([
                                        string(credentialsId: pipelinex.PackagesRepo.ARTIFACTORY_IGUAZIO[2], variable: 'PACKAGES_ARTIFACTORY_PASSWORD')
                                ]) {
                                    common.upload_file_to_artifactory(pipelinex.PackagesRepo.ARTIFACTORY_IGUAZIO[0], pipelinex.PackagesRepo.ARTIFACTORY_IGUAZIO[1], PACKAGES_ARTIFACTORY_PASSWORD, "iguazio-devops/k8s", "v3ctl-${github.TAG_VERSION}-linux-amd64")
                                }
                            }
                        },
                        'upload darwin binaries': {
                            container('jnlp') {
                                github.upload_asset(git_project, git_project_user, "v3ctl-${github.TAG_VERSION}-darwin-amd64", RELEASE_ID, GIT_TOKEN)
                            }
                        },
                        'upload windows binaries': {
                            container('jnlp') {
                                github.upload_asset(git_project, git_project_user, "v3ctl-${github.TAG_VERSION}-windows-amd64", RELEASE_ID, GIT_TOKEN)
                            }
                        }
                    )
                }
            }
        }
    }
}

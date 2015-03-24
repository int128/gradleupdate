package service

import groovy.transform.TupleConstructor
import groovyx.gaelyk.GaelykBindings
import infrastructure.GitHub
import model.GitHubRepository

@TupleConstructor
@GaelykBindings
class GitHubRepositoryService {

    final GitHub gitHub = new GitHub()

    List<GitHubRepository> listPullRequestOnStableRelease() {
        datastore.execute {
            select all from 'GitHubRepository'
            where pullRequestOnStableRelease == true
        }.collect {
            it as GitHubRepository
        }
    }

    GitHubRepository query(String fullName) {
        assert fullName

        def repo = gitHub.getRepository(fullName)
        if (repo.permissions?.admin) {
            GitHubRepository.get(fullName) ?: new GitHubRepository(fullName: fullName)
        } else {
            null
        }
    }

    GitHubRepository save(GitHubRepository entity) {
        assert entity.fullName

        def repo = gitHub.getRepository(entity.fullName)
        if (repo.permissions?.admin) {
            entity.save()
            entity
        } else {
            null
        }
    }

}

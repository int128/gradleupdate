package service

import groovy.transform.TupleConstructor
import infrastructure.GitHub
import model.GitHubRepository

@TupleConstructor
class GitHubRepositoryService {

    final GitHub gitHub = new GitHub()

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

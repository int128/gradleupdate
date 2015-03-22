package service

import groovy.transform.Immutable

@Immutable
class GitHubRepositoryMetadata {
    String fullName
    boolean gradleProject
    String gradleVersion
    boolean autoUpdate
}

package service

import groovy.transform.Immutable

@Immutable
class GitHubRepositoryMetadata {
    String fullName
    boolean admin

    boolean gradleProject
    String gradleVersion

    // null if user has no permission
    Boolean autoUpdate
}

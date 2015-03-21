package model

import groovyx.gaelyk.datastore.Entity
import groovyx.gaelyk.datastore.Indexed
import groovyx.gaelyk.datastore.Key

@Entity
class GitHubRepository {
    @Key String fullName
    @Indexed boolean autoUpdate
}

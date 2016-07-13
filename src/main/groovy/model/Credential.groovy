package model

import groovyx.gaelyk.datastore.Entity
import groovyx.gaelyk.datastore.Key
import groovyx.gaelyk.datastore.Unindexed

@Entity
class Credential {

    @Key String service
    @Unindexed String secret

    static enum CredentialKey {
        GitHubToken,
        GitHubClientId,
        GitHubClientKey,
        GitHubWebHookSecret,
    }

    static Credential get(CredentialKey key) {
        final credential = Credential.get(key.name())
        if (credential) {
            credential
        } else {
            throw new IllegalStateException("Credential should be set by initial setup: ${key.name()}")
        }
    }

}

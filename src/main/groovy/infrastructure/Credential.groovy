package infrastructure

import groovy.transform.PackageScope

@PackageScope
class Credential {

    static githubToken
    static githubClientId
    static githubClientSecret

    static {
        def properties = new Properties()
        properties.load(Credential.getResourceAsStream('/credential.properties'))
        properties.each { k, v -> Credential."$k" = "$v" }
    }

}

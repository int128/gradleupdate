package service

import groovyx.gaelyk.GaelykBindings
import model.Credential

@GaelykBindings
class CredentialRepository {

    Credential find(String serviceName) {
        def credential = Credential.get(serviceName)
        if (credential == null) {
            credential = new Credential()
            credential.service = serviceName
            credential.save()
        }
        credential
    }

}

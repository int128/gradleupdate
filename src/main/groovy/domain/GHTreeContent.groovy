package domain

import groovy.transform.Immutable

@Immutable
class GHTreeContent {

    String path
    String mode
    String base64encoded

}

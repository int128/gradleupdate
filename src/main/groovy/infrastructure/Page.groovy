package infrastructure

import groovy.transform.Canonical
import wslite.rest.Response

@Canonical
class Page {

    final Map<String, String> rel
    final List current

    def Page(Response response) {
        rel = parseLinkHeader(response.headers.link as String)
        current = response.json
    }

    static Map<String, String> parseLinkHeader(String linkHeader) {
        linkHeader?.findAll(~/<(.+?)>; rel="(.+?)"/) { [it[2], it[1]] }?.collectEntries() ?: [:]
    }

}

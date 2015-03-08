['/feed/all', '/feed/stable'].each { uri ->
    memcache.clearCacheForUri(uri)
    log.info("Purged the page cache of $uri")
}

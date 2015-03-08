import com.google.appengine.api.memcache.MemcacheService
import com.google.appengine.api.memcache.MemcacheServiceFactory
import groovy.transform.CompileDynamic
import groovy.transform.CompileStatic
import groovy.xml.MarkupBuilder

@CompileStatic
class FeedService {

    static enum CacheKey {
        all,
        stable
    }

    final GradleService service = new GradleService()

    final MemcacheService memcache = MemcacheServiceFactory.memcacheService

    String getAllVersions() {
        getCacheOrFetch(CacheKey.all) {
            feed(service.fetchAllVersions(), 'Gradle Releases')
        }
    }

    String getStableVersions() {
        getCacheOrFetch(CacheKey.stable) {
            feed(service.fetchStableVersions(), 'Gradle Releases (Stable)')
        }
    }

    void purgeCache() {
        memcache.deleteAll([CacheKey.all, CacheKey.stable])
    }

    private <T> T getCacheOrFetch(Object key, Closure<T> fetcher) {
        def cached = memcache.get(key)
        if (cached) {
            cached as T
        } else {
            def fetched = fetcher()
            memcache.put(key, fetched)
            fetched
        }
    }

    @CompileDynamic
    static feed(List versions, String titleOfFeed) {
        def writer = new StringWriter()
        new MarkupBuilder(writer).feed {
            title(titleOfFeed)
            link(href: 'https://gradleupdate.appspot.com')
            id('https://gradleupdate.appspot.com')
            author('Gradle Update')
            updated()

            versions.each { version ->
                entry {
                    title(version.version)
                    link(href: version.downloadUrl)
                    id(version.downloadUrl)
                    updated(datetime(version.buildTime))
                    summary("Gradle $version.version")

                    raw {
                        buildTime(version.buildTime)
                        current(version.current)
                        snapshot(version.snapshot)
                        nightly(version.nightly)
                        activeRc(version.activeRc)
                        rcFor(version.rcFor)
                        broken(version.broken)
                    }
                }
            }
        }
        writer.toString()
    }

    static datetime(String datetimeInJson) {
        // 20120912104602+0000
        // 2005-07-31T12:29:29Z
        Date.parse('yyyyMMddHHmmssZ', datetimeInJson).format("yyyy-MM-dd'T'HH:mm:ss'Z'", TimeZone.getTimeZone('UTC'))
    }

}

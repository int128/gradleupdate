import groovy.transform.CompileDynamic
import groovy.transform.CompileStatic
import groovy.transform.TupleConstructor
import groovy.xml.MarkupBuilder

import javax.servlet.http.HttpServletRequest
import javax.servlet.http.HttpServletResponse

@CompileStatic
@TupleConstructor
class FeedController {

    HttpServletRequest request
    HttpServletResponse response

    def service = new GradleService()

    void stableVersions() {
        feed(service.fetchStableVersions(), 'Gradle Stable Versions')
    }

    @CompileDynamic
    void feed(List versions, String titleOfFeed) {
        response.setContentType('text/xml')
        response.setCharacterEncoding('UTF-8')
        response.writer.withWriter { writer ->
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
        }
    }

    static datetime(String datetimeInJson) {
        // 20120912104602+0000
        // 2005-07-31T12:29:29Z
        Date.parse('yyyyMMddHHmmssZ', datetimeInJson).format("yyyy-MM-dd'T'HH:mm:ss'Z'", TimeZone.getTimeZone('UTC'))
    }

}

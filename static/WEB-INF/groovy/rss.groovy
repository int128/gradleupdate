import gradle.VersionWatcher

final watcher = new VersionWatcher()
final releases = watcher.rcReleasesWithFixedIssues()

response.contentType = 'application/xml'

html.rss(version: '2.0') {
    channel {
        title('Gradle Releases')
        link('https://gradleupdate.appspot.com')
        description('RSS feed of Gradle releases with fixed issues')
        lastBuildDate(formatTime(new Date()))

        releases.each { release ->
            final buildTime = formatTime(parseTime(release.buildTime))
            item {
                title("Gradle $release.version")
                pubDate(buildTime.toString())
                guid(release.version)
                description {
                    p("Gradle $release.version on $buildTime")
                    ul {
                        release.fixedIssues?.each { issue ->
                            li {
                                a(href: issue.link, issue.key)
                                span(issue.summary)
                                span("($issue.type)")
                            }
                        }
                    }
                }
            }
        }
    }
}

static parseTime(String datetimeInJson) {
    // 20120912104602+0000
    Date.parse('yyyyMMddHHmmssZ', datetimeInJson)
}

static formatTime(Date time) {
    // RFC822
    time.format("EEE, dd MMM yyyy HH:mm:ss 'Z'", TimeZone.getTimeZone('UTC'))
}

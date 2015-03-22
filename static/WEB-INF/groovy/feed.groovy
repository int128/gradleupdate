import infrastructure.GradleRegistry

final registry = new GradleRegistry()

switch (params.filter) {
    case 'stable':
        feed(versions: registry.getStableReleasesWithFixedIssues(), title: 'Gradle Releases')
        break

    case 'rc':
        feed(versions: registry.getReleaseCandidateReleasesWithFixedIssues(), title: 'Gradle Releases including Candidates')
        break

    default:
        response.sendError(404)
        break
}

def feed(Map data) {
    response.contentType = 'application/xml'

    html.feed(xmlns: 'http://www.w3.org/2005/Atom') {
        title(data.title)
        link(href: 'https://gradleupdate.appspot.com')
        id("https://gradleupdate.appspot.com/feed/${params.filter}")
        author {
            name('gradleupdate')
        }
        updated(formatTime(new Date()))

        data.versions.each { version ->
            entry {
                title("Gradle $version.version")
                link(href: version.downloadUrl)
                id(version.downloadUrl)
                updated(formatTime(parseTime(version.buildTime)))
                summary("Gradle $version.version build $version.buildTime")

                content(type: 'xhtml') {
                    div {
                        p("Gradle $version.version build $version.buildTime")
                        ul {
                            version.fixedIssues?.each { issue ->
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
}

static parseTime(String datetimeInJson) {
    // 20120912104602+0000
    Date.parse('yyyyMMddHHmmssZ', datetimeInJson)
}

static formatTime(Date time) {
    // 2005-07-31T12:29:29Z
    time.format("yyyy-MM-dd'T'HH:mm:ss'Z'", TimeZone.getTimeZone('UTC'))
}

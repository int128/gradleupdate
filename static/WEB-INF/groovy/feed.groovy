final service = new GradleService()

switch (params.filter) {
    case 'stable':
        feed(versions: service.fetchStableVersions(), title: 'Gradle Releases (Stable)')
        break

    case 'all':
        feed(versions: service.fetchAllVersions(), title: 'Gradle Releases')
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
                title(version.version)
                link(href: version.downloadUrl)
                id(version.downloadUrl)
                updated(formatTime(parseTime(version.buildTime)))
                summary("Gradle $version.version")

                content(type: 'application/xml') {
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

static parseTime(String datetimeInJson) {
    // 20120912104602+0000
    Date.parse('yyyyMMddHHmmssZ', datetimeInJson)
}

static formatTime(Date time) {
    // 2005-07-31T12:29:29Z
    time.format("yyyy-MM-dd'T'HH:mm:ss'Z'", TimeZone.getTimeZone('UTC'))
}

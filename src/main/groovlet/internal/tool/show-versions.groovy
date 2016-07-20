import domain.GHSession
import domain.GradleWrapperTemplateRepository
import entity.CurrentGradleVersion
import infrastructure.GradleRegistry

final gradleRegistry = new GradleRegistry()

final template = GradleWrapperTemplateRepository.get(GHSession.defaultToken())

html.html {
    body {
        h1('Gradle Versions')

        h2('Gradle Registry')
        p(gradleRegistry.fetchCurrentStableRelease().toString())
        p(gradleRegistry.fetchCurrentReleaseCandidateRelease().toString())

        h2('Datastore (polling every 1 hour)')
        CurrentGradleVersion.findAll().each { entity ->
            p(entity.toString())
        }

        h2('Template Repository')
        p { a(href: template.repository.rawJson.html_url, template.repository.fullName) }
        p(template.repository.defaultBranch.name)
        p(template.gradleWrapper.version)
    }
}

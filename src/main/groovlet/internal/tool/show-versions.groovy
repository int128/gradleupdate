import gradle.TemplateRepository
import infrastructure.GradleRegistry
import model.CurrentGradleVersion

final gradleRegistry = new GradleRegistry()

final templateRepository = new TemplateRepository()

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
        p { a(href: templateRepository.htmlUrl, templateRepository.fullName) }
        p('master')
        p(templateRepository.fetchGradleWrapperVersion('master'))
    }
}

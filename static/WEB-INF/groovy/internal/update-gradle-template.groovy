import gradle.TemplateRepository
import infrastructure.GitHub

assert params.gradleVersion

final gitHub = new GitHub()
final templateRepository = new TemplateRepository(gitHub)
templateRepository.bumpTemplate(params.gradleVersion)

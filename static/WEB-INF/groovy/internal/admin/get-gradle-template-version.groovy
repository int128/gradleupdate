import gradle.TemplateRepository
import infrastructure.GitHub

final gitHub = new GitHub()
final templateRepository = new TemplateRepository(gitHub)

response.contentType = 'text/plain'
println templateRepository.queryGradleWrapperVersion()

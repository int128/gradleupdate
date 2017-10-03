package org.hidetake.gradleupdate.app

import org.hidetake.gradleupdate.domain.RepositoryPath
import org.springframework.http.HttpHeaders
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.servlet.ModelAndView
import javax.servlet.http.HttpServletResponse

@Controller
@RequestMapping("/{owner}/{repo}")
class RepositoryController(val service: GradleUpdateService) {
    @GetMapping("status.svg")
    fun badge(@PathVariable owner: String, @PathVariable repo: String, response: HttpServletResponse) =
        BadgeSvg.render(service.getGradleWrapperVersionStatus(RepositoryPath(owner, repo))).also {
            // https://docs.spring.io/spring-security/site/docs/current/reference/html/headers.html#headers-cache-control
            response.setHeader(HttpHeaders.CACHE_CONTROL, "public, max-age=30")
        }

    @GetMapping("status")
    fun html(@PathVariable owner: String, @PathVariable repo: String) =
        ModelAndView("status", mapOf(
            "repository" to service.getRepository(RepositoryPath(owner, repo)),
            "pullRequest" to service.findPullRequestForUpdate(RepositoryPath(owner, repo))
        ))

    @PostMapping("update")
    fun update(@PathVariable owner: String, @PathVariable repo: String): String {
        service.createPullRequestForLatestGradleWrapper(RepositoryPath(owner, repo))
        return "redirect:status"
    }
}

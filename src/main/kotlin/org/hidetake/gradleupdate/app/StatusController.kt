package org.hidetake.gradleupdate.app

import org.hidetake.gradleupdate.domain.RepositoryPath
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.servlet.ModelAndView

@Controller
class StatusController(val service: GradleUpdateService) {
    @GetMapping("/{owner}/{repo}/status.svg")
    fun badge(@PathVariable owner: String, @PathVariable repo: String) =
        BadgeSvg.render(service.getGradleWrapperVersionStatus(RepositoryPath(owner, repo)))

    @GetMapping("/{owner}/{repo}/status")
    fun html(@PathVariable owner: String, @PathVariable repo: String) =
        ModelAndView("status", mapOf(
            "repository" to service.getRepository(RepositoryPath(owner, repo)),
            "pullRequest" to service.findPullRequestForUpdate(RepositoryPath(owner, repo))
        ))
}

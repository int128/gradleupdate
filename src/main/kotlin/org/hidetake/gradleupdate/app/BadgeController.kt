package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable

@Controller
class BadgeController(val service: GradleUpdateService) {
    @GetMapping("/{owner}/{repo}/status.svg")
    fun get(@PathVariable owner: String, @PathVariable repo: String) =
        BadgeSvg.render(service.getGradleWrapperVersionStatus("$owner/$repo"))
}

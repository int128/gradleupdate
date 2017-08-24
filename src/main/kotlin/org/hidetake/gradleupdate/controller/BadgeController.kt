package org.hidetake.gradleupdate.controller

import org.hidetake.gradleupdate.service.GradleUpdateService
import org.hidetake.gradleupdate.view.BadgeSvg
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import java.io.FileNotFoundException

@Controller
class BadgeController(val gradleUpdateService: GradleUpdateService) {
    @GetMapping("/{owner}/{repo}/status.svg")
    fun get(@PathVariable owner: String, @PathVariable repo: String) =
        BadgeSvg.render(gradleUpdateService.getStatus("$owner/$repo"))

    @ExceptionHandler(FileNotFoundException::class)
    fun notFound() = BadgeSvg.notFound()
}

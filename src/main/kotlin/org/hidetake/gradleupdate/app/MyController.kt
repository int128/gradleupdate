package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.servlet.ModelAndView

@Controller
@RequestMapping("/my")
class MyController(private val service: UserService) {
    @GetMapping
    fun get() = ModelAndView(
        "/my",
        mapOf(
            "user" to service.getLoginUser()
        )
    )
}

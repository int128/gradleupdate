package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.web.servlet.ModelAndView
import org.springframework.web.servlet.support.ServletUriComponentsBuilder

@Controller
class LoginController(private val service: LoginService) {
    @GetMapping("/login")
    fun authorize() = ModelAndView(
        "redirect:${service.getRedirectURL()}",
        service.createAuthorizationParameters(
            ServletUriComponentsBuilder.fromCurrentContextPath()
                .path("/login/auth")
                .build().toUriString()))

    @GetMapping("/login/auth")
    fun continueAuthorization(@RequestParam state: String, @RequestParam code: String): String {
        service.continueAuthorization(state, code)
        return "redirect:/my"
    }
}

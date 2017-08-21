package example

import org.springframework.stereotype.Controller
import org.springframework.ui.Model
import org.springframework.web.bind.annotation.GetMapping
import java.time.LocalDateTime

@Controller
class HomeController {
    @GetMapping("/hello")
    fun index(model: Model): String {
        model.addAttribute("currentDateTime", LocalDateTime.now())
        return "hello"
    }
}

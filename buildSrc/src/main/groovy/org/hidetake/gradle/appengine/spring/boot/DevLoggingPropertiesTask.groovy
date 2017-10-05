package org.hidetake.gradle.appengine.spring.boot

import org.gradle.api.DefaultTask
import org.gradle.api.tasks.Input
import org.gradle.api.tasks.OutputFile
import org.gradle.api.tasks.TaskAction

/**
 * A task to configure logging format of the Dev Server.
 */
class DevLoggingPropertiesTask extends DefaultTask {
  @Input
  String pattern = '%1$tY-%1$tm-%1$td %1$tH:%1$tM:%1$tS.%1$tL %4$s --- %3$s : %5$s %6$s%n'

  @Input
  String level = 'INFO'

  @OutputFile
  File loggingProperties = new File(temporaryDir, 'logging.properties')

  @TaskAction
  void createFile() {
    loggingProperties.text = """\
handlers=java.util.logging.ConsoleHandler
.level=$level
java.util.logging.ConsoleHandler.level=$level
java.util.logging.ConsoleHandler.formatter=java.util.logging.SimpleFormatter
java.util.logging.SimpleFormatter.format=$pattern
"""
  }
}

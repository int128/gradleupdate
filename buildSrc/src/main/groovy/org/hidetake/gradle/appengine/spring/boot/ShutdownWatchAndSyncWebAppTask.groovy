package org.hidetake.gradle.appengine.spring.boot

import org.gradle.api.DefaultTask
import org.gradle.api.tasks.Input
import org.gradle.api.tasks.TaskAction

class ShutdownWatchAndSyncWebAppTask extends DefaultTask {
  @Input
  WatchAndSyncWebAppTask watchAndSyncWebAppTask

  @TaskAction
  void shutdownThreads() {
    watchAndSyncWebAppTask.threads*.interrupt()
  }
}

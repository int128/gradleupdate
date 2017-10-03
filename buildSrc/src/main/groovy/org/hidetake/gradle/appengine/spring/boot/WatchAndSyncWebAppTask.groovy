package org.hidetake.gradle.appengine.spring.boot

import org.gradle.api.DefaultTask
import org.gradle.api.tasks.TaskAction
import org.hidetake.gradle.appengine.spring.boot.WatchAndSync.SyncSpec

class WatchAndSyncWebAppTask extends DefaultTask {
  final List<Thread> threads = []

  @TaskAction
  void startThread() {
    threads << new Thread(new WatchAndSync(*forResources(), forWebApp())).start()
    project.gradle.buildFinished {
      threads*.interrupt()
    }
  }

  private SyncSpec forWebApp() {
    assert project.appengine.stage.sourceDirectory
    new SyncSpec(project.webAppDir, project.file(project.appengine.stage.sourceDirectory))
  }

  private List<SyncSpec> forResources() {
    assert project.appengine.stage.sourceDirectory
    final webInfClasses = project.file("${project.appengine.stage.sourceDirectory}/WEB-INF/classes")
    project.sourceSets.main.resources.srcDirs.collect { File srcDir ->
      new SyncSpec(srcDir, webInfClasses)
    }
  }
}

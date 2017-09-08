package org.hidetake.gradle.appengine.spring.boot

class DotEnv {
  static Properties loadOrEmpty(File file) {
    final properties = new Properties()
    if (file?.canRead()) {
      file.withReader { reader ->
        properties.load(reader)
      }
    }
    properties
  }
}

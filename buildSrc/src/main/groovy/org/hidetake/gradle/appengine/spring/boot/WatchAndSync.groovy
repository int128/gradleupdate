package org.hidetake.gradle.appengine.spring.boot

import groovy.transform.Canonical

import java.nio.file.*
import java.nio.file.WatchEvent.Kind

import static java.nio.file.StandardCopyOption.REPLACE_EXISTING
import static java.nio.file.StandardWatchEventKinds.*

@Canonical
class WatchAndSync implements Runnable {
  final Path sourceDirectory
  final Path targetDirectory
  final Map<WatchKey, Path> watchKeyPathMap = [:]

  @Override
  void run() {
    final watchService = FileSystems.default.newWatchService()
    registerWatchService(watchService)

    while (true) {
      try {
        final watchKey = watchService.take()
        watchKey.pollEvents().each { event ->
          if (event.kind() != OVERFLOW) {
            final directory = watchKeyPathMap.get(watchKey)
            try {
              processEvent(event.kind(), directory, event.context() as Path)
            } catch (IOException e) {
              log(e.toString())
            }
          }
        }
        watchKey.reset()
      } catch (InterruptedException ignore) {
        break
      }
    }
  }

  private void registerWatchService(WatchService watchService) {
    final watchKey = sourceDirectory.register(watchService, ENTRY_CREATE, ENTRY_MODIFY, ENTRY_DELETE)
    watchKeyPathMap.put(watchKey, Paths.get(''))
    log("Watching directory: $sourceDirectory")

    sourceDirectory.eachDirRecurse { Path directory ->
      final childWatchKey = directory.register(watchService, ENTRY_CREATE, ENTRY_MODIFY, ENTRY_DELETE)
      final relativePath = sourceDirectory.relativize(directory)
      watchKeyPathMap.put(childWatchKey, relativePath)
      log("Watching directory: $directory -> $relativePath")
    }
  }

  private void processEvent(Kind kind, Path directory, Path path) {
    final sourceFile = sourceDirectory.resolve(directory).resolve(path).toFile()
    final targetFile = targetDirectory.resolve(directory).resolve(path).toFile()

    switch (kind) {
      case ENTRY_CREATE:
        if (sourceFile.directory) {
          log("Creating directory $targetFile")
          targetFile.mkdirs()
        } else {
          log("Creating file $targetFile")
          Files.copy(sourceFile.toPath(), targetFile.toPath(), REPLACE_EXISTING)
        }
        break

      case ENTRY_MODIFY:
        if (sourceFile.file) {
          log("Updating file $targetFile")
          Files.copy(sourceFile.toPath(), targetFile.toPath(), REPLACE_EXISTING)
        }
        break

      case ENTRY_DELETE:
        if (targetFile.file) {
          log("Removing file $targetFile")
          targetFile.delete()
        }
        break
    }
  }

  private static void log(String message) {
    println("[${Thread.currentThread().name}] $message")
  }
}

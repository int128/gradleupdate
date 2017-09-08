package org.hidetake.gradle.appengine.spring.boot

import groovy.transform.Immutable

import java.nio.file.*

import static java.nio.file.StandardCopyOption.REPLACE_EXISTING
import static java.nio.file.StandardWatchEventKinds.*

class WatchAndSync implements Runnable {
  final List<SyncSpec> syncSpecs
  private final Map<WatchKey, WatchContext> watchContextMap = [:]

  WatchAndSync(SyncSpec... syncSpecs) {
    this.syncSpecs = syncSpecs
  }

  @Immutable(knownImmutableClasses = [File])
  static class SyncSpec {
    final File sourceDirectory
    final File targetDirectory

    Path getSourcePath() {
      sourceDirectory.toPath()
    }
    Path getTargetPath() {
      targetDirectory.toPath()
    }
  }

  @Immutable(knownImmutableClasses = [Path])
  static class WatchContext {
    final SyncSpec syncSpec
    final Path relativePath
  }

  @Override
  void run() {
    final watchService = FileSystems.default.newWatchService()
    syncSpecs.each { syncSpec ->
      registerWatchService(watchService, syncSpec)
    }

    while (true) {
      try {
        final watchKey = watchService.take()
        watchKey.pollEvents().findAll { it.kind() != OVERFLOW }.each { event ->
          try {
            final context = watchContextMap.get(watchKey)
            processWatchEvent(event as WatchEvent<Path>, context)
          } catch (IOException e) {
            log(e.toString())
          }
        }
        watchKey.reset()
      } catch (InterruptedException ignore) {
        break
      }
    }
  }

  private void registerWatchService(WatchService watchService, SyncSpec syncSpec) {
    final watchKey = syncSpec.sourcePath.register(watchService, ENTRY_CREATE, ENTRY_MODIFY, ENTRY_DELETE)
    watchContextMap.put(watchKey, new WatchContext(syncSpec, Paths.get('')))
    log("Watching directory: ${syncSpec.sourceDirectory}")

    syncSpec.sourcePath.eachDirRecurse { Path directory ->
      final childWatchKey = directory.register(watchService, ENTRY_CREATE, ENTRY_MODIFY, ENTRY_DELETE)
      final relativePath = syncSpec.sourcePath.relativize(directory)
      watchContextMap.put(childWatchKey, new WatchContext(syncSpec, relativePath))
      log("Watching directory: $directory -> $relativePath")
    }
  }

  private static void processWatchEvent(WatchEvent<Path> event, WatchContext context) {
    final sourceFile = context.syncSpec.sourcePath.resolve(context.relativePath).resolve(event.context()).toFile()
    final targetFile = context.syncSpec.targetPath.resolve(context.relativePath).resolve(event.context()).toFile()

    switch (event.kind()) {
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

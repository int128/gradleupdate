import static util.RequestUtil.relativePath

defaultQueue.add(url: relativePath(request, '1-stable.groovy'))
defaultQueue.add(url: relativePath(request, '1-rc.groovy'))
